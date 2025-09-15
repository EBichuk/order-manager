package cache

import (
	"order-manager/internal/models"
	"sync"
)

type Cache struct {
	orders    []string
	cacheList map[string]models.Order
	size      int
	mu        *sync.RWMutex
}

func NewCache(size int) *Cache {
	return &Cache{
		orders:    make([]string, 0, size),
		cacheList: make(map[string]models.Order),
		size:      size,
		mu:        &sync.RWMutex{},
	}
}

func (c *Cache) SetOrder(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, found := c.cacheList[order.OrderUID]; found {
		c.cacheList[order.OrderUID] = order
		return
	}

	if len(c.orders) >= c.size {
		old := c.orders[0]
		delete(c.cacheList, old)
		c.orders = c.orders[1:]
	}

	c.cacheList[order.OrderUID] = order
	c.orders = append(c.orders, order.OrderUID)
}

func (c *Cache) GetOrder(orderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, found := c.cacheList[orderUID]
	return order, found
}
