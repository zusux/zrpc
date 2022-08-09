package zetcd

// Balancer 平衡模式
type Balancer interface {
	GetPoint(count int) (int,error)
}

