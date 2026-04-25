package client

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/NeverENG/BanKV/internal/network/ziface"
	"github.com/NeverENG/BanKV/internal/network/znet"
)

// Client 客户端结构
type Client struct {
	conn    net.Conn
	addr    string
	packer  ziface.IDataPack
}

// NewClient 创建客户端实例
func NewClient(addr string) *Client {
	return &Client{
		addr:   addr,
		packer: znet.NewDataPack(),
	}
}

// Connect 连接到服务端
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("connect to %s failed: %v", c.addr, err)
	}
	c.conn = conn
	fmt.Printf("Connected to %s\n", c.addr)
	return nil
}

// SendRequest 发送请求
func (c *Client) SendRequest(msgID uint32, data []byte) ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	// 创建消息
	msg := znet.NewMessage(msgID, data)

	// 打包消息
	pkgData, err := c.packer.Pack(msg)
	if err != nil {
		return nil, fmt.Errorf("pack message failed: %v", err)
	}

	// 发送消息
	_, err = c.conn.Write(pkgData)
	if err != nil {
		return nil, fmt.Errorf("send message failed: %v", err)
	}

	// 接收响应
	headData := make([]byte, c.packer.GetHeadLen())
	_, err = c.conn.Read(headData)
	if err != nil {
		return nil, fmt.Errorf("read head failed: %v", err)
	}

	// 解包头部
	msgHead, err := c.packer.UnPack(headData)
	if err != nil {
		return nil, fmt.Errorf("unpack head failed: %v", err)
	}

	// 读取数据部分
	dataLen := msgHead.GetMsgLen()
	if dataLen > 0 {
		respData := make([]byte, dataLen)
		_, err = c.conn.Read(respData)
		if err != nil {
			return nil, fmt.Errorf("read data failed: %v", err)
		}
		return respData, nil
	}

	return nil, nil
}

// SendPut 发送 PUT 请求
func (c *Client) SendPut(key, value []byte) error {
	// 构建数据：key_len + key + value_len + value
	data := make([]byte, 8+len(key)+len(value))
	
	// 写入 key 长度
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(key)))
	// 写入 key
	copy(data[4:4+len(key)], key)
	// 写入 value 长度
	binary.LittleEndian.PutUint32(data[4+len(key):8+len(key)], uint32(len(value)))
	// 写入 value
	copy(data[8+len(key):], value)

	// 发送请求
	resp, err := c.SendRequest(1, data)
	if err != nil {
		return err
	}

	// 检查响应
	if len(resp) > 0 && resp[0] != 0 {
		return fmt.Errorf("server error")
	}

	return nil
}

// SendGet 发送 GET 请求
func (c *Client) SendGet(key []byte) ([]byte, error) {
	// 构建数据：key_len + key
	data := make([]byte, 4+len(key))
	
	// 写入 key 长度
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(key)))
	// 写入 key
	copy(data[4:], key)

	// 发送请求
	resp, err := c.SendRequest(2, data)
	if err != nil {
		return nil, err
	}

	// 检查响应
	if len(resp) == 0 {
		return nil, fmt.Errorf("no response")
	}

	if resp[0] != 0 {
		return nil, fmt.Errorf("server error")
	}

	// 解析响应：状态 + value_len + value
	if len(resp) < 5 {
		return nil, fmt.Errorf("invalid response")
	}

	valueLen := binary.LittleEndian.Uint32(resp[1:5])
	if len(resp) < 5+int(valueLen) {
		return nil, fmt.Errorf("invalid response length")
	}

	return resp[5 : 5+int(valueLen)], nil
}

// SendDelete 发送 DELETE 请求
func (c *Client) SendDelete(key []byte) error {
	// 构建数据：key_len + key
	data := make([]byte, 4+len(key))
	
	// 写入 key 长度
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(key)))
	// 写入 key
	copy(data[4:], key)

	// 发送请求
	resp, err := c.SendRequest(3, data)
	if err != nil {
		return err
	}

	// 检查响应
	if len(resp) > 0 && resp[0] != 0 {
		return fmt.Errorf("server error")
	}

	return nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
