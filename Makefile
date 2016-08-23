all:
	mkdir -p out
	GOOS=darwin  GOARCH=386   go build -o clifx     main.go
	tar -acf out/mac-x86.tar.gz clifx
	GOOS=darwin  GOARCH=amd64 go build -o clifx     main.go
	tar -acf out/mac-x64.tar.gz clifx
	GOOS=linux   GOARCH=386   go build -o clifx     main.go
	tar -acf out/linux-x86.tar.gz clifx
	GOOS=linux   GOARCH=amd64 go build -o clifx     main.go
	tar -acf out/linux-x64.tar.gz clifx
	GOOS=linux   GOARCH=arm   go build -o clifx     main.go
	tar -acf out/linux-arm.tar.gz clifx
	GOOS=linux   GOARCH=arm64 go build -o clifx     main.go
	tar -acf out/linux-arm64.tar.gz clifx
	GOOS=windows GOARCH=386   go build -o clifx.exe main.go
	zip out/windows-x86.zip clifx.exe
	GOOS=windows GOARCH=amd64 go build -o clifx.exe main.go
	zip out/windows-x64.zip clifx.exe
	rm clifx
	rm clifx.exe
