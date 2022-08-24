package currencies

import (
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	creppkg "github.com/wolframdeus/exchange-rates-backend/internal/repositories/currencies"
	"sync"
	"time"
)

// Хранит в себе список всех кешированных валют.
type cache struct {
	// Дата последнего обновления.
	updatedAt time.Time
	// Кеш, в котором содержится последние полученные данные.
	cache map[models.CurrencyId]*models.Currency
	rep   *creppkg.Currencies
	mu    sync.Mutex
}

// FindByIds возвращает валюты по их ключам.
func (c *cache) FindByIds(ids []models.CurrencyId) ([]*models.Currency, error) {
	err := c.sync()
	res := make([]*models.Currency, 0, len(ids))

	for _, id := range ids {
		if item, ok := c.cache[id]; ok {
			res = append(res, item)
		}
	}

	return res, err
}

// FindById возвращает валюту по её ключу.
func (c *cache) FindById(id models.CurrencyId) (*models.Currency, error) {
	err := c.sync()
	cur, ok := c.cache[id]
	if ok {
		return cur, err
	}
	return nil, err
}

// FindAll возвращает все валюты.
func (c *cache) FindAll() ([]*models.Currency, error) {
	err := c.sync()
	res := make([]*models.Currency, 0, len(c.cache))

	for _, cur := range c.cache {
		res = append(res, cur)
	}
	return res, err
}

// В случае необходимости, синхронизирует текущий кеш с БД.
func (c *cache) sync() error {
	// Разрешаем работать с этим методом лишь 1 горутине в один момент времени.
	c.mu.Lock()

	// Не забываем разблокировать доступ к этому методу.
	defer c.mu.Unlock()

	// Кеш валиден лишь в течение 5 минут.
	if c.updatedAt.Add(5 * time.Minute).After(time.Now()) {
		return nil
	}

	// Получаем список всех валют.
	currencies, err := c.rep.FindAll()
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
func newCache(rep *creppkg.Currencies) *cache {
	return &cache{
		rep:   rep,
		cache: map[models.CurrencyId]*models.Currency{},
	}
}
