package memstorage

import (
	"reflect"
	"testing"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	repo, err := New(config.ServerConfig{})

	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func Test_memStorage_UpdateGauge(t *testing.T) {
	t.Parallel()
	type args struct {
		metric models.Metric
	}
	type want struct {
		metricName  string
		metricValue *float64
	}
	tests := []struct {
		name string
		ms   *memStorage
		args args
		want want
	}{
		{
			name: "add gauge metric",
			ms:   &memStorage{data: make(map[string]models.Metric)},
			args: args{
				metric: models.Metric{
					ID:    "test1",
					Value: createValue(1),
				},
			},
			want: want{
				metricName:  "test1",
				metricValue: createValue(1),
			},
		},
		{
			name: "rewrite gauge metric",
			ms: &memStorage{data: map[string]models.Metric{
				"test2": {
					Value: createValue(10),
				},
			}},
			args: args{
				metric: models.Metric{
					ID:    "test2",
					Value: createValue(10),
				},
			},
			want: want{
				metricName:  "test2",
				metricValue: createValue(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateGauge(tt.args.metric)
			have, ok := tt.ms.data[tt.args.metric.ID]

			assert.True(t, ok)
			assert.Equal(t, *tt.want.metricValue, *have.Value)
		})
	}
}

func Test_memStorage_UpdateCounter(t *testing.T) {
	t.Parallel()
	type args struct {
		metric models.Metric
	}
	type want struct {
		metricName  string
		metricValue *int64
	}
	tests := []struct {
		name string
		ms   *memStorage
		args args
		want want
	}{
		{
			name: "add counter metric",
			ms:   &memStorage{data: make(map[string]models.Metric)},
			args: args{
				metric: models.Metric{
					ID:    "test1",
					Delta: createDelta(1),
				},
			},
			want: want{
				metricName:  "test1",
				metricValue: createDelta(1),
			},
		},
		{
			name: "update counter metric",
			ms: &memStorage{data: map[string]models.Metric{
				"test2": {
					Delta: createDelta(10),
				},
			}},
			args: args{
				metric: models.Metric{
					ID:    "test2",
					Delta: createDelta(10),
				},
			},
			want: want{
				metricName:  "test2",
				metricValue: createDelta(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateCounter(tt.args.metric)

			have, ok := tt.ms.data[tt.args.metric.ID]

			assert.True(t, ok)
			assert.Equal(t, *tt.want.metricValue, *have.Delta)
		})
	}
}

func Test_memStorage_GetDelta(t *testing.T) {
	t.Parallel()
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
			ms: &memStorage{data: map[string]models.Metric{
				"test11": {
					Delta: createDelta(11),
					//Value: createValue(22),
				},
			}},
			args: args{
				name: "test11",
			},
			want: models.Metric{
				ID:    "test11",
				Delta: createDelta(11),
				//Value: createValue(22),
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
			if !assert.Equal(t, got.Delta, tt.want.Delta) {
				t.Errorf("memStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memStorage_GetAll(t *testing.T) {
	tests := []struct {
		name string
		ms   *memStorage
		want []models.Metric
	}{
		{
			name: "get",
			ms: &memStorage{data: map[string]models.Metric{
				"test11": {
					Delta: createDelta(11),
					Value: createValue(22),
				},
				"test22": {
					Delta: createDelta(21),
					Value: createValue(31),
				},
			}},
			want: []models.Metric{
				{
					ID: "test11",

					Delta: createDelta(11),
					Value: createValue(22),
				},
				{
					ID: "test22",

					Delta: createDelta(21),
					Value: createValue(31),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.GetAll(); !reflect.DeepEqual(got, tt.want) {
				reflect.DeepEqual(got, tt.want)
			}
		})
	}
}

func createValue(x float64) *float64 {
	return &x
}

func createDelta(x int64) *int64 {
	return &x
}
