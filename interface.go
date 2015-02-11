package redislog

import (
	"github.com/0studio/redisapi"
)

type RedisLogger interface {
	AddLog(key string, value interface{})
	FlushAll()
}

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
