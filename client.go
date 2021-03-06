package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}
	client.conn = conn

	return client
}

func (c *Client) DealResponse() {
	// 一旦有数据, 就拷贝到stdout到标准输出上,永久阻塞
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>请输入合法数字!<<<<<<")
		return false
	}
	return true
}

func (c *Client) UpdateName() bool {
	fmt.Println("请输入用户名:")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (c *Client) PublicChat() {
	// 提示用户输入聊天内容
	var chatMsg string
	fmt.Println("请输入聊天内容, exit退出")
	fmt.Scanln(&chatMsg)
	for {
		if chatMsg != "exit" {
			// 发送给服务器
			if len(chatMsg) != 0 {
				sendMsg := chatMsg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn,Write err: ", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入聊天内容, exit退出")
			fmt.Scanln(&chatMsg)

		}
	}
}

func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

func (c *Client) PrivateChat() {
	// 查询当前谁在线
	c.SelectUsers()
	fmt.Println("请输入聊天的对象的用户名, exit退出")
	var remoteName string
	var chatMsg string
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println("请输入消息内容")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			// 发送给服务器
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn,Write err: ", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入消息内容")
			fmt.Scanln(&chatMsg)
		}

		c.SelectUsers()
		fmt.Println("请输入聊天的对象的用户名, exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() == true {
			// 根据不同模式处理不同业务
			switch c.flag {
			case 1:
				c.PublicChat()
				break
			case 2:
				c.PrivateChat()
				break
			case 3:
				c.UpdateName()
				break
			}
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址,默认地址是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器PORT端口,默认地址端口是8888")
}

func main() {
	// 命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>链接服务器成功!")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>>>>>链接服务器成功!!!!")
	client.Run()
}
