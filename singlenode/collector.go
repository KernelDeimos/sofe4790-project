package singlenode

type Collector interface {
	Collect(ev Event)
}

type HandleHereCollector struct {
	e *Emitter
}

type DelegateToPeerCollector struct {
	//
}

func (c *HandleHereCollector) Collect(ev Event) {
	c.e.Emit(ev)
}
