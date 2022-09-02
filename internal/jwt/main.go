package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
)

// Информация о том, какие поля могут быть в JWT:
// https://ru.wikipedia.org/wiki/JSON_Web_Token#Полезная_нагрузка

// Способ подписания пэйлоада который мы везде используем.
var signMethod = jwt.SigningMethodHS256

// SignUserAccessToken подписывает пользовательский токен доступа.
func SignUserAccessToken(uid models.UserId) (string, error) {
	return sign(&UserAccessToken{Uid: uid})
}

// DecodeUserAccessToken декодирует пользовательский токен доступа.
func DecodeUserAccessToken(token string) (*UserAccessToken, error) {
	return decode[UserAccessToken](token)
}

// Подписывает указанный payload.
func sign(payload interface{}) (string, error) {
	// Сериализуем пэйлоад в массив байт.
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Конвертируем массив байт в карту.
	var pmap map[string]interface{}
	if err := json.Unmarshal(bytes, &pmap); err != nil {
		return "", err
	}

	return jwt.
		NewWithClaims(signMethod, jwt.MapClaims(pmap)).
		SignedString(configs.Jwt.Secret)
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
