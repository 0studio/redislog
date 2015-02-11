# Demo
```

type Student struct{
name string
age int
}

logger:=NewRedisLogger(redisclient,100,30)
logger.AddLog("key","value")
logger.AddLog("key2",Student{"joseph",11})
logger.FlushAll() // should call this , before shutdown
```
