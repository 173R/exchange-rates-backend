package models

import (
	"database/sql/driver"
	"encoding/json"
)

type ImageJsonbSetItem struct {
	// Ширина изображения.
	Width uint `json:"width"`
	// Высота изображения.
	Height uint `json:"height"`
	// Ссылка на изображение.
	Url string `json:"url"`
	// Увеличение изображения.
	Scale uint `json:"scale"`
}

type ImageJsonb struct {
	// Список изображений.
	Set []ImageJsonbSetItem `json:"set"`
}

func (v *ImageJsonb) Scan(value any) error {
	return ScanTo(value, v)
}

func (v ImageJsonb) Value() (driver.Value, error) {
	return json.Marshal(v)
}
