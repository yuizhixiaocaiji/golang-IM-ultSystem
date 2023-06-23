package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Message 消息
type Message struct {
	gorm.Model
	FromId   int64  //发送者
	TargetId int64  //接受者
	Type     int    //发送类型 1.群聊 2.私聊 3.心跳
	Media    int    //消息类型 1.文字 2.图片 3.音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn          *websocket.Conn
	Addr          string //客户端地址
	FirstTime     uint64 //首次连接时间
	HeartbeatTime uint64 //心跳时间
	LoginTime     uint64 //登录时间
	DataQueue     chan []byte
	GroupSets     set.Interface
}

//映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

//读写锁
var rwLocker sync.RWMutex

// Chat 通信服务：需要：1.发送者ID， 接受者ID，消息类型，发送的内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并校验 token等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalida := true //checkToken()
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2.获取conn
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}

	//3.用户关系
	//4.userId 跟 node绑定 并且加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送的逻辑
	go sendProc(node)
	//6.完成接受的逻辑
	go recvProc(node)

	//sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] sendProc >>>> msg: ", string(data))
			err :=
				node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		//心跳检测 msg.Media == -1 || msg.Type == 3
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			broadMsg(data) //todo 将消息广播到局域网
			fmt.Println("[ws] recvProc <<<<< ", string(data))
		}
	}
}

// Heartbeat 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine")
}

// 完成udp数据发送携程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc data : ", string(data))
			_, err :=
				conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接受携程
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()
	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

//后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch data : ", string(data))
		sendMsg(msg.TargetId, data)
		//case 2://群发
		//	sendGroupMsg()
		//case 3://广播
		//	sendAllMsg
		//case 4:
		//
	}
}

func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg >>> userID: ", userId, " msg: ", string(msg))
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
