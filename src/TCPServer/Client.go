package TCPServer

import (
	"Crypt"
	"Logger"
	"bufio"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type TCPClient struct {
	conn    net.Conn
	Logger  Logger.Interface
	Host    string
	Port    int
	Channel chan interface{}
	Key     string
	Iv      string
	Id      string
}

type Request struct {
	Command    string            `json:"command"`
	Proc       string            `json:"proc"`
	Connection string            `json:"connection"`
	Data       map[string]string `json:"data"`
}

func NewClient() *TCPClient {
	return &TCPClient{}
}

func (c *TCPClient) Servers(host string, port int) *TCPClient {
	c.Host = host
	c.Port = port
	c.Channel = make(chan interface{}, 2)
	return c
}

func (c *TCPClient) SetId(serverId string) *TCPClient {
	c.Id = serverId
	return c
}

func (c *TCPClient) Secrets(key string, iv string) *TCPClient {
	c.Key = key
	c.Iv = iv
	return c
}

func (c *TCPClient) SetLogger(logg Logger.Interface) *TCPClient {
	c.Logger = logg
	return c
}

func (c *TCPClient) Connect() (*TCPClient, error) {
	addr := strings.Join([]string{c.Host, strconv.Itoa(c.Port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return c, err
	}
	c.conn = conn
	return c, nil
}

func (c *TCPClient) IsConnected() (bool, error) {
	one := make([]byte, 1)
	err := c.conn.SetReadDeadline(time.Now())
	if err != nil {
		return false, err
	}
	if _, err := c.conn.Read(one); err == io.EOF {
		err := c.conn.Close()
		if err != nil {
			return false, err
		}
		c.conn = nil
		return false, err
	} else {
		err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		if err != nil {
			return false, err
		}
	}
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return false, neterr
	}
	return true, nil
}

func (c *TCPClient) SendAndRead(request Request) string {
	clientCrypter, _ := Crypt.NewCrypter([]byte(c.Key), []byte(c.Iv))
	command := clientCrypter.EncodeToString(request.Command)
	_, errF := c.conn.Write([]byte(command))
	if errF != nil {
		c.Logger.Error(10002, "conn.Write. Error: %v", errF)
		return ""
	}
	reader := bufio.NewReader(c.conn)
	result, _ := reader.ReadString('\n')
	if len(result) > 0 {
		return clientCrypter.DecodeToString(result)
	}
	return ""
}

func (c *TCPClient) SendJson(request Request) string {
	clientCrypter, _ := Crypt.NewCrypter([]byte(c.Key), []byte(c.Iv))
	marsh, _ := json.Marshal(request)
	//fmt.Println("marsh: " + string(marsh))
	command := clientCrypter.EncodeToString(string(marsh))
	//cmd := clientCrypter.DecodeToString(command)
	//fmt.Println("command: " + command)
	//fmt.Println("cmd    : " + cmd)
	_, errF := c.conn.Write([]byte(command + "\n"))
	if errF != nil {
		c.Logger.Error(10002, "conn.Write. Error: %v", errF)
		return ""
	}

	reader := bufio.NewReader(c.conn)
	result, _ := reader.ReadString('\n')
	//c.conn.Close()
	if len(result) > 0 {
		rest, _ := clientCrypter.Decode(result)
		return rest
	}
	return ""
}

func (c *TCPClient) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}
