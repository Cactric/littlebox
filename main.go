package main

import (
    "fmt"
    "log"
    "os"
    "net/http"
    "strconv"
    "io"
    "io/ioutil"
    "html/template"
    "time"
)

func staticFile(w http.ResponseWriter, r *http.Request) {
    filename := ""
    mimeType := "text/plain"
    fmt.Println("request for " + r.URL.Path)
    switch r.URL.Path {
        case "/upload":
            filename = "uploaded.html"
            mimeType = "text/html"
        case "/icon.svg":
            filename = "icon.svg"
            mimeType = "image/svg+xml"
        case "/littlebox_style.css":
            filename = "littlebox_style.css"
            mimeType = "text/css"
        default:
            filename = "homepage.html"
            mimeType = "text/html"
    }
    if mimeType == "text/html" {
        t, err := template.ParseFiles(filename)
        if err != nil {
            log.Fatal(err)
        }
        w.Header().Set("Content-Type", mimeType)
        w.WriteHeader(http.StatusOK)
        err = t.ExecuteTemplate(w, filename, nil)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        w.Header().Set("Content-Type", mimeType)
        w.WriteHeader(http.StatusOK)
        content, err := ioutil.ReadFile(filename)
        if err != nil {
            log.Fatal(err)
        }
        _, err = w.Write(content)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func recvFile(w http.ResponseWriter, r *http.Request) {
    receivingFile, receivingFileHeader, err := r.FormFile("FileToUpload")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Receiving file's size: " + strconv.Itoa(int(receivingFileHeader.Size)))
    outputFolder := "uploads"
    filename := strconv.Itoa(int(time.Now().UnixNano())) + "_" + receivingFileHeader.Filename
    file, err := os.Create(outputFolder + "/" + filename)
    if err != nil {
        log.Print(err)
        w.WriteHeader(500)
        t, t_err := template.ParseFiles("errorpages/500.html")
        if t_err != nil {
            log.Fatal(t_err)
        }
        // Send an error page to the user
        err_data := struct {
            Error_description string
        } {
            Error_description: err.Error(),
        }
        t.ExecuteTemplate(w, "500.html", err_data)
    } else {
        defer file.Close()
        io.Copy(file, receivingFile)
        
        // Successful - send back the page
        staticFile(w, r)
    }
}

func main() {
    http.HandleFunc("/", staticFile)
    http.HandleFunc("/upload", recvFile)
    portString := strconv.Itoa(8000)
    fmt.Println("About to serve Littlebox on port " + portString)
    log.Fatal(http.ListenAndServe(":" + portString, nil))
}
