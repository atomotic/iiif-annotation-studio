# IIIF Annotation Studio

This is the [Mirador IIIF Viewer](https://projectmirador.org) packaged as a desktop app with an embedded annotation endpoint that saves annotations to a local sqlite database.

![](https://gist.githubusercontent.com/atomotic/f621ab4364a0bd520dd1066bc4d0e439/raw/1ae8d554576bbe0ecae6b5630dbaa262e26d8615/iiif-annotation-studio.png)

This is a proof of concept, slightly usable at the moment.
A lot of things are missing:

- [x] annotation api: create
- [x] annotation api: update and delete
- [ ] UI: a page listing all manifests and relative annotations
- [ ] export (and import) of annotations
- [ ] builds for linux (AppImage) and Windows
- [ ] share annotations (maybe using DAT)

## Build

Clone

    git clone https://github.com/atomotic/iiif-annotation-studio
    cd iiif-annotation-studio

Install go-bindata

    go get github.com/jteeuwen/go-bindata/...
    go get github.com/elazarl/go-bindata-assetfs/...

Install modules (requires [dep](https://github.com/golang/dep))

    dep ensure

Package static assets

    go-bindata-assetfs static/...

— Build for macos

    go build -o "build/macos/IIIF Annotation Studio.app/Contents/MacOS/IIIF Annotation Studio.app"

Run `IIIF Annotation Studio.app` from build/macos  
(if you trust me there is a ready made build in [releases](https://github.com/atomotic/iiif-annotation-studio/releases))

— Build for linux. GTK3 and GtkWebkit2 required, not yet tested.

    go build
    ./iiif-annotation-studio

## Backup

The first run creates the sqlite database `$HOME/.annotations/annotations.db` if not existing. Just backup this file.

    ➜  ~ sqlite3 ~/.annotations/annotations.db .schema
    CREATE TABLE annotations (
    		id INTEGER PRIMARY KEY,
    		annoid VARCHAR,
    		created_at DATETIME,
    		target VARCHAR,
    		manifest VARCHAR,
    		body TEXT);
    CREATE UNIQUE INDEX annotation_id ON annotations (annoid);

## Disclaimer

Is this another Electron app? NO, is made in golang with [zserge/webview](https://github.com/zserge/webview).

## See also

[Mirador Desktop](https://github.com/ProjectMirador/mirador-desktop): Electron based, annotations on LocalStorage

## Acknowledgement

[demetsiiify](https://github.com/jbaiter/demetsiiify) by [Johannes Baiter](https://github.com/jbaiter) contains an extremely simple annotation server that is probably the best piece of code to learn how Mirador interacts with annotations stores.
