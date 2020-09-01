// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusremotewriteexporter

import (
	"go.opentelemetry.io/collector/consumer/pdata"
	common "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	"strconv"
	"testing"

	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"

	//common "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	otlp "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
)

var (
	validIntGauge = "valid_IntGauge"
	validDoubleGauge = "valid_DoubleGauge"
	validIntSum = "valid_IntSum"
	validDoubleSum ="valid_DoubleSum"
	validIntHistogram ="valid_IntHistogram"
	validDoubleHistogram = "valid_DoubleHistogram"

	validIntGaugeDirty = "*valid_IntGauge$"

	unmatchedBoundBucketIntHist = "unmatchedBoundBucketIntHist"
	unmatchedBoundBucketDoubleHist = "unmatchedBoundBucketDoubleHist"

	// valid metrics as input should not return error
	validMetrics1   = map[string]*otlp.Metric{
		validIntGauge:{
			Name: validIntGauge,
			Data:
			&otlp.Metric_IntGauge{
				IntGauge: &otlp.IntGauge{
					DataPoints: []*otlp.IntDataPoint{
						getIntDataPoint(lbs1,intVal1,time1),
						nil,
					},
				},
			},
		},
		validDoubleGauge:{
			Name: validDoubleGauge,
			Data:
			&otlp.Metric_DoubleGauge{
				DoubleGauge: &otlp.DoubleGauge{
					DataPoints: []*otlp.DoubleDataPoint{
						getDoubleDataPoint(lbs1,floatVal1,time1),
						nil,
					},
				},
			},
		},
		validIntSum:{
			Name: validIntSum,
			Data:
			&otlp.Metric_IntSum{
				IntSum: &otlp.IntSum{
					DataPoints: []*otlp.IntDataPoint{
						getIntDataPoint(lbs1,intVal1,time1),
						nil,
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validDoubleSum: {
			Name: validDoubleSum,
			Data:
			&otlp.Metric_DoubleSum{
				DoubleSum: &otlp.DoubleSum{
					DataPoints: []*otlp.DoubleDataPoint{
						getDoubleDataPoint(lbs1,floatVal1,time1),
						nil,
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validIntHistogram: {
			Name: validIntHistogram,
			Data:
			&otlp.Metric_IntHistogram{
				IntHistogram: &otlp.IntHistogram{
					DataPoints: []*otlp.IntHistogramDataPoint{
						getIntHistogramDataPoint(lbs1, time1, floatVal1, uint64(intVal1), bounds, buckets),
						nil,
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validDoubleHistogram:{
			Name: validDoubleHistogram,
			Data:
			&otlp.Metric_DoubleHistogram{
				DoubleHistogram: &otlp.DoubleHistogram{
					DataPoints: []*otlp.DoubleHistogramDataPoint{
						getDoubleHistogramDataPoint(lbs1, time1, floatVal1, uint64(intVal1), bounds, buckets),
						nil,
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
	}
	validMetrics2   = map[string]*otlp.Metric{
		validIntGauge:{
			Name: validIntGauge,
			Data:
			&otlp.Metric_IntGauge{
				IntGauge: &otlp.IntGauge{
					DataPoints: []*otlp.IntDataPoint{
						getIntDataPoint(lbs2,intVal2,time2),
					},
				},
			},
		},
		validDoubleGauge:{
			Name: validDoubleGauge,
			Data:
			&otlp.Metric_DoubleGauge{
				DoubleGauge: &otlp.DoubleGauge{
					DataPoints: []*otlp.DoubleDataPoint{
						getDoubleDataPoint(lbs2,floatVal2,time2),
					},
				},
			},
		},
		validIntSum:{
			Name: validIntSum,
			Data:
			&otlp.Metric_IntSum{
				IntSum: &otlp.IntSum{
					DataPoints: []*otlp.IntDataPoint{
						getIntDataPoint(lbs2,intVal2,time2),
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validDoubleSum: {
			Name: validDoubleSum,
			Data:
			&otlp.Metric_DoubleSum{
				DoubleSum: &otlp.DoubleSum{
					DataPoints: []*otlp.DoubleDataPoint{
						getDoubleDataPoint(lbs2,floatVal2,time2),
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validIntHistogram: {
			Name: validIntHistogram,
			Data:
			&otlp.Metric_IntHistogram{
				IntHistogram: &otlp.IntHistogram{
					DataPoints: []*otlp.IntHistogramDataPoint{
						getIntHistogramDataPoint(lbs2, time2, floatVal2, uint64(intVal2), bounds, buckets),
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validDoubleHistogram:{
			Name: validDoubleHistogram,
			Data:
			&otlp.Metric_DoubleHistogram{
				DoubleHistogram: &otlp.DoubleHistogram{
					DataPoints: []*otlp.DoubleHistogramDataPoint{
						getDoubleHistogramDataPoint(lbs2, time2, floatVal2, uint64(intVal2), bounds, buckets),
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		validIntGaugeDirty:{
			Name: validIntGaugeDirty,
			Data:
			&otlp.Metric_IntGauge{
				IntGauge: &otlp.IntGauge{
					DataPoints: []*otlp.IntDataPoint{
						getIntDataPoint(lbs1,intVal1,time1),
						nil,
					},
				},
			},
		},
		unmatchedBoundBucketIntHist: {
			Name: unmatchedBoundBucketIntHist,
			Data:
			&otlp.Metric_IntHistogram{
				IntHistogram: &otlp.IntHistogram{
					DataPoints: []*otlp.IntHistogramDataPoint{
						{
							ExplicitBounds:[]float64{0.1,0.2,0.3},
							BucketCounts:[]uint64{1,2},
						},
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		unmatchedBoundBucketDoubleHist: {
			Name:unmatchedBoundBucketDoubleHist,
			Data:
			&otlp.Metric_DoubleHistogram{
				DoubleHistogram: &otlp.DoubleHistogram{
					DataPoints: []*otlp.DoubleHistogramDataPoint{
						{
							ExplicitBounds:[]float64{0.1,0.2,0.3},
							BucketCounts:[]uint64{1,2},
						},
					},
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
	}

	nilMetric = "nil"
	empty = "empty"

	// Category 1: type and data field doesn't match
	notMatchIntGauge = "noMatchIntGauge"
	notMatchDoubleGauge = "notMatchDoubleGauge"
	notMatchIntSum = "notMatchIntSum"
	notMatchDoubleSum = "notMatchDoubleSum"
	notMatchIntHistogram = "notMatchIntHistogram"
	notMatchDoubleHistogram = "notMatchDoubleHistogram"

	// Category 2: invalid type and temporality combination
	invalidIntSum = "invalidIntSum"
	invalidDoubleSum = "invalidDoubleSum"
	invalidIntHistogram = "invalidIntHistogram"
	invalidDoubleHistogram = "invalidDoubleHistogram"

	//Category 3: nil data points
	nilDataPointIntGauge = "nilDataPointIntGauge"
	nilDataPointDoubleGauge = "nilDataPointDoubleGauge"
	nilDataPointIntSum = "nilDataPointIntSum"
	nilDataPointDoubleSum = "nilDataPointDoubleSum"
	nilDataPointIntHistogram = "nilDataPointIntHistogram"
	nilDataPointDoubleHistogram = "nilDataPointDoubleHistogram"

	// different metrics that will not pass validate metrics
	invalidMetrics = map[string]*otlp.Metric{
		// nil
		nilMetric: nil,
		// Data = nil
		empty: {},
		notMatchIntGauge: {
			Name: notMatchIntGauge,
			Data: &otlp.Metric_IntGauge{},
		},
		notMatchDoubleGauge: {
			Name: notMatchDoubleGauge,
			Data: &otlp.Metric_DoubleGauge{},
		},
		notMatchIntSum: {
			Name: notMatchIntSum,
			Data: &otlp.Metric_IntSum{},
		},
		notMatchDoubleSum: {
			Name: notMatchDoubleSum,
			Data: &otlp.Metric_DoubleSum{},
		},
		notMatchIntHistogram: {
			Name: notMatchIntHistogram,
			Data: &otlp.Metric_IntHistogram{},
		},
		notMatchDoubleHistogram: {
			Name: notMatchDoubleHistogram,
			Data: &otlp.Metric_DoubleHistogram{},
		},
		invalidIntSum: {
			Name: invalidIntSum,
			Data: &otlp.Metric_IntSum{
				IntSum:
					&otlp.IntSum{
					AggregationTemporality: otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
				},
			},
		},
		invalidDoubleSum: {
			Name: invalidDoubleSum,
			Data: &otlp.Metric_DoubleSum{
				DoubleSum:
				&otlp.DoubleSum{
					AggregationTemporality: otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
				},
			},
		},
		invalidIntHistogram: {
			Name: invalidIntHistogram,
			Data: &otlp.Metric_IntHistogram{
				IntHistogram:
				&otlp.IntHistogram{
					AggregationTemporality: otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
				},
			},
		},
		invalidDoubleHistogram: {
			Name: invalidDoubleHistogram,
			Data: &otlp.Metric_DoubleHistogram{
				DoubleHistogram:
				&otlp.DoubleHistogram{
					AggregationTemporality: otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
				},
			},
		},
	}

	// different metrics that will cause the exporter to return an error
	errorMetrics = map[string]*otlp.Metric{

		nilDataPointIntGauge: {
			Name: nilDataPointIntGauge,
			Data: &otlp.Metric_IntGauge{
				IntGauge: &otlp.IntGauge{DataPoints:nil},
			},
		},
		nilDataPointDoubleGauge: {
			Name: nilDataPointDoubleGauge,
			Data: &otlp.Metric_DoubleGauge{
				DoubleGauge: &otlp.DoubleGauge{DataPoints:nil},
			},
		},
		nilDataPointIntSum: {
			Name: nilDataPointIntSum,
			Data: &otlp.Metric_IntSum{
				IntSum: &otlp.IntSum{
					DataPoints:nil,
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		nilDataPointDoubleSum: {
			Name: nilDataPointDoubleSum,
			Data: &otlp.Metric_DoubleSum{
				DoubleSum: &otlp.DoubleSum{
					DataPoints:nil,
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		nilDataPointIntHistogram: {
			Name: nilDataPointIntHistogram,
			Data: &otlp.Metric_IntHistogram{
				IntHistogram: &otlp.IntHistogram{
					DataPoints:nil,
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},
		nilDataPointDoubleHistogram: {
			Name: nilDataPointDoubleHistogram,
			Data: &otlp.Metric_DoubleHistogram{
				DoubleHistogram: &otlp.DoubleHistogram{
					DataPoints:nil,
					AggregationTemporality:otlp.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				},
			},
		},

	}
	)

// Test_validateMetrics checks validateMetrics return true if a type and temporality combination is valid, false
// otherwise.
func Test_validateMetrics(t *testing.T) {

	// define a single test
	type combTest struct {
		name string
		metric *otlp.Metric
		want bool
	}

	tests := []combTest{}

	// append true cases
	for k, validMetric := range validMetrics1 {
		name := "valid_" + k

		tests = append(tests, combTest{
			name,
			validMetric,
			true,
		})
	}

	// append nil case
	tests = append(tests, combTest{"invalid_nil", nil, false})

	for k, invalidMetric := range invalidMetrics {
		name := "valid_" + k

		tests = append(tests, combTest{
			name,
			invalidMetric,
			false,
		})
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateMetrics(tt.metric)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Test_addSample checks addSample updates the map it receives correctly based on the sample and Label
// set it receives.
// Test cases are two samples belonging to the same TimeSeries,  two samples belong to different TimeSeries, and nil
// case.
func Test_addSample(t *testing.T) {
	type testCase struct {
		metric   *otlp.Metric
		sample prompb.Sample
		labels []prompb.Label
	}

	tests := []struct {
		name     string
		orig     map[string]*prompb.TimeSeries
		testCase []testCase
		want     map[string]*prompb.TimeSeries
	}{
		{
			"two_points_same_ts_same_metric",
			map[string]*prompb.TimeSeries{},
			[]testCase{
				{validMetrics1[validDoubleGauge],
					getSample(floatVal1, msTime1),
					promLbs1,
				},
				{
					validMetrics1[validDoubleGauge],
					getSample(floatVal2, msTime2),
					promLbs1,
				},
			},
			twoPointsSameTs,
		},
		{
			"two_points_different_ts_same_metric",
			map[string]*prompb.TimeSeries{},
			[]testCase{
				{validMetrics1[validIntGauge],
					getSample(float64(intVal1), msTime1),
					promLbs1,
				},
				{validMetrics1[validIntGauge],
					getSample(float64(intVal1), msTime2),
					promLbs2,
				},
			},
			twoPointsDifferentTs,
		},
	}
	t.Run("nil_case", func(t *testing.T) {
		tsMap := map[string]*prompb.TimeSeries{}
		addSample(tsMap, nil, nil, nil)
		assert.Exactly(t, tsMap, map[string]*prompb.TimeSeries{})
	})
	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addSample(tt.orig, &tt.testCase[0].sample, tt.testCase[0].labels, tt.testCase[0].metric)
			addSample(tt.orig, &tt.testCase[1].sample, tt.testCase[1].labels, tt.testCase[1].metric)
			assert.Exactly(t, tt.want, tt.orig)
		})
	}
}

// Test_timeSeries checks timeSeriesSignature returns consistent and unique signatures for a distinct label set and
// metric type combination.
func Test_timeSeriesSignature(t *testing.T) {
	tests := []struct {
		name string
		lbs  []prompb.Label
		metric *otlp.Metric
		want string
	}{
		{
			"int64_signature",
			promLbs1,
			validMetrics1[validIntGauge],
			strconv.Itoa(int(pdata.MetricDataTypeIntGauge)) + lb1Sig,
		},
		{
			"histogram_signature",
			promLbs2,
			validMetrics1[validIntHistogram],
			strconv.Itoa(int(pdata.MetricDataTypeIntHistogram)) + lb2Sig,
		},
		{
			"unordered_signature",
			getPromLabels(label22, value22, label21, value21),
			validMetrics1[validIntHistogram],
			strconv.Itoa(int(pdata.MetricDataTypeIntHistogram)) + lb2Sig,
		},
		// descriptor type cannot be nil, as checked by validateMetrics
		{
			"nil_case",
			nil,
			validMetrics1[validIntHistogram],
			strconv.Itoa(int(pdata.MetricDataTypeIntHistogram)),
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.want, timeSeriesSignature(tt.metric, &tt.lbs))
		})
	}
}

// Test_createLabelSet checks resultant label names are sanitized and label in extra overrides label in labels if
// collision happens. It does not check whether labels are not sorted
func Test_createLabelSet(t *testing.T) {
	tests := []struct {
		name   string
		orig   []*common.StringKeyValue
		extras []string
		want   []prompb.Label
	}{
		{
			"labels_clean",
			lbs1,
			[]string{label31, value31, label32, value32},
			getPromLabels(label11, value11, label12, value12, label31, value31, label32, value32),
		},
		{
			"labels_duplicate_in_extras",
			lbs1,
			[]string{label11, value31},
			getPromLabels(label11, value31, label12, value12),
		},
		{
			"labels_dirty",
			lbs1Dirty,
			[]string{label31 + dirty1, value31, label32, value32},
			getPromLabels(label11+"_", value11, "key_"+label12, value12, label31+"_", value31, label32, value32),
		},
		{
			"no_original_case",
			nil,
			[]string{label31, value31, label32, value32},
			getPromLabels(label31, value31, label32, value32),
		},
		{
			"empty_extra_case",
			lbs1,
			[]string{"", ""},
			getPromLabels(label11, value11, label12, value12, "", ""),
		},
		{
			"single_left_over_case",
			lbs1,
			[]string{label31, value31, label32},
			getPromLabels(label11, value11, label12, value12, label31, value31),
		},
	}
	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, createLabelSet(tt.orig, tt.extras...))
		})
	}
}

// Tes_getPromMetricName checks if OTLP metric names are converted to Cortex metric names correctly.
// Test cases are empty namespace, monotonic metrics that require a total suffix, and metric names that contains
// invalid characters.
func Test_getPromMetricName(t *testing.T) {
	tests := []struct {
		name string
		metric *otlp.Metric
		ns   string
		want string
	}{
		{
			"nil_case",
			nil,
			ns1,
			"",
		},
		{
			"normal_case",
			validMetrics1[validDoubleGauge],
			ns1,
			"test_ns_" + validDoubleGauge,
		},
		{
			"empty_namespace",
			validMetrics1[validDoubleGauge],
			"",
			validDoubleGauge,
		},
		{
			"total_suffix",
			validMetrics1[validIntSum],
			ns1,
			"test_ns_" + validIntSum + delimeter + totalStr,
		},
		{
			"dirty_string",
			validMetrics2[validIntGaugeDirty],
			"7" + ns1,
			"key_7test_ns__"+ validIntGauge + "_",
		},
	}
	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getPromMetricName(tt.metric, tt.ns))
		})
	}
}