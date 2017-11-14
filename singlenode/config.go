package singlenode

type Config struct {
	Sources   ConfigSources   `yaml:"sources"`
	Endpoints ConfigEndpoints `yaml:"endpoints"`
	Triggers  []ConfigTrigger `yaml:"triggers"`
}

type ConfigSources struct {
	Stdin bool `yaml:"stdin"`
}

type ConfigEndpoints struct {
	Appendlog []ConfigEndAppendLog `yaml:"appendlog"`
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

type ConfigTrigger struct {
	Source   string            `yaml:"when"`
	Endpoint string            `yaml:"do"`
	DataMap  map[string]string `yaml:"map"`
}
