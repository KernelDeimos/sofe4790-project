package singlenode

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Source interface {
	Start(name string, e *Emitter) (*sync.WaitGroup, error)
}

type Endpoint interface {
	Start(name string) (*sync.WaitGroup, error)
	Catch(source string, data map[string]string)
}

type Application struct {
	Sources   map[string]Source
	Endpoints map[string]Endpoint
	emitter   *Emitter
}

func (app *Application) Start() {
	for name, end := range app.Endpoints {
		_, err := end.Start(name)
		if err != nil {
			logrus.Error(err)
		}
	}
	for name, src := range app.Sources {
		_, err := src.Start(name, app.emitter)
		if err != nil {
			logrus.Error(err)
		}
	}
}

type ApplicationBuilder struct {
	config *Config
}

func (builder *ApplicationBuilder) Build() *Application {
	sources := map[string]Source{}
	endpoints := map[string]Endpoint{}

	emitter := NewEmitter()

	builder.buildSources(sources)
	builder.buildEndpoints(endpoints)
	builder.attachTriggers(emitter, endpoints)

	return &Application{sources, endpoints, emitter}
}

func (builder *ApplicationBuilder) buildSources(sources map[string]Source) {
	c := builder.config
	if c.Sources.Stdin {
		src := NewStdinSource()
		sources["stdin"] = src
	}
}

func (builder *ApplicationBuilder) buildEndpoints(endpoints map[string]Endpoint) {
	c := builder.config
	for _, item := range c.Endpoints.Appendlog {
		logrus.Infof("Adding endpoint '%s'", item.Name)
		end := NewStreamEndpoint(item.Path)
		endpoints[item.Name] = end
	}
	for _, item := range c.Endpoints.Tweet {
		logrus.Infof("Adding endpoint '%s'", item.Name)
		end := NewTweetEndpoint(item.Auth)
		endpoints[item.Name] = end
	}
}

func (builder *ApplicationBuilder) attachTriggers(
	e *Emitter,
	endpoints map[string]Endpoint,
) {
	c := builder.config

	for _, trig := range c.Triggers {
		endpoint, exists := endpoints[trig.Endpoint]
		if !exists {
			panic("Invalid configuration: " + trig.Endpoint)
		}
		trigger := NewEndpointTrigger(endpoint, trig.DataMap)
		e.AddTrigger(trig.Source, trigger)
	}
}

func RunApplication() {
	var err error

	configBytes, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	builder := ApplicationBuilder{&config}
	app := builder.Build()

	app.Start()

	time.Sleep(time.Second * 60)
}
