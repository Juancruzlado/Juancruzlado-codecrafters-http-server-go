[![progress-banner](https://backend.codecrafters.io/progress/http-server/f4c7b30e-a5bb-44bf-bccb-0e6eab4cb094)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# T.R.A.S.H -> HTTP Server in Golang

- Trash stands for Traditional Raw Amateur Server for HTTP
Why Traditional, because it is based on [RFC 2616 HTTP/1.1](https://datatracker.ietf.org/doc/html/rfc2616)
- Raw because it's not looking to be perfect and is full of things to work on and improve
- Amateur because it's my first HTTP Server ever :) 
- Server for 
- HTTP is pretty self-explanatory. 

# Why?
It's built as a hobby project also to try out codecrafters platform.  
# How?
It's built entirely in the Go Language. You can run it with:
```bash
go run trash.go
``` 

- It listens in Port 4221 for TCP connection with the client
- It can function as an 200 ok page if you curl the main page like this:
```bash
curl -i localhost:4221/
```
- It returns you your own string that you send via url if you curl this route:
```bash
 localhost:4221/echo/hello 
```
  in this example i pick the word hello, use whichever you want.
- It returns your User-Agent if you go to 
```bash
localhost:4221/user-agent
```
- It can handle concurrent connections with the use of goroutines
- You can GET a file from directory with
```bash
curl -i localhost:4221/ --directory
```
- You can POST a file [- To DO -]