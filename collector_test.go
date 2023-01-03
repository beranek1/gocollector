package gocollector

import (
	"testing"
	"time"
)

type colSrc struct {
	Count uint
	Keys  []string
}

func (s *colSrc) Get(key string) any {
	s.Count += 1
	println("Get " + key)
	return s.Count
}

func (s *colSrc) List() []string {
	println("List")
	return s.Keys
}

type colDes struct {
	Count uint
}

func (d *colDes) Put(key string, data any) {
	println("Put " + key)
	d.Count += 1
}

func TestCreate(t *testing.T) {
	k := []string{"0", "1", "2", "3", "4"}
	s := &colSrc{Keys: k}
	d := &colDes{}
	i := 100 * time.Millisecond
	c := Create(s, d, i)
	if c.Interval() != i {
		t.Error("Collector has wrong interval.")
	}
	if s.Count != 0 {
		t.Error("Collector has already read from source.")
	}
	if d.Count != 0 {
		t.Error("Collector has already written to destination.")
	}
	c.Start()
	time.Sleep(1300 * time.Millisecond)
	c.Stop()
	if s.Count != 65 {
		t.Error("Collector has read from source insufficient times: ", s.Count)
	}
	if d.Count != 65 {
		t.Error("Collector has written to destination insufficient times: ", d.Count)
	}
	time.Sleep(100 * time.Millisecond)
	if s.Count != 65 || d.Count != 65 {
		t.Error("Collector didn't stop.")
	}
}
