package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"partyframe/thirdparty/cellnet"
	"partyframe/thirdparty/cellnet/extend"
)

type wsAcceptor struct {
	*wsPeer
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func (self *wsAcceptor) Start(address string) cellnet.Peer {

	if self.IsRunning() {
		return self
	}

	self.SetRunning(true)

	url, err := url.Parse(address)

	if err != nil {
		logger.Errorln(err, address)
		return self
	}

	if url.Path == "" {
		logger.Errorln("websocket: expect path in url to listen", address)
		return self
	}

	self.SetAddress(address)

	http.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf("websocket err:%v", err)
			return
		}
		//fmt.Println(fmt.Sprintf("%v", r.Header))

		ses := newSession(c, self)
		ses.SetForwardedFor(r.Header.Values("X-Forwarded-For"))
		// 添加到管理器
		self.Add(ses)

		// 断开后从管理器移除
		ses.OnClose = func() {
			self.Remove(ses)
		}

		ses.run()

		// 通知逻辑
		extend.PostSystemEvent(ses, cellnet.Event_Accepted, self.ChainListRecv(), cellnet.Result_OK)

	})

	go func() {

		err = http.ListenAndServe(url.Host, nil)

		if err != nil {
			logger.Errorln(err)
		}

		self.SetRunning(false)

	}()

	return self
}

func (self *wsAcceptor) Stop() {
	if !self.IsRunning() {
		return
	}
}

func NewAcceptor(q cellnet.EventQueue) cellnet.Peer {

	self := &wsAcceptor{
		wsPeer: newPeer(q, cellnet.NewSessionManager()),
	}

	return self
}
