package helpful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const jsonPathDelimiter = "."

func NewJsonCfg(fileName string) (Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cant read config file: %v", err)
	}
	return newJsonCfgFromBytes(data)
}

func NewJsonCfgFromURL(url string) (Config, error) {
	return newJsonCfgFromURL(url)
}

func NewJsonCfgFromURLWithRefresh(url string) (Config, func() error, error) {
	return newJsonCfgWithRefresh(func() ([]byte, error) {
		return getConfigDataByURL(url)
	})
}

func newJsonCfgFromURL(url string) (*jsonConfig, error) {
	data, err := getConfigDataByURL(url)
	if err != nil {
		return nil, err
	}
	return newJsonCfgFromBytes(data)
}

func newJsonCfgFromBytes(data []byte) (*jsonConfig, error) {
	cfg, err := parseMap(data)
	if err != nil {
		return nil, err
	}
	return &jsonConfig{
		cfg: cfg,
	}, nil
}

func newJsonCfgWithRefresh(dataFunc func() ([]byte, error)) (*jsonConfig, func() error, error) {
	res := &jsonConfig{}
	f := refreshFunc(dataFunc, res)
	err := f()
	if err != nil {
		return nil, nil, err
	}
	return res, f, nil
}

func refreshFunc(dataFunc func() ([]byte, error), res *jsonConfig) func() error {
	f := func() error {
		data, err := dataFunc()
		if err != nil {
			return err
		}
		m, err := parseMap(data)
		if err != nil {
			return err
		}
		res.Lock()
		defer res.Unlock()
		res.cfg = m
		return nil
	}
	return f
}

func getConfigDataByURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cant get config data from url %v: %v", url, err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cand read data from response for url %v: %v", url, err)
	}
	return data, nil
}

func parseMap(data []byte) (map[string]interface{}, error) {
	cfg := make(map[string]interface{})

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cant unmarshal config: %v", err)
	}

	return cfg, nil
}

type jsonConfig struct {
	cfg        map[string]interface{}
	pathPrefix string
	parent     *jsonConfig

	sync.RWMutex
}

func (j *jsonConfig) Child(path string) Config {
	//if j.pathPrefix != "" {
	//	path = j.pathPrefix + jsonPathDelimiter + path
	//}
	return &jsonConfig{
		cfg:        j.cfg,
		pathPrefix: path,
		parent:     j,
	}
}

func (j *jsonConfig) GetArray(path string) ([]Config, error) {
	val, err := j.getValByPath(path)
	if err != nil {
		return nil, err
	}

	switch arr := val.(type) {
	case []interface{}:
		res := make([]Config, 0)
		for i, v := range arr {
			switch v.(type) {
			case map[string]interface{}:
				res = append(res, &jsonConfig{
					parent:     j,
					pathPrefix: path + jsonPathDelimiter + strconv.Itoa(i),
				})
			default:
				return nil, fmt.Errorf("element with index %v of array by path %v is not a json object", i, path)
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("value by path %v is not an array", path)
	}
}

func (j *jsonConfig) GetInterfaceArray(path string) ([]interface{}, error) {
	val, err := j.getValByPath(path)
	if err != nil {
		return nil, err
	}
	switch arr := val.(type) {
	case []interface{}:
		return arr, nil
	default:
		return nil, fmt.Errorf("value by path %v is not an array", path)
	}
}

func (j *jsonConfig) GetInterface(path string) (interface{}, error) {
	return j.getValByPath(path)
}

func (j *jsonConfig) getValByPath(path string) (interface{}, error) {
	if j.pathPrefix != "" {
		path = j.pathPrefix + jsonPathDelimiter + path
	}
	if j.parent != nil {
		return j.parent.getValByPath(path)
	}
	j.RLock()
	defer j.RUnlock()
	names := strings.Split(path, jsonPathDelimiter)
	var v interface{} = j.cfg
	processedPath := ""
	delim := ""
	for _, name := range names {
		switch m := v.(type) {
		case map[string]interface{}:
			var exists bool
			v, exists = m[name]
			if !exists {
				return nil, fmt.Errorf("cant get value for %v, element %v doesn't exist", path, name)
			}
		case []interface{}:
			i, err := strconv.Atoi(name)
			if err != nil {
				return nil, fmt.Errorf("cant get value for %v, element %v is an array, but given %v key isn't a number", path, processedPath, name)
			}
			if i > len(m)-1 {
				return nil, fmt.Errorf("cant get value for %v, array %v has %v elements, but %v key has given", path, processedPath, len(m), i)
			}
			v = m[i]
		default:
			return nil, fmt.Errorf("cant get value for %v, element %v doesn't exist", path, name)
		}
		processedPath += delim + name
		delim = jsonPathDelimiter
	}
	return v, nil
}

func (j *jsonConfig) GetInt(path string) (int, error) {
	val, err := j.getValByPath(path)
	if err != nil {
		return 0, err
	}

	if res, ok := val.(int); ok {
		return res, nil
	}
	switch k := reflect.ValueOf(val).Kind(); k {
	case reflect.Float32:
		return int(val.(float32)), nil
	case reflect.Float64:
		return int(val.(float64)), nil
	case reflect.String:
		val, err := strconv.Atoi(val.(string))
		if err != nil {
			return 0, err
		}
		return val, nil
	default:
		return 0, fmt.Errorf("value of %v is %v type and has unconvertable value %v", path, reflect.TypeOf(val), val)
	}
}

func (j *jsonConfig) GetString(path string) (string, error) {
	val, err := j.getValByPath(path)
	if err != nil {
		return "", err
	}
	switch t := val.(type) {
	case string:
		return val.(string), nil
	default:
		return "", fmt.Errorf("value %v is %v type", path, t)
	}
}
