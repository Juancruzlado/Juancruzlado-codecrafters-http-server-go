package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Log for the server")
	directory := flag.String("directory", ".", "Directory to serve files from")
	flag.Parse()

	// Stage 1 Bind del puerto 4221 que va a ocupar el servidor HTTP
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221 to enable listening", err.Error())
		os.Exit(1)
	}

	// Stage 6 Hacerlo concurrente para multiples conexiones al mismo tiempo con goroutines
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleRequest(conn, directory)
	}
}

// Metodo que responde con el HTTP Response y el string que quiere imprimir el cliente
func responseEcho(conn net.Conn, path string, acceptEncoding string) {
	msg := strings.Split(path, "/")[2]
	if acceptEncoding == "gzip" {
		resp := "HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg
		conn.Write([]byte(resp))
	} else {
		resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg
		conn.Write([]byte(resp))
	}
}

// Metodo que responde con el HTTP Response y le concatena el Header User-Agent extraido del Request del Cliente 
func responseUserAgent(conn net.Conn, content string) {
	lines := strings.Split(content, "\r\n")
	userAgent := strings.Split(lines[2], ": ")[1]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
	conn.Write([]byte(resp))
}

// Metodo que maneja las Request segun sus casos y determina que va a responder el Server.
func HandleRequest(conn net.Conn, directory *string) {

	buffer := make([]byte, 1024)
	contentLength, err := conn.Read(buffer)
	if err != nil {
		fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
		conn.Close()
		return
	}

	content := string(buffer[:contentLength])
	httpRequest := strings.Split(content, "\r\n")
	startLine := strings.Split(httpRequest[0], " ")
	method := startLine[0]
	path := strings.ReplaceAll(startLine[1], " ", "")
	target := startLine[1]
	reqParts := strings.Split(target, "/")

	acceptEncoding := getHeaderValue(httpRequest, "Accept-Encoding")

	if method == "POST" && len(reqParts) > 2 && reqParts[1] == "files" {
		// Stage: Manejo de solicitudes POST para subir archivos
		filename := reqParts[2]
		body := extractRequestBody(content)
		err := ioutil.WriteFile(*directory+"/"+filename, []byte(body), 0644)
		if err != nil {
			fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
			conn.Close()
			return
		}
		fmt.Fprintf(conn, "HTTP/1.1 201 Created\r\n\r\n")
		conn.Close()
		return
	}

	// Stage 2 y 3 si el path es / devuelve 200 OK
	if path == "/" {
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
		conn.Close()
		return
	} else if strings.HasPrefix(path, "/echo/") {
		// Stage 4 si el path empieza con /echo/ entonces devuelve el string
		responseEcho(conn, path, acceptEncoding)
		conn.Close()
		return
	} else if path == "/user-agent" {
		// Stage 5 Si el path es user agent devuelve el header
		responseUserAgent(conn, content)
		conn.Close()
		return
	} else if len(reqParts) > 2 && reqParts[1] == "files" && checkFileExists(*directory+"/"+reqParts[2]) {
		// Stage 7 Agregando la capacidad de hacerle un GET a un archivo.
		fileContents := getFileContents(*directory + "/" + reqParts[2])
		response := ([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprint(len(fileContents)) + "\r\n\r\n" + fileContents))
		conn.Write(response)
		conn.Close()
		return
	}

	fmt.Fprintf(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	conn.Close()
}

// Metodo que se fija si existe o no el archivo y retorna un booleano de 0 o 1
func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Metodo que lee los contenidos del archivo y los devuelve en un string.
func getFileContents(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file: ", err.Error())
		os.Exit(1)
	}
	return string(data)
}

// Metodo que extrae el cuerpo de la solicitud POST
func extractRequestBody(content string) string {
	parts := strings.Split(content, "\r\n\r\n")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// Metodo que obtiene el valor de una cabecera específica
func getHeaderValue(httpRequest []string, headerName string) string {
	headerName = strings.ToLower(headerName) + ":"
	for _, line := range httpRequest {
		if strings.HasPrefix(strings.ToLower(line), headerName) {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return ""
}
