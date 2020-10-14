package labelsprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"

	v11 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	otlpmetrics "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
	v1 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/resource/v1"
)

func TestAttachLabels(t *testing.T) {

	// First setting up config and parameters for tests
	cfg := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "cluster", Value: "c1"},
			{Key: "__replica__", Value: "r1"},
		},
	}

	metricAllDataTypes := []*otlpmetrics.ResourceMetrics{
		{
			Resource: &v1.Resource{
				Attributes:             []*v11.KeyValue{},
				DroppedAttributesCount: 0,
			},
			InstrumentationLibraryMetrics: []*otlpmetrics.InstrumentationLibraryMetrics{
				{
					InstrumentationLibrary: &v11.InstrumentationLibrary{},
					Metrics: []*otlpmetrics.Metric{
						{
							Name:        "counter-int",
							Description: "counter-int",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntSum{
								IntSum: &otlpmetrics.IntSum{
									DataPoints: []*otlpmetrics.IntDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             123,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             456,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "counter-double",
							Description: "counter-double",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleSum{
								DoubleSum: &otlpmetrics.DoubleSum{
									DataPoints: []*otlpmetrics.DoubleDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             1.23,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             4.56,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "double-histogram",
							Description: "double-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleHistogram{
								DoubleHistogram: &otlpmetrics.DoubleHistogram{
									DataPoints: []*otlpmetrics.DoubleHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{1},
											Exemplars: []*otlpmetrics.DoubleExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
						{
							Name:        "int-histogram",
							Description: "int-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntHistogram{
								IntHistogram: &otlpmetrics.IntHistogram{
									DataPoints: []*otlpmetrics.IntHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{0},
											Exemplars: []*otlpmetrics.IntExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
						{
							Name:        "int-gauge",
							Description: "int-gauge",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntGauge{
								IntGauge: &otlpmetrics.IntGauge{
									DataPoints: []*otlpmetrics.IntDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             123,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             456,
										},
									},
								},
							},
						},
						{
							Name:        "double-gauge",
							Description: "double-gauge",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleGauge{
								DoubleGauge: &otlpmetrics.DoubleGauge{
									DataPoints: []*otlpmetrics.DoubleDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             1.23,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             4.56,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	expectedMetricAllDataTypes := []*otlpmetrics.ResourceMetrics{
		{
			Resource: &v1.Resource{
				Attributes:             []*v11.KeyValue{},
				DroppedAttributesCount: 0,
			},
			InstrumentationLibraryMetrics: []*otlpmetrics.InstrumentationLibraryMetrics{
				{
					InstrumentationLibrary: &v11.InstrumentationLibrary{},
					Metrics: []*otlpmetrics.Metric{
						{
							Name:        "counter-int",
							Description: "counter-int",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntSum{
								IntSum: &otlpmetrics.IntSum{
									DataPoints: []*otlpmetrics.IntDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             123,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             456,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "counter-double",
							Description: "counter-double",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleSum{
								DoubleSum: &otlpmetrics.DoubleSum{
									DataPoints: []*otlpmetrics.DoubleDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             1.23,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             4.56,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "double-histogram",
							Description: "double-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleHistogram{
								DoubleHistogram: &otlpmetrics.DoubleHistogram{
									DataPoints: []*otlpmetrics.DoubleHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{1},
											Exemplars: []*otlpmetrics.DoubleExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
						{
							Name:        "int-histogram",
							Description: "int-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntHistogram{
								IntHistogram: &otlpmetrics.IntHistogram{
									DataPoints: []*otlpmetrics.IntHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{0},
											Exemplars: []*otlpmetrics.IntExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
						{
							Name:        "int-gauge",
							Description: "int-gauge",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntGauge{
								IntGauge: &otlpmetrics.IntGauge{
									DataPoints: []*otlpmetrics.IntDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             123,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             456,
										},
									},
								},
							},
						},
						{
							Name:        "double-gauge",
							Description: "double-gauge",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleGauge{
								DoubleGauge: &otlpmetrics.DoubleGauge{
									DataPoints: []*otlpmetrics.DoubleDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             1.23,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             4.56,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	metricDuplicateLabels := []*otlpmetrics.ResourceMetrics{
		{
			Resource: &v1.Resource{
				Attributes:             []*v11.KeyValue{},
				DroppedAttributesCount: 0,
			},
			InstrumentationLibraryMetrics: []*otlpmetrics.InstrumentationLibraryMetrics{
				{
					InstrumentationLibrary: &v11.InstrumentationLibrary{},
					Metrics: []*otlpmetrics.Metric{
						{
							Name:        "counter-int",
							Description: "counter-int",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntSum{
								IntSum: &otlpmetrics.IntSum{
									DataPoints: []*otlpmetrics.IntDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             123,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             456,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "counter-double",
							Description: "counter-double",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleSum{
								DoubleSum: &otlpmetrics.DoubleSum{
									DataPoints: []*otlpmetrics.DoubleDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             1.23,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             4.56,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "double-histogram",
							Description: "double-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleHistogram{
								DoubleHistogram: &otlpmetrics.DoubleHistogram{
									DataPoints: []*otlpmetrics.DoubleHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{1},
											Exemplars: []*otlpmetrics.DoubleExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
						{
							Name:        "int-histogram",
							Description: "int-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntHistogram{
								IntHistogram: &otlpmetrics.IntHistogram{
									DataPoints: []*otlpmetrics.IntHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "label-value-1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{0},
											Exemplars: []*otlpmetrics.IntExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
					},
				},
			},
		},
	}

	expectedMetricDuplicateLabels := []*otlpmetrics.ResourceMetrics{
		{
			Resource: &v1.Resource{
				Attributes:             []*v11.KeyValue{},
				DroppedAttributesCount: 0,
			},
			InstrumentationLibraryMetrics: []*otlpmetrics.InstrumentationLibraryMetrics{
				{
					InstrumentationLibrary: &v11.InstrumentationLibrary{},
					Metrics: []*otlpmetrics.Metric{
						{
							Name:        "counter-int",
							Description: "counter-int",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntSum{
								IntSum: &otlpmetrics.IntSum{
									DataPoints: []*otlpmetrics.IntDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             123,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             456,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "counter-double",
							Description: "counter-double",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleSum{
								DoubleSum: &otlpmetrics.DoubleSum{
									DataPoints: []*otlpmetrics.DoubleDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             1.23,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Value:             4.56,
										},
									},
									AggregationTemporality: 2,
									IsMonotonic:            true,
								},
							},
						},
						{
							Name:        "double-histogram",
							Description: "double-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_DoubleHistogram{
								DoubleHistogram: &otlpmetrics.DoubleHistogram{
									DataPoints: []*otlpmetrics.DoubleHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{1},
											Exemplars: []*otlpmetrics.DoubleExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
						{
							Name:        "int-histogram",
							Description: "int-histogram",
							Unit:        "1",
							Data: &otlpmetrics.Metric_IntHistogram{
								IntHistogram: &otlpmetrics.IntHistogram{
									DataPoints: []*otlpmetrics.IntHistogramDataPoint{
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-1", Value: "label-value-1"},
												{Key: "label-3", Value: "label-value-3"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
										},
										{
											Labels: []*v11.StringKeyValue{
												{Key: "label-2", Value: "label-value-2"},
												{Key: "cluster", Value: "c1"},
												{Key: "__replica__", Value: "r1"},
											},
											StartTimeUnixNano: 1581452772000000321,
											TimeUnixNano:      1581452773000000789,
											Count:             1,
											Sum:               15,
											BucketCounts:      []uint64{0, 1},
											ExplicitBounds:    []float64{0},
											Exemplars: []*otlpmetrics.IntExemplar{
												{
													FilteredLabels: []*v11.StringKeyValue{
														{Key: "exemplar-attachment", Value: "exemplar-attachment-value"},
													},
													TimeUnixNano: 1581452773000000123,
													Value:        15,
													SpanId:       v11.SpanID{},
													TraceId:      v11.TraceID{},
												},
											},
										},
									},
									AggregationTemporality: 2,
								},
							},
						},
					},
				},
			},
		},
	}

	// Set up tests
	tests := []struct {
		name                   string
		config                 *Config
		incomingMetric         []*otlpmetrics.ResourceMetrics
		expectedOutgoingMetric []*otlpmetrics.ResourceMetrics
	}{
		{
			name:                   "test_add_labels_on_all_metrics_types",
			config:                 cfg,
			incomingMetric:         metricAllDataTypes,
			expectedOutgoingMetric: expectedMetricAllDataTypes,
		},
		{
			name:                   "test_adding_duplicate_labels",
			config:                 cfg,
			incomingMetric:         metricDuplicateLabels,
			expectedOutgoingMetric: expectedMetricDuplicateLabels,
		},
	}

	// Execute tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory()
			tmn := &testMetricsConsumer{}
			rmp, err := factory.CreateMetricsProcessor(context.Background(), component.ProcessorCreateParams{}, cfg, tmn)
			require.NoError(t, err)
			assert.Equal(t, true, rmp.GetCapabilities().MutatesConsumedData)

			sourceMetricData := pdata.MetricsFromOtlp(tt.incomingMetric)
			wantMetricData := pdata.MetricsFromOtlp(tt.expectedOutgoingMetric)
			err = rmp.ConsumeMetrics(context.Background(), sourceMetricData)
			require.NoError(t, err)
			assert.EqualValues(t, wantMetricData, tmn.md)
		})
	}

}

type testMetricsConsumer struct {
	md pdata.Metrics
}

func (tmn *testMetricsConsumer) ConsumeMetrics(_ context.Context, md pdata.Metrics) error {
	// simply store md in struct so that we can compare the result to the input
	tmn.md = md
	return nil
}
