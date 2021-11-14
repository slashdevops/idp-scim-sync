package hash

import (
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{value: ""},
			want: Get(""),
		},
		{
			name: "slice of string",
			args: args{
				value: []string{"a", "b", "c"},
			},
			want: Get([]string{"a", "b", "c"}),
		},
		{
			name: "array of string",
			args: args{
				value: [3]string{"a", "b", "c"},
			},
			want: Get([3]string{"a", "b", "c"}),
		},
		{
			name: "array of integers",
			args: args{
				value: [3]int{1, 2, 3},
			},
			want: Get([3]int{1, 2, 3}),
		},
		{
			name: "array of ordered structs",
			args: args{
				value: [3]struct {
					A int
					B string
					C bool
				}{
					{1, "test 1", true},
					{2, "test 2", false},
					{3, "test 3", true},
				},
			},
			want: Get([3]struct {
				A int
				B string
				C bool
			}{
				{1, "test 1", true},
				{2, "test 2", false},
				{3, "test 3", true},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.value); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetV2(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{value: ""},
			want: GetV2(""),
		},
		{
			name: "slice of string",
			args: args{
				value: []string{"a", "b", "c"},
			},
			want: GetV2([]string{"a", "b", "c"}),
		},
		{
			name: "array of string",
			args: args{
				value: [3]string{"a", "b", "c"},
			},
			want: GetV2([3]string{"a", "b", "c"}),
		},
		{
			name: "array of integers",
			args: args{
				value: [3]int{1, 2, 3},
			},
			want: GetV2([3]int{1, 2, 3}),
		},
		{
			name: "array of ordered structs",
			args: args{
				value: [3]struct {
					A int
					B string
					C bool
				}{
					{1, "test 1", true},
					{2, "test 2", false},
					{3, "test 3", true},
				},
			},
			want: GetV2([3]struct {
				A int
				B string
				C bool
			}{
				{1, "test 1", true},
				{2, "test 2", false},
				{3, "test 3", true},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetV2(tt.args.value); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
