package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mkishere/goup"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	verbose   = kingpin.Flag("verbose", "Prints verbose messages.").Short('v').Bool()
	incBeta   = kingpin.Flag("beta", "Include Beta in list of consideration. True if local version is beta.").Short('b').Bool()
	incRC     = kingpin.Flag("rc", "Include Release Candidate in list of consideration. True if local version is RC.").Short('c').Bool()
	autoUpd   = kingpin.Flag("silent", "Auto download and upgrade local Go without confirmation.").Short('s').Bool()
	goExePath = kingpin.Arg("path", "Path to Go executable. If omitted, will use\n1. go executable on $PATH\n2. Go default installation path").String()
	jumpVer   = kingpin.Flag("upgrade", "Jump to latest version if available. If not set, will only update to latest build.").Short('u').Bool()
)

func main() {
	kingpin.Parse()
	goExeFullPath := filepath.Join(*goExePath, "go")

	printVerbose("Running command \"%v version\"\n", goExeFullPath)
	localVer, platform, arch, err := goup.LocalGoInfo(goExeFullPath)
	if err != nil {
		// Try default path
		printVerbose("Trying default installation directory %s", goup.DefaultInstallDir)
		goExeFullPath = filepath.Join(goup.DefaultInstallDir, "go")
		localVer, platform, arch, err = goup.LocalGoInfo(goExeFullPath)
		if err != nil {
			fmt.Println("Error when getting local Go infomration", err)
			return
		}
	}

	printVerbose("Running command \"%v env\"\n", goExeFullPath)
	gopath, err := goup.GoPath(goExeFullPath)
	if err != nil {
		fmt.Println("Error when getting local Go infomration", err)
		return
	}

	printVerbose("Local Go Info:(Version:%v, OS:%v, Arch:%v, GoHome:%v)\n", localVer, platform, arch, gopath)

	availVerList, err := goup.LatestVersionInfo()
	if err != nil {
		fmt.Println("Cannot retrieve version information", err)
	}

	// Assume user will like beta and RC if they are already using beta/RC
	if localVer.Beta {
		*incBeta = true
		*incRC = true
	}
	if localVer.RC {
		*incRC = true
	}

	var latestVer goup.VersionInfo
	for i := range availVerList {
		if !*incRC && availVerList[i].RC {
			continue
		}
		if !*incBeta && availVerList[i].Beta {
			continue
		}
		if !*jumpVer && (availVerList[i].Major != localVer.Major || availVerList[i].Minor != localVer.Minor) {
			continue
		}
		latestVer = availVerList[i]
		break
	}

	fmt.Printf("Latest version is %v\n", latestVer)
	if latestVer == localVer {
		fmt.Println("Your Go is at latest version. Exiting...")
		return
	}

	if !*autoUpd {
		var input string
	loop:
		for {
			fmt.Print("Do you want to download and upgrade now (Y/n):")
			fmt.Scanln(&input)
			switch input {
			case "Y", "y":
				break loop
			case "n", "N":
				return
			}
		}
	}

	latestGoBin, err := ioutil.TempFile("", "go"+latestVer.String()+arch+platform)
	if err != nil {
		fmt.Println("Cannot create temporary file:", err)
		return
	}
	defer latestGoBin.Close()

	dlUrl := goup.DownloadUrl(latestVer, platform, arch)
	fmt.Printf("Downloading from %s\n", dlUrl)
	fileSize, err := goup.DownloadPackage(dlUrl,
		func(totalSize int64, src io.Reader) error {
			// Create a progress bar in console for download
			bar := pb.New(int(totalSize)).SetUnits(pb.U_BYTES)
			bar.Start()
			reader := bar.NewProxyReader(src)
			defer reader.Close()
			_, err := io.Copy(latestGoBin, reader)
			if err != nil {
				fmt.Println("\nError occured while downloading:", err)
				return err
			}
			bar.FinishPrint("Download completed")
			return nil
		})

	if err != nil {
		fmt.Println("Cannot download file: ", err)
		return
	}

	// Backup current Go installation to temp directory
	fmt.Println("Backing up current Go to temporary directory")
	backupDir, err := ioutil.TempDir("", "gobackup-"+localVer.String()+"-")
	printVerbose("Backup location: %s\n", backupDir)
	if err != nil {
		fmt.Println("Error creating backup in temporary directory:", err)
		return
	}
	err = goup.RecursiveCopyDir(gopath, backupDir)
	if err != nil {
		fmt.Println("Error backing up in temporary directory:", err)
		return
	}

	// Remove current Go installation
	err = os.RemoveAll(gopath)
	printVerbose("Removing %s\n", gopath)
	if err != nil {
		fmt.Println("Error removing existing Go directory. Make sure goup runs with elevated permissions:", err)
		return
	}
	// Extract archive
	fmt.Printf("Extracting latest Go to %s\n", gopath)
	err = goup.ExtractArchive(latestGoBin, fileSize, gopath, printVerbose)
	if err != nil {
		printVerbose("Error: %v\n", err)
		fmt.Println("Error extracting new Go package, restoring...")
		err = restore(backupDir, gopath)
		if err != nil {
			fmt.Println("Unrecoverable error, please consider reinstall Go manually ", err)
		}
		return
	}

	// Verify
	newLocalVer, _, _, err := goup.LocalGoInfo(goExeFullPath)
	if err != nil || newLocalVer != latestVer {
		printVerbose("Error: %v\n", err)
		err = restore(backupDir, gopath)
		if err != nil {
			fmt.Println("Unrecoverable error, please consider reinstall Go manually ", err)
			return
		}
	}
}

func restore(backupPath, goPath string) (err error) {
	err = os.RemoveAll(goPath)
	if err != nil {
		return
	}
	err = os.MkdirAll(goPath, 0755)
	if err != nil {
		return
	}
	err = goup.RecursiveCopyDir(backupPath, goPath)
	return
}

func printVerbose(format string, a ...interface{}) {
	if *verbose {
		fmt.Printf(format, a...)
	}
}
