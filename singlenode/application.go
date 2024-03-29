package singlenode

import (
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/davecgh/go-spew/spew"

	"github.com/KernelDeimos/sofe4790/network"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Source interface {
	Start(name string, c Collector) (*sync.WaitGroup, error)
}

type Endpoint interface {
	Start(name string) (*sync.WaitGroup, error)
	Catch(source string, data map[string]string)
}

type Application struct {
	Sources   map[string]Source
	Endpoints map[string]Endpoint
	emitter   *Emitter  //struct
	collector Collector //interface
}

func (app *Application) Start() {
	for name, end := range app.Endpoints {
		_, err := end.Start(name)
		if err != nil {
			logrus.Error(err)
		}
	}
	for name, src := range app.Sources {
		_, err := src.Start(name, app.collector)
		if err != nil {
			logrus.Error(err)
		}
	}
}

type ApplicationBuilder struct {
	config *Config
}

func (builder *ApplicationBuilder) Build(n *network.Network) *Application {
	sources := map[string]Source{}
	endpoints := map[string]Endpoint{}

	emitter := NewEmitter()

	collector := &DelegateToPeerCollector{emitter, n}

	builder.buildSources(sources)
	builder.buildEndpoints(endpoints)
	builder.attachTriggers(emitter, endpoints)

	return &Application{sources, endpoints, emitter, collector}
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

func RunApplication(host string, port, id int, leader int) {
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

	n := network.NewDefaultNetwork(host, id, port, false)

	nAdded := 0

	// Generate list of peer nodes
	for _, peer := range config.Peers {
		if peer.ID == id {
			//continue // Do not include self as a peer
		}

		n.AddPeer(network.NewPeer(peer.Host, peer.Port, peer.ID))
		nAdded++
	}

	logrus.Infof("Added %d peers to network list.", nAdded)

	n.SetLeader(leader)
	if leader == id {
		logrus.Info("This is the leader node.")
	} else {
		logrus.Info("This is a gateway node.")
	}

	builder := ApplicationBuilder{&config}
	app := builder.Build(n)

	r := gin.Default()
	n.Attach(r, func(data interface{}) {
		spew.Dump(data)
	})

	go n.StartService()
	app.Start()

	r.Run(":" + strconv.Itoa(port))
}
