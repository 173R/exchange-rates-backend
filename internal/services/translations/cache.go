package services

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories"
	"sync"
	"time"
)

// Хранит в себе список всех кешированных переводов.
type translationsCache struct {
	// Дата последнего обновления.
	updatedAt time.Time
	// Кеш, в котором содержится последние полученные данные.
	cache map[models.TranslationId]*models.Translation
	rep   *repositories.Translations
	mu    sync.Mutex
}

// FindByIds возвращает переводы по их ключам.
func (t *translationsCache) FindByIds(ids []models.TranslationId) ([]models.Translation, error) {
	var err error

	if err = t.sync(); err != nil {
		if t.cache == nil {
			return nil, err
		}
	}

	res := make([]models.Translation, 0, len(ids))
	for _, id := range ids {
		if tr, ok := t.cache[id]; ok {
			res = append(res, *tr)
		}
	}

	return res, err
}

// В случае необходимости, синхронизирует текущий кеш с БД.
func (t *translationsCache) sync() error {
	// Разрешаем работать с этим методом лишь 1 горутине в один момент времени.
	t.mu.Lock()

	// Не забываем разблокировать доступ к этому методу.
	defer t.mu.Unlock()

	// Кеш валиден лишь в течение 5 минут.
	if t.updatedAt.Add(5 * time.Minute).After(time.Now()) {
		return nil
	}

	// Получаем список всех переводов.
	translations, err := t.rep.FindAll()
	if err != nil {
		return err
	}

	// Создаем новый кеш.
	t.cache = make(map[models.TranslationId]*models.Translation, len(translations))

	for _, tr := range translations {
		tVal := tr
		t.cache[tr.Id] = &tVal
	}

	// Устанавливаем новую дату обновления кеша.
	t.updatedAt = time.Now()

	return nil
}
