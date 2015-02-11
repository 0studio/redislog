package redislog

import (
	"github.com/0studio/redisapi"
)

type RedisLogger interface {
	AddLog(key string, value interface{})
	FlushAll()
}

// maxCacheCnt 某个key 的队列长度达到此值后会强制批量把其刷到redis中(如果想每次AddLog 直接刷到redis，达设置此值为1即可)
// idleFlushIntervalSeconds 空闲时间达到此值后会强制批量把其刷到redis中
func NewRedisLogger(client redisapi.Redis, maxCacheCnt int, idleFlushIntervalSeconds int) RedisLogger {
	logger := &RedisLoggerImpl{client: client,
		logMgrMap:                make(map[string]*LogMgr),
		chanLogEntity:            make(chan LogEntity),
		flushAllChan:             make(chan bool),
		idleFlushAllChan:         make(chan string),
		maxCacheCnt:              maxCacheCnt,
		idleFlushIntervalSeconds: idleFlushIntervalSeconds,
	}
	logger.startRedisLogger()
	return logger
}
