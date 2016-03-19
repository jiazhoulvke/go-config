package goconfig

import (
	"testing"
)

const (
	TestFile string = "test.conf"
)

type cfg struct {
	Host     string  `cfgname:"HOST"`
	DBName   string  `default:"go"`
	Username string  `cfgname:"USERNAME" default:"jiazhoulvke"`
	Port     int     `cfgname:"PORT" default:"80"`
	Version  float64 `cfgname:"VERSION" default:"1.0"`
}

var config *Config
var c cfg

func Test_Config(t *testing.T) {
	var err error
	config, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if err = config.Parse("test.conf"); err != nil {
		t.Error(err)
	}
	if err = config.Init(&c); err != nil {
		t.Error(err)
	}
	ts, section, err := config.Search("tstring")
	if err != nil {
		t.Error(err)
	}
	if ts != "333" || section != "abc" {
		t.Error("Search fail")
	}
	if c.DBName != "go" || c.Host != "localhost" || c.Port != 1984 || c.Username != "jiazhoulvke" || c.Version != 1.1 {
		t.Error("value error")
	}
}
