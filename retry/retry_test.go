package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/jrjaro18/retry/config"
)

func TestRetryNormalSuccess(t *testing.T) {
	conf := config.Config{
		RetryMethod: config.Normal,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return nil
	}

	results := Retry(conf, fn)

	result := <-results
	if !result.Success {
		t.Errorf("Expected success, got failure with error: %v", result.Error)
	}
}

func TestRetryNormalFailure(t *testing.T) {
	conf := config.Config{
		RetryMethod: config.Normal,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return errors.New("test error")
	}

	results := Retry(conf, fn)

	var lastResult RetryResult
	for result := range results {
		lastResult = result
	}

	if lastResult.Success {
		t.Errorf("Expected failure, got success")
	}
	if lastResult.Error == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRetryExponentialSuccess(t *testing.T) {
	conf := config.Config{
		RetryMethod: config.Exponential,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return nil
	}

	results := Retry(conf, fn)

	result := <-results
	if !result.Success {
		t.Errorf("Expected success, got failure with error: %v", result.Error)
	}
}

func TestRetryExponentialFailure(t *testing.T) {
	conf := config.Config{
		RetryMethod: config.Exponential,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return errors.New("test error")
	}

	results := Retry(conf, fn)

	var lastResult RetryResult
	for result := range results {
		lastResult = result
	}

	if lastResult.Success {
		t.Errorf("Expected failure, got success")
	}
	if lastResult.Error == nil {
		t.Errorf("Expected error, got nil")
	}
}
