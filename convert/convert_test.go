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
