package gocollector

import "time"

type CollectorSource interface {
	Get(string) any
	List() []string
}

type CollectorDestination interface {
	Put(string, any)
}

type Collector struct {
	source      CollectorSource
	destination CollectorDestination
	interval    time.Duration
	ticker      *time.Ticker
	done        chan bool
}

func Create(source CollectorSource, destination CollectorDestination, interval time.Duration) Collector {
	return Collector{source: source, destination: destination, interval: interval, done: make(chan bool)}
}

func (c *Collector) collect() {
	for _, key := range c.source.List() {
		value := c.source.Get(key)
		c.destination.Put(key, value)
	}
}

func (c *Collector) Interval() time.Duration {
	return c.interval
}

func (c *Collector) Update(interval time.Duration) {
	c.interval = interval
	if c.ticker != nil {
		c.ticker.Reset(c.interval)
	}
}

func (c *Collector) Start() {
	if c.ticker == nil {
		c.ticker = time.NewTicker(c.interval)
	} else {
		c.ticker.Reset(c.interval)
	}
	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-c.ticker.C:
				c.collect()
			}
		}
	}()
}

func (c *Collector) Stop() {
	c.ticker.Stop()
	c.done <- true
}
