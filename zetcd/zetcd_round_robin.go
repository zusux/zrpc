package zetcd

import (
	"errors"
	"sync/atomic"
)

// NewRoundRobin returns a load balancer that returns services in sequence.
func NewRoundRobin() *roundRobin {
	return &roundRobin{
		c: 0,
	}
}

type roundRobin struct {
	c uint64
}

func (r *roundRobin) GetPoint(count int) (int,error) {
	if count <= 0 {
		return -1, errors.New("endpoint count less 0")
	}
	old := atomic.AddUint64(&r.c, 1) - 1
	idx := old % uint64(count)
	return int(idx), nil
}
