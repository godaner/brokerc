pgk:pkglinuxamd64 pkglinux386 pkgwindowsamd64 pkgwindows386
# build
build:
	CGO_ENABLED=0 go build

# linux amd64 build and package
pkglinuxamd64:buildlinuxamd64 upxlinuxamd64
	tar -czf ./bin/brokerc_linux-amd64.tar.gz ./bin/brokerc
buildlinuxamd64:
	-rm ./bin/brokerc
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/brokerc
upxlinuxamd64:
	-upx -9 ./bin/brokerc

# linux 386 build and package
pkglinux386:buildlinux386 upxlinux386
	tar -czf ./bin/brokerc_linux-386.tar.gz ./bin/brokerc
buildlinux386:
	-rm ./bin/brokerc
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/brokerc
upxlinux386:
	-upx -9 ./bin/brokerc

# windwos amd64 build and package
pkgwindowsamd64:buildwindowsamd64 upxwindowsamd64
	tar -czf ./bin/brokerc_windows-amd64.tar.gz ./bin/brokerc.exe
buildwindowsamd64:
	-rm ./bin/brokerc.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/brokerc.exe
upxwindowsamd64:
	-upx -9 ./bin/brokerc.exe

# windwos 386 build and package
pkgwindows386:buildwindows386 upxwindows386
	tar -czf ./bin/brokerc_windows-386.tar.gz ./bin/brokerc.exe
buildwindows386:
	-rm ./bin/brokerc.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/brokerc.exe
upxwindows386:
	-upx -9 ./bin/brokerc.exe