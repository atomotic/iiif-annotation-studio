bindata:
	go-bindata-assetfs static/...

build-macos:
	go build -o "build/macos/IIIF Annotation Studio.app/Contents/MacOS/IIIF Annotation Studio.app"