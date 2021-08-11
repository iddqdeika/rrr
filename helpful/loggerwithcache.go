package helpful

import (
	"fmt"
	"sync"
)

func LoggerWithCache(size int) (defaultLogger, *cachedPrinter) {
	cp := &cachedPrinter{maxSize: size}
	return defaultLogger{
		p:     cp,
		level: 0,
	}, cp
}

type cachedPrinter struct {
	cache   []string
	maxSize int

	sync.RWMutex
}

func (c *cachedPrinter) GetLastLogs() []string {
	c.RLock()
	defer c.RUnlock()
	return pickRight(c.cache, c.maxSize)
}

func (c *cachedPrinter) Printf(format string, a ...interface{}) (n int, err error) {
	c.add(fmt.Sprintf(format, a...))
	return fmt.Printf(format, a...)
}

func (c *cachedPrinter) add(val string) {
	c.Lock()
	defer c.Unlock()
	c.cache = append(c.cache, val)
	if len(c.cache) > c.maxSize*2 {
		c.cache = pickRight(c.cache, c.maxSize)
	}
}

func pickRight(slice []string, size int) []string {
	if len(slice) <= size {
		return slice
	}
	res := make([]string, 0, size)
	res = append(res, slice[len(slice)-size:]...)
	return res
}
