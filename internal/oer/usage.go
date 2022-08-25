package oer

import (
	"time"
)

type usagePlan struct {
	// Частота обновления данных.
	UpdateFrequencyRaw string `mapstructure:"update_frequency"`
}

type usage struct {
	// Информация о подписке.
	Plan *usagePlan `mapstructure:"plan"`
}

// UpdateFrequency возвращает частоту обновлений данных в количестве секунд.
func (u *usagePlan) UpdateFrequency() time.Duration {
	d, err := time.ParseDuration(u.UpdateFrequencyRaw)
	if err != nil {
		return 5 * time.Minute
	}
	return d
}
