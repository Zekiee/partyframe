package conn

import (
	"partyframe/codec"
	"partyframe/thirdparty/cellnet"
	"time"
)

// 新建一个连接
func NewConn(Uid int64, SocketSession cellnet.Session) *Conn {
	conn := &Conn{
		LastLoginTime:  time.Now().Unix(),
		Uid:            Uid,
		SocketSession:  SocketSession,
		LastUpdateTime: time.Now().Unix(),
		Next:           nil,
		Prev:           nil,
	}
	//conn.Next = conn
	//conn.Prev = conn
	return conn
}

// 连接的struct
// 基本信息, ScoketSession
type Conn struct {
	// 用户ID
	Uid int64

	// 最后登陆时间
	LastLoginTime int64

	// socketSession, 发送消息
	SocketSession cellnet.Session

	//socketId 区分不同连接
	SocketId string

	//最后向usercenter发送时间
	Last2UserCenterTime int64

	// 最后tick的更新时间
	LastUpdateTime int64
	// 双向链表
	Next *Conn
	Prev *Conn
}

// 向当前连接发送消息
func (c *Conn) Send(res *codec.Message) {
	//res.ResUid = c.Uid
	c.SocketSession.Send(res)
}
