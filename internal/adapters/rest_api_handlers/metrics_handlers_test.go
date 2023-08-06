package handlers

import (
	"net/http"
	"testing"
)

func Test_metricsHandlers_UpdateMetric(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		mh   *metricsHandlers
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mh.UpdateMetric(tt.args.w, tt.args.r)
		})
	}
}
