package goconfig

import (
	"testing"
)

const (
	TestFile string = "test.conf"
)

type Cfg struct {
	Host     string  `cfgname:"HOST"`
	DBName   string  `default:"go"`
	Username string  `cfgname:"USERNAME" default:"jiazhoulvke"`
	Port     int     `cfgname:"PORT" default:"80"`
	Version  float64 `cfgname:"VERSION" default:"1.0"`
}

func Test_Int(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	port, err := cfg.Int("PORT")
	if err != nil {
		t.Error(err)
	}
	if port != 1984 {
		t.Error("Int error")
	} else {
		t.Log("Int success")
	}
}

func Test_IntDefault(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	dv := cfg.IntDefault("NOKEY", 1234)
	if dv != 1234 {
		t.Error("IntDefault error")
	} else {
		t.Log("IntDefault success")
	}
}

func Test_Int64(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	port, err := cfg.Int64("PORT")
	if err != nil {
		t.Error(err)
	}
	if port != 1984 {
		t.Error("Int64 error")
	} else {
		t.Log("Int64 success")
	}
}

func Test_Int64Default(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	dv := cfg.Int64Default("NOKEY", 1234)
	if dv != 1234 {
		t.Error("Int64Default error")
	} else {
		t.Log("Int64Default success")
	}
}

func Test_Float(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	port, err := cfg.Float("VERSION")
	if err != nil {
		t.Error(err)
	}
	if port != 1.1 {
		t.Error("Float error")
	} else {
		t.Log("Float success")
	}
}

func Test_FloatDefault(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	dv := cfg.FloatDefault("NOKEY", 12.34)
	if dv != 12.34 {
		t.Error("FloatDefault error")
	} else {
		t.Log("FloatDefault success")
	}
}

func Test_String(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	username, err := cfg.String("USERNAME")
	if err != nil {
		t.Error(err)
	}
	if username != "jiazhoulvke" {
		t.Error("String error")
	} else {
		t.Log("String success")
	}
}

func Test_StringDefault(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	dv := cfg.StringDefault("NOKEY", "hello,world!")
	if dv != "hello,world!" {
		t.Error("StringDefault error")
	} else {
		t.Log("StringDefault success")
	}
}

func Test_Init(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	err = cfg.Parse(TestFile)
	if err != nil {
		t.Error(err)
	}
	var config Cfg
	err = cfg.Init(&config)
	if config.DBName != "go" || config.Username != "jiazhoulvke" || config.Port != 1984 || config.Version != 1.1 {
		t.Error("Init error")
	} else {
		t.Log("Init success")
	}
}

func Test_Set(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Error(err)
	}
	cfg.Set("ABC", 123).Set("DEF", 456).Set("GHI", 7.89)
	abc, err := cfg.Int("ABC")
	if err != nil {
		t.Error(err)
	}
	def, err := cfg.Int64("DEF")
	if err != nil {
		t.Error(err)
	}
	ghi, err := cfg.Float("GHI")
	if err != nil {
		t.Error(err)
	}
	if abc != 123 || def != 456 || ghi != 7.89 {
		t.Error("Set error")
	} else {
		t.Log("Set success")
	}
}
