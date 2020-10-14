package labelsprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/consumer/pdata"
	v11 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	v1 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
)

type labelMetricProcessor struct {
	cfg *Config
}

func newLabelMetricProcessor(cfg *Config) (*labelMetricProcessor, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &labelMetricProcessor{cfg: cfg}, nil
}

func validateConfig(cfg *Config) error {
	// Ensure no empty keys/values exist
	for _, elem := range cfg.Labels {
		if elem.Key == "" || elem.Value == "" {
			return fmt.Errorf("Labels Processor configuration contains an empty key or value")
		}
	}

	//Ensure no duplicate keys exist
	keys := make(map[string]bool)
	for _, elem := range cfg.Labels {
		_, value := keys[elem.Key]
		if value {
			return fmt.Errorf("Labels Processor configuration contains duplicate keys")
		}
		keys[elem.Key] = true
	}

	return nil
}

func (lp *labelMetricProcessor) ProcessMetrics(_ context.Context, md pdata.Metrics) (pdata.Metrics, error) {

	otlpMetrics := pdata.MetricsToOtlp(md)

	for _, otlpMetric := range otlpMetrics {
		for _, instrMetric := range otlpMetric.GetInstrumentationLibraryMetrics() {
			for _, metric := range instrMetric.GetMetrics() {

				// Multiple types of Data Points exists, since there is no way to determine this beforehand we initialize variables for all
				var intDataPoint []*v1.IntDataPoint
				var doubleDataPoint []*v1.DoubleDataPoint
				var intHistogramDataPoint []*v1.IntHistogramDataPoint
				var doubleHistogramDataPoint []*v1.DoubleHistogramDataPoint

				if metric.GetDoubleGauge() != nil {
					doubleDataPoint = metric.GetDoubleGauge().GetDataPoints()
				} else if metric.GetDoubleHistogram() != nil {
					doubleHistogramDataPoint = metric.GetDoubleHistogram().GetDataPoints()
				} else if metric.GetDoubleSum() != nil {
					doubleDataPoint = metric.GetDoubleSum().GetDataPoints()
				} else if metric.GetIntGauge() != nil {
					intDataPoint = metric.GetIntGauge().GetDataPoints()
				} else if metric.GetIntHistogram() != nil {
					intHistogramDataPoint = metric.GetIntHistogram().GetDataPoints()
				} else if metric.GetIntSum() != nil {
					intDataPoint = metric.GetIntSum().GetDataPoints()
				}

				// Note only 1 of these variables will get populated at a time, hence the loops for the remaining variables do nothing
				for _, label := range lp.cfg.Labels {
					for _, dataPoint := range intDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
					for _, dataPoint := range doubleDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
					for _, dataPoint := range intHistogramDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
					for _, dataPoint := range doubleHistogramDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
				}

			}
		}
	}

	return md, nil
}

// This processor will always by default update existing label values. Also assumes duplicate labels do not already exist in the metric
func deDuplicateAndAppend(labels *[]*v11.StringKeyValue, key string, value string) {
	// If the key already exists, overwrite it
	for _, elem := range *labels {
		if elem.Key == key {
			elem.Value = value
			return
		}
	}
	// If it does not exist, append it
	*labels = append(*labels, &v11.StringKeyValue{Key: key, Value: value})
}
