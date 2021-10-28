package version

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name    string
		prepare func()
		want    string
	}{
		{
			name: "empty",
			prepare: func() {
				Version = ""
			},
			want: "0.0.0",
		},
		{
			name: "fixed version",
			prepare: func() {
				Version = "1.1.1"
			},
			want: "1.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			if got := GetVersion(); got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVersionInfo(t *testing.T) {
	tests := []struct {
		name    string
		prepare func()
		want    string
	}{
		{
			name: "empty",
			prepare: func() {
				Version = ""
				Revision = ""
				Branch = ""
			},
			want: "(version=0.0.0, revision=0, branch=unknown)",
		},
		{
			name: "fixed version",
			prepare: func() {
				Version = "1.1.1"
				Revision = "1"
				Branch = "devel"
			},
			want: "(version=1.1.1, revision=1, branch=devel)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			if got := GetVersionInfo(); got != tt.want {
				t.Errorf("GetVersionInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVersionInfoExtended(t *testing.T) {
	tests := []struct {
		name    string
		prepare func()
		want    string
	}{
		{
			name: "empty",
			prepare: func() {
				Version = ""
				Revision = ""
				Branch = ""
				BuildDate = ""
				BuildUser = ""
				GoVersion = runtime.Version()
			},
			want: fmt.Sprintf("(version=0.0.0, revision=0, branch=unknown, go=%s, user=unknown, date=unknown)", runtime.Version()),
		},
		{
			name: "fixed version",
			prepare: func() {
				Version = "1.1.1"
				Revision = "1"
				Branch = "devel"
				BuildDate = "2021-11-01"
				BuildUser = "jhon"
				GoVersion = runtime.Version()
			},
			want: fmt.Sprintf("(version=1.1.1, revision=1, branch=devel, go=%s, user=jhon, date=2021-11-01)", runtime.Version()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			if got := GetVersionInfoExtended(); got != tt.want {
				t.Errorf("GetVersionInfoExtended() = %v, want %v", got, tt.want)
			}
		})
	}
}
