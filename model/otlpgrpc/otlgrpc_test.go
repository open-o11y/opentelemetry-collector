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

package otlpgrpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/model/pdata"
)

func TestUnmarshalLogsRequest_error(t *testing.T) {
	r := NewLogsRequest()
	_, err := r.UnmarshalLogsRequest([]byte("#123"))
	assert.Error(t, err)
}

func TestUnmarshalMetricsRequest_error(t *testing.T) {
	r := NewMetricsRequest()
	_, err := r.UnmarshalMetricsRequest([]byte("#123"))
	assert.Error(t, err)
}

func TestUnmarshalTracesRequest_error(t *testing.T) {
	r := NewTracesRequest()
	_, err := r.UnmarshalTracesRequest([]byte("#123"))
	assert.Error(t, err)
}

func TestUnmarshalLogsRequest(t *testing.T) {
	r := NewLogsRequest()
	_, err := r.UnmarshalLogsRequest([]byte(""))
	assert.NoError(t, err)
}

func TestUnmarshalMetricsRequest(t *testing.T) {
	r := NewMetricsRequest()
	_, err := r.UnmarshalMetricsRequest([]byte(""))
	assert.NoError(t, err)
}

func TestUnmarshalTracesRequest(t *testing.T) {
	r := NewTracesRequest()
	_, err := r.UnmarshalTracesRequest([]byte(""))
	assert.NoError(t, err)
}

func TestMarshalLogsRequest(t *testing.T) {
	r := NewLogsRequest()
	td := pdata.NewLogs()
	_, err := r.MarshalLogsRequest(td)
	assert.NoError(t, err)
}

func TestMarshalMetricsRequest(t *testing.T) {
	r := NewMetricsRequest()
	td := pdata.NewMetrics()
	_, err := r.MarshalMetricsRequest(td)
	assert.NoError(t, err)
}

func TestMarshalTracesRequest(t *testing.T) {
	r := NewTracesRequest()
	td := pdata.NewTraces()
	_, err := r.MarshalTracesRequest(td)
	assert.NoError(t, err)
}

func TestUnmarshalLogsResponse_error(t *testing.T) {
	r := NewLogsResponse()
	err := r.UnmarshalLogsResponse([]byte("#123"))
	assert.Error(t, err)
}

func TestUnmarshalMetricsResponse_error(t *testing.T) {
	r := NewMetricsResponse()
	err := r.UnmarshalMetricsResponse([]byte("#123"))
	assert.Error(t, err)
}

func TestUnmarshalTracesResponse_error(t *testing.T) {
	r := NewTracesResponse()
	err := r.UnmarshalTracesResponse([]byte("#123"))
	assert.Error(t, err)
}

func TestUnmarshalLogsResponse(t *testing.T) {
	r := NewLogsResponse()
	err := r.UnmarshalLogsResponse([]byte(""))
	assert.NoError(t, err)
}

func TestUnmarshalMetricsResponse(t *testing.T) {
	r := NewMetricsResponse()
	err := r.UnmarshalMetricsResponse([]byte(""))
	assert.NoError(t, err)
}

func TestUnmarshalTracesResponse(t *testing.T) {
	r := NewTracesResponse()
	err := r.UnmarshalTracesResponse([]byte(""))
	assert.NoError(t, err)
}

func TestMarshalLogsResponse(t *testing.T) {
	r := NewLogsResponse()
	_, err := r.MarshalLogsResponse()
	assert.NoError(t, err)
}

func TestMarshalMetricsResponse(t *testing.T) {
	r := NewMetricsResponse()
	_, err := r.MarshalMetricsResponse()
	assert.NoError(t, err)
}

func TestMarshalTracesResponse(t *testing.T) {
	r := NewTracesResponse()
	_, err := r.MarshalTracesResponse()
	assert.NoError(t, err)
}

func BenchmarkLogsRequest_LogsToProtobuf(b *testing.B) {
	lr := NewLogsRequest()
	logs := generateBenchmarkLogs(128)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf, err := lr.MarshalLogsRequest(logs)
		require.NoError(b, err)
		assert.NotEqual(b, 0, len(buf))
	}
}

