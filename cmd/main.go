package main

import (
	"flag"
	"fmt"
	"github.com/andykuszyk/zipkin-load-test/pkg/loadtest"
	"net/http"
)

func main() {
	sampleSize := flag.Int("n", 1000, "The sample size to use.")
	test := flag.String("test", "sqs", "Which test implementation to run (zipkin or sqs)")
	sqsRegion := flag.String("sqsRegion", "eu-west-1", "Which region to use (sqs only)")
	sqsEndpoint := flag.String("sqsEndpoint", "http://localhost:4100", "Which endpoint to use (sqs only)")
	flag.Parse()

	fmt.Printf("Running load test %s\n", *test)
	var loadTest loadtest.LoadTest

	if *test == "zipkin" {
		loadTest = &loadtest.ZipkinTest{
			Client: http.DefaultClient,
		}
	} else if *test == "sqs" {
		loadTest = &loadtest.SqsTest{
			SQSEndpoint: sqsEndpoint,
			SQSRegion: sqsRegion,
		}
	}

	testRunner := loadtest.TestRunner{
		SampleSize: *sampleSize,
		Test:       loadTest,
	}
	testRunner.Execute()
}