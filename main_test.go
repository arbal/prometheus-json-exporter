package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func refString(s string) *string {
	return &s
}

func refFloat64(f float64) *float64 {
	return &f
}

func refMetricType(mt dto.MetricType) *dto.MetricType {
	return &mt
}

func TestWalkJSON(t *testing.T) {
	testData := []struct {
		name     string
		bytes    []byte
		expected []*dto.MetricFamily
	}{
		{
			name:  "float value",
			bytes: []byte(`{"x": 1.0}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "int value",
			bytes: []byte(`{"x": 1}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "bool value",
			bytes: []byte(`{"x": true}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
					},
				},
			},
		},
		{
			name:     "string value",
			bytes:    []byte(`{"x": "ok"}`),
			expected: nil,
		},
		{
			name:     "null value",
			bytes:    []byte(`{"x": null}`),
			expected: nil,
		},
		{
			name:  "array value",
			bytes: []byte(`{"x": [1, 2, 3]}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x::array_0"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("0"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("1"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(2.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("2"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(3.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "nested value",
			bytes: []byte(`{"x": {"y": 1}}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x::y"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "nested^2 value",
			bytes: []byte(`{"x": {"y": {"z": 1}}}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x::y::z"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "array in nested value",
			bytes: []byte(`{"x": {"y": [1, 2, 3]}}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x::y::array_0"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("0"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("1"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(2.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("2"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(3.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "array in array value",
			bytes: []byte(`{"x": [[1, 2], [3, 4]]}`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("x::array_0::array_1"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("0"),
								},
								&dto.LabelPair{
									Name:  refString("array_1_index"),
									Value: refString("0"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("0"),
								},
								&dto.LabelPair{
									Name:  refString("array_1_index"),
									Value: refString("1"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(2.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("1"),
								},
								&dto.LabelPair{
									Name:  refString("array_1_index"),
									Value: refString("0"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(3.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("1"),
								},
								&dto.LabelPair{
									Name:  refString("array_1_index"),
									Value: refString("1"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(4.0),
							},
						},
					},
				},
			},
		},
		{
			name:  "array at root",
			bytes: []byte(`[1, 2, 3]`),
			expected: []*dto.MetricFamily{
				&dto.MetricFamily{
					Name: refString("array_0"),
					Help: refString("Retrieved value"),
					Type: refMetricType(dto.MetricType_GAUGE),
					Metric: []*dto.Metric{
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("0"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(1.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("1"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(2.0),
							},
						},
						&dto.Metric{
							Label: []*dto.LabelPair{
								&dto.LabelPair{
									Name:  refString("array_0_index"),
									Value: refString("2"),
								},
							},
							Gauge: &dto.Gauge{
								Value: refFloat64(3.0),
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var jsonData interface{}
			err := json.Unmarshal(tt.bytes, &jsonData)
			if err != nil {
				t.Errorf("Error: %v", err)
			}

			registry := prometheus.NewRegistry()

			doWalkJSON("", jsonData, registry)
			actual, err := registry.Gather()
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Got: %+v, expected: %+v", actual, tt.expected)
			}
		})
	}
}
