# bldver
Building in git version information into a Go language program

## Introduction

This demonstration Golang program demonstrates various techniques related to application versioning within the Go language environment. The references at the end also provide useful information for versioning and its use in Go programs.

## Building

The following build commands are for building this program on Linux or Windows. It should be straight forward to make something like this work your favorite Go program or module. **NB:** the variables main.Ver, main.Dat and main.Githash correspond to package scope variables in the main package. If version information was needed in an imported package, with this kind of version data being linked in the package would need to define these variables and use them in package code accordingly.  

### Linux
``` sh
	HEAD=`git rev-parse HEAD` go build -ldflags "-X main.Ver=`git describe
		--tags $HEAD` -X main.Dat=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
		-X main.Githash=`git rev-parse HEAD`" bldver.go
```

### Windows
``` bat
	git rev-parse HEAD > temp.txt
	set /p GITHASH=<temp.txt
	git describe --tags %GITHASH% > temp.txt
	set /p VER=<temp.txt
	go build -ldflags "-X main.Ver=%VER% -X appversion.Dat=%DATE%_%TIME%
		-X main.Githash=%GITHASH%" bldver.go
```

## Git Tags

Displaying the git hash for the HEAD
``` sh
git rev-parse HEAD
```

Listing all the tags
``` sh
git tag --list
```

List the last tag. This method requires two commands, which does not work well when passed to the linler by go build.
``` sh
git tag --sort=taggerdate | tail -1
```

List the githash for the HEAD
``` sh
$ HEAD=`git rev-list --tags --max-count=1` && echo $HEAD
3d1cc50d8f96d5a4a08e6c679a30ea16cb0f878d
```
Listing the githash for the HEAD and using it to display the associated tag.
``` sh
$ HEAD=`git rev-list --tags --max-count=1` && git describe --tags $HEAD
v0.0.1
```

Add a tag to the HEAD 
``` sh
git tag v0.0.10
```

Pushing tags
``` sh
git push <tag-name>
git push --tags
```

Deleting tags
``` sh
git tag -d v1.0.0
git push --delete <tag-name>
```


## References

- How to display modules version from inside of code
	https://stackoverflow.com/questions/62009264/golang-how-to-display-modules-version-from-inside-of-code
- Go Runtime - Returning the build information embedded in the running binary.
	https://pkg.go.dev/runtime/debug#ReadBuildInfo
- Golang Auto Build Versioning
	https://www.atatus.com/blog/golang-auto-build-versioning/
- Module version numbering
	https://go.dev/doc/modules/version-numbers#pseudo-version-number
- Semantic Versioning
	https://semver.org/
- Git Tags	
	https://initialcommit.com/blog/git-tag
