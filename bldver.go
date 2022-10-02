package main

/*
Program:	bldver.go
Version:	v0.0.10
Date:		01Oct2022
Author:		nzkiwi1g@gmail.com
Purpose:	Demonstration of git versioning in the Go language

Note:		The build command below wont include any hash if the repo has not
			changed. So its a good idea to change a file in the repo before
			the build is done.

Building:
	The following build commands are for building this program on Linux or
	Windows. It should be straight forward to make this work for your favorite
	Go program or module. NB: the variables main.Ver, main.Dat and main.Githash
	correspond to package scope variables in the main package. If an imported
	module was being built with this data being linked in the package would need
	to define these variables and use them in package code accordingly.

	On Linux use this pattern to build this Go program:
		HEAD=`git rev-parse HEAD` go build -ldflags "-X main.Ver=`git describe
			--tags $HEAD` -X main.Dat=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
			-X main.Githash=`git rev-parse HEAD`" bldver.go

	an old way:
		HEAD=`git rev-list --tags --max-count=1` go build -ldflags
			"-X main.Ver=`git describe --tags $HEAD`
			-X main.Dat=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
			-X main.Githash=$HEAD" bldver.go

	On Windows:
		git rev-parse HEAD > temp.txt
		set /p GITHASH=<temp.txt
		git describe --tags %GITHASH% > temp.txt
		set /p VER=<temp.txt
		go build -ldflags "-X main.Ver=%VER% -X appversion.Dat=%DATE%_%TIME%
			-X main.Githash=%GITHASH%" bldver.go

References:
	https://stackoverflow.com/questions/62009264/golang-how-to-display-modules-version-from-inside-of-code
	https://pkg.go.dev/runtime/debug#ReadBuildInfo
	https://www.atatus.com/blog/golang-auto-build-versioning/
	https://go.dev/doc/modules/version-numbers#pseudo-version-number
	https://initialcommit.com/blog/git-tag

Output:
	$ ./bldver
	bldver version v0.0.10 started
	bldver ended

	$ ./bldver -v
	bldver ver: v0.0.10 at: 2022-10-02_04:27:12AM githash: 36d2f5de81003116ceb92f3fa0fd9468a654e25c

	Build Info:
	go      go1.19.1
	path    command-line-arguments
	dep     github.com/mattn/go-colorable   v0.1.12 h1:jF+Du6AlPIjs2BiUiQlKOX0rt3SujHxPnksPKZbaA40=
	dep     github.com/mattn/go-isatty      v0.0.14 h1:yVuAays6BHfxijgZPzw+3Zlu5yQgKGP2/hcQbHb7S9Y=
	dep     github.com/rs/zerolog   v1.28.0 h1:MirSo27VyNi7RJYP3078AA1+Cyzd2GB66qy3aUHvsWY=
	dep     golang.org/x/sys        v0.0.0-20210927094055-39ccf1dd6fa6      h1:foEbQz/B0Oz6YIqu/69kfXPYeFQAuuMYFkjaqXzl5Wo=
	build   -compiler=gc
	build   -ldflags="-X main.Ver=v0.0.10 -X main.Dat=2022-10-02_04:27:12AM -X main.Githash=36d2f5de81003116ceb92f3fa0fd9468a654e25c"
	build   CGO_ENABLED=1
	build   CGO_CFLAGS=
	build   CGO_CPPFLAGS=
	build   CGO_CXXFLAGS=
	build   CGO_LDFLAGS=
	build   GOARCH=amd64
	build   GOOS=linux
	build   GOAMD64=v1


	Dependencies:
			&debug.Module{Path:"github.com/mattn/go-colorable", Version:"v0.1.12", Sum:"h1:jF+Du6AlPIjs2BiUiQlKOX0rt3SujHxPnksPKZbaA40=", Replace:(*debug.Module)(nil)}
			&debug.Module{Path:"github.com/mattn/go-isatty", Version:"v0.0.14", Sum:"h1:yVuAays6BHfxijgZPzw+3Zlu5yQgKGP2/hcQbHb7S9Y=", Replace:(*debug.Module)(nil)}
			&debug.Module{Path:"github.com/rs/zerolog", Version:"v1.28.0", Sum:"h1:MirSo27VyNi7RJYP3078AA1+Cyzd2GB66qy3aUHvsWY=", Replace:(*debug.Module)(nil)}
			&debug.Module{Path:"golang.org/x/sys", Version:"v0.0.0-20210927094055-39ccf1dd6fa6", Sum:"h1:foEbQz/B0Oz6YIqu/69kfXPYeFQAuuMYFkjaqXzl5Wo=", Replace:(*debug.Module)(nil)}

	$ ./bldver -b
	HEAD=`git rev-parse HEAD` go build -ldflags "-X main.Ver=`git describe --tags $HEAD` -X main.Dat=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" bldver.go

*/

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/rs/zerolog"
)

const (
	bldCmdLinux = "HEAD=`git rev-parse HEAD` go build -ldflags \"-X main.Ver=`" +
		"git describe --tags $HEAD` -X main.Dat=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" +
		" -X main.Githash=`git rev-parse HEAD`\" bldver.go"
	bldCmdWindows = "git rev-parse HEAD > temp.txt\n" +
		"set /p GITHASH=<temp.txt\n" +
		"git describe --tags %GITHASH% > temp.txt\n" +
		"set /p VER=<temp.txt\n" +
		"go build -ldflags \"-X main.Ver=%VER% -X appversion.Dat=%DATE%_%TIME% " +
		"-X main.Githash=%GITHASH%\" bldver.go"
)

var (
	showVer, showHelp bool
	showBld           bool
	log               zerolog.Logger
)

func init() {
	log = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()

	flag.BoolVar(&showVer, "version", false, "Show the current version")
	flag.BoolVar(&showVer, "v", false, "Show the current version")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showBld, "b", false, "Show build command for this program")
}

var (
	// program name set by code developer at program creation time
	Pgm = "bldver"

	// date of last build replaced by go build command (see above)
	Dat = "2022-10-01_00:00:00AM"

	// version of program replaced by go build command from most
	// recent git tag
	Ver = "v0.0.10"

	// githash for latest revision of program replaced by go build
	// command from git hash for HEAD
	Githash = "0000000000000000000000000000000000000000"
)

func main() {
	flag.Parse()
	if showHelp {
		flag.PrintDefaults()
		return
	}
	if showVer {
		fmt.Printf("%s ver: %s at: %s githash: %s\n\n", Pgm, Ver, Dat, Githash)

		info, ok := debug.ReadBuildInfo()
		fmt.Printf("Build Info: \n%v\n\n", info)

		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			log.Printf("Failed to read build info")
			return
		}
		fmt.Printf("Dependencies:\n")
		for _, dep := range buildInfo.Deps {
			fmt.Printf("\t%#v\n", dep)
		}
		return
	}
	if showBld {
		if runtime.GOOS == "linux" {
			fmt.Println(bldCmdLinux)
		} else if runtime.GOOS == "windows" {
			fmt.Println(bldCmdWindows)
		} else {
			log.Printf("System type %q show build command not supported", runtime.GOOS)
		}
		return
	}
	fmt.Printf("%s version %s started\n", Pgm, Ver)
	defer fmt.Printf("%s ended\n", Pgm)
	// . . .
}
