package provider

import (
	"fmt"
	"strings"
	"sync"
)

type Memory struct {
	kv sync.Map
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Set(key string, value any) error {
	m.kv.Store(key, value)
	return nil
}

func (m *Memory) Get(key string) (string, error) {
	value, ok := m.kv.Load(key)
	if !ok {
		return "", nil
	}
	return value.(string), nil
}

func (m *Memory) Delete(key string) error {
	m.kv.Delete(key)
	return nil
}

//func (m *Memory) GetAll() (map[string]string, error) {
//	keyValues := make(map[string]string)
//	m.kv.Range(func(key, value any) bool {
//		//keyValues[strings.Split(key.(string), "::")[2]] = value.(string)
//		keyValues[key.(string)] = value.(string)
//		return true
//	})
//	return keyValues, nil
//}

func (m *Memory) GetByNameSpaceAndProfile(namespace, profile string) (map[string]string, error) {
	keyValues := make(map[string]string)
	m.kv.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), fmt.Sprintf("%s::%s", namespace, profile)) {
			//keyValues[strings.Split(key.(string), "::")[2]] = value.(string)
			keyValues[key.(string)] = value.(string)
		}
		return true
	})
	return keyValues, nil
}
