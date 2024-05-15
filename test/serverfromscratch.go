package main 

import(
		"fmt"
		"net"
		"strings"
		"os"
)

func handleRequest(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Fprintf(conn, "HTTP/1.1 500 Internal Server Error")
	}
	httpReq := strings.Split(string(buffer), "\r\n\r\n")
	line := strings.Split(httpReq[0], " ")
	path := strings.ReplaceAll(line[1], "", " ")

	if path == "/" {
		fmt.Fprintf(conn, "HTTP/1.1 200 OK \r\n\r\n")
	}
	fmt.Fprintf(conn, "HTTP/1.1 404 Not Found \r\n\r\n")

}

func main() {
	fmt.Println("log server")
	l, err := net.Listen("tcp","0.0.0.0:2441")
	if err != nil {
		fmt.Println("Error al bindear puerto", err.Error())
		os.Exit(1)
	}

	conn, err = net.Accept()
	if err != nil {
		fmt.Println("Error al aceptar listner p:2241", err.Error())
		os.Exit(1)
	}

	handleRequest(conn)
}