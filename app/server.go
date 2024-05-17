package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Log for the server")
	directory := flag.String("directory", ".", "Directory to serve files from")
	flag.Parse()

	// Stage 1
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221 to enable listening", err.Error())
		os.Exit(1)
	}

	// Stage 6 Hacerlo concurrente para multiples conexiones al mismo tiempo
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleRequest(conn, directory)
	}
	// stage 2 estaba aca, pero se saca el http 200 ok response
}

func responseEcho(conn net.Conn, path string) {
	msg := strings.Split(path, "/")[2]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg
	conn.Write([]byte(resp))
}

func responseUserAgent(conn net.Conn, content string) {
	lines := strings.Split(content, "\r\n")
	userAgent := strings.Split(lines[2], ": ")[1]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
	conn.Write([]byte(resp))
}

func HandleRequest(conn net.Conn, directory *string) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	contentLength, err := conn.Read(buffer)
	if err != nil {
		fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
	}

	content := string(buffer[:contentLength])
	httpRequest := strings.Split(string(buffer), "\r\n")
	startLine := strings.Split(httpRequest[0], " ")
	path := strings.ReplaceAll(startLine[1], " ", "")
	target := startLine[1]
	req_parts := strings.Split(target, "/")
	// Stage 3 si el path es / devuelve 200 OK
	if path == "/" {
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
		return
	} else if strings.HasPrefix(path, "/echo/") {
		// Stage 4 si el path empieza con /echo/ entonces devuelve el string
		responseEcho(conn, path)
	} else if path == "/user-agent" {
		// Stage 5 si es user agent devuelve el header
		responseUserAgent(conn, content)
	} else if len(req_parts) > 2 && req_parts[1] == "files" && checkFileExists(*directory+req_parts[2]) {
		file_contents := getFileContents(*directory + req_parts[2])
		response := ([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprint(len(file_contents)) + "\r\n\r\n" + file_contents))
		conn.Write(response)
	}

	fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	conn.Close()
}

func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getFileContents(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file: ", err.Error())
		os.Exit(1)
	}
	return string(data)
}
