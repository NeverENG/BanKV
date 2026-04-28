package banNet

import (
	"fmt"
	"net"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/network/banIface"
)

type Server struct {
	IP        string
	Port      int
	Name      string
	IPVersion string
	MsgHandle banIface.IMsgHandle
	ConnMgr   banIface.IConnManager

	ConnStartFunc func(conn banIface.IConnect)
	ConnStopFunc  func(conn banIface.IConnect)
}

func (s *Server) AddRouter(msgId uint32, router banIface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
}

/*
基础回显函数 v 0.2
func CallBackClient(conn *net.TCPConn, buf []byte, cnt int) error {
	fmt.Println("[CallBack]START")
	if _, err := conn.Write(buf[:cnt]); err != nil {
		fmt.Println("[CallBack]Write err:", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}
*/

func NewServer() banIface.IServer {
	return &Server{
		IPVersion: "tcp4",
		IP:        config.G.Host,
		Name:      config.G.Name,
		Port:      config.G.Port,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
}

func (s *Server) GetConnMgr() banIface.IConnManager {
	return s.ConnMgr
}

func (s *Server) Start() {
	fmt.Printf("[START]BanKVNetWork:%s ip: %s port:%d \n", s.Name, s.IP, s.Port)

	go func() {

		s.MsgHandle.StartWorkerPool()

		TcPAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[ERROR] Get the Tcp Addr err :", err)
		}
		listener, err := net.ListenTCP(s.IPVersion, TcPAddr)
		if err != nil {
			fmt.Println("[ERROR] ListenTcp err :", err)
		}
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[ERROR] Accept err :", err)
				continue
			}

			if s.ConnMgr.Len() >= config.G.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConnection(conn, cid, s.MsgHandle, s)
			fmt.Println("链接启动中")
			go dealConn.Start()
			cid++
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP]server listener at IP : " + s.IP)
	// 处理副作用并安全推出
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) SetConnStartFunc(f func(conn banIface.IConnect)) {
	s.ConnStartFunc = f
}
func (s *Server) SetConnStopFunc(f func(conn banIface.IConnect)) {
	s.ConnStopFunc = f
}
func (s *Server) CallConnStartFunc(conn banIface.IConnect) {
	if s.ConnStartFunc == nil {
		fmt.Println("[ERROR] CallConnStartFunc is nil!")
		return
	}
	s.ConnStartFunc(conn)
}

func (s *Server) CallConnStopFunc(conn banIface.IConnect) {
	if s.ConnStopFunc == nil {
		fmt.Println("[ERROR] CallConnStopFunc is nil!")
		return
	}
	s.ConnStopFunc(conn)
}
