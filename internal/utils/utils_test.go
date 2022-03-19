package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		name      string
		args      args
		wantError bool
		want      []byte
	}{
		{
			name: "nil argument",
			args: args{
				stc: nil,
			},
			wantError: false,
			want:      []byte(""),
		},
		{
			name: "empty string argument",
			args: args{
				stc: "",
			},
			wantError: false,
			want:      []byte(""),
		},
		{
			name: "bad argument",
			args: args{
				stc: map[string]interface{}{
					"this will fail when is serialize": make(chan int),
				},
			},
			wantError: true,
			want:      nil,
		},
		{
			name: "empty string argument",
			args: args{
				stc: s{
					A: "my string",
					B: 0,
				},
			},
			wantError: false,
			want: []byte(`{
  "A": "my string",
  "B": 0,
  "C": false
}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantError {
				assert.Panics(t, func() {
					ToJSON(tt.args.stc)
				})
			} else {
				if got := ToJSON(tt.args.stc); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ToJSON() = %s, want %v", string(got), string(tt.want))
				}
			}
		})
	}
}

func TestToYAML(t *testing.T) {
	type s struct {
		A string
		B int
		C bool
	}
	type args struct {
		stc interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		want      []byte
	}{
		{
			name: "nil argument",
			args: args{
				stc: nil,
			},
			wantError: false,
			want:      []byte(""),
		},
		{
			name: "empty string argument",
			args: args{
				stc: "",
			},
			wantError: false,
			want:      []byte(""),
		},
		{
			name: "bad argument",
			args: args{
				stc: map[string]interface{}{
					"this will fail when is serialize": make(chan int),
				},
			},
			wantError: true,
			want:      nil,
		},
		{
			name: "empty string argument",
			args: args{
				stc: s{
					A: "my string",
					B: 0,
				},
			},
			wantError: false,
			want: []byte(`a: my string
b: 0
c: false
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantError {
				assert.Panics(t, func() {
					ToYAML(tt.args.stc)
				})
			} else {
				if got := ToYAML(tt.args.stc); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ToYAML() = %s, want %v", string(got), string(tt.want))
				}
			}
		})
	}
}
