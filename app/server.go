package main

import(
        "fmt"
        "net"
        "strings"
        "os"
)
 
func responseEcho(conn net.Conn, path string) {
        msg := strings.Split(path, "/")[2]
        resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg
        conn.Write([]byte(resp))
}
 
func responseUserAgent(conn net.Conn, content string) {
        lines := strings.Split(content, "\r\n")
        fmt.Println("lines here", len(lines), lines)
        userAgent := strings.Split(lines[2], ": ")[1]
        resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
        conn.Write([]byte(resp))
}
 
func HandleRequest(conn net.Conn) {
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
 
fmt.Printf("path: `%s`\n", path)
        if path == "/" {
                fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
                return
        } else if strings.HasPrefix(path, "/echo/") {
                responseEcho(conn, path)
        } else if path == "/user-agent" {
                responseUserAgent(conn, content)
        }
        fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
}
 
func main(){
        fmt.Println("Log for the server")
        l, err := net.Listen("tcp", "0.0.0.0:4221")
        if err != nil {
                fmt.Println("Failed to bind to port 4221 to enable listening", err.Error())
                os.Exit(1)
        }
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleRequest(conn)
        }
        HandleRequest(conn)
}
        