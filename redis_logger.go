package redislog

import (
	"encoding/json"
	"github.com/0studio/redisapi"
	log "github.com/cihub/seelog"
	"time"
)

type RedisLoggerImpl struct {
	client redisapi.Redis
	// key=redis.key
	logMgrMap                map[string]*LogMgr
	chanLogEntity            chan LogEntity
	flushAllChan             chan bool
	idleFlushAllChan         chan string // key
	maxCacheCnt              int         // 如果log 积累到100条 则flush 一次
	idleFlushIntervalSeconds int         // 如果60秒没flush 过 ，则执行一次flush 操作
}

func (logger *RedisLoggerImpl) startRedisLogger() {
	go func() {
		for {
			select {
			case logEntity := <-logger.chanLogEntity:
				if logEntity.flush {
					logMgr, _ := logger.logMgrMap[logEntity.key]
					logMgr.flush()
				} else {
					logger.addLog(logEntity.key, logEntity.value)
				}
			case <-logger.flushAllChan:
				logger.flushAll()
			case key := <-logger.idleFlushAllChan:
				logMgr, ok := logger.logMgrMap[key]
				if ok {
					logMgr.flushIfIdle()
				}
			}
		}
	}()
}

func (instance *RedisLoggerImpl) FlushAll() {
	instance.flushAllChan <- true
}

func (instance *RedisLoggerImpl) flushAll() {
	for key, _ := range instance.logMgrMap {
		instance.logMgrMap[key].flush()
	}
}

func (instance *RedisLoggerImpl) AddLog(key string, value interface{}) {
	instance.chanLogEntity <- LogEntity{key: key, value: value}
}
func (instance *RedisLoggerImpl) addLog(key string, value interface{}) bool {
	logMgr, ok := instance.logMgrMap[key]
	if !ok {
		logMgr = defaultLogEntity(key, instance, instance.maxCacheCnt, instance.idleFlushIntervalSeconds)
		instance.logMgrMap[key] = logMgr
		instance.startTimerFlush(key)
	}

	return logMgr.addLog(value)
}

func (instance *RedisLoggerImpl) startTimerFlush(key string) {
	var interval int = instance.idleFlushIntervalSeconds / 5
	if interval == 0 {
		interval = instance.idleFlushIntervalSeconds
	}

	go func() {

		for {
			select {
			case <-time.After(time.Duration(interval) * time.Second):
				instance.idleFlushAllChan <- key
			}
		}
	}()
}
func getJsonByte(value interface{}) []byte {
	redisMsg, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return redisMsg
}

func (logger *RedisLoggerImpl) pushRedis(key string, values LogList) bool {
	if values == nil || values == nil || len(values) == 0 {
		return true
	}
	err := logger.client.Lpush(key, values.ToJsonByte())
	if err != nil {
		log.Error("add redis log err:", err, key, values)
		return false
	}
	return true
}
