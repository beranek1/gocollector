package gocollector

import (
	"time"

	"github.com/beranek1/godatainterface"
)

type Collector struct {
	source      godatainterface.DataSource
	destination godatainterface.DataDestination
	interval    time.Duration
	ticker      *time.Ticker
	done        chan bool
}

func Create(source godatainterface.DataSource, destination godatainterface.DataDestination, interval time.Duration) Collector {
	return Collector{source: source, destination: destination, interval: interval, done: make(chan bool)}
}

func (c *Collector) collect() {
	keys, err := c.source.List()
	if err != nil {
		return
	}
	for _, key := range keys {
		value, err := c.source.Get(key)
		if err == nil {
			c.destination.Put(key, value)
		}
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
