package prometheus

import (
	prometheuspb "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheus"
	"github.com/prometheus/common/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func buildQueryResponse(vector model.Vector) *prometheuspb.QueryResponse {
	samples := vectorToSamples(vector)
	res := prometheuspb.QueryResponse_builder{
		Data: samples,
	}
	return res.Build()
}

func buildStreamQueryResponse(vector model.Vector) *prometheuspb.StreamQueryResponse {
	samples := vectorToSamples(vector)
	res := prometheuspb.StreamQueryResponse_builder{
		Data: samples,
	}
	return res.Build()
}

func buildQueryRangeResponse(matrix model.Matrix) *prometheuspb.QueryRangeResponse {
	sampleStreams := matrixToSampleStreams(matrix)
	res := prometheuspb.QueryRangeResponse_builder{
		Data: sampleStreams,
	}
	return res.Build()
}

func buildStreamQueryRangeResponse(matrix model.Matrix) *prometheuspb.StreamQueryRangeResponse {
	sampleStreams := matrixToSampleStreams(matrix)
	res := prometheuspb.StreamQueryRangeResponse_builder{
		Data: sampleStreams,
	}
	return res.Build()
}

func vectorToSamples(vector model.Vector) []*prometheuspb.Sample {
	samples := make([]*prometheuspb.Sample, vector.Len())

	for i, sample := range vector {
		samples[i] = prometheuspb.Sample_builder{
			Metric: metricToMap(sample.Metric),
			Value: prometheuspb.Value_builder{
				Timestamp: timestamppb.New(sample.Timestamp.Time()),
				Value:     float64(sample.Value),
			}.Build(),
		}.Build()
	}

	return samples
}

func matrixToSampleStreams(matrix model.Matrix) []*prometheuspb.SampleStream {
	sampleStreams := make([]*prometheuspb.SampleStream, matrix.Len())

	for i, sampleStream := range matrix {
		values := make([]*prometheuspb.Value, len(sampleStream.Values))
		for j, v := range sampleStream.Values {
			values[j] = prometheuspb.Value_builder{
				Timestamp: timestamppb.New(v.Timestamp.Time()),
				Value:     float64(v.Value),
			}.Build()
		}
		sampleStreams[i] = prometheuspb.SampleStream_builder{
			Metric: metricToMap(sampleStream.Metric),
			Values: values,
		}.Build()
	}

	return sampleStreams
}

func metricToMap(metric model.Metric) map[string]string {
	metricMap := make(map[string]string)
	for k, v := range metric {
		metricMap[string(k)] = string(v)
	}
	return metricMap
}
