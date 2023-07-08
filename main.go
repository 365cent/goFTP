package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"

	"github.com/jlaffaye/ftp"
)

type Page struct {
	Title string
	Files []string
}

var (
	ftpAddress  = "ftp.example.com:21"
	ftpUsername = "ftpUsername"
	ftpPassword = "ftpPassword"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/remote-download", remoteDownloadHandler)
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/" // Default to root if no path is provided
	}

	c, err := ftp.Dial(ftpAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Login(ftpUsername, ftpPassword)
	if err != nil {
		fmt.Println(err)
		return
	}

	entries, err := c.List(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var files []string
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFolder {
			continue
		}
		files = append(files, entry.Name)
	}

	p := &Page{Title: "FTP File List", Files: files}
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, p)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is missing.", http.StatusBadRequest)
		return
	}

	c, err := ftp.Dial(ftpAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.Login(ftpUsername, ftpPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader, err := c.Retr(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(fileName)))

	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to write file to response.", http.StatusInternalServerError)
		return
	}
}

func remoteDownloadHandler(w http.ResponseWriter, r *http.Request) {
	sourceURL := r.FormValue("sourceURL")

	if sourceURL == "" {
		http.Error(w, "Missing parameters.", http.StatusBadRequest)
		return
	}

	// Parse source URL and get the filename
	u, err := url.Parse(sourceURL)
	if err != nil {
		http.Error(w, "Invalid URL.", http.StatusBadRequest)
		return
	}

	filename := path.Base(u.Path)
	targetPath := "/" + filename // set to root directory

	// Get file from source URL
	resp, err := http.Get(sourceURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Connect to FTP server
	conn, err := ftp.Dial(ftpAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = conn.Login(ftpUsername, ftpPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload file to FTP server
	err = conn.Stor(targetPath, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File downloaded from %s and uploaded to %s successfully.", sourceURL, targetPath)
}
