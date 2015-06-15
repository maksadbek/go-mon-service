package rcache

type ConcurrentRedis struct {
	p *redis.Pool
}

func (rc *ConcurrentRedis) Start(host string) error {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", host)
		return
	}

	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}

	c := rc.p.Get()
	defer c.Close()
	return c.Err()
}

func (rc *ConcurrentRedis) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := rc.p.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}
