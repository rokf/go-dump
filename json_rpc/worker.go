package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/url"
)

// Calculator is the struct which contains the
// methods available via RPC calls
type Calculator struct{}

// Args represents the two arguments
// for the Add function
type Args struct {
	Int1, Int2 int
}

// Reply contains the result of the
// RPC computation
type Reply struct {
	Result int
}

var quitChan chan bool

// Add is the only
// method available
// to the RPC client
func (c *Calculator) Add(args *Args, reply *Reply) error {
	reply.Result = args.Int1 + args.Int2
	return nil
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")

	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil // success
}

func runServer(port int) {
	calc := new(Calculator)
	server := rpc.NewServer()
	server.Register(calc)

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatal("listen error:", err)
		quitChan <- true
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("accept error:", err)
			quitChan <- true
		}

		log.Println("accepted connection")

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	port, err := getFreePort()

	fmt.Println(port)

	if err != nil {
		log.Fatal("port error:", err)
	}
	go runServer(port)
	// send port to master
	resp, err := http.PostForm("http://127.0.0.1:8888/register", url.Values{"port": {fmt.Sprintf("%d", port)}})

	body, err := ioutil.ReadAll(resp.Body)
	log.Print(string(body))

	if err != nil {
		log.Fatal("registration error:", err)
	}

	resp.Body.Close()

	<-quitChan
}
