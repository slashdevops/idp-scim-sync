package convert

import (
	"reflect"
	"testing"
)

func TestToJSON(t *testing.T) {
	type args struct {
		stc interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil",
			args: args{
				stc: nil,
			},
			want: []byte(""),
		},
		{
			name: "empty",
			args: args{
				stc: "",
			},
			want: []byte(""),
		},
		{
			name: "string",
			args: args{
				stc: "test",
			},
			want: []byte("\"test\""),
		},
		{
			name: "int",
			args: args{
				stc: 1,
			},
			want: []byte("1"),
		},
		{
			name: "int64",
			args: args{
				stc: int64(1),
			},
			want: []byte("1"),
		},
		{
			name: "struct",
			args: args{
				stc: struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				}{
					Name: "test",
					Age:  1,
				},
			},
			want: []byte("{\"name\":\"test\",\"age\":1}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSON(tt.args.stc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJSONString(t *testing.T) {
	type args struct {
		stc   interface{}
		ident []bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{
				stc: nil,
			},
			want: "",
		},
		{
			name: "empty",
			args: args{
				stc: "",
			},
			want: "",
		},
		{
			name: "string",
			args: args{
				stc: "test",
			},
			want: "\"test\"",
		},
		{
			name: "int",
			args: args{
				stc: 1,
			},
			want: "1",
		},
		{
			name: "struct",
			args: args{
				stc: struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				}{
					Name: "test",
					Age:  1,
				},
			},
			want: `{"name":"test","age":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSONString(tt.args.stc, tt.args.ident...); got != tt.want {
				t.Errorf("ToJSONString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToYAML(t *testing.T) {
	type args struct {
		stc interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil",
			args: args{
				stc: nil,
			},
			want: []byte(""),
		},
		{
			name: "empty",
			args: args{
				stc: "",
			},
			want: []byte(""),
		},
		{
			name: "string",
			args: args{
				stc: "test",
			},
			want: []byte("test\n"),
		},
		{
			name: "int",
			args: args{
				stc: 1,
			},
			want: []byte("1\n"),
		},
		{
			name: "int64",
			args: args{
				stc: int64(1),
			},
			want: []byte("1\n"),
		},
		{
			name: "struct",
			args: args{
				stc: struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				}{
					Name: "test",
					Age:  1,
				},
			},
			want: []byte("name: test\nage: 1\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToYAML(tt.args.stc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}
