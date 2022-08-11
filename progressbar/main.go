package main

import (
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
)

func main() {
	golangPkg := "https://golang.google.cn/dl/go1.16.4.darwin-amd64.pkg"
	req, _ := http.NewRequest("GET", golangPkg, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile("go1.16.4.darwin-amd64.pkg", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"正在下载",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
