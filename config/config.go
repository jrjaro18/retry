package config

import "time"

// can be either Normal (for uniformly distributed retries) or Expo (for retries with exponentially increasing duration)
type RetryMethod int

const (
	// uniformly distributed retries until max retry is reached
	Normal RetryMethod = iota
	// retries with exponentially increasing duration until max retry is reached
	Exponential
)

// configuration for retries
type Config struct {
	// method for retrying it can either be normal or exponential
	RetryMethod RetryMethod
	// interval between each subsequent request
	Interval    time.Duration
	// total number of times to retry after the first failure, 
	//
	// e.g., if MaxRetries = 10, then the given function can execute maximum of 11 times
	MaxRetries  uint16
}

type configFunc func(*Config)

// adds a method for retrying
func WithRetryMethod(m RetryMethod) configFunc {
    return func(c *Config) {
        if m != Normal && m != Exponential {
            panic("invalid retry method")
        }
        c.RetryMethod = m
    }
}


// adds the interval between each subsequent request
func WithInterval(i time.Duration) configFunc {
	return func(c *Config) {
		c.Interval = i
	}
}

// adds the number of retries allowed after an inital try
//
// MaxRetries means total number of times to retry after the first failure
func WithMaxRetries(r uint16) configFunc {
	return func(c *Config) {
		c.MaxRetries = r
	}
}

// when no arguments are passed returns the default configuration: 
//
// Config {
// 	RetryMethod: Normal, 
// 	Interval: 5 * time.Second, 
// 	MaxRetries: 5
// }
//
//to use custom configurations use methods:
//
// WithInterval(x) WithMaxRetries(x) WithRetryMethod(x)
func NewConfig(configs ...configFunc) Config {
	finalConfig := Config{
		RetryMethod: Normal,
		Interval:    5 * time.Second,
		MaxRetries:  5,
	}
	for _, config := range configs {
		config(&finalConfig);
	}
	return finalConfig
}
