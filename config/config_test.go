package config

import (
	"testing"
	"time"
)

func TestNewConfig_Default(t *testing.T) {
	config := NewConfig()

	if config.RetryMethod != Normal {
		t.Errorf("expected RetryMethod to be Normal, got %v", config.RetryMethod)
	}
	if config.Interval != 5*time.Second {
		t.Errorf("expected Interval to be 5s, got %v", config.Interval)
	}
	if config.MaxRetries != 5 {
		t.Errorf("expected MaxRetries to be 5, got %v", config.MaxRetries)
	}
}

func TestNewConfig_WithRetryMethod(t *testing.T) {
	config := NewConfig(WithRetryMethod(Exponential))

	if config.RetryMethod != Exponential {
		t.Errorf("expected RetryMethod to be Exponential, got %v", config.RetryMethod)
	}
}

func TestNewConfig_WithInterval(t *testing.T) {
	interval := 10 * time.Second
	config := NewConfig(WithInterval(interval))

	if config.Interval != interval {
		t.Errorf("expected Interval to be %v, got %v", interval, config.Interval)
	}
}

func TestNewConfig_WithMaxRetries(t *testing.T) {
	maxRetries := uint16(10)
	config := NewConfig(WithMaxRetries(maxRetries))

	if config.MaxRetries != maxRetries {
		t.Errorf("expected MaxRetries to be %v, got %v", maxRetries, config.MaxRetries)
	}
}

func TestNewConfig_CustomConfig(t *testing.T) {
	interval := 10 * time.Second
	maxRetries := uint16(10)
	config := NewConfig(
		WithRetryMethod(Exponential),
		WithInterval(interval),
		WithMaxRetries(maxRetries),
	)

	if config.RetryMethod != Exponential {
		t.Errorf("expected RetryMethod to be Exponential, got %v", config.RetryMethod)
	}
	if config.Interval != interval {
		t.Errorf("expected Interval to be %v, got %v", interval, config.Interval)
	}
	if config.MaxRetries != maxRetries {
		t.Errorf("expected MaxRetries to be %v, got %v", maxRetries, config.MaxRetries)
	}
}

func TestWithRetryMethod_Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid retry method")
		}
	}()

	NewConfig(WithRetryMethod(RetryMethod(999)))
}
