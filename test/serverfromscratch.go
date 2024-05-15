// NO HELP NO NOTHING STRAIGHT FROM THE CAP
package main

import(
		"fmt"
		"net"
		"strings"
		"os"
)

func HandleRequest(Conn conn.Accept){
	defer conn.Close()
	buffer := make([]byte, 1024)
	// ERROR en poner = en vez de  :=, además de no agregar como parametro al Read() el argumento buffer que declaramosen la linea anterior
	// _, err = conn.Read()
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error \r\n\r\n")
	}
	// esto no se nada!!!
	httpRequest := strings.Split(string(buffer), "\r\n")
	startLine := strings.Split(httpRequest[0], " ")
	path := strings.ReplaceAll(startLine[1], " ", "")

	if path == "/" {
		fmt.Fprintf("HTTP/1.1 200 OK \r\n\r\n")
		// lo olvide el return
		return
	}
	fmt.Fprintf("HTTP/1.1 404 Not Found \r\n\r\n")
}


func main() {
	fmt.Println("Iniciando el servidor HTTP")
	l, err := net.Listen("tcp", "0.0.0.0:2441")
	if err != nil {
		fmt.Println("Error al bindear puerto 2441", err.Error())
		os.Exit(1)
	}
	conn, err = l.Accept()
	if err != nil {
		fmt.Println("Error al aceptar conec en listener port 2441", err.Error())
		os.Exit(1)
	}
	// fmt.Fprintf(conn, "HTTP/1.1 200 OK \r\n\r\n")
	HandleRequest(conn)
}
