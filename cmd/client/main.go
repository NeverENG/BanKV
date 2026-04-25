package main

import (
	"fmt"
	"os"

	"github.com/NeverENG/BanKV/cmd/client"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := os.Args[1]
	addr := "localhost:8080" // 默认服务端地址

	// 创建客户端
	c := client.NewClient(addr)

	// 连接服务端
	err := c.Connect()
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer c.Close()

	switch cmd {
	case "put":
		if len(os.Args) < 4 {
			fmt.Println("Usage: client put <key> <value>")
			return
		}
		key := os.Args[2]
		value := os.Args[3]
		err := c.SendPut([]byte(key), []byte(value))
		if err != nil {
			fmt.Printf("Error putting key-value: %v\n", err)
			return
		}
		fmt.Printf("Put successful: %s = %s\n", key, value)

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: client get <key>")
			return
		}
		key := os.Args[2]
		value, err := c.SendGet([]byte(key))
		if err != nil {
			fmt.Printf("Error getting value: %v\n", err)
			return
		}
		fmt.Printf("Get successful: %s = %s\n", key, string(value))

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: client delete <key>")
			return
		}
		key := os.Args[2]
		err := c.SendDelete([]byte(key))
		if err != nil {
			fmt.Printf("Error deleting key: %v\n", err)
			return
		}
		fmt.Printf("Delete successful: %s\n", key)

	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  client put <key> <value> - Put a key-value pair")
	fmt.Println("  client get <key>        - Get value by key")
	fmt.Println("  client delete <key>     - Delete a key")
}
