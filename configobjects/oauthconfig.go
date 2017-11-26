package configobjects

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/oauth1"
)

type OAuth struct {
	ConsumerKey    string `yaml:"consumer-key"`
	ComsumerSecret string `yaml:"consumer-secret"`
	Token          string `yaml:"token"`
	TokenSecret    string `yaml:"token-secret"`
}

func (conf OAuth) GetOAuth1() (config *oauth1.Config, token *oauth1.Token) {
	spew.Dump(conf)
	config = oauth1.NewConfig(conf.ConsumerKey, conf.ComsumerSecret)
	token = oauth1.NewToken(conf.Token, conf.TokenSecret)
	return
}
