package memstorage

import (
	"reflect"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/models"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	repo := New()

	assert.NotNil(t, repo)
}

func Test_memStorage_UpdateGauge(t *testing.T) {
	type args struct {
		metric models.Metric
	}
	type want struct {
		metricName  string
		metricValue models.Gauge
	}
	tests := []struct {
		name string
		ms   *memStorage
		args args
		want want
	}{
		{
			name: "add gauge metric",
			ms:   &memStorage{data: make(map[string]models.MetricValue)},
			args: args{
				metric: models.Metric{
					Name: "test1",
					MetricValue: models.MetricValue{
						Gauge: 1,
					},
				},
			},
			want: want{
				metricName:  "test1",
				metricValue: models.Gauge(1),
			},
		},
		{
			name: "rewrite gauge metric",
			ms: &memStorage{data: map[string]models.MetricValue{
				"test2": {
					Gauge: 1,
				},
			}},
			args: args{
				metric: models.Metric{
					Name: "test2",
					MetricValue: models.MetricValue{
						Gauge: 10,
					},
				},
			},
			want: want{
				metricName:  "test2",
				metricValue: models.Gauge(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateGauge(tt.args.metric)
			have, ok := tt.ms.data[tt.args.metric.Name]

			assert.True(t, ok)
			assert.Equal(t, tt.want.metricValue, have.Gauge)
		})
	}
}

func Test_memStorage_UpdateCounter(t *testing.T) {
	type args struct {
		metric models.Metric
	}
	type want struct {
		metricName  string
		metricValue models.Counter
	}
	tests := []struct {
		name string
		ms   *memStorage
		args args
		want want
	}{
		{
			name: "add counter metric",
			ms:   &memStorage{data: make(map[string]models.MetricValue)},
			args: args{
				metric: models.Metric{
					Name: "test1",
					MetricValue: models.MetricValue{
						Counter: 1,
					},
				},
			},
			want: want{
				metricName:  "test1",
				metricValue: models.Counter(1),
			},
		},
		{
			name: "update counter metric",
			ms: &memStorage{data: map[string]models.MetricValue{
				"test2": {
					Counter: 1,
				},
			}},
			args: args{
				metric: models.Metric{
					Name: "test2",
					MetricValue: models.MetricValue{
						Counter: 10,
					},
				},
			},
			want: want{
				metricName:  "test2",
				metricValue: models.Counter(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateCounter(tt.args.metric)

			have, ok := tt.ms.data[tt.args.metric.Name]

			assert.True(t, ok)
			assert.Equal(t, tt.want.metricValue, have.Counter)
		})
	}
}

func Test_memStorage_Get(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		ms      *memStorage
		args    args
		want    models.Metric
		wantErr bool
	}{
		{
			name: "return metric",
			ms: &memStorage{data: map[string]models.MetricValue{
				"test11": {
					Counter: 11,
					Gauge:   22,
				},
			}},
			args: args{
				name: "test11",
			},
			want: models.Metric{
				Name: "test11",
				MetricValue: models.MetricValue{
					Counter: 11,
					Gauge:   22,
				},
			},
			wantErr: false,
		},
		{
			name: "return error",
			ms:   &memStorage{},
			args: args{
				name: "test11",
			},
			want:    models.Metric{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ms.Get(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("memStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
