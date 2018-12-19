# goup
Tool that helps you to retrieve latest Go update from Google

# Overview
Goup is a little utility that helps you to check and upgrade your local non-container Go version. It bascially does the following step-by-step:

1. Run `go version` to determine local Go version and `go env` for `$GOPATH`. The location of `go` is determined by supplied `-p` param, `$PATH` variable or default installation location (`/usr/local/go` or `C:\Go`)
2. Check https://go.googlesource.com/go/+refs to see if there is a version (tags starts with `go`), compare it against local version retrieved in (1). Can include beta, RC and latest major/minor version for comparison with `-b`, `-rc` and `-u`.
3. If there is a new version available, download it to temporary directory.
4. Backup existing Go installtion to temp.
5. Extract new Go archive to `$GOROOT`.
6. In case of an error, reverse backup to `$GOROOT`.

# Compile and run
```
go get -u github.com/mkishere/goup
go build cmd/goup.go

goup
```