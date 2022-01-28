package main

import "time"

type config struct {
	workerNum int
	timeout   time.Duration
}

type Option interface {
	apply(*config)
} // 这里设置为一个接口， 方便自行设计config，只暴露这个类

type optionFunc func(*config) // 内部实现方式

func (o optionFunc) apply(c *config) {
	o(c)
} // 一个函数就自实现了 option

// 设定默认值，遍历设置选项

func NewConfig(opts ...Option) *config {
	defaultConfig := &config{
		timeout:   3 * time.Second,
		workerNum: 10,
	}
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

func WithWorkerNum(num int) Option {
	return optionFunc(func(c *config) {
		c.workerNum = num
	})
}

func WithTimeout(t time.Duration) Option {
	return optionFunc(func(c *config) {
		c.timeout = t
	})
}

var _ = Option(optionFunc(func(c *config) {})) // 隐式一种实现方式
