package retry

import (
	"time"

	"github.com/jrjaro18/retry/config"
)

// RetryResult holds the result of a retry attempt, including success status and any error encountered.
type RetryResult struct {
	Success bool        // Indicates if the retry attempt was successful
	Error   error       // Contains any error encountered during the retry attempt
}

// Retry initiates the retry process based on the given configuration and retry function.
//
// It returns a channel that will receive RetryResult values.
func Retry(conf config.Config, fn func() error) <-chan RetryResult {
	// Create a channel to send RetryResult
	chRR := make(chan RetryResult)

	// Start the retry process in a goroutine based on the selected retry method
	if conf.RetryMethod == config.Normal {
		go normalRetry(conf.Interval, conf.MaxRetries, chRR, fn)
	} else {
		go exponentialRetry(conf.Interval, conf.MaxRetries, chRR, fn)
	}

	// Return the channel so the caller can receive retry results
	return chRR
}

func normalRetry(interval time.Duration, maxRetries uint16, chRR chan<- RetryResult, fn func() error) {
	defer close(chRR)

	for i := uint16(0); i <= maxRetries; i++ {
		err := fn()
		if err == nil {
			chRR <- RetryResult{
				Success: true,
				Error:   nil,
			}
			break
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

func exponentialRetry(interval time.Duration, maxRetries uint16, chRR chan<- RetryResult, fn func() error) {
	defer close(chRR)

	for i := uint16(0); i <= maxRetries; i++ {
		err := fn()
		if err == nil {
			chRR <- RetryResult{
				Success: true,
				Error:   nil,
			}
			break
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