func BenchmarkLogsRequest_LogsFromProtobuf(b *testing.B) {
	lr := NewLogsRequest()
	baseLogs := generateBenchmarkLogs(128)
	buf, err := lr.MarshalLogsRequest(baseLogs)
	require.NoError(b, err)
	assert.NotEqual(b, 0, len(buf))
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		lr = NewLogsRequest()
		logs, err := lr.UnmarshalLogsRequest(buf)
		require.NoError(b, err)
		assert.Equal(b, baseLogs.ResourceLogs().Len(), logs.ResourceLogs().Len())
	}
}

func BenchmarkMetricsRequest_MetricsToProtobuf(b *testing.B) {
	mr := NewMetricsRequest()
	metrics := generateBenchmarkMetrics(128)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf, err := mr.MarshalMetricsRequest(metrics)
		require.NoError(b, err)
		assert.NotEqual(b, 0, len(buf))
	}
}

func BenchmarkMetricsRequest_MetricsFromProtobuf(b *testing.B) {
	mr := NewMetricsRequest()
	baseMetrics := generateBenchmarkMetrics(128)
	buf, err := mr.MarshalMetricsRequest(baseMetrics)
	require.NoError(b, err)
	assert.NotEqual(b, 0, len(buf))
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		mr = NewMetricsRequest()
		metrics, err := mr.UnmarshalMetricsRequest(buf)
		require.NoError(b, err)
		assert.Equal(b, baseMetrics.ResourceMetrics().Len(), metrics.ResourceMetrics().Len())
	}
}

func BenchmarkTracesRequest_TracesToProtobuf(b *testing.B) {
	tr := NewTracesRequest()
	traces := generateBenchmarkTraces(128)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf, err := tr.MarshalTracesRequest(traces)
		require.NoError(b, err)
		assert.NotEqual(b, 0, len(buf))
	}
}

func BenchmarkTracesRequest_TracesFromProtobuf(b *testing.B) {
	tr := NewTracesRequest()
	baseTraces := generateBenchmarkTraces(128)
	buf, err := tr.MarshalTracesRequest(baseTraces)
	require.NoError(b, err)
	assert.NotEqual(b, 0, len(buf))
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tr = NewTracesRequest()
		traces, err := tr.UnmarshalTracesRequest(buf)
		require.NoError(b, err)
		assert.Equal(b, baseTraces.ResourceSpans().Len(), traces.ResourceSpans().Len())
	}
}

func generateBenchmarkLogs(logsCount int) pdata.Logs {
	endTime := pdata.NewTimestampFromTime(time.Now())
	md := pdata.NewLogs()
	ilm := md.ResourceLogs().AppendEmpty().InstrumentationLibraryLogs().AppendEmpty()
	ilm.Logs().EnsureCapacity(logsCount)
	for i := 0; i < logsCount; i++ {
		im := ilm.Logs().AppendEmpty()
		im.SetName("test_name")
		im.SetTimestamp(endTime)
	}
	return md
}

func generateBenchmarkMetrics(metricsCount int) pdata.Metrics {
	now := time.Now()
	startTime := pdata.NewTimestampFromTime(now.Add(-10 * time.Second))
	endTime := pdata.NewTimestampFromTime(now)

	md := pdata.NewMetrics()
	ilm := md.ResourceMetrics().AppendEmpty().InstrumentationLibraryMetrics().AppendEmpty()
	ilm.Metrics().EnsureCapacity(metricsCount)
	for i := 0; i < metricsCount; i++ {
		im := ilm.Metrics().AppendEmpty()
		im.SetName("test_name")
		im.SetDataType(pdata.MetricDataTypeSum)
		idp := im.Sum().DataPoints().AppendEmpty()
		idp.SetStartTimestamp(startTime)
		idp.SetTimestamp(endTime)
		idp.SetIntVal(123)
	}
	return md
}

func generateBenchmarkTraces(metricsCount int) pdata.Traces {
	now := time.Now()
	startTime := pdata.NewTimestampFromTime(now.Add(-10 * time.Second))
	endTime := pdata.NewTimestampFromTime(now)

	md := pdata.NewTraces()
	ilm := md.ResourceSpans().AppendEmpty().InstrumentationLibrarySpans().AppendEmpty()
	ilm.Spans().EnsureCapacity(metricsCount)
	for i := 0; i < metricsCount; i++ {
		im := ilm.Spans().AppendEmpty()
		im.SetName("test_name")
		im.SetStartTimestamp(startTime)
		im.SetEndTimestamp(endTime)
	}
	return md
}
