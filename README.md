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
## Sample Session
```
# saves a copy of the output to /tmp/dupes.txt while showing it live in your terminal session
$ ~/bin/qwkdupfinder /mnt/media /mnt/restore /home/admin/go /home/admin/work_files/ | tee /tmp/dupes.txt

=================================================================================================================
DUPLICATES REPORT ON BASIS OF FILESIZE AND SAMPLING CONTENTS -- ✓ indicates verified duplicate (for files <100KB)
=================================================================================================================
MATCH:1  	SIZE:1153525978	"/mnt/media/home-movies/298329884.mov"
MATCH:1  	SIZE:1153525978	"/mnt/restore/0000000238.MOV"

MATCH:2         SIZE:2335727    "/mnt/media/pics//IMAG0107.jpg"
MATCH:2         SIZE:2335727    "/mnt/restore/pic-00000019238.jpg"

MATCH:3 ✓       SIZE:4764       "/home/admin/work_files/2019/best_app_evarrr/go/pkg/mod/github.com/go-sql-driver/mysql@v1.5.0/rows.go"
MATCH:3 ✓       SIZE:4764       "/home/admin/go/src/github.com/go-sql-driver/mysql/rows.go"
```
To use the output to help you remove duplicates, run this
```
$ egrep "(^$|^MATCH:)" /tmp/dupes.txt | sed -e 's/^[^"]*//' > /tmp/dupes.raw
```
and edit `/tmp/dupes.raw` to REMOVE ALL FILES **YOU WANT TO KEEP**.  
Once you're confident that the only files left are the ones you want to disappear, which may resemble my edits --
```
# file /tmp/dupes.raw

"/mnt/restore/0000000238.MOV"

"/mnt/restore/pic-00000019238.jpg"

"/home/admin/work_files/2019/best_app_evarrr/go/pkg/mod/github.com/go-sql-driver/mysql@v1.5.0/rows.go"

```
-- run this and re-read your list of files that are about to go away forever:
```
$ sed -e 's/^"/rm "/' /tmp/dupes.raw
rm "/mnt/restore/0000000238.MOV"

rm "/mnt/restore/pic-00000019238.jpg"

rm "/home/admin/work_files/2019/best_app_evarrr/go/pkg/mod/github.com/go-sql-driver/mysql@v1.5.0/rows.go"
```
If the above `rm` commands look good to you, complete the dedupliation:
```
# sed -e 's/^"/rm "/' /tmp/dupes.raw | sh
```
(remove the comment character to actually perform the deletion.)

```
