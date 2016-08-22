// go build -o chat // chatは自由に決めていい

package main

import (
	// "fmt"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once // 一度実行するともうよべない
	filename string
	temp1    *template.Template
}

// ServeHTTPメソッドをもつとhttp.Handleに渡せる
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		// templateを一回だけコンパイルする
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, r) // templateからHTMLを出力する際にhttp.Requestに含まれるデータを参照できる
}

func main() {
	// 起動時に ./chan -addr=":3000"などと渡せる
	var addr = flag.String("addr", ":8080", "アプリのアドレス") // key, default, desc
	flag.Parse()

	// routeにたいし、viewを返している
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){w.Write([]byte(`<html>`))})
	http.Handle("/", &templateHandler{filename: "temp1.html"})

	http.Handle("/about", &templateHandler{filename: "about.html"})

	// 		     request					  requestからprefixを取り除く         assetsまでのパス
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	r := newRoom()
	// /roomへのリクエストはrが受け取る。
	// rは無限ループでr.join, r.leave, r.clientsをずっと監視している
	http.Handle("/room", r)
	go r.run()

	// サーバ起動
	log.Println("webサーバーを起動します. port: ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
