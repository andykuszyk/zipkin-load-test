package loadtest

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type TestHttpClient struct {
	DoCount int
}

func (c *TestHttpClient) Do(r *http.Request) (*http.Response, error) {
	c.DoCount += 1
	resp := http.Response{StatusCode:202}
	return &resp, nil
}

func TestRun_ShouldSendNSpans(t *testing.T) {
	httpClient := &TestHttpClient{}
	testRunner := TestRunner{
		Client: httpClient,
		SampleSize: 1000,
	}
	testRunner.Run()
	assert.Equal(t, 1000, httpClient.DoCount)
}
