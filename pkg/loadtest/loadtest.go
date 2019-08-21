package loadtest

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"time"
)

type LoadTest interface {
	Setup()
	Run()
	Teardown()
}

type TestRunner struct {
	SampleSize int
	Test       LoadTest
}

func (r *TestRunner) Execute() {
	s := metrics.NewUniformSample(r.SampleSize)
	h := metrics.NewHistogram(s)
	metrics.Register("histo", h)

	r.Test.Setup()

	for i := 0; i < r.SampleSize; i++ {
		fmt.Printf("\rProcessing iteration %d/%d", i + 1, r.SampleSize)
		start := time.Now()

		r.Test.Run()

		h.Update(time.Since(start).Nanoseconds())

	}

	r.Test.Teardown()

	fmt.Println("")
	ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	fmt.Printf("  count:       %9d\n", h.Count())
	fmt.Printf("  min:         %12.2fms\n", float64(h.Min()) / 1000000)
	fmt.Printf("  max:         %12.2fms\n", float64(h.Max()) / 1000000)
	fmt.Printf("  mean:        %12.2fms\n", h.Mean() / 1000000)
	fmt.Printf("  stddev:      %12.2fms\n", h.StdDev() / 1000000)
	fmt.Printf("  median:      %12.2fms\n", ps[0] / 1000000)
	fmt.Printf("  75%%:         %12.2fms\n", ps[1] / 1000000)
	fmt.Printf("  95%%:         %12.2fms\n", ps[2] / 1000000)
	fmt.Printf("  99%%:         %12.2fms\n", ps[3] / 1000000)
	fmt.Printf("  99.9%%:       %12.2fms\n", ps[4] / 1000000)
}
