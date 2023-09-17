package websocket

import (
	"github.com/gorilla/websocket"
	"partyframe/thirdparty/cellnet"
	"partyframe/thirdparty/cellnet/extend"
)

type wsSession struct {
	OnClose func() // 关闭函数回调

	id int64

	xForwardedFor interface{}
	p             cellnet.Peer

	conn *websocket.Conn

	tag interface{}

	sendChan chan *cellnet.Event

	// 推送队列消息索引
	pushMsgIdx uint16
}

func (c *wsSession) Tag() interface{} {
	return c.tag
}
func (c *wsSession) SetTag(tag interface{}) {
	c.tag = tag
}

func (c *wsSession) ID() int64 {
	return c.id
}

func (c *wsSession) SetID(id int64) {
	c.id = id
}

func (c *wsSession) FromPeer() cellnet.Peer {
	return c.p
}

func (c *wsSession) SetForwardedFor(xF interface{}) {
	c.xForwardedFor = xF
}

func (c *wsSession) GetForwardedFor() interface{} {
	return c.xForwardedFor
}

func (c *wsSession) RemoteAddr() interface{} {
	return c.conn.RemoteAddr().String()
}

func (c *wsSession) Close() {
	select {
	case c.sendChan <- nil:
		return
	default:
		logger.Error("send message to sendChan for close error.")
	}
}

func (c *wsSession) Send(data interface{}) {

	ev := cellnet.NewEvent(cellnet.Event_Send, c)
	ev.Msg = data

	if ev.ChainSend == nil {
		ev.ChainSend = c.p.ChainSend()
	}

	c.RawSend(ev)
}

func (c *wsSession) RawSend(ev *cellnet.Event) {
	ev.Ses = c
	if ev.ChainSend != nil {
		ev.ChainSend.Call(ev)
	}

	// 发送日志
	cellnet.MsgLog(ev)

	// 放入发送队列
	select {
	case c.sendChan <- ev:
		return
	default:
		logger.Error("send message to sendChan error.")
	}
}

func (c *wsSession) RawConn() interface{} {
	return c.conn
}

func (c *wsSession) GetNextPushMsg(move bool) *extend.PushMessage {
	if move {
		c.pushMsgIdx++
	}
	return extend.PushMsgMgr().GetMsg(c.pushMsgIdx)
}

func (c *wsSession) sendThread() {
	defer func() {
		if r := recover(); r != nil {
			//nerr.ProcessError(0, r)
		}

		c.conn.Close()
		c.OnClose()
	}()

	pm := c.GetNextPushMsg(false)
	for {
		var ev *cellnet.Event
		endFlag := false
		select {
		case ev = <-c.sendChan:
			if ev == nil {
				endFlag = true
			}
		case <-pm.C:
			if pm.IsForMe == nil || pm.IsForMe(c.tag) {
				ev = pm.Msg
				if ev != nil {
					ev.ChainSend = c.p.ChainSend()
					ev.ChainSend.Call(ev)
				}
			}
			pm = c.GetNextPushMsg(true)
		}

		if endFlag {
			c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			break
		}

		if ev == nil {
			continue
		}

		meta := cellnet.MessageMetaByID(ev.MsgID)

		if meta == nil {
			// TODO
			//logger.Error("websocket sendThread Result_CodecError")
			ev.SetResult(cellnet.Result_CodecError)
			continue
		}

		// 组websocket包
		raw := composeBinaryPacket(meta.Name, ev.Data)
		c.conn.WriteMessage(websocket.BinaryMessage, raw)
	}
}

func (c *wsSession) ReadPacket() (msgid uint32, data []byte, result cellnet.Result) {

	// 读超时
	t, raw, err := c.conn.ReadMessage()
	if err != nil {
		logger.Error("websocket err:", err.Error())
		return 0, nil, errToResult(err)
	}

	switch t {
	case websocket.TextMessage:

		msgName, userdata := parsePacket(raw)

		data = userdata

		if msgName != "" {

			meta := cellnet.MessageMetaByName(msgName)

			if meta == nil || meta.Codec == nil {
				logger.Error("websocket CodecError ", err.Error())
				return 0, nil, cellnet.Result_CodecError
			}

			msgid = meta.ID

		}

	case websocket.BinaryMessage:
		_, userData := parseBinaryPacket(raw)
		return 3634688514, userData, cellnet.Result_OK

	case websocket.CloseMessage:
		return 0, nil, cellnet.Result_RequestClose
	}

	return msgid, data, cellnet.Result_OK
}

func (c *wsSession) recvThread() {

	defer func() {
		if r := recover(); r != nil {
			//nerr.ProcessError(0, r)
		}
	}()

	for {

		msgid, data, result := c.ReadPacket()

		chainList := c.p.ChainListRecv()

		if result != cellnet.Result_OK {
			extend.PostSystemEvent(c, cellnet.Event_Closed, chainList, result)
			break
		}

		ev := cellnet.NewEvent(cellnet.Event_Recv, c)
		ev.MsgID = msgid
		ev.Data = data

		// 接收日志
		cellnet.MsgLog(ev)

		chainList.Call(ev)

		if ev.Result() != cellnet.Result_OK {
			// TODO
			//if ev.Result() != cellnet.Result_RequestClose {
			//	logger.Errorf("websocket recvThread error2, code:%d", ev.Result())
			//}
			extend.PostSystemEvent(ev.Ses, cellnet.Event_Closed, chainList, ev.Result())
			break
		}
	}
}

func (c *wsSession) run() {

	go c.recvThread()

	go c.sendThread()
}

func newSession(c *websocket.Conn, p cellnet.Peer) *wsSession {

	self := &wsSession{
		p:          p,
		conn:       c,
		sendChan:   make(chan *cellnet.Event, 256),
		pushMsgIdx: extend.PushMsgMgr().GetIdx(),
	}

	return self
}
