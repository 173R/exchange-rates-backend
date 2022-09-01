package currencies

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
	"sync"
	"time"
)

// Хранит в себе список всех кешированных валют.
type cache struct {
	// Дата последнего обновления.
	updatedAt time.Time
	// Кеш, в котором содержится последние полученные данные.
	cache map[models.CurrencyId]*models.Currency
	rep   *currencies.Repository
	mu    sync.Mutex
}

// FindByIds возвращает валюты по их ключам.
func (c *cache) FindByIds(ctx context.Context, ids []models.CurrencyId) ([]*models.Currency, error) {
	err := c.sync(ctx)
	res := make([]*models.Currency, 0, len(ids))

	for _, id := range ids {
		if item, ok := c.cache[id]; ok {
			res = append(res, item)
		}
	}

	return res, err
}

// FindById возвращает валюту по её ключу.
func (c *cache) FindById(ctx context.Context, id models.CurrencyId) (*models.Currency, error) {
	err := c.sync(ctx)
	cur, ok := c.cache[id]
	if ok {
		return cur, err
	}
	return nil, err
}

// FindAll возвращает все валюты.
func (c *cache) FindAll(ctx context.Context) ([]*models.Currency, error) {
	err := c.sync(ctx)
	res := make([]*models.Currency, 0, len(c.cache))

	for _, cur := range c.cache {
		res = append(res, cur)
	}
	return res, err
}

// В случае необходимости, синхронизирует текущий кеш с БД.
func (c *cache) sync(ctx context.Context) error {
	// Разрешаем работать с этим методом лишь 1 горутине в один момент времени.
	c.mu.Lock()

	// Не забываем разблокировать доступ к этому методу.
	defer c.mu.Unlock()

	// Кеш валиден лишь в течение 5 минут.
	if c.updatedAt.Add(5 * time.Minute).After(time.Now()) {
		return nil
	}

	// Получаем список всех валют.
	currencies, err := c.rep.FindAll(ctx)
	if err != nil {
		return err
	}

	// Создаем новый кеш.
	c.cache = make(map[models.CurrencyId]*models.Currency, len(currencies))

	for _, cur := range currencies {
		c.cache[cur.Id] = cur
	}

	// Устанавливаем новую дату обновления кеша.
	c.updatedAt = time.Now()

	return nil
}

// Создает новый экземпляр cache.
func newCache(rep *currencies.Repository) *cache {
	return &cache{
		rep:   rep,
		cache: map[models.CurrencyId]*models.Currency{},
	}
}
