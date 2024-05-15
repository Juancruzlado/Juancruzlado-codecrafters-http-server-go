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
        } else if strings.HasPrefix(path, "/echo/") {
		responseEcho(conn, path)
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

        conn, err := l.Accept()
        if err != nil {
                fmt.Println("Failed to accept incoming client connection", err.Error())
                os.Exit(1)
        }
	HandleRequest(conn)
        // agregarle logica de si mandan a ruta echo un string asi /echo/unstring que devuelva un body con el string enviado
}
