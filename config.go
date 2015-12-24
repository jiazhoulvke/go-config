package goconfig

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrNotString  = errors.New("not string")
	ErrNotInt     = errors.New("not int")
	ErrNotInt64   = errors.New("not int64")
	ErrNotFloat   = errors.New("not float")
	ErrNotUintptr = errors.New("not uintptr")
)

type Config struct {
	data map[string]interface{}
}

func New() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

func (p *Config) Set(key string, value interface{}) *Config {
	p.data[key] = value
	return p
}

func (p *Config) Init(obj interface{}) error {
	objT := reflect.TypeOf(obj).Elem()
	objV := reflect.ValueOf(obj).Elem()
	for i := 0; i < objT.NumField(); i++ {
		elem := objT.Field(i)
		if !objV.FieldByName(elem.Name).CanSet() {
			continue
		}
		cfgName := elem.Name
		if len(elem.Tag.Get("cfgname")) > 0 {
			cfgName = elem.Tag.Get("cfgname")
		}
		originDefaultValue := elem.Tag.Get("default")
		if len(originDefaultValue) > 0 {
			switch elem.Type.String() {
			case "string":
				objV.FieldByName(elem.Name).SetString(p.StringDefault(cfgName, originDefaultValue))
			case "int":
				defaultValue, err := strconv.Atoi(originDefaultValue)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetInt(int64(p.IntDefault(cfgName, defaultValue)))
			case "int64":
				defaultValue, err := strconv.ParseInt(originDefaultValue, 10, 64)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetInt(p.Int64Default(cfgName, defaultValue))
			case "float64":
				defaultValue, err := strconv.ParseFloat(originDefaultValue, 64)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetFloat(p.FloatDefault(cfgName, defaultValue))
			}
		} else {
			switch elem.Type.String() {
			case "string":
				value, err := p.GetString(cfgName)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetString(value)
			case "int":
				value, err := p.GetInt(cfgName)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetInt(int64(value))
			case "int64":
				value, err := p.GetInt64(cfgName)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetInt(value)
			case "float64":
				value, err := p.GetFloat(cfgName)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetFloat(value)
			}
		}
	}
	return nil
}

func Parse(filename string) (*Config, error) {
	cfg := New()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	bomb, err := buf.Peek(3)
	if err != nil {
		return nil, err
	}
	if err == nil && bytes.Equal(bomb, []byte{239, 187, 191}) {
		buf.Read(make([]byte, 3))
	}
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, []byte{}) {
			continue
		}
		if line[0] == '#' || line[0] == ';' {
			continue
		}
		splitedLine := bytes.SplitN(line, []byte{'='}, 2)
		if len(splitedLine) != 2 {
			continue
		}
		key, value := bytes.TrimSpace(splitedLine[0]), bytes.TrimSpace(splitedLine[1])
		if len(key) == 0 || len(value) == 0 {
			continue
		}
		cfg.data[string(key)] = string(value)
	}
	return cfg, nil
}

func (p *Config) GetInt(key string) (int, error) {
	if v, ok := p.data[key]; ok {
		switch value := v.(type) {
		case int:
			return value, nil
		case string:
			valueInt, err := strconv.Atoi(value)
			if err != nil {
				return 0, ErrNotInt
			}
			return valueInt, nil
		default:
			return 0, ErrNotInt
		}
	}
	return 0, ErrNotFound
}

func (p *Config) IntDefault(key string, defaultValue int) int {
	value, err := p.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (p *Config) GetInt64(key string) (int64, error) {
	if v, ok := p.data[key]; ok {
		switch value := v.(type) {
		case int:
			return int64(value), nil
		case int64:
			return value, nil
		case string:
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, ErrNotInt64
			}
			return valueInt64, nil
		default:
			return 0, ErrNotInt64
		}
	}
	return 0, ErrNotFound
}

func (p *Config) Int64Default(key string, defaultValue int64) int64 {
	value, err := p.GetInt64(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (p *Config) GetFloat(key string) (float64, error) {
	if v, ok := p.data[key]; ok {
		switch value := v.(type) {
		case float64:
			return value, nil
		case string:
			valueFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, ErrNotFloat
			}
			return valueFloat, nil
		default:
			return 0, ErrNotFloat
		}
	}
	return 0, ErrNotFound
}

func (p *Config) FloatDefault(key string, defaultValue float64) float64 {
	value, err := p.GetFloat(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (p *Config) GetString(key string) (string, error) {
	if v, ok := p.data[key]; ok {
		if value, isString := v.(string); isString {
			return value, nil
		} else {
			return "", ErrNotString
		}
	}
	return "", ErrNotFound
}

func (p *Config) StringDefault(key string, defaultValue string) string {
	value, err := p.GetString(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (p *Config) GetInterface(key string) (interface{}, error) {
	if v, ok := p.data[key]; ok {
		return v, nil
	}
	return nil, ErrNotFound
}
