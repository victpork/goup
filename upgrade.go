package goup

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	RelVerURL              = "https://go.googlesource.com/go/+refs"
	DownloadURLWithPattern = "https://dl.google.com/go/go[version].[os]-[arch].[ext]"
)

type VersionInfo struct {
	Major       int
	Minor       int
	Build       int
	RC          bool
	RCVersion   int
	Beta        bool
	BetaVersion int
}

func (vi VersionInfo) String() string {
	if vi.Beta {
		return fmt.Sprintf("%d.%dbeta%d", vi.Major, vi.Minor, vi.BetaVersion)
	} else if vi.RC {
		return fmt.Sprintf("%d.%drc%d", vi.Major, vi.Minor, vi.RCVersion)
	} else {
		return fmt.Sprintf("%d.%d.%d", vi.Major, vi.Minor, vi.Build)
	}
}

func DownloadUrl(version VersionInfo, os, arch string) string {
	replacer := strings.NewReplacer("[version]", version.String(),
		"[arch]", arch,
		"[ext]", Format,
		"[os]", os)

	return replacer.Replace(DownloadURLWithPattern)
}

// DownloadPackage downloads Go compiled binaries from dl.google.com
// version: Go version
// arch: Go architecture
func DownloadPackage(url string, dlCallback func(totalSize int64, src io.Reader) error) (size int64, err error) {

	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	err = dlCallback(resp.ContentLength, resp.Body)
	return resp.ContentLength, err
}

// LatestVersionInfo returns all version (defined in GoogleSource) available in a slice
func LatestVersionInfo() (versionInfo []VersionInfo, err error) {
	resp, err := http.Get(RelVerURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Error code: " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	verList := make([]VersionInfo, 0, 10)
	doc.Find(".RefList-item").Each(func(i int, s *goquery.Selection) {
		verStr := s.Find("a").Text()
		if strings.HasPrefix(verStr, "go") {
			verInfo, err := ExtractVersionInfo(verStr[2:])
			if err == nil {
				verList = append(verList, verInfo)
			}
		}
	})

	// Sort all version with latest go first
	// Standard build > RC > Beta
	sort.Slice(verList, func(i, j int) bool {
		if verList[i].Major < verList[j].Major {
			return false
		}
		if verList[i].Major > verList[j].Major {
			return true
		}
		if verList[i].Minor < verList[j].Minor {
			return false
		}
		if verList[i].Minor > verList[j].Minor {
			return true
		}
		if verList[i].Build < verList[j].Build {
			return false
		}
		if verList[i].Build > verList[j].Build {
			return true
		}
		if (verList[i].Beta || verList[i].RC) && !(verList[j].RC || verList[j].Beta) {
			return false
		}
		if !(verList[i].Beta || verList[i].RC) && (verList[j].RC || verList[j].Beta) {
			return true
		}
		if verList[i].Beta && verList[j].RC {
			return false
		}
		if verList[i].RC && verList[j].Beta {
			return true
		}
		if verList[i].RCVersion < verList[j].RCVersion {
			return false
		}
		if verList[i].RCVersion > verList[j].RCVersion {
			return true
		}
		if verList[i].BetaVersion < verList[j].BetaVersion {
			return false
		}
		if verList[i].BetaVersion > verList[j].BetaVersion {
			return true
		}

		return true
	})

	return verList, nil
}

// LocalGoInfo returns local Go version numbers, OS and Arch
func LocalGoInfo(exePath string) (ver, os, arch string, err error) {
	verCmd := exec.Command(exePath, "version")
	out, err := verCmd.Output()
	if err != nil {
		return "", "", "", err
	}
	if len(out) <= 0 {
		return "", "", "", errors.New("go version returns improper")
	}
	goVerInfo := strings.Split(string(out), " ")
	archOSInfo := strings.Split(goVerInfo[3], "/")

	return strings.TrimSpace(goVerInfo[2][2:]), strings.TrimSpace(archOSInfo[0]), strings.TrimSpace(archOSInfo[1]), nil
}

// GoPath extract GOPATH path from `go env` command
func GoPath(exePath string) (string, error) {
	verCmd := exec.Command(exePath, "env")
	out, err := verCmd.Output()
	if err != nil {
		return "", err
	}
	envStr := string(out)
	scanner := bufio.NewScanner(strings.NewReader(envStr))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "GOROOT") {
			varPair := strings.Split(scanner.Text(), "=")
			return strings.Trim(varPair[1], "\""), nil
		}
	}
	return "", errors.New("GOROOT not found")
}

func ExtractVersionInfo(version string) (versionInfo VersionInfo, err error) {
	verArr := strings.Split(version, ".")
	if len(verArr) != 3 && len(verArr) != 2 {
		return VersionInfo{}, errors.New("Cannot parse version correctly")
	}

	versionInfo.Major, err = strconv.Atoi(verArr[0])
	if err != nil {
		return VersionInfo{}, errors.New("Cannot parse major version")
	}

	if strings.Contains(verArr[1], "beta") {
		versionInfo.Beta = true
		betaInfo := strings.Split(verArr[1], "beta")
		versionInfo.Minor, err = strconv.Atoi(betaInfo[0])
		if err != nil {
			return VersionInfo{}, errors.New("Cannot parse minor version")
		}
		versionInfo.BetaVersion, err = strconv.Atoi(betaInfo[1])
		if err != nil {
			return VersionInfo{}, errors.New("Cannot parse minor version")
		}
	} else if strings.Contains(verArr[1], "rc") {
		versionInfo.RC = true
		rcInfo := strings.Split(verArr[1], "rc")
		versionInfo.Minor, err = strconv.Atoi(rcInfo[0])
		if err != nil {
			return VersionInfo{}, errors.New("Cannot parse minor version")
		}
		versionInfo.RCVersion, err = strconv.Atoi(rcInfo[1])
		if err != nil {
			return VersionInfo{}, errors.New("Cannot parse minor version")
		}
	} else {
		versionInfo.Minor, err = strconv.Atoi(verArr[1])
		if err != nil {
			return VersionInfo{}, errors.New("Cannot parse minor version")
		}
	}

	if len(verArr) == 2 {
		return
	}

	versionInfo.Build, err = strconv.Atoi(verArr[2])
	if err != nil {
		return VersionInfo{}, errors.New("Cannot parse build version")
	}

	return
}
