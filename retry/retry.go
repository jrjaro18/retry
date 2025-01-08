package retry

import (
	"time"

	"github.com/jrjaro18/retry/config"
)

type RetryResult struct {
	Success bool
	Error   error
}

func Retry(conf config.Config, fn func() error) <-chan RetryResult {
	chRR := make(chan RetryResult)

	if conf.RetryMethod == config.Normal {
		go normalRetry(conf.Interval, conf.MaxRetries, chRR, fn)
	} else {
		go exponentialRetry(conf.Interval, conf.MaxRetries, chRR, fn)
	}

	return chRR
}

func normalRetry(interval time.Duration, maxRetries uint16, chRR chan<- RetryResult, fn func() error) {
	defer close(chRR)
	for i := uint16(0); i <= maxRetries; i++ {
		err := fn()
		if err == nil {
			chRR <- RetryResult{
				Success: true,
				Error: nil,
			}
			break
		}
		chRR <- RetryResult{
			Success: false,
			Error: err,
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
				Error: nil,
			}
            break
        }
        chRR <- RetryResult{
			Success: false,
			Error: err,
		}
        if i < maxRetries {
            time.Sleep(interval)
            interval *= 2 
        }
    }
}
