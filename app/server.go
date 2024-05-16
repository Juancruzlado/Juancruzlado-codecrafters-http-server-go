package main

import(
        "flag"
        "fmt"
        "net"
        "strings"
        "os"
)

type Config struct {
	port      int
	directory string
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
 
func HandleRequest(cfg Config, conn net.Conn) {
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
}  else if strings.HasPrefix(path, "/files/") {
        // Stage 7 GET file 
		fileName, _ := strings.CutPrefix(path, "/files/")
		filePath := fmt.Sprintf("%s/%s", cfg.directory, fileName)
		if _, err := os.Stat(filePath); err == nil {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
                                fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
			} else {
                                resp := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprint(len(fileContent)) + "\r\n\r\n"
                                conn.Write([]byte(resp))
			}
		} else {
			fmt.Printf("File %s does not exist\n", filePath)
                        fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		}
	}
	conn.Write([]byte(response.String()))
	conn.Close()
        
fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
}
 
func main(){
        fmt.Println("Log for the server")
        // Stage 1
        l, err := net.Listen("tcp", "0.0.0.0:4221")
        if err != nil {
                fmt.Println("Failed to bind to port 4221 to enable listening", err.Error())
                os.Exit(1)
        }

        var cfg Config
	flag.IntVar(&cfg.port, "port", 4221, "TCP server port")
	flag.StringVar(&cfg.directory, "directory", "", "Directory to serve files from")
        flag.Parse()

        // Stage 6 Hacerlo concurrente para multiples conexiones al mismo tiempo
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleRequest(cfg, conn)
        }
        // stage 2 estaba aca, pero se saca el http 200 ok response
}
        