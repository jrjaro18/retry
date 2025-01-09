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
		return testFunction("rohan")
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
