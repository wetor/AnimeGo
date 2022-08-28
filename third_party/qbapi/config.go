package qbapi

import "time"

type Config struct {
	Username string
	Password string
	Host     string
	Timeout  time.Duration
}

type Option func(c *Config)

func WithAuth(user string, pwd string) Option {
	return func(c *Config) {
		c.Username = user
		c.Password = pwd
	}
}

func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

func WithTimeout(duration time.Duration) Option {
	return func(c *Config) {
		c.Timeout = duration
	}
}
