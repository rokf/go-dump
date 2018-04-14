package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

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

var channel chan *rpc.Client

func main() {
	channel = make(chan *rpc.Client, 10)

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		a := r.FormValue("a")
		b := r.FormValue("b")

		if len(a) == 0 || len(b) == 0 {
			w.Write([]byte("error: missing one or more params (a, b)"))
			return
		}

		aInt, _ := strconv.Atoi(a)
		bInt, _ := strconv.Atoi(b)

		var c *rpc.Client

		var reply Reply
		var args *Args

		args = &Args{aInt, bInt}

	loop:
		for {
			c = <-channel
			err := c.Call("Calculator.Add", args, &reply)
			if err != nil {
				log.Printf("error: calling remote functon didn't work (%s)", err.Error())
			} else {
				w.Write([]byte(fmt.Sprintf("result: %d\n", reply.Result)))
				channel <- c
				break loop
			}
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Print("error in registration handler form parsing:", err)
			return
		}
		p := r.Form.Get("port")
		if len(p) == 0 {
			log.Print("missing port parameter in request")
			return
		}
		w.Write([]byte(fmt.Sprintf("registered port %s", p)))

		// now start RPC client
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", p))
		if err != nil {
			log.Print("tcp dial error", err)
		}

		c := jsonrpc.NewClient(conn)

		channel <- c
	})
	log.Fatal(http.ListenAndServe(":8888", nil))
}
