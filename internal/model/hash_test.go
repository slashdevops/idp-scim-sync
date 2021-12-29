package model

import (
	"bytes"
	"encoding/gob"
	"testing"
)

type CustomStructArray struct {
	More []string
}

type CustomStruct struct {
	Name    string
	Age     int
	Friends []string
	Things  []CustomStructArray
}

type CustomStructSerialized struct {
	Name    string
	Age     int
	Friends []string
	Things  []CustomStructArray
}

func (cs CustomStructSerialized) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(cs.Name); err != nil {
		panic(err)
	}
	if err := enc.Encode(cs.Age); err != nil {
		panic(err)
	}
	if err := enc.Encode(cs.Friends); err != nil {
		panic(err)
	}
	if err := enc.Encode(cs.Things); err != nil {
		panic(err)
	}
	return buf.Bytes(), nil
}

func TestHash(t *testing.T) {
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
			want: Hash(""),
		},
		{
			name: "slice of string",
			args: args{
				value: []string{"a", "b", "c"},
			},
			want: Hash([]string{"a", "b", "c"}),
		},
		{
			name: "array of string",
			args: args{
				value: [3]string{"a", "b", "c"},
			},
			want: Hash([3]string{"a", "b", "c"}),
		},
		{
			name: "array of integers",
			args: args{
				value: [3]int{1, 2, 3},
			},
			want: Hash([3]int{1, 2, 3}),
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
			want: Hash([3]struct {
				A int
				B string
				C bool
			}{
				{1, "test 1", true},
				{2, "test 2", false},
				{3, "test 3", true},
			}),
		},
		{
			name: "CustomStruct type",
			args: args{
				value: CustomStruct{
					Name:    "John",
					Age:     30,
					Friends: []string{"a", "b", "c"},
					Things: []CustomStructArray{
						{
							More: []string{"a", "b", "c"},
						},
						{
							More: []string{"1", "2", "3"},
						},
					},
				},
			},
			want: Hash(CustomStruct{
				Name:    "John",
				Age:     30,
				Friends: []string{"a", "b", "c"},
				Things: []CustomStructArray{
					{
						More: []string{"a", "b", "c"},
					},
					{
						More: []string{"1", "2", "3"},
					},
				},
			}),
		},
		{
			name: "Pointer to CustomStruct type",
			args: args{
				value: &CustomStruct{
					Name:    "John",
					Age:     30,
					Friends: []string{"a", "b", "c"},
					Things: []CustomStructArray{
						{
							More: []string{"a", "b", "c"},
						},
						{
							More: []string{"1", "2", "3"},
						},
					},
				},
			},
			want: Hash(&CustomStruct{
				Name:    "John",
				Age:     30,
				Friends: []string{"a", "b", "c"},
				Things: []CustomStructArray{
					{
						More: []string{"a", "b", "c"},
					},
					{
						More: []string{"1", "2", "3"},
					},
				},
			}),
		},
		{
			name: "CustomStructSerialized type",
			args: args{
				value: CustomStructSerialized{
					Name:    "John",
					Age:     30,
					Friends: []string{"a", "b", "c"},
					Things: []CustomStructArray{
						{
							More: []string{"a", "b", "c"},
						},
						{
							More: []string{"1", "2", "3"},
						},
					},
				},
			},
			want: Hash(CustomStructSerialized{
				Name:    "John",
				Age:     30,
				Friends: []string{"a", "b", "c"},
				Things: []CustomStructArray{
					{
						More: []string{"a", "b", "c"},
					},
					{
						More: []string{"1", "2", "3"},
					},
				},
			}),
		},
		{
			name: "Pointer to CustomStructSerialized type",
			args: args{
				value: &CustomStructSerialized{
					Name:    "John",
					Age:     30,
					Friends: []string{"a", "b", "c"},
					Things: []CustomStructArray{
						{
							More: []string{"a", "b", "c"},
						},
						{
							More: []string{"1", "2", "3"},
						},
					},
				},
			},
			want: Hash(&CustomStructSerialized{
				Name:    "John",
				Age:     30,
				Friends: []string{"a", "b", "c"},
				Things: []CustomStructArray{
					{
						More: []string{"a", "b", "c"},
					},
					{
						More: []string{"1", "2", "3"},
					},
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.value); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
