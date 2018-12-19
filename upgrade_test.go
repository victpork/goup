package goup

import (
	"reflect"
	"testing"
)

func TestVersionInfo_String(t *testing.T) {
	tests := []struct {
		name string
		vi   VersionInfo
		want string
	}{
		{
			"TestCase 1",
			VersionInfo{
				Major: 1,
				Minor: 10,
				Build: 3,
			},
			"1.10.3",
		}, {
			"TestCase 2",
			VersionInfo{
				Major:       1,
				Minor:       11,
				Build:       0,
				Beta:        true,
				BetaVersion: 2,
			},
			"1.11beta2",
		}, {
			"TestCase 3",
			VersionInfo{
				Major:     2,
				Minor:     12,
				Build:     0,
				RC:        true,
				RCVersion: 4,
			},
			"2.12rc4",
		}, {
			"TestCase 4",
			VersionInfo{
				Major: 1,
				Minor: 9,
				Build: 0,
			},
			"1.9",
		},
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
	// The following test will not pass for 100%: it will always fail in one of the case
	// as the archive extension is on a platform-variant constant
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
		{
			"TestCase 1",
			args{
				VersionInfo{
					Major: 1,
					Minor: 10,
					Build: 3,
				},
				"windows",
				"amd64",
			},
			"https://dl.google.com/go/go1.10.3.windows-amd64.zip",
		}, {
			"TestCase 2",
			args{
				VersionInfo{
					Major:       1,
					Minor:       12,
					Build:       0,
					Beta:        true,
					BetaVersion: 1,
				},
				"linux",
				"arm64",
			},
			"https://dl.google.com/go/go1.12beta1.linux-arm64.tar.gz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DownloadUrl(tt.args.version, tt.args.os, tt.args.arch); got != tt.want {
				t.Errorf("DownloadUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatestVerInfo(t *testing.T) {
	verInfo, err := LatestVersionInfo()
	if err != nil {
		t.Errorf("Error in LatestVersionInfo(): %v", err)
	}
	if len(verInfo) == 0 {
		t.Errorf("Empty version array returned")
	}
}

func TestExtractVersionInfo(t *testing.T) {
	tests := []struct {
		name            string
		args            string
		wantVersionInfo VersionInfo
		wantErr         bool
	}{
		{
			"TestCase 1",
			"1.10.1",
			VersionInfo{
				Major: 1,
				Minor: 10,
				Build: 1,
				Beta:  false,
				RC:    false,
			},
			false,
		}, {
			"TestCase 2",
			"1.11rc2",
			VersionInfo{
				Major:     1,
				Minor:     11,
				Build:     0,
				Beta:      false,
				RC:        true,
				RCVersion: 2,
			},
			false,
		}, {
			"TestCase 3",
			"1.13beta3",
			VersionInfo{
				Major:       1,
				Minor:       13,
				Build:       0,
				Beta:        true,
				RC:          false,
				BetaVersion: 3,
			},
			false,
		}, {
			"TestCase 4",
			"1.18",
			VersionInfo{
				Major: 1,
				Minor: 18,
				Build: 0,
				Beta:  false,
				RC:    false,
			},
			false,
		}, {
			"TestCase 5",
			"1.18rc",
			VersionInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersionInfo, err := ExtractVersionInfo(tt.args)
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
