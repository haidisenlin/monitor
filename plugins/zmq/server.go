package zmq

import (
	"fmt"
	"github.com/zeromq/goczmq"
)

func main() {

	fmt.Println("开始监听")
	// create a router channeler
	router := goczmq.NewRouterChanneler("tcp://192.168.37.20:5555")
	defer router.Destroy()

	for {
		// receive the hello message
		request := <-router.RecvChan
		fmt.Println("我写的" + string(request[0]) + ":" + string(request[1]) + ":" + string(request[2]))

		// first frame is identity of client - let's append 'World'
		// to the message and route it back
		request = append(request, []byte("World"))

		// send the reply
		router.SendChan <- request

		fmt.Println("发送完成")
	}
}
