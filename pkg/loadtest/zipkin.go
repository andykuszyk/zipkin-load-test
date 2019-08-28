package loadtest

import (
	"bytes"
	"github.com/openzipkin/zipkin-go/idgenerator"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"log"
	"net/http"
	"time"
)

func NewSpan() model.SpanModel {
	idGen := idgenerator.NewRandom64()
	traceId := idGen.TraceID()
	return model.SpanModel{
		SpanContext:    model.SpanContext{
			TraceID:  traceId,
			ID:       idGen.SpanID(traceId),
		},
		Name:           "foo",
		Kind:           model.Client,
		Timestamp:      time.Now(),
	}
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ZipkinTest struct {
	Client     HttpClient
}

func (r *ZipkinTest) Setup() {
	r.Client = &http.Client{}
}

func (r *ZipkinTest) Run() {
	span := NewSpan()
	batch := []*model.SpanModel{&span}
	body, err := reporter.JSONSerializer{}.Serialize(batch)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:9411/api/v2/spans", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)

	if err != nil {
		log.Fatalf("\r\nError when posting span: %s", err)
	}
	if resp.StatusCode != 202 {
		log.Fatalf("\r\nUnexpected status code when posting span: %d", resp.StatusCode)
	}
}

func (r *ZipkinTest) Teardown() {
}