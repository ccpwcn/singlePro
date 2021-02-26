package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static
var embedFiles embed.FS

func main() {
	useOs := len(os.Args) > 1 && os.Args[1] == "live"
	http.Handle("/", http.FileServer(getFileSys(useOs)))
	if err := http.ListenAndServe(":9100", nil); err != nil {
		panic(err)
	}
}

func getFileSys(useOs bool) http.FileSystem {
	if useOs {
		log.Println("使用常规文件系统模式，页面静态资源从工作目录的static下加载...")
		return http.FS(os.DirFS("static"))
	} else {
		log.Println("使用内嵌文件模式，页面静态资源从模拟的static下加载...")
		if fSys, err := fs.Sub(embedFiles, "static"); err != nil {
			panic(err)
		} else {
			return http.FS(fSys)
		}
	}
}
