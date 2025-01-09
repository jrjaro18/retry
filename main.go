package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jrjaro18/retry/config"
	"github.com/jrjaro18/retry/retry"
)

func main() {
	// arg 1
	config := config.NewConfig(config.WithInterval(1*time.Second), config.WithMaxRetries(10), config.WithRetryMethod(config.Normal))
	// arg 2
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// arg 3
	fn := func() error {
		return testFunction("jrjaro18")
	}

	chRR := retry.Retry(ctx, config, fn)
	for r := range chRR {
		log.Printf("%+v\n", r)
	}
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
