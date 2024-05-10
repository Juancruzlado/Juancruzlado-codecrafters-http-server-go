package main

import(
        "fmt"
        "net"
        "strings"
        "os"
)

func HandleRequest(conn net.Conn) {
defer conn.Close()

buffer := make([]byte, 1024)
_, err := conn.Read(buffer)
if err != nil {
        fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
}
httpRequest := strings.Split(string(buffer), "\r\n")
startLine := strings.Split(httpRequest[0], " ")
path := strings.ReplaceAll(startLine[1], " ", "")

fmt.Printf("path: `%s`\n", path)
if path == "/" {
        fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
        return
}

fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
}

func main(){
        fmt.Println("Log for the server")
        l, err := net.Listen("tcp", "0.0.0.0:4221")
        if err != nil {
                fmt.Println("Failed to bind to port 4221 to enable listening")
                os.Exit(1)
        }

        defer l.Close()
        conn, err := l.Accept()
        if err != nil {
			fmt.Println("Failed to accept incoming client connection")
			os.Exit(1)
        }
        defer conn.Close()
        conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
        conn.Close()
	HandleRequest(conn)

}
