package server

import (
	pb "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheus_proxy/v1"
	"github.com/prometheus/common/model"
)

func VectorToQueryResponse(vector model.Vector) *pb.QueryResponse {
	data := make([]*pb.Sample, vector.Len())

	for i, sample := range vector {
		data[i] = pb.Sample_builder{
			Metric: metricToMap(sample.Metric),
			Value: pb.Value_builder{
				Timestamp: int64(sample.Timestamp),
				Value:     float64(sample.Value),
			}.Build(),
		}.Build()
	}

	res := pb.QueryResponse_builder{
		Data: data,
	}

	return res.Build()
}

func MatrixToQueryRangeResponse(matrix model.Matrix) *pb.QueryRangeResponse {
	data := make([]*pb.SampleStream, matrix.Len())

	for i, sampleStream := range matrix {
		values := make([]*pb.Value, len(sampleStream.Values))
		for j, v := range sampleStream.Values {
			values[j] = pb.Value_builder{
				Timestamp: int64(v.Timestamp),
				Value:     float64(v.Value),
			}.Build()
		}
		data[i] = pb.SampleStream_builder{
			Metric: metricToMap(sampleStream.Metric),
			Values: values,
		}.Build()
	}

	res := pb.QueryRangeResponse_builder{
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
