
// for debug
go-bindata-assetfs -debug static/... templates

// for release
go-bindata-assetfs -nomemcopy static/... templates

env GOOS=linux GOARCH=arm GOARM=6 go build