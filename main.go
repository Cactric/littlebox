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
    "flag"
)

var globalUploadDir string
var globalResourcesDir string

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
        t, err := template.ParseFiles(globalResourcesDir + filename)
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
        content, err := ioutil.ReadFile(globalResourcesDir + filename)
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
    outputFolder := globalUploadDir
    filename := strconv.Itoa(int(time.Now().UnixNano())) + "_" + receivingFileHeader.Filename
    file, err := os.Create(outputFolder + "/" + filename)
    if err != nil {
        log.Print(err)
        w.WriteHeader(500)
        t, t_err := template.ParseFiles(globalResourcesDir + "errorpages/500.html")
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
    var port int
    resourcesDir := "./"
    // globalUploadDir is used
    // globalResourcesDir is used
    flag.IntVar(&port, "p", 8000, "Specify port to listen on, default 8000")
    flag.StringVar(&globalUploadDir, "d", "./uploads", "Specify folder to put uploaded files, default is ./uploads")
    flag.StringVar(&resourcesDir, "r", "./", "Specify directory with Littlebox’s resources (HTML files, etc.) Default is the current directory")
    
    flag.Parse()
    
    // Going to put resourcesDir in globalResourcesDir
    // If the last character of the resourcesDir is not a /, add one to it before putting it in globalResourcesDir
    if !(resourcesDir[len(resourcesDir) - 1:] == "/") {
        globalResourcesDir = resourcesDir + "/"
    } else {
        globalResourcesDir = resourcesDir
    }
    
    http.HandleFunc("/", staticFile)
    http.HandleFunc("/upload", recvFile)
    portString := strconv.Itoa(port)
    fmt.Println("About to serve Littlebox on port " + portString)
    fmt.Println("Uploads will go into the following directory: " + globalUploadDir)
    fmt.Println("Resources are being loaded from " + globalResourcesDir)
    log.Fatal(http.ListenAndServe(":" + portString, nil))
}
