package proxy

import (
	"sync"
)

type Config struct {
	Interceptor []Route     `json:"interceptor"`
	AuthConfig  AuthConfig  `json:"authConfig"`
	Limiter     LimitConfig `json:"limiter"`
}

type LimitConfig struct {
	Enable bool `json:"enable"`
	Qps    int  `json:"qps"`
}

type Server struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Timeout int64  `json:"timeout"`
}

type ConfigHandler struct {
	ch   chan Config
	once sync.Once
	wg   sync.WaitGroup
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		ch: make(chan Config),
	}
}

func (ch *ConfigHandler) Put(conf Config) {
	ch.ch <- conf
}

func (ch *ConfigHandler) Close() {
	ch.once.Do(func() {
		ch.wg.Done()
		close(ch.ch)
	})
}

func (ch *ConfigHandler) Wait() {
	ch.wg.Wait()
}

func (ch *ConfigHandler) Listen(fn func(conf Config)) {
	for {
		select {
		case conf, ok := <-ch.ch:
			if !ok {
				return
			}
			fn(conf)
		}
	}
}
