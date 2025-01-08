package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	cf "github.com/jrjaro18/retry/config"
	re "github.com/jrjaro18/retry/retry"
)

func main() {
	config := cf.NewConfig(cf.WithInterval(100 * time.Millisecond), cf.WithMaxRetries(10), cf.WithRetryMethod(cf.Normal))
	chRR := re.Retry(config, func() error {
		return testFunction("rohan")
	})
	for r := range(chRR) {
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
	if (x == 4) {
		return true, nil
	}
	return false, fmt.Errorf("faliure in service!!!, %v != %v", x, 4)
}