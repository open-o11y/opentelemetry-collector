package labelsprocessor

import (
	v11 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	otlpmetrics "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
	v1 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/resource/v1"
)

var metricDuplicateLabels []*otlpmetrics.ResourceMetrics = []*otlpmetrics.ResourceMetrics{
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

var expectedMetricDuplicateLabels []*otlpmetrics.ResourceMetrics = []*otlpmetrics.ResourceMetrics{
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
