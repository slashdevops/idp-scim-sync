package deepcopy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSliceOfPointers_any(t *testing.T) {
	tests := []struct {
		name   string
		toTest []*any
		want   []*any
	}{
		{
			name:   "empty slice",
			toTest: []*any{},
			want:   []*any{},
		},
		{
			name:   "nil slice",
			toTest: nil,
			want:   nil,
		},
		{
			name:   "slice of nil pointers",
			toTest: []*any{nil, nil, nil},
			want:   []*any{nil, nil, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceOfPointers(tt.toTest)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("TestSliceOfPointers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSliceOfPointers_string(t *testing.T) {
	tests := []struct {
		name   string
		toTest []*string
		want   []*string
	}{
		{
			name:   "empty slice",
			toTest: []*string{},
			want:   []*string{},
		},
		{
			name:   "nil slice",
			toTest: nil,
			want:   nil,
		},
		{
			name:   "slice of nil pointers",
			toTest: []*string{nil, nil, nil},
			want:   []*string{nil, nil, nil},
		},
		{
			name:   "slice of pointers",
			toTest: []*string{new(string), new(string), new(string)},
			want:   []*string{new(string), new(string), new(string)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceOfPointers(tt.toTest)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("TestSliceOfPointers() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSliceOfPointers_struct(t *testing.T) {
	tests := []struct {
		name   string
		toTest []*struct{}
		want   []*struct{}
	}{
		{
			name:   "empty slice",
			toTest: []*struct{}{},
			want:   []*struct{}{},
		},
		{
			name:   "nil slice",
			toTest: nil,
			want:   nil,
		},
		{
			name:   "slice of nil pointers",
			toTest: []*struct{}{nil, nil, nil},
			want:   []*struct{}{nil, nil, nil},
		},
		{
			name:   "slice of pointers",
			toTest: []*struct{}{new(struct{}), new(struct{}), new(struct{})},
			want:   []*struct{}{new(struct{}), new(struct{}), new(struct{})},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceOfPointers(tt.toTest)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("TestSliceOfPointers() (-want +got):\n%s", diff)
			}
		})
	}
}
