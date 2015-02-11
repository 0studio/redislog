package redislog

import (
	"time"
)

var (
	UninitedTime time.Time
)

type LogEntity struct {
	key   string
	value interface{}
	flush bool
}
type LogList []interface{}

func (l LogList) ToJsonByte() (bytes [][]byte) {
	bytes = make([][]byte, len(l))
	for idx, _ := range l {
		bytes[idx] = getJsonByte(l[idx])
	}
	return
}

type LogMgr struct {
	logger *RedisLoggerImpl
	// lock sync.Mutex
	key                      string
	list                     LogList
	lastFlushTime            time.Time
	idleFlushIntervalSeconds int //如果60秒没flush 过 ，则执行一次flush 操作
}

func defaultLogEntity(key string, logger *RedisLoggerImpl, maxCacheCnt int, idleFlushIntervalSeconds int) (e *LogMgr) {
	e = &LogMgr{
		logger: logger,
		key:    key,
		list:   make(LogList, 0, maxCacheCnt),
		idleFlushIntervalSeconds: idleFlushIntervalSeconds,
	}
	return
}

func (l *LogMgr) addLog(value interface{}) bool { // not goroutine safe
	if len(l.list) == cap(l.list)-1 {
		l.list = append(l.list, value)
		l.flush()
		return true
	}
	l.list = append(l.list, value)
	return true
}

func (l *LogMgr) clear() {
	l.list = l.list[0:0]
}
func (l *LogMgr) flush() {
	if l.logger.pushRedis(l.key, l.list) {
		l.clear()
		l.lastFlushTime = time.Now()
	}
}

func (l *LogMgr) flushIfIdle() {
	// 如果1min内 没有flush 过 则执行flush()刷缓存至redis
	if l.lastFlushTime == UninitedTime || time.Now().After(l.lastFlushTime.Add(time.Second*time.Duration(l.idleFlushIntervalSeconds))) {
		l.flush()
	}
}
