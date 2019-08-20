package main

import (
	"flag"
	"github.com/andykuszyk/zipkin-load-test/pkg/loadtest"
	"net/http"
)

func main() {
	sampleSize := flag.Int("n", 1000, "The sample size to use.")
	flag.Parse()
	testRunner := loadtest.TestRunner{
		Client:     &http.Client{},
		SampleSize: *sampleSize,
	}
	testRunner.Run()
}