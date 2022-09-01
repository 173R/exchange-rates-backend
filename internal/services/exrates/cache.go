package exrates

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/exrates"
	"sync"
	"time"
)

// Хранит в себе список всех кешированных курсов обмена.
type cache struct {
	// Дата последнего обновления.
	updatedAt time.Time
	// Кеш, в котором содержатся последние курсы обменов.
	latest map[models.CurrencyId]*models.ExchangeRate
	// Кеш, в котором содержатся курсы обменов за предыдущий день.
	prevDay map[models.CurrencyId]*models.ExchangeRate
	rep     *exrates.Repository
	mu      sync.Mutex
}

// FindLatestById возвращает актуальный курс указанный валюты.
func (c *cache) FindLatestById(ctx context.Context, id models.CurrencyId) (*models.ExchangeRate, error) {
	err := c.sync(ctx)
	return c.latest[id], err
}

// FindPrevDayById возвращает курс указанный валюты за предыдущий день.
func (c *cache) FindPrevDayById(ctx context.Context, id models.CurrencyId) (*models.ExchangeRate, error) {
	err := c.sync(ctx)
	return c.prevDay[id], err
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

	// Получаем самые свежие курсы обменов.
	latest, err := c.rep.FindLatest(ctx)
	if err != nil {
		return err
	}

	// Получаем данные за предыдущий день.
	// Мы приводим текущую дату к началу дня по UTC и используем в дальнейшем
	// запросе.
	ts := time.Now().In(time.UTC).Truncate(24 * time.Hour)

	prevDay, err := c.rep.FindByTimestamp(ctx, ts)
	if err != nil {
		return err
	}

	// Обновляем кеши.
	cacheRates(latest, &c.latest)
	cacheRates(prevDay, &c.prevDay)

	// Устанавливаем новую дату обновления кеша.
	c.updatedAt = time.Now()

	return nil
}

// Кеширует указанные курсы обменов и помещает их по переданному указателю.
func cacheRates(
	rates []*models.ExchangeRate,
	dest *map[models.CurrencyId]*models.ExchangeRate,
) {
	res := make(map[models.CurrencyId]*models.ExchangeRate, len(rates))
	for _, rate := range rates {
		res[rate.CurrencyId] = rate
	}

	*dest = res
}

// Создает новый экземпляр cache.
func newCache(rep *exrates.Repository) *cache {
	return &cache{
		rep:     rep,
		latest:  map[models.CurrencyId]*models.ExchangeRate{},
		prevDay: map[models.CurrencyId]*models.ExchangeRate{},
	}
}
