//读取配置文件的工具包
package goconfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrNotString = errors.New("not string")
	ErrNotInt    = errors.New("not int")
	ErrNotInt64  = errors.New("not int64")
	ErrNotFloat  = errors.New("not float")
)

const (
	DEFAILT_SECTION string = "_DEFAULT_"
)

type Config struct {
	Section     string
	sectionList []string
	data        map[string]map[string]string
}

func New(objs ...interface{}) (*Config, error) {
	config := &Config{
		Section:     DEFAILT_SECTION,
		sectionList: []string{DEFAILT_SECTION},
		data:        make(map[string]map[string]string),
	}
	config.SetSection(DEFAILT_SECTION)
	for _, obj := range objs {
		err := config.Init(obj)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
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
		originDefaultVal := elem.Tag.Get("default")
		if len(originDefaultVal) > 0 {
			switch elem.Type.String() {
			case "string":
				objV.FieldByName(elem.Name).SetString(p.StringDefault(cfgName, originDefaultVal))
			case "int":
				defaultVal, err := strconv.Atoi(originDefaultVal)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetInt(int64(p.IntDefault(cfgName, defaultVal)))
			case "int64":
				defaultVal, err := strconv.ParseInt(originDefaultVal, 10, 64)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetInt(p.Int64Default(cfgName, defaultVal))
			case "float64":
				defaultVal, err := strconv.ParseFloat(originDefaultVal, 64)
				if err != nil {
					return err
				}
				objV.FieldByName(elem.Name).SetFloat(p.FloatDefault(cfgName, defaultVal))
			}
		} else {
			switch elem.Type.String() {
			case "string":
				value, err := p.String(cfgName)
				if err == nil {
					objV.FieldByName(elem.Name).SetString(value)
				}
			case "int":
				value, err := p.Int(cfgName)
				if err == nil {
					objV.FieldByName(elem.Name).SetInt(int64(value))
				}
			case "int64":
				value, err := p.Int64(cfgName)
				if err == nil {
					objV.FieldByName(elem.Name).SetInt(value)
				}
			case "float64":
				value, err := p.Float(cfgName)
				if err == nil {
					objV.FieldByName(elem.Name).SetFloat(value)
				}
			}
		}
	}
	return nil
}

//设置当前section
func (p *Config) SetSection(section string) *Config {
	p.Section = section
	if _, ok := p.data[section]; !ok {
		p.sectionList = append(p.sectionList, section)
		p.data[section] = make(map[string]string)
	}
	return p
}

//设置字段
func (p *Config) Set(key string, value interface{}) *Config {
	if _, ok := p.data[p.Section]; ok {
		p.data[p.Section][key] = fmt.Sprintf("%v", value)
	} else {
		p.data[p.Section] = map[string]string{key: fmt.Sprintf("%v", value)}
	}
	return p
}

//读取配置文件
func (p *Config) Parse(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	bomb, err := buf.Peek(3)
	if err != nil {
		return err
	}
	if err == nil && bytes.Equal(bomb, []byte{239, 187, 191}) {
		buf.Read(make([]byte, 3))
	}
	for {
		line, _, err := buf.ReadLine()
		//如果文件结束
		if err == io.EOF {
			break
		}
		line = bytes.TrimSpace(line)
		//如果是空行
		if bytes.Equal(line, []byte{}) {
			continue
		}
		//如果是注释
		if line[0] == '#' || line[0] == ';' {
			continue
		}
		//如果是section头
		if line[0] == '[' && line[len(line)-1] == ']' && len(line) > 2 {
			section := string(line[1 : len(line)-1])
			p.SetSection(section)
		} else {
			//是字段
			splitedLine := bytes.SplitN(line, []byte{'='}, 2)
			if len(splitedLine) != 2 {
				continue
			}
			key, value := bytes.TrimSpace(splitedLine[0]), bytes.TrimSpace(splitedLine[1])
			if len(key) == 0 {
				continue
			}
			p.Set(string(key), string(value))
		}
	}
	p.SetSection(DEFAILT_SECTION)
	return nil
}

//在所有section中搜索
func (p *Config) Search(key string) (string, string, error) {
	for _, section := range p.sectionList {
		if _, ok := p.data[section]; ok {
			if v, ok := p.data[section][key]; ok {
				return v, section, nil
			}
		}
	}
	return "", "", ErrNotFound
}

//获取String
func (p *Config) String(key string) (string, error) {
	if _, ok := p.data[p.Section]; ok {
		if v, ok := p.data[p.Section][key]; ok {
			return v, nil
		}
	}
	return "", ErrNotFound
}

//获取String，如果找不到或出错则返回给予的默认值
func (p *Config) StringDefault(key string, defaultVal string) string {
	value, err := p.String(key)
	if err != nil {
		return defaultVal
	}
	return value
}

//获取Int
func (p *Config) Int(key string) (int, error) {
	v, err := p.String(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return value, nil
}

//获取Int，如果找不到或出错则返回给予的默认值
func (p *Config) IntDefault(key string, defaultVal int) int {
	value, err := p.Int(key)
	if err != nil {
		return defaultVal
	}
	return value
}

//获取Int64
func (p *Config) Int64(key string) (int64, error) {
	v, err := p.String(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

//获取Int64，如果找不到或出错则返回给予的默认值
func (p *Config) Int64Default(key string, defaultVal int64) int64 {
	value, err := p.Int64(key)
	if err != nil {
		return defaultVal
	}
	return value
}

//获取Float
func (p *Config) Float(key string) (float64, error) {
	v, err := p.String(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

//获取Float，如果找不到或出错则返回给予的默认值
func (p *Config) FloatDefault(key string, defaultVal float64) float64 {
	value, err := p.Float(key)
	if err != nil {
		return defaultVal
	}
	return value
}

//获取List
func (p *Config) List(key string, delimiterArg ...string) ([]string, error) {
	value, err := p.String(key)
	if err != nil {
		return nil, err
	}
	delimiter := " "
	if len(delimiterArg) > 0 {
		delimiter = delimiterArg[0]
	}
	return strings.Split(value, delimiter), nil
}

//获取List，如果找不到或出错则返回给予的默认值
func (p *Config) ListDefault(key string, defaultVal []string, delimiterArg ...string) []string {
	delimiter := " "
	if len(delimiterArg) > 0 {
		delimiter = delimiterArg[0]
	}
	value, err := p.List(key, delimiter)
	if err != nil {
		return defaultVal
	}
	return value
}
