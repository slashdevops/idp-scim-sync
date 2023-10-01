package aws

import (
	"reflect"
	"testing"
)

func Test_toJSON(t *testing.T) {
	type args struct {
		stc   interface{}
		ident []bool
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
			if got := toJSON(tt.args.stc, tt.args.ident...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
