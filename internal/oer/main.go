package oer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
)

//// Структура со всеми возможными полями ответа.
//type response[D interface{}] struct {
//	Status      int    `json:"status"`
//	Data        *D     `json:"data"`
//	Error       bool   `json:"error"`
//	Message     string `json:"message"`
//	Description string `json:"description"`
//}

// OER позволяет работать с API Open Exchange Rates.
type OER struct {
	// Ключ для доступа к API.
	appId string
}

// FetchUsage получает информацию об использовании текущего ключа.
func (o *OER) FetchUsage() (*usage, error) {
	return sendHttpRequest[usage]("usage.json", o.appId)
}

// FetchLatest получает актуальную информацию о текущем обмене курса валют.
func (o *OER) FetchLatest() (*latest, error) {
	return sendHttpRequest[latest]("latest.json", o.appId)
}

func sendHttpRequest[T interface{}](path string, appId string) (result *T, err error) {
	// Отправляем запрос.
	res, err := http.Get(
		fmt.Sprintf("https://openexchangerates.org/api/%s?app_id=%s", path, appId),
	)
	if err != nil {
		return nil, err
	}

	// Не забываем закрыть поток.
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Println(closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// Пытаемся обработать запрос как успешный.
	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// Наличие поля "status" говорит о том, что это шаблонный ответ.
	if s, ok := r["status"]; ok {
		if status, ok := s.(float64); ok {
			// Возвращаем ошибку, если статус отличен от успешного.
			if status != 200 {
				msg := fmt.Sprintf("unknown error occurred: %f", status)

				if d, ok := r["description"]; ok {
					if desc, ok := d.(string); ok {
						msg = desc
					}
				}
				return nil, errors.New(msg)
			}

			// Извлекаем данные из поля "data".
			if d, ok := r["data"]; ok {
				var data T
				if err := mapstructure.Decode(d, &data); err == nil {
					return &data, nil
				}
			}
		}
	} else {
		var data T
		if err := mapstructure.Decode(r, &data); err == nil {
			return &data, nil
		}
	}

	// FIXME
	return nil, errors.New("TODO")
}

//func sendHttpRequest[T interface{}](path string, appId string) (result *T, err error) {
//	// Отправляем запрос.
//	res, err := http.Get(
//		fmt.Sprintf("https://openexchangerates.org/api/%s?app_id=%s", path, appId),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	// Не забываем закрыть поток.
//	defer func() {
//		if closeErr := res.Body.Close(); closeErr != nil {
//			if err == nil {
//				err = closeErr
//			}
//		}
//	}()
//
//	// Пытаемся обработать запрос как успешный.
//	var r response[T]
//
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		return nil, err
//	}
//
//	// Возвращаем ошибку, если статус отличен от успешного.
//	if r.Status != 200 {
//		if r.Message == "" {
//			return nil, errors.New(fmt.Sprintf("unknown error occurred: %d", r.Status))
//		}
//		return nil, errors.New(fmt.Sprintf("%s: %s", r.Message, r.Description))
//	}
//	return r.Data, nil
//}

// New создает новый экземпляр OER.
func New(appId string) *OER {
	return &OER{appId: appId}
}
