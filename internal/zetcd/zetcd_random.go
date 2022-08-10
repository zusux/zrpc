package zetcd

import (
	"errors"
	"math/rand"
)

// NewRandom returns a load balancer that selects services randomly.
func NewRandom(seed int64) *random {
	return &random{
		r: rand.New(rand.NewSource(seed)),
	}
}

type random struct {
	r *rand.Rand
}

func (r *random) GetPoint(count int) (int,error) {

	if count <= 0 {
		return -1, errors.New("endpoint count less 0")
	}
	return r.r.Intn(count), nil
}
