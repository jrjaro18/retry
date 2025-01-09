# Retry Package

The `retry` package provides a simple mechanism to retry a function that may fail intermittently. It allows you to configure retries using either a normal or exponential backoff strategy. The retry logic supports context cancellation, so retries can be stopped at any time.

## Features

- Retry failed function calls using **Normal** or **Exponential** backoff strategies.
- Handle retries with configurable intervals and maximum retry counts.
- Support for context cancellation to stop retries immediately when needed.
- Returns a result indicating success or failure for each retry attempt.

## Installation

To use this package in your project, you can import it as follows:

```go
import "github.com/jrjaro18/retry"
```
## Example Usage

```
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	cf "github.com/jrjaro18/retry/config"
	re "github.com/jrjaro18/retry/retry"
)

func main() {
	config := cf.NewConfig(cf.WithInterval(1*time.Second), cf.WithMaxRetries(10), cf.WithRetryMethod(cf.Normal))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chRR := re.Retry(ctx, config, func() error {
		return testFunction("jrjaro18")
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range chRR {
			log.Printf("%+v\n", r)
		}
	}()
	wg.Wait()
}

// user function
func testFunction(name string) error {
	_, err := testExternalService()
	if err != nil {
		return err
	}
	return nil
}

// supposed external service
func testExternalService() (bool, error) {
	time.Sleep(100 * time.Millisecond)
	x := rand.Intn(7)
	if x == 4 {
		return true, nil
	}
	return false, fmt.Errorf("faliure in service!!!, %v != %v", x, 4)
}
```

## Retry Configuration

The retry mechanism is highly configurable using the Config struct. It allows you to define the retry method, interval between retries, and the maximum number of retries. Here's a breakdown of each configuration option:

- `RetryMethod`: Defines how retries will be handled. The available methods are:
  - `Normal`: Uniform retry intervals (fixed duration between retries).
  - `Exponential`: The retry interval increases exponentially (doubles after each failure).

- `Interval`: The duration between retry attempts. This is the amount of time to wait between each retry.

- `MaxRetries`: Specifies the total number of retries allowed after the first attempt. For example, if MaxRetries = 10, the function can execute a maximum of 11 times (initial try + 10 retries).

## Default Configuration
If no custom configuration is provided, the package will use the following default values:

```
Config{
    RetryMethod: Normal,     // Uniform retries
    Interval:    5 * time.Second,  // 5 seconds between retries
    MaxRetries:  5,          // 5 retries allowed after the first failure
}
```

## With Methods
The retry package provides a set of "with" methods that allow you to customize the retry configuration when creating a new Config. These methods are chainable, allowing you to set each option individually.

- `WithRetryMethod`: Used to set the retry method (either Normal or Exponential). Example:

```
WithRetryMethod(Exponential) // Sets retry method to Exponential backoff
```

- `WithInterval`: Used to set the interval between retries. Example:

```
WithInterval(2 * time.Second) // Sets interval to 2 seconds between retries
```

- `WithMaxRetries`: Used to set the maximum number of retries allowed. Example:

```
WithMaxRetries(3) // Sets max retries to 3
```

