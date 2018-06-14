package models

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"

	"tsEngine/tsMail"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	Conn *websocket.Conn
}

var connManage = make(map[string]*WsClient)

//连接服务器
func WsConn(wss string) error {

	conn, _, err := websocket.DefaultDialer.Dial(wss, nil)
	if err != nil {
		return err
	}

	var cli WsClient
	cli.Conn = conn
	connManage[wss] = &cli
	return nil

}

//关闭服务器
func WsClose(wss string) {
	conn := connManage[wss].Conn
	conn.Close()
}

//读取数据
func WsRead(wss string) ([]byte, error) {
	conn := connManage[wss].Conn
	_, content, err := conn.ReadMessage()
	if err != nil {
		beego.Error("read:", err)
		go SendMail(wss + "平台无法接受数据，连接已断开，请及时处理~~")
	}
	return content, err
}

//发送数据
func WsSend(wss, content string) {
	conn := connManage[wss].Conn
	err := conn.WriteMessage(websocket.TextMessage, []byte(content))
	if err != nil {
		beego.Error("send:", err)
	}

}

//gzip压缩算法
func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {

		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {

			return nil, err
		}
		return undatas, nil
	}
}

//邮件发送
func SendMail(body string) {

	user := "43943634@qq.com"
	password := "nzyjoxsojksxcbef"
	host := "smtp.qq.com:25"
	subject := "提款机系统报警"

	go tsMail.SendMail(user, password, host, "18732065@qq.com", subject, body, "text")
	go tsMail.SendMail(user, password, host, "43943634@qq.com", subject, body, "text")
}
