package singlenode

import "github.com/KernelDeimos/sofe4790/network"

type Collector interface {
	Collect(ev Event)
}

type HandleHereCollector struct {
	e *Emitter
}

type DelegateToPeerCollector struct {
	e *Emitter
	n *network.Network
}

func (c *HandleHereCollector) Collect(ev Event) {
	c.e.Emit(ev)
}

func (c *DelegateToPeerCollector) Collect(ev Event) {
	if !c.n.SendToLeader(ev) {
		c.e.Emit(ev)
	}
}
