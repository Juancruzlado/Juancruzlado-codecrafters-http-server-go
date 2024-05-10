package main

import(
        "fmt"
        "net"
        "os"
)

func main(){
        fmt.Println("Log for the server")
        l, err := net.Listen("tcp", "0.0.0.0:4221")
        if err != nil {
                fmt.Println("Failed to bind to port 4221")
                os.Exit(1)
        }

        conn, err := l.Accept()
        defer conn.Close()
        if err != nil {
                fmt.Println("Failed to connect")
                os.Exit(1)
        }
        conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
        fmt.Println("exiting")
}
~     