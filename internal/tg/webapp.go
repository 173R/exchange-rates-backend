package tg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/wolframdeus/exchange-rates-backend/internal/language"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type UserId int64

type InitData struct {
	QueryId      string `mapstructure:"query_id"`
	User         *User  `mapstructure:"user"`
	Receiver     *User  `mapstructure:"receiver"`
	StartParam   string `mapstructure:"start_param"`
	CanSendAfter int    `mapstructure:"can_send_after"`
	AuthDateRaw  int    `mapstructure:"auth_date"`
	Hash         string `mapstructure:"hash"`
}

func (d *InitData) AuthDate() time.Time {
	return time.Unix(int64(d.AuthDateRaw), 0)
}

type User struct {
	Id           UserId        `mapstructure:"id"`
	IsBot        bool          `mapstructure:"is_bot"`
	FirstName    string        `mapstructure:"first_name"`
	LastName     string        `mapstructure:"last_name"`
	Username     string        `mapstructure:"username"`
	LanguageCode language.Lang `mapstructure:"language_code"`
	IsPremium    bool          `mapstructure:"is_premium"`
	PhotoUrl     string        `mapstructure:"photo_url"`
}

// Срок действия параметров запуска.
const paramsExpireIn = 10 * 24 * time.Hour

// ValidateInitData валидирует параметры запуска, которые были переданы из
// клиентского приложения. Ожидается, что initData будет равен значению
// window.Telegram.WebApp.initData
// https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app
func ValidateInitData(initData, token string) (bool, error) {
	// Ожидаем, что переданные данные это query параметры.
	params, err := url.ParseQuery(initData)
	if err != nil {
		return false, err
	}

	// Храним тут список пар ключ-значение.
	pairs := make([]string, 0, len(params))

	// Храним найденный хеш и дату создания параметров.
	var hash string
	var authDate time.Time

	// Пробегаемся по всем полям и добавляем их в pairs.
	for k, v := range params {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err == nil {
				authDate = time.Unix(int64(i), 0)
			}
		}
		pairs = append(pairs, k+"="+v[0])
	}

	// Хеш обязателен.
	if hash == "" {
		return false, errors.New("hash is empty")
	}

	// Дата создания параметров запуска обязательна.
	if authDate.IsZero() {
		return false, errors.New("auth_date is empty")
	}

	// Параметры запуска валидны лишь в течение определенного времени.
	if authDate.Add(paramsExpireIn).Before(time.Now()) {
		return false, errors.New("init data is expired")
	}

	// Сортируем по возрастанию ключа.
	sort.Strings(pairs)

	// Создаем корректный data check string.
	imploded := strings.Join(pairs, "\n")

	// Создаем подписи HMAC_SHA256.
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(token))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(imploded))

	return hex.EncodeToString(impHmac.Sum(nil)) == hash, nil
}

// ParseInitData конвертирует параметры запуска представленные в виде
// query-параметров в известную структуру - InitData.
func ParseInitData(initData string) (*InitData, error) {
	// Ожидаем, что переданные данные это query параметры.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{}, len(q))
	for k, v := range q {
		val := v[0]

		if val == "true" || val == "false" {
			m[k] = val == "true"
		} else if asInt, err := strconv.ParseInt(val, 10, 64); err == nil {
			m[k] = asInt
		} else {
			// Сначала пытаемся спарсить строку как JSON.
			var j map[string]interface{}
			if err := json.Unmarshal([]byte(val), &j); err == nil {
				m[k] = j
			} else {
				// В противном случае это просто строка.
				m[k] = val
			}
		}
	}

	// Создаем экземпляр InitData из карты.
	var mapped InitData
	if err := mapstructure.Decode(m, &mapped); err != nil {
		return nil, err
	}

	// Если язык не является нам известным, мы поставим язык по умолчанию.
	if !mapped.User.LanguageCode.Known() {
		mapped.User.LanguageCode = language.Default
	}
	return &mapped, nil
}
