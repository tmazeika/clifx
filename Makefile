all:
	mkdir -p out
	GOOS=darwin  GOARCH=386   go build -o clifx     main.go
	tar -acf out/clifx-mac-x86.tar.gz clifx
	GOOS=darwin  GOARCH=amd64 go build -o clifx     main.go
	tar -acf out/clifx-mac-x64.tar.gz clifx
	GOOS=linux   GOARCH=386   go build -o clifx     main.go
	tar -acf out/clifx-linux-x86.tar.gz clifx
	GOOS=linux   GOARCH=amd64 go build -o clifx     main.go
	tar -acf out/clifx-linux-x64.tar.gz clifx
	GOOS=linux   GOARCH=arm   go build -o clifx     main.go
	tar -acf out/clifx-linux-arm.tar.gz clifx
	GOOS=linux   GOARCH=arm64 go build -o clifx     main.go
	tar -acf out/clifx-linux-arm64.tar.gz clifx
	GOOS=windows GOARCH=386   go build -o clifx.exe main.go
	zip out/clifx-windows-x86.zip clifx.exe
	GOOS=windows GOARCH=amd64 go build -o clifx.exe main.go
	zip out/clifx-windows-x64.zip clifx.exe
	rm clifx
	rm clifx.exe
