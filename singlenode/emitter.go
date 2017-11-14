package singlenode

import "github.com/sirupsen/logrus"

type Event struct {
	Key  string
	Data map[string]string
}

type Trigger interface {
	Invoke(e *Emitter, ev Event)
}

type ConnectionTrigger struct {
	OutputKey string
	DataMap   map[string]string
}

type EndpointTrigger struct {
	Endpoint Endpoint
	DataMap  map[string]string
}

func NewEndpointTrigger(end Endpoint, dataMap map[string]string) *EndpointTrigger {
	return &EndpointTrigger{end, dataMap}
}

func (trigger *EndpointTrigger) Invoke(e *Emitter, ev Event) {
	trigger.Endpoint.Catch(ev.Key, ev.Data)
}

type Emitter struct {
	triggers map[string]Trigger
}

func NewEmitter() *Emitter {
	return &Emitter{
		map[string]Trigger{},
	}
}

func (e *Emitter) Emit(ev Event) {
	logrus.Debugf("Emitting event from '%s'", ev.Key)
	for key, trig := range e.triggers {
		if key == ev.Key {
			logrus.Debugf("Found trigger for '%s'", ev.Key)
			trig.Invoke(e, ev)
		}
	}
}

func (e *Emitter) AddTrigger(key string, trig Trigger) {
	e.triggers[key] = trig
}
