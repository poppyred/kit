package generator

import (
	"testing"
)

func TestNewNewService(t *testing.T) {
	setDefaults()
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{name: "test service", args: args{name: "test"}, want: nil},
		{name: "test2 service", args: args{name: "test2"}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNewService(tt.args.name).Generate(); got != tt.want {
				t.Errorf("NewNewService() = %v, want %v", got, tt.want)
			}

		})
	}
}
