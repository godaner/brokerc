build:buildlinux upxlinux buildwin64 upxwind64 clearupx
buildlinux:
	-rm ./bin/brokerclinux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/brokerclinux
upxlinux:
	-upx -9 ./bin/brokerclinux
buildwin64:
	-rm ./bin/brokercwin64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/brokercwin64.exe
upxwind64:
	-upx -9 ./bin/brokercwin64.exe
clearupx:
	-rm *.upx
