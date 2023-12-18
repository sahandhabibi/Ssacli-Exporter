package main

import (
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"os"
	"testing"
)

func TestDiskStatus_Successful(t *testing.T) {
	// given
	testOutput, _ := os.ReadFile("test/ssacli_output_ok.txt")

	// when
	parse(string(testOutput))

	// then
	metricsChan := make(chan prometheus.Metric, 100)
	disk_status.Collect(metricsChan)

	if len(metricsChan) != 3 {
		t.Fatalf(`Expected 3 metrics, but collected %d`, len(metricsChan))
	}

	assertValue(t, disk_status, "box 6 bay 1 type SAS", 1.0)
}

func TestDiskStatus_Rebuilding(t *testing.T) {
	// given
	testOutput, _ := os.ReadFile("test/ssacli_output_rebuilding.txt")

	// when
	parse(string(testOutput))

	// then
	assertValue(t, disk_status, "box 6 bay 2 type SAS", 0.0)
}

func assertValue(t *testing.T, gaugeVec *prometheus.GaugeVec, name string, expectedVal float64) {
	gaugeVal, gaugeName := extractValue(gaugeVec, name)
	if gaugeVal != expectedVal {
		t.Fatalf(`Expected %s to be %f, but was %f`, gaugeName, expectedVal, gaugeVal)
	}
}


func extractValue(gauge *prometheus.GaugeVec, name string) (float64, string) {
	metricsChan := make(chan prometheus.Metric, 1)
	gauge.With(prometheus.Labels{"physicaldrive": name}).Collect(metricsChan)
	metric := <-metricsChan
	metricDto := io_prometheus_client.Metric{}
	metric.Write(&metricDto)
	return *metricDto.Gauge.Value, *metricDto.Label[0].Value
}
