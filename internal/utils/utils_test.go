package utils

import (
	"reflect"
	"testing"
)

func TestToJSON(t *testing.T) {
	type s struct {
		A string
		B int
		C bool
	}

	type args struct {
		stc interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil argument",
			args: args{
				stc: nil,
			},
			want: []byte(""),
		},
		{
			name: "empty string argument",
			args: args{
				stc: "",
			},
			want: []byte(""),
		},
		{
			name: "empty string argument",
			args: args{
				stc: s{
					A: "my string",
					B: 0,
				},
			},
			want: []byte(`{
  "A": "my string",
  "B": 0,
  "C": false
}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSON(tt.args.stc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %s, want %v", string(got), string(tt.want))
			}
		})
	}
}
