package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

type FileHandler struct {
	Dir string
}

func (fh *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		filePath := fh.Dir + "/" + header.Filename + time.Now().String()
		err = ioutil.WriteFile(filePath, data, 0777)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "File %s has been successfully uploaded", header.Filename)
	default:
		http.Error(w, "Wrong method", http.StatusBadRequest)
	}
}

func main() {
	handler := &FileHandler{Dir: "upload"}
	http.Handle("/upload", handler)

	http.HandleFunc("/list", func(w http.ResponseWriter, req *http.Request) {
		files, err := ioutil.ReadDir(handler.Dir)
		if err != nil {
			fmt.Fprintf(w, "Unable to read file list")
			return
		}

		filter := req.Header.Get("Filter")
		hasFilter := len(filter) > 0

		for _, file := range files {
			if hasFilter && filepath.Ext(file.Name()) == filter {
				fmt.Fprintf(w, "%s - %d\n", file.Name(), file.Size())
			} else if !hasFilter {
				fmt.Fprintf(w, "%s - %d\n", file.Name(), file.Size())
			}

		}
	})

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	dirToServe := http.Dir(handler.Dir)
	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fs.ListenAndServe()
	srv.ListenAndServe()
}
