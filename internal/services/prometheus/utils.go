package prometheus

import (
	prometheusv1 "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheus/v1"
	"github.com/prometheus/common/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func VectorToQueryResponse(vector model.Vector) *prometheusv1.QueryResponse {
	data := make([]*prometheusv1.Sample, vector.Len())

	for i, sample := range vector {
		data[i] = prometheusv1.Sample_builder{
			Metric: metricToMap(sample.Metric),
			Value: prometheusv1.Value_builder{
				Timestamp: timestamppb.New(sample.Timestamp.Time()),
				Value:     float64(sample.Value),
			}.Build(),
		}.Build()
	}

	res := prometheusv1.QueryResponse_builder{
		Data: data,
	}

	return res.Build()
}

func MatrixToQueryRangeResponse(matrix model.Matrix) *prometheusv1.QueryRangeResponse {
	data := make([]*prometheusv1.SampleStream, matrix.Len())

	for i, sampleStream := range matrix {
		values := make([]*prometheusv1.Value, len(sampleStream.Values))
		for j, v := range sampleStream.Values {
			values[j] = prometheusv1.Value_builder{
				Timestamp: timestamppb.New(v.Timestamp.Time()),
				Value:     float64(v.Value),
			}.Build()
		}
		data[i] = prometheusv1.SampleStream_builder{
			Metric: metricToMap(sampleStream.Metric),
			Values: values,
		}.Build()
	}

	res := prometheusv1.QueryRangeResponse_builder{
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
