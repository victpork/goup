package goup

import (
	"io"
	"reflect"
	"testing"
)

func TestVersionInfo_String(t *testing.T) {
	tests := []struct {
		name string
		vi   VersionInfo
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vi.String(); got != tt.want {
				t.Errorf("VersionInfo.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadUrl(t *testing.T) {
	type args struct {
		version VersionInfo
		os      string
		arch    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DownloadUrl(tt.args.version, tt.args.os, tt.args.arch); got != tt.want {
				t.Errorf("DownloadUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadPackage(t *testing.T) {
	type args struct {
		url        string
		dlCallback func(totalSize int64, src io.Reader) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadPackage(tt.args.url, tt.args.dlCallback); (err != nil) != tt.wantErr {
				t.Errorf("DownloadPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLatestVersionInfo(t *testing.T) {
	tests := []struct {
		name            string
		wantVersionInfo []VersionInfo
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersionInfo, err := LatestVersionInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("LatestVersionInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVersionInfo, tt.wantVersionInfo) {
				t.Errorf("LatestVersionInfo() = %v, want %v", gotVersionInfo, tt.wantVersionInfo)
			}
		})
	}
}

func TestLocalGoInfo(t *testing.T) {
	type args struct {
		exePath string
	}
	tests := []struct {
		name     string
		args     args
		wantVer  string
		wantOs   string
		wantArch string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVer, gotOs, gotArch, err := LocalGoInfo(tt.args.exePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalGoInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVer != tt.wantVer {
				t.Errorf("LocalGoInfo() gotVer = %v, want %v", gotVer, tt.wantVer)
			}
			if gotOs != tt.wantOs {
				t.Errorf("LocalGoInfo() gotOs = %v, want %v", gotOs, tt.wantOs)
			}
			if gotArch != tt.wantArch {
				t.Errorf("LocalGoInfo() gotArch = %v, want %v", gotArch, tt.wantArch)
			}
		})
	}
}

func TestGoPath(t *testing.T) {
	type args struct {
		exePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoPath(tt.args.exePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractVersionInfo(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name            string
		args            args
		wantVersionInfo VersionInfo
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersionInfo, err := ExtractVersionInfo(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVersionInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVersionInfo, tt.wantVersionInfo) {
				t.Errorf("ExtractVersionInfo() = %v, want %v", gotVersionInfo, tt.wantVersionInfo)
			}
		})
	}
}
