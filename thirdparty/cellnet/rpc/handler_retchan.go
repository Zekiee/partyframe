package rpc

import "partyframe/thirdparty/cellnet"

type RetChanHandler struct {
	ret chan interface{}
}

func (self *RetChanHandler) Call(ev *cellnet.Event) {

	self.ret <- ev.Msg
}

func NewRetChanHandler(ret chan interface{}) cellnet.EventHandler {
	return &RetChanHandler{
		ret: ret,
	}

}
