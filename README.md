# qwkdupefinder
A command-line app to identify duplicate files by comparing samples of files.  Especially useful for large sets of media files.

## Building and Running

### with your golang environment
Assuming you have a working Go v1.14+ setup,
```
cd $GOPATH/src
git clone https://github.com/feromax/qwkdupefinder
cd qwkdupefinder
go build -o qwkdupefinder cmd/qwkdupefinder.go
./qwkdupefinder -h
```

### don't have golang?  if you have docker...
These steps will use a prepacked golang v1.14 container image to build the binary.  Assuming you'll want the binary in `~/bin`,

1. Download and start the golang container.
	```
	$ cd
	$ mkdir -p bin
	$ docker run -it --rm -v `pwd`/bin:/build golang:1.14.7-alpine
	```
2. Once the container image has downloaded and started, clone the source code repositoroy and build.
	```
	/go # cd src
	/go/src # apk add --no-cache git
	/go/src # git clone https://github.com/feromax/qwkdupefinder
	/go/src # go build -o /build/qwkdupefinder qwkdupefinder/cmd/qwkdupefinder.go
	/go/src # exit
	```
3. Magically, the built artifact (i.e., the binary) will be in `~/bin/`.  Run via 
	```
	./bin/qwkdupfinder -h
	```

