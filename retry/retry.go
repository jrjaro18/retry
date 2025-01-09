package retry

import (
	"context"
	"time"

	"github.com/jrjaro18/retry/config"
)

// RetryResult holds the result of a retry attempt, including success status and any error encountered.
type RetryResult struct {
	Success bool  // Indicates if the retry attempt was successful
	Error   error // Contains any error encountered during the retry attempt
}

// Retry initiates the retry process based on the given configuration and retry function.
//
// It returns a channel that will receive RetryResult values.
func Retry(ctx context.Context, conf config.Config, fn func() error) <-chan RetryResult {
	// Create a channel to send RetryResult
	chRR := make(chan RetryResult)

	// Start the retry process in a goroutine based on the selected retry method
	if conf.RetryMethod == config.Normal {
		go normalRetry(ctx, conf.Interval, conf.MaxRetries, chRR, fn)
	} else {
		go exponentialRetry(ctx, conf.Interval, conf.MaxRetries, chRR, fn)
	}

	// Return the channel so the caller can receive retry results
	return chRR
}

func normalRetry(ctx context.Context, interval time.Duration, maxRetries uint16, chRR chan<- RetryResult, fn func() error) {
	defer close(chRR)

	for i := uint16(0); i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			chRR <- RetryResult{
				Success: false,
				Error:   ctx.Err(),
			}
			return
		default:
			err := fn()
			if err == nil {
				chRR <- RetryResult{
					Success: true,
					Error:   nil,
				}
				return
			}

			chRR <- RetryResult{
				Success: false,
				Error:   err,
			}

			if i < maxRetries {
				time.Sleep(interval)
			}
		}
	}
}

func exponentialRetry(ctx context.Context, interval time.Duration, maxRetries uint16, chRR chan<- RetryResult, fn func() error) {
	defer close(chRR)

	for i := uint16(0); i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			chRR <- RetryResult{
				Success: false,
				Error:   ctx.Err(),
			}
			return
		default:
			err := fn()
			if err == nil {
				chRR <- RetryResult{
					Success: true,
					Error:   nil,
				}
				return
			}

			chRR <- RetryResult{
				Success: false,
				Error:   err,
			}

			if i < maxRetries {
				time.Sleep(interval)
				interval *= 2
			}
		}
	}
}
