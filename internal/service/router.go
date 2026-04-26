package service

import (
	"encoding/binary"

	"github.com/NeverENG/BanKV/internal/network/ziface"
)

// Router 基础路由处理器
type Router struct {
	kv *KVServer

	// 前置处理函数
	preHandleFunc func(request ziface.IRequest)
	// 后置处理函数
	postHandleFunc func(request ziface.IRequest)
}

// NewRouter 创建新的路由处理器
func NewRouter(kv *KVServer) *Router {
	return &Router{
		kv: kv,
	}
}

// SetPreHandle 设置前置处理函数
func (r *Router) SetPreHandle(f func(request ziface.IRequest)) {
	r.preHandleFunc = f
}

// SetPostHandle 设置后置处理函数
func (r *Router) SetPostHandle(f func(request ziface.IRequest)) {
	r.postHandleFunc = f
}

// PreHandle 前置处理
func (r *Router) PreHandle(request ziface.IRequest) {
	if r.preHandleFunc != nil {
		r.preHandleFunc(request)
	}
}

// Handle 处理请求
func (r *Router) Handle(request ziface.IRequest) {
	// 获取消息类型和数据
	msgID := request.GetMsgID()
	data := request.GetMsgData()

	switch msgID {
	case 1: // PUT 操作
		r.handlePut(data, request)
	case 2: // GET 操作
		r.handleGet(data, request)
	case 3: // DELETE 操作
		r.handleDelete(data, request)
	}
}

// handlePut 处理 PUT 操作
func (r *Router) handlePut(data []byte, request ziface.IRequest) {
	// 解析数据格式：key_len + key + value_len + value
	if len(data) < 8 {
		return
	}

	// 使用 LittleEndian 解析长度字段，与客户端保持一致
	keyLen := int(binary.LittleEndian.Uint32(data[0:4]))
	valueLen := int(binary.LittleEndian.Uint32(data[4:8]))

	if len(data) < 8+keyLen+valueLen {
		return
	}

	key := data[8 : 8+keyLen]
	value := data[8+keyLen : 8+keyLen+valueLen]

	// 创建命令并通过 Raft 追加日志
	cmd := Command{
		Type:  "Put",
		Key:   key,
		Value: value,
	}

	index, err := r.kv.AppendEntry(cmd)
	if err != nil {
		// 发送错误响应
		response := []byte{0x01} // 错误标志
		request.GetConnection().SendMsg(5, response)
		return
	}

	// 等待 Raft 提交确认
	if err := r.kv.WaitForCommit(index); err != nil {
		// 发送错误响应
		response := []byte{0x01} // 错误标志
		request.GetConnection().SendMsg(5, response)
		return
	}

	// 发送成功响应
	response := []byte{0x00} // 成功标志
	request.GetConnection().SendMsg(4, response)
}

// handleGet 处理 GET 操作
func (r *Router) handleGet(data []byte, request ziface.IRequest) {
	// 解析数据格式：key_len + key
	if len(data) < 4 {
		return
	}

	// 使用 LittleEndian 解析长度字段，与客户端保持一致
	keyLen := int(binary.LittleEndian.Uint32(data[0:4]))

	if len(data) < 4+keyLen {
		return
	}

	key := data[4 : 4+keyLen]

	// 从存储获取值
	value, err := r.kv.Get(key)
	if err != nil {
		// 发送错误响应
		response := []byte{0x01} // 错误标志
		request.GetConnection().SendMsg(5, response)
		return
	}

	// 发送成功响应：状态 + value_len + value
	response := make([]byte, 1+4+len(value))
	response[0] = 0x00 // 成功标志

	// 使用 LittleEndian 编码长度字段，与客户端保持一致
	binary.LittleEndian.PutUint32(response[1:5], uint32(len(value)))

	// 写入 value 数据
	copy(response[5:], value)

	request.GetConnection().SendMsg(4, response)
}

// handleDelete 处理 DELETE 操作
func (r *Router) handleDelete(data []byte, request ziface.IRequest) {
	// 解析数据格式：key_len + key
	if len(data) < 4 {
		return
	}

	// 使用 LittleEndian 解析长度字段，与客户端保持一致
	keyLen := int(binary.LittleEndian.Uint32(data[0:4]))

	if len(data) < 4+keyLen {
		return
	}

	key := data[4 : 4+keyLen]

	// 创建命令并通过 Raft 追加日志
	cmd := Command{
		Type: "Delete",
		Key:  key,
	}

	index, err := r.kv.AppendEntry(cmd)
	if err != nil {
		// 发送错误响应
		response := []byte{0x01} // 错误标志
		request.GetConnection().SendMsg(5, response)
		return
	}

	// 等待 Raft 提交确认
	if err := r.kv.WaitForCommit(index); err != nil {
		// 发送错误响应
		response := []byte{0x01} // 错误标志
		request.GetConnection().SendMsg(5, response)
		return
	}

	// 发送成功响应
	response := []byte{0x00} // 成功标志
	request.GetConnection().SendMsg(4, response)
}

// PostHandle 后置处理
func (r *Router) PostHandle(request ziface.IRequest) {
	if r.postHandleFunc != nil {
		r.postHandleFunc(request)
	}
}

// GetFSM 获取 FSM 实例
func (r *Router) GetFSM() *KVServer {
	return r.kv
}
