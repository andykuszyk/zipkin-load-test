package loadtest

import (
	"bytes"
	"fmt"
	"github.com/openzipkin/zipkin-go/idgenerator"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/rcrowley/go-metrics"
	"log"
	"net/http"
	"time"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type TestRunner struct {
	Client     HttpClient
	SampleSize int
}

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

func (r *TestRunner) Run() {
	s := metrics.NewUniformSample(r.SampleSize)
	h := metrics.NewHistogram(s)
	metrics.Register("histo", h)

	for i := 0; i < r.SampleSize; i++ {
		fmt.Printf("\rProcessing span %d/%d", i + 1, r.SampleSize)
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

		start := time.Now()
		resp, err := r.Client.Do(req)
		h.Update(time.Since(start).Nanoseconds())


		if err != nil {
			log.Fatalf("\r\nError when posting span: %s", err)
		}
		if resp.StatusCode != 202 {
			log.Fatalf("\r\nUnexpected status code when posting span: %d", resp.StatusCode)
		}
	}
	fmt.Println("")
	ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	fmt.Printf("  count:       %9d\n", h.Count())
	fmt.Printf("  min:         %9d\n", h.Min() / 1000000)
	fmt.Printf("  max:         %9d\n", h.Max() / 1000000)
	fmt.Printf("  mean:        %12.2f\n", h.Mean() / 1000000)
	fmt.Printf("  stddev:      %12.2f\n", h.StdDev() / 1000000)
	fmt.Printf("  median:      %12.2f\n", ps[0] / 1000000)
	fmt.Printf("  75%%:         %12.2f\n", ps[1] / 1000000)
	fmt.Printf("  95%%:         %12.2f\n", ps[2] / 1000000)
	fmt.Printf("  99%%:         %12.2f\n", ps[3] / 1000000)
	fmt.Printf("  99.9%%:       %12.2f\n", ps[4] / 1000000)
}
