package prometheusproxy

import (
	prometheusproxyv1 "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheusproxy/v1"
	"github.com/prometheus/common/model"
)

func VectorToQueryResponse(vector model.Vector) *prometheusproxyv1.QueryResponse {
	data := make([]*prometheusproxyv1.Sample, vector.Len())

	for i, sample := range vector {
		data[i] = prometheusproxyv1.Sample_builder{
			Metric: metricToMap(sample.Metric),
			Value: prometheusproxyv1.Value_builder{
				Timestamp: int64(sample.Timestamp),
				Value:     float64(sample.Value),
			}.Build(),
		}.Build()
	}

	res := prometheusproxyv1.QueryResponse_builder{
		Data: data,
	}

	return res.Build()
}

func MatrixToQueryRangeResponse(matrix model.Matrix) *prometheusproxyv1.QueryRangeResponse {
	data := make([]*prometheusproxyv1.SampleStream, matrix.Len())

	for i, sampleStream := range matrix {
		values := make([]*prometheusproxyv1.Value, len(sampleStream.Values))
		for j, v := range sampleStream.Values {
			values[j] = prometheusproxyv1.Value_builder{
				Timestamp: int64(v.Timestamp),
				Value:     float64(v.Value),
			}.Build()
		}
		data[i] = prometheusproxyv1.SampleStream_builder{
			Metric: metricToMap(sampleStream.Metric),
			Values: values,
		}.Build()
	}

	res := prometheusproxyv1.QueryRangeResponse_builder{
		Data: data,
	}

	return res.Build()
}

func metricToMap(metric model.Metric) map[string]string {
	metricMap := make(map[string]string)
	for k, v := range metric {
		metricMap[string(k)] = string(v)
	}
	return metricMap
}
