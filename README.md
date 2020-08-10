# qwkdupefinder
A command-line app to identify duplicate files by comparing samples of files.  Especially useful for large sets of media files.

## Building and Running
Assuming you have a working Go v1.14+ setup,
```
cd $GOPATH/src
git clone https://github.com/feromax/qwkdupefinder
cd qwkdupefinder
go build -o qwkdupefinder cmd/qwkdupefinder.go
./qwkdupefinder -h
```

