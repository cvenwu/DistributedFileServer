package main

import (
	"fileSystem/handler"
	"log"
	"net/http"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/14 7:14 下午
 * @Desc:
 */

func main() {
	//建立路由规则
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Failed to start server，err:%s\n", err)
	}
}
