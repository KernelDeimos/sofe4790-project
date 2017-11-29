package singlenode

import "github.com/KernelDeimos/sofe4790/configobjects"

type Config struct {
	Peers     []ConfigPeer    `yaml:"peers"`
	Sources   ConfigSources   `yaml:"sources"`
	Endpoints ConfigEndpoints `yaml:"endpoints"`
	Triggers  []ConfigTrigger `yaml:"triggers"`
}

type ConfigSources struct {
	Stdin bool `yaml:"stdin"`
}

type ConfigEndpoints struct {
	Appendlog []ConfigEndAppendLog `yaml:"appendlog"`
	Tweet     []ConfigEndTweet     `yaml:"tweet"`
}

type ConfigPeer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

/*
type ConfigEnd struct {
	Name string `yaml:"name"`
}

func (confEnd ConfigEnd) GetName() string {
	return confEnd.Name
}
*/

type ConfigEndAppendLog struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type ConfigEndTweet struct {
	Name string              `yaml:"name"`
	Auth configobjects.OAuth `yaml:"auth"`
}

type ConfigTrigger struct {
	Source   string            `yaml:"when"`
	Endpoint string            `yaml:"do"`
	DataMap  map[string]string `yaml:"map"`
}
