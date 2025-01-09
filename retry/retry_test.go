package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jrjaro18/retry/config"
	"github.com/jrjaro18/retry/retry"
)

func TestRetryNormalSuccess(t *testing.T) {
	ctx := context.Background()
	conf := config.Config{
		RetryMethod: config.Normal,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return nil
	}

	chRR := retry.Retry(ctx, conf, fn)
	result := <-chRR

	if !result.Success {
		t.Errorf("expected success, got failure with error: %v", result.Error)
	}
}

func TestRetryNormalFailure(t *testing.T) {
	ctx := context.Background()
	conf := config.Config{
		RetryMethod: config.Normal,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return errors.New("retry error")
	}

	chRR := retry.Retry(ctx, conf, fn)
	var result retry.RetryResult
	for result = range chRR {
	}

	if result.Success {
		t.Errorf("expected failure, got success")
	}
	if result.Error == nil || result.Error.Error() != "retry error" {
		t.Errorf("expected error 'retry error', got: %v", result.Error)
	}
}

func TestRetryExponentialSuccess(t *testing.T) {
	ctx := context.Background()
	conf := config.Config{
		RetryMethod: config.Exponential,
		Interval:    100 * time.Millisecond,
		MaxRetries:  3,
	}

	fn := func() error {
		return nil
	}

	chRR := retry.Retry(ctx, conf, fn)
	result := <-chRR

	if !result.Success {
		t.Errorf("expected success, got failure with error: %v", result.Error)
	}
}