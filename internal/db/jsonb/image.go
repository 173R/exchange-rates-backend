package jsonb

import (
	"database/sql/driver"
	"encoding/json"
)

type ImageSetItem struct {
	// Ширина изображения.
	Width uint `json:"width"`
	// Высота изображения.
	Height uint `json:"height"`
	// Ссылка на изображение.
	Url string `json:"url"`
	// Увеличение изображения.
	Scale uint `json:"scale"`
}

type Image struct {
	// Список изображений.
	Set []ImageSetItem `json:"set"`
}

func (v *Image) Scan(value any) error {
	return ScanTo(value, v)
}

func (v Image) Value() (driver.Value, error) {
	return json.Marshal(v)
}
