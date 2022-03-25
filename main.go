// This example only illustrates the assembly of chunks back into a file again.
// In reality, we would send each chunk to Azure directly as blocks without saving the file to the filesystem.
// https://docs.microsoft.com/en-us/rest/api/storageservices/put-block
// https://docs.microsoft.com/en-us/rest/api/storageservices/put-block-list
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func uploadChunkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chunk, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var file *os.File

	if file == nil {
		file, err = os.OpenFile(fmt.Sprintf("%s/%s", "./uploads", fileHeader.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("Error creating file on the filesystem: " + err.Error())
		}
	}
	defer file.Close()

	bytes, err := io.Copy(file, chunk)
	if err != nil {
		panic("Error during chunk write:" + err.Error())
	}

	fmt.Printf("Upload bytes written %d\n", bytes)

	length, _ := strconv.Atoi(r.Header.Get("Upload-Length"))
	offset, _ := strconv.Atoi(r.Header.Get("Upload-Offset"))
	maxSize, _ := strconv.Atoi(r.Header.Get("Upload-Max-Chunk-Size"))

	fmt.Println("Upload Offset", offset)
	fmt.Println("Upload length", length)
	fmt.Println("Upload Chunks", r.Header.Get("Upload-Chunks"))
	fmt.Println("Upload Chunk-Size", maxSize)

	if length-offset <= maxSize {
		fmt.Println("Last Upload Chunk Detected")
	}
	fmt.Fprintf(w, "Upload successful")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadChunkHandler)

	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
