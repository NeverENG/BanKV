package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"zinx/src/ziface"
)

func ClientTest1() {
	fmt.Println("[Client] client start")

	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("[Client] client dial error:", err)
	}

	dp := NewDataPack()
	fmt.Println("[Client] client dial ok")
	for {
		msg := NewMassage(1, []byte("I am MsgId 1"))
		data, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("[Client1] client pack error:", err)
			continue
		}
		_, err = conn.Write(data)

		if err != nil {
			fmt.Println("[Client] client write error:", err)
		}
		Handerdata := make([]byte, 8)
		io.ReadFull(conn, Handerdata)

		msgHead, err := dp.UnPack(Handerdata)
		msg = msgHead.(*Message)
		msg.Data = make([]byte, msg.GetMsgLen())
		_, err = io.ReadFull(conn, msg.Data)
		if err != nil {
			fmt.Println("[Client] client write error:", err)
		}
		fmt.Println("[Client] client write ok", string(msg.GetData()))

		if err != nil {
			fmt.Println("[Client] client write error:", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func ClientTest2() {
	fmt.Println("[Client] client start")

	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("[Client] client dial error:", err)
	}

	dp := NewDataPack()
	fmt.Println("[Client2] client dial ok")
	for {
		msg := NewMassage(2, []byte("I am MsgId 2"))
		data, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("[Client] client pack error:", err)
			continue
		}
		_, err = conn.Write(data)

		if err != nil {
			fmt.Println("[Client] client write error:", err)
		}
		Handerdata := make([]byte, 8)
		io.ReadFull(conn, Handerdata)

		msgHead, err := dp.UnPack(Handerdata)
		msg = msgHead.(*Message)
		msg.Data = make([]byte, msg.GetMsgLen())
		_, err = io.ReadFull(conn, msg.Data)
		if err != nil {
			fmt.Println("[Client] client write error:", err)
		}
		fmt.Println("[Client] client write ok", string(msg.GetData()))

		if err != nil {
			fmt.Println("[Client] client write error:", err)
		}

		time.Sleep(1 * time.Second)
	}
}

type Router1 struct {
	BaseRouter
}

func (r *Router1) Handle(request ziface.IRequest) {

	err := request.GetConnection().SendMsg(1, []byte(request.GetMsgData()))
	if err != nil {
		fmt.Println("[Router Handle] client write error:", err)
	}
	time.Sleep(1 * time.Second)
}

type Router2 struct {
	BaseRouter
}

func (r *Router2) Handle(request ziface.IRequest) {

	err := request.GetConnection().SendMsg(2, []byte(request.GetMsgData()))
	if err != nil {
		fmt.Println("[Router Handle] client write error:", err)
	}
	time.Sleep(1 * time.Second)
	request.GetConnection().Stop()
}

func Hock1(conn ziface.IConnect) {

	data := []byte("Hock1启动,启动")
	err := conn.SendMsg(1, data)
	if err != nil {
		fmt.Println("[Router Handle] client write error:", err)
		return
	}
}

func Hock2(conn ziface.IConnect) {
	data := []byte("Hock2启动,完全关闭")
	err := conn.SendMsg(1, data)
	if err != nil {
		fmt.Println("[Router Handle] client write error:", err)
		return
	}
	time.Sleep(1 * time.Second)
}

func TestServer(t *testing.T) {
	fmt.Println("[Server] server start")
	s := NewServer()
	s.SetConnStopFunc(Hock2)
	s.SetConnStartFunc(Hock1)
	go ClientTest1()
	go ClientTest2()
	s.AddRouter(1, &Router1{BaseRouter{}})

	s.AddRouter(2, &Router2{BaseRouter{}})
	s.Serve()
}
