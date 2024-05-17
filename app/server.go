package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"io/ioutil" 
)

func main() {
	fmt.Println("Log for the server")
	directory := flag.String("directory", ".", "Directory to serve files from")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Printf("Failed to bind to port 4221: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleRequest(conn, *directory)
	}
}

func handleRequest(conn net.Conn, directory string) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	contentLength, err := conn.Read(buffer)
	if err != nil {
		writeResponse(conn, "500 Internal Server Error", "text/plain", "Internal Server Error")
		return
	}

	request := string(buffer[:contentLength])
	requestLines := strings.Split(request, "\r\n")
	startLine := strings.Fields(requestLines[0])
	if len(startLine) < 2 {
		writeResponse(conn, "400 Bad Request", "text/plain", "Bad Request")
		return
	}

	method, path := startLine[0], startLine[1]
	if method == "POST" && strings.HasPrefix(path, "/files/") { // Added to handle POST /files/<filename>
		handleFileUpload(conn, directory, path, buffer[contentLength:]) // Added to handle file upload
		return
	} else if method != "GET" {
		writeResponse(conn, "405 Method Not Allowed", "text/plain", "Method Not Allowed")
		return
	}

	switch {
	case path == "/":
		writeResponse(conn, "200 OK", "text/plain", "")
	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path)
	case path == "/user-agent":
		handleUserAgent(conn, requestLines)
	case strings.HasPrefix(path, "/files/"):
		handleFileRequest(conn, directory, path)
	default:
		writeResponse(conn, "404 Not Found", "text/plain", "Not Found")
	}
}

func handleEcho(conn net.Conn, path string) {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 3 {
		writeResponse(conn, "400 Bad Request", "text/plain", "Bad Request")
		return
	}
	message := parts[2]
	writeResponse(conn, "200 OK", "text/plain", message)
}

func handleUserAgent(conn net.Conn, requestLines []string) {
	for _, line := range requestLines {
		if strings.HasPrefix(line, "User-Agent:") {
			userAgent := strings.TrimSpace(strings.TrimPrefix(line, "User-Agent:"))
			writeResponse(conn, "200 OK", "text/plain", userAgent)
			return
		}
	}
	writeResponse(conn, "400 Bad Request", "text/plain", "User-Agent header not found")
}

func handleFileRequest(conn net.Conn, directory, path string) {
	filePath := directory + strings.TrimPrefix(path, "/files/")
	if !fileExists(filePath) {
		writeResponse(conn, "404 Not Found", "text/plain", "File Not Found")
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		writeResponse(conn, "500 Internal Server Error", "text/plain", "Internal Server Error")
		return
	}

	writeResponse(conn, "200 OK", "application/octet-stream", string(content))
}

func handleFileUpload(conn net.Conn, directory, path string, body []byte) {
	filename := strings.TrimPrefix(path, "/files/")
	filePath := directory + "/" + filename

	err := ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		fmt.Printf("Error writing file %s: %v\n", filePath, err)
		writeResponse(conn, "500 Internal Server Error", "text/plain", "Internal Server Error")
		return
	}

	writeResponse(conn, "201 Created", "text/plain", "File created")
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func writeResponse(conn net.Conn, status, contentType, body string) {
	response := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", status, contentType, len(body), body)
	conn.Write([]byte(response))
}
