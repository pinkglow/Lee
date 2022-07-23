package lee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		log.Printf("[Lee] [%s] %d | %s | %v", c.Method, c.StatusCode, c.Path, time.Since(start))
	}
}
