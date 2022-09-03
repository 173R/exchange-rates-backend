package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"time"
)

// Информация о том, какие поля могут быть в JWT:
// https://ru.wikipedia.org/wiki/JSON_Web_Token#Полезная_нагрузка

type signResult struct {
	// Токен доступа.
	Token string
	// Дата истечения срока годности токена.
	ExpiresAt time.Time
}

type UserToken struct {
	AccessToken  *signResult
	RefreshToken *signResult
}

// Способ подписания пэйлоада который мы везде используем.
var signMethod = jwt.SigningMethodHS256

// CreateUserToken создает токен для использования пользователем.
func CreateUserToken(uid models.UserId, lang language.Lang) (*UserToken, error) {
	// Генерируем токен доступа.
	accessToken, err := sign(&UserAccessToken{
		Uid:      uid,
		Language: lang,
	}, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	// Генерируем токен для обновления токена доступа.
	refreshToken, err := sign(map[string]interface{}{}, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &UserToken{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// DecodeUserAccessToken декодирует пользовательский токен доступа.
func DecodeUserAccessToken(token string) (*UserAccessToken, error) {
	return decode[UserAccessToken](token)
}

// DecodeRefreshToken декодирует токен обновления.
func DecodeRefreshToken(token string) error {
	_, err := decode[interface{}](token)
	return err
}

// Подписывает указанный payload.
func sign(payload interface{}, expIn time.Duration) (*signResult, error) {
	// Сериализуем пэйлоад в массив байт.
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Конвертируем массив байт в карту.
	var pmap map[string]interface{}
	if err := json.Unmarshal(bytes, &pmap); err != nil {
		return nil, err
	}

	expAt := time.Now().Add(expIn)

	// Проставляем дату выдачи токена.
	pmap["iat"] = time.Now().Unix()
	// Добавляем дату истечения срока годности токена.
	pmap["exp"] = expAt.Unix()

	// Подписываем пэйлоад.
	t, err := jwt.NewWithClaims(signMethod, jwt.MapClaims(pmap)).SignedString(configs.Jwt.Secret)
	if err != nil {
		return nil, err
	}

	return &signResult{Token: t, ExpiresAt: expAt}, nil
}

// Проверяет токен на валидность, а также извлекает из него контент.
func decode[T interface{}](seq string) (*T, error) {
	// Парсим токен.
	token, err := jwt.Parse(seq, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм шифрования на соответствие нашему.
		alg, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if alg != signMethod {
			return nil, errors.New("unexpected alg found")
		}
		return configs.Jwt.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	// Если токен невалиден, вернуть ошибку.
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("payload is invalid")
	}

	// Приводим payload к ожидаемому типу.
	var data T
	if err := mapstructure.Decode(payload, &data); err != nil {
		return nil, errors.New("payload had unexpected format")
	}
	return &data, nil
}
