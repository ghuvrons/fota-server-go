package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	giotgo "github.com/ghuvrons/g-IoT-Go"
)

func main() {
	fmt.Println("Start")

	var server = giotgo.NewServer()

	setCmdHandlers(server)

	server.ClientAuth(func(username, password string) bool {
		fmt.Println(username, password)
		return true
	})
	go server.Serve("0.0.0.0:2000")

	http.HandleFunc("/", routeSubmitPost)
	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := r.FormValue("alias")

	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename

	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}

	fileLocation := filepath.Join(dir, "files", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("done"))
}
