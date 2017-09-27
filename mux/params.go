package mux

import "fmt"

const maxParam = 128

var (
	maxParamsErr = fmt.Errorf("Params length must less than %d ", maxParam)
)

type Params []*UrlParam

type UrlParam struct {
	Key   string
	Value string
}

func (p Params) Get(key string) (string, bool) {
	for _, entry := range p {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

func (p Params) Set(key, val string) error {
	if len(p) >= maxParam {
		return maxParamsErr
	}
	for _, entry := range p {
		if entry.Key == key {
			return fmt.Errorf("Params \"%s\" is exist ", entry.Key)
		}
	}
	p = append(p, &UrlParam{Key: key, Value: val})
	return nil
}

func NewParams() Params {
	return Params(make([]*UrlParam, 0, maxParam))
}
