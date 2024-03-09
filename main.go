package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var mimeIcons = map[string]string{
	"text/plain":                    "txt.png",
	"text/html":                     "html.png",
	"application/pdf":               "pdf.png",
	"application/msword":            "doc.png",
	"application/vnd.ms-excel":      "xls.png",
	"application/vnd.ms-powerpoint": "ppt.png",
	"image/jpeg":                    "jpg.png",
	"image/png":                     "png.png",
	"image/gif":                     "gif.png",
	"video/mp4":                     "mp4.png",
	"audio/mpeg":                    "mp3.png",
	"application/zip":               "zip.png",
	"application/x-7z-compressed":   "7z.png",
	"application/x-rar-compressed":  "rar.png",
	"application/x-tar":             "tar.png",
	"application/gzip":              "gz.png",
	"application/x-bzip2":           "bz2.png",
}

func main() {
	http.HandleFunc("/", listFiles)
	http.HandleFunc("/dir/", handleDirectory)
	http.HandleFunc("/download/", downloadFile)
	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "<html><head>")
	fmt.Fprintf(w, "<title>wExplorer - Files List</title>")
	fmt.Fprintf(w, `
		<style>
			body {
				font-family: Arial, sans-serif;
			}
		</style>
	`)
	fmt.Fprintf(w, "</head><body>")

	printNavButtons(w, "", r.URL.Path)
	fmt.Fprintf(w, "<h1>Files</h1><ul>")
	err := listDir(".", w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

func listDir(dir string, w http.ResponseWriter) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		icon := getFileIcon(path, file.IsDir())
		if file.IsDir() {
			fmt.Fprintf(w, "<li>%s <a href='/dir/%s'>%s</a></li>", icon, path, file.Name())
		} else {
			fmt.Fprintf(w, "<li>%s <a href='/download/%s'>%s</a></li>", icon, path, file.Name())
		}
	}
	return nil
}

func handleDirectory(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Path[len("/dir/"):]
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><head>")
	fmt.Fprintf(w, "<title>wExplorer - %s</title>", dirPath)
	fmt.Fprintf(w, `
		<style>
			body {
				font-family: Arial, sans-serif;
			}
		</style>
	`)
	fmt.Fprintf(w, "</head><body>")
	printNavButtons(w, dirPath, r.URL.Path)
	fmt.Fprintf(w, "<h1>%s</h1><ul>", dirPath)
	err := listDir(dirPath, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/download/"):]
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func printNavButtons(w http.ResponseWriter, currentDir, currentPath string) {
	fmt.Fprintf(w, "<div>")
	fmt.Fprintf(w, "<a href='/'>[DIR] Root</a> | ")

	if currentDir != "" {
		parentDir := filepath.Dir(currentDir)
		fmt.Fprintf(w, "<a href='/dir/%s'>[DIR] Parent</a> | ", parentDir)
	}

	if currentPath != "" {
		prevPath := filepath.Dir(currentPath)
		if prevPath != "/" {
			fmt.Fprintf(w, "<a href='%s'>[PREV] Previous</a> | ", prevPath)
		}
	}

	fmt.Fprintf(w, "</div>")
}

func getFileIcon(path string, isDir bool) string {
	if isDir {
		return "<img src='/download/assets/dir.png' alt='Directory' width='16' height='16'>"
	}

	ext := filepath.Ext(path)
	mimeType := getMimeType(ext)
	iconFile, ok := mimeIcons[mimeType]
	if ok {
		return fmt.Sprintf("<img src='/download/assets/%s' alt='%s' width='16' height='16'>", iconFile, mimeType)
	}

	return "<img src='/download/assets/file.png' alt='File' width='16' height='16'>"
}

func getMimeType(ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	switch ext {
	case "txt":
		return "text/plain"
	case "html", "htm":
		return "text/html"
	case "pdf":
		return "application/pdf"
	case "doc", "docx":
		return "application/msword"
	case "xls", "xlsx":
		return "application/vnd.ms-excel"
	case "ppt", "pptx":
		return "application/vnd.ms-powerpoint"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "mp4":
		return "video/mp4"
	case "mp3":
		return "audio/mpeg"
	case "zip":
		return "application/zip"
	case "7z":
		return "application/x-7z-compressed"
	case "rar":
		return "application/x-rar-compressed"
	case "tar":
		return "application/x-tar"
	case "gz":
		return "application/gzip"
	case "bz2":
		return "application/x-bzip2"
	default:
		return "application/octet-stream"
	}
}
