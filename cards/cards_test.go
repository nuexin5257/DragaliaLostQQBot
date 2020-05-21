package cards

import (
	"reflect"
	"testing"
)

func TestFindCardIndex(t *testing.T) {
	type args struct {
		name []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"basic", args{[]string{"龙龙"}}, []int{146}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindCardIndex(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindCardIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
