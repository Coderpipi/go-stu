package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"net/http"
)

type h struct {
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(len(r.URL.Query()))
	buf := bytes.NewBuffer(nil)
	content := []byte("[{\"id\":\"footer\",\"title\":\"页尾菜单\",\"category\":2,\"children\":[{\"id\":4054469,\"title\":\"子商品\",\"origin_type\":-1,\"type\":\"web\",\"source_id\":\"0\",\"url\":\"https://zyxm.preview.shoplazza.com/collections/%E4%B9%A6%E7%B1%8D/products/%E9%9D%92%E6%98%A5%E7%8C%AA%E5%A4%B4%E5%B0%91%E5%A5%B3%E4%B8%8D%E4%BC%9A%E6%A2%A6%E5%88%B0%E5%85%94%E5%A5%B3%E9%83%8E%E5%AD%A6%E5%A7%90?spm=..collection.collection_1.3\\u0026spm_prev=..p\",\"pid\":0,\"status\":0,\"url_obj\":{\"type\":\"web\",\"source_id\":\"0\",\"url\":\"https://zyxm.preview.shoplazza.com/collections/%E4%B9%A6%E7%B1%8D/products/%E9%9D%92%E6%98%A5%E7%8C%AA%E5%A4%B4%E5%B0%91%E5%A5%B3%E4%B8%8D%E4%BC%9A%E6%A2%A6%E5%88%B0%E5%85%94%E5%A5%B3%E9%83%8E%E5%AD%A6%E5%A7%90?spm=..collection.collection_1.3\\u0026spm_prev=..p\",\"title\":\"\"},\"expanded\":true,\"children\":[{\"id\":4054470,\"title\":\"子商品\",\"origin_type\":-1,\"type\":\"web\",\"source_id\":\"0\",\"url\":\"https://zyxm.preview.shoplazza.com/collections/%E4%B9%A6%E7%B1%8D/products/%E9%9D%92%E6%98%A5%E7%8C%AA%E5%A4%B4%E5%B0%91%E5%A5%B3%E4%B8%8D%E4%BC%9A%E6%A2%A6%E5%88%B0%E5%85%94%E5%A5%B3%E9%83%8E%E5%AD%A6%E5%A7%90?spm=..collection.collection_1.3\\u0026spm_prev=..p\",\"pid\":4054469,\"status\":0,\"url_obj\":{\"type\":\"web\",\"source_id\":\"0\",\"url\":\"https://zyxm.preview.shoplazza.com/collections/%E4%B9%A6%E7%B1%8D/products/%E9%9D%92%E6%98%A5%E7%8C%AA%E5%A4%B4%E5%B0%91%E5%A5%B3%E4%B8%8D%E4%BC%9A%E6%A2%A6%E5%88%B0%E5%85%94%E5%A5%B3%E9%83%8E%E5%AD%A6%E5%A7%90?spm=..collection.collection_1.3\\u0026spm_prev=..p\",\"title\":\"\"},\"expanded\":true}]}]}]")
	dict, _ := flate.NewWriter(buf, flate.BestCompression)
	dict.Write(content)
	dict.Flush()
	w.Header().Set("Content-Encoding", "deflate")
	res, _ := ioutil.ReadAll(buf)
	w.Write(res)
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":80", nil)
}
