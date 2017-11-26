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
	triggers map[string][]Trigger
}

func NewEmitter() *Emitter {
	return &Emitter{
		map[string][]Trigger{},
	}
}

func (e *Emitter) Emit(ev Event) {
	logrus.Debugf("Emitting event from '%s'", ev.Key)
	if trigList, exists := e.triggers[ev.Key]; exists {
		logrus.Debugf("%d triggers for '%s'", len(trigList), ev.Key)
		for _, trig := range trigList {
			trig.Invoke(e, ev)
		}
	} else {
		logrus.Debugf("No trigger for '%s'", ev.Key)
	}
}

func (e *Emitter) AddTrigger(key string, trig Trigger) {
	if _, exists := e.triggers[key]; !exists {
		e.triggers[key] = []Trigger{}
	}
	e.triggers[key] = append(e.triggers[key], trig)
}
