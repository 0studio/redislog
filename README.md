# Demo
```
package main

import (
	"github.com/0studio/redisapi"
	"github.com/0studio/redislog"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	client, _ = redisapi.InitRedisClient("127.0.0.1:6379", 1, 1, true)


	logger := redislog.NewRedisLogger(client, 1, 5)
	logger.AddLog("key", "value")
	logger.AddLog("key2", Student{"joseph", 11})
	logger.FlushAll() // should call this , before shutdown
}
```
