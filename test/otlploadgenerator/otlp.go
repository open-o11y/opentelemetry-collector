package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"time"

	service "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	metrics "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	"google.golang.org/grpc"
)

type sender struct {
	client service.MetricsServiceClient
}

var (
	path      = "./data.txt" 				// data file path
	item      = 1000         				// total number of metrics / lines in output file
	metric    = "metricName" 				// base metricName. output file has only unique metricName with a number suffix
	gauge     = "gauge"
	counter   = "counter"
	histogram = "histogram"
	summary   = "summary"
	types     = []string{					// types of metrics generatedq
		counter,
		gauge,
		histogram,
		summary,
	}
	labels = []string{ 						// each metric will have from 1 to 4 sets of labels
		"label1 value1",
		"label2 value2",
		"label3 value3",
		"label4 value4",
	}
	delimeter  = ","                        // separate name, type, labels and metric value
	space      = " "                        // separate a set of label values or metric values
	valueBound = 5000                       // metric values are [0, valueBound)
	bounds     = []float64{0.01, 0.5, 0.99} // fixed quantile/buckets

	endpoint = "localhost:55680"
	waitTime = 5 * time.Second				// wait time between two sends
)

func main() {
	// Writes metrics in the following format to a text file:
	// 		name, type, label1 labelvalue1 , value1 value2 value3 value4 value5
	// gauge and counter has only one value
	generateData()
	// wait for collector to start
	time.Sleep(time.Second * 5)
	// connect to the Collector
	clientConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := service.NewMetricsServiceClient(clientConn)
	s := &sender{
		client,
	}
	// read from file and send metrics
	s.createAndSendMetricsFromFile()
}

// createAndSendMetricsFromFile reads a text file, parse each line to build the corresponding otlp metric, then send the
// metric to the Collector
func (s *sender) createAndSendMetricsFromFile() {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// parse each line and build metric
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), space)
		params := strings.Split(line, delimeter)

		// get metric name and labels
		name := strings.Trim(params[0], space)
		labelSet := getLabels(strings.Split(strings.Trim(params[2], space), space)...)

		mTpye := params[1]
		values := params[3]
		var m *metrics.Metric
		// build metrics
		switch mTpye {
		case gauge, counter:
			m = buildScalarMetric(name, labelSet, parseNumber(values), intComb)
		case histogram:
			m = buildHistogramMetric(name, labelSet, parseuUInt64Slice(values))
		case summary:
			m = buildSummaryMetric(name, labelSet, parseuUInt64Slice(values))
		default:
			log.Println("Invalid metric type")
			continue
		}
		log.Println(m.String())
		s.sendMetric(m)
	}
}

func (s *sender) sendMetric(m *metrics.Metric) {
	request := service.ExportMetricsServiceRequest{

		ResourceMetrics: []*metrics.ResourceMetrics{
			{
				InstrumentationLibraryMetrics: []*metrics.InstrumentationLibraryMetrics{
					{
						Metrics: []*metrics.Metric{
							m,
						},
					},
				},
			},
		},
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err := s.client.Export(ctx, &request)
	time.Sleep(waitTime)
	if err != nil {
		log.Fatal(err)
	}
}
