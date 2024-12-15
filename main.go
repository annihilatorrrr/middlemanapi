package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var domain = os.Getenv("URL")

type jSONyResponse struct {
	Title    string `json:"title"`
	Thumb    string `json:"thumb"`
	Dlurl    string `json:"dlurl"`
	Duration int    `json:"duration"`
	Size     int64  `json:"size"`
}

func getResponse(input string) *jSONyResponse {
	jdata := &jSONyResponse{}
	get, err := http.Get(fmt.Sprintf("%sydl?key=%s&q=%s", domain, os.Getenv("key"), input))
	if err != nil {
		return jdata
	}
	defer get.Body.Close()
	byteee, _ := io.ReadAll(get.Body)
	_ = json.Unmarshal(byteee, &jdata)
	return jdata
}

func HandletyDL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}
	ipAddress := r.Header.Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	qurl := r.URL.Query().Get("q")
	if qurl == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}
	log.Printf("Requested by: %s - %s!", ipAddress, qurl)
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(getResponse(qurl)); err != nil {
		http.Error(w, "Error encoding response!\nContact github.com/annihilatorrrr ASAP!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(buf.Bytes())
}

func handleico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./favicon.ico")
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}
	_, _ = fmt.Fprintf(w, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Service Status</title>
    <style>
        body {
            background-color: #121212;
            color: #e0e0e0;
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
        }
        h1 {
            color: #bb86fc;
        }
        pre {
            background-color: #1e1e1e;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <h1>I'm alive!</h1>
    <p><b>Supporting:</b></p>
    <pre>
> Video from Yt.
    </pre>
    <p>Paid API Key can be obtained from <a href="https://t.me/annihilatorrrr" style="color: #bb86fc;">t.me/annihilatorrrr</a>!</p>
    <p>Go Version: %s</p>
    <p>Go Routines: %d</p>
    <p>Usage:</p>
    <pre>
GET /ydl?key=YOURKEY&q=QUERY or Link (YT)
    </pre>
<p>Not Working As Accepted:</p>
<pre>/ydl</pre>
    <p>Source Code: <a href="https://runurl.in/GetSourceCode" style="color: #bb86fc;">https://runurl.in/GetSourceCode</a></p>
</body>
</html>`, runtime.Version(), runtime.NumGoroutine())
}

func main() {
	if domain == "" {
		domain = "https://f76f-20-157-218-109.ngrok-free.app/"
	}
	if !strings.HasSuffix(domain, "/") {
		domain = domain + "/"
	}
	router := http.NewServeMux()
	router.HandleFunc("/", HandleHome)
	router.HandleFunc("/ydl", HandletyDL)
	router.HandleFunc("/favicon.ico", handleico)
	httphandler := http.Handler(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      httphandler,
	}
	go func() {
		time.Sleep(time.Second * 21600)
		self, err := os.Executable()
		if err != nil {
			log.Println(err.Error())
			return
		}
		_ = syscall.Exec(self, os.Args, os.Environ())
	}()
	log.Println("Started!")
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	log.Fatal("Bye!")
}
