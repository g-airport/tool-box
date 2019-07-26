package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const maxUploadSize = 10 * 1024 * 1024 // 10 MB
const uploadPath = "/Users/tqll/work/go/src/github.com/g-airport/tool-box/upload/tmp"

func main() {
	log.Print("current path v1:", getCurrentPath())
	log.Print("current path v2:", getCurrentPathV2())
	log.Print("current path v3:", getCurrentPathV3())

	http.HandleFunc("/upload", uploadFileHandler())

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Print("Server started on localhost:8080, use /upload for uploading files and /files/{fileName} for downloading")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "file too big", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		fileType := r.PostFormValue("type")
		file, _, err := r.FormFile("upload_file")
		if err != nil {
			renderError(w, "invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "invalid file", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		//fileTyp := http.DetectContentType(fileBytes)
		//switch fileTyp {
		//case "image/jpeg", "image/jpg":
		//case "image/gif", "image/png":
		//case "application/pdf":
		//	break
		//default:
		//	renderError(w, "invalid file type", http.StatusBadRequest)
		//	return
		//}

		//fileName := randToken(12)
		//fileEndings, err := mime.ExtensionsByType(fileType)
		//if err != nil {
		//	renderError(w, "can't read file type", http.StatusInternalServerError)
		//	return
		//}
		//newPath := filepath.Join(uploadPath, fileName+fileEndings[0])

		// set save file name eg example.csv
		newPath := filepath.Join(uploadPath, "example.csv")
		log.Printf("file_type: %s, file: %s\n", fileType, newPath)
		_, err = os.Stat(uploadPath)
		if err != nil {
			_ = os.MkdirAll(uploadPath, 0755)
		}
		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "can't write file", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "can't write file", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte("success"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal("get current path failed", err)
	}
	s = strings.Replace(s, "\\", "/", -1)
	s = strings.Replace(s, "\\\\", "/", -1)
	i := strings.LastIndex(s, "/")
	path := string(s[0 : i+1])
	return path
}

func getCurrentPathV2() string {
	//configfile := filepath.Join(filepath.Dir(execPath), "./config.yml")
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("get current path failed", err)
	}
	return execPath
}

func getCurrentPathV3() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("get current path failed", err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
