package extend

import (
	"partyframe/thirdparty/cellnet"
	"sync"
)

type PushMessage struct {
	C   chan struct{}
	Msg *cellnet.Event
	IsForMe func(tag interface{}) bool
}

type PushMessageManager struct {
	idx  uint16
	lock sync.RWMutex
	msg  [65536]*PushMessage
}

var (
	once     sync.Once
	instance *PushMessageManager
)

func init() {
	PushMsgMgr()
}

func PushMsgMgr() *PushMessageManager {
	once.Do(func() {
		instance = &PushMessageManager {
			idx: 0,
		}
		instance.msg[0] = &PushMessage{
			C: make(chan struct{}),
		}
	})
	return instance
}

func (this *PushMessageManager) GetIdx() uint16 {
	var idx uint16
	this.lock.RLock()
	idx = this.idx
	this.lock.RUnlock()
	return idx
}

func (this *PushMessageManager) Send(data interface{}, isForMe func(tag interface{}) bool) {
	ev := cellnet.NewEvent(cellnet.Event_Send, nil)
	ev.Msg = data
	this.lock.Lock()
	msg := this.msg[this.idx]
	this.idx++
	this.msg[this.idx] = &PushMessage{
		C: make(chan struct{}),
	}
	this.lock.Unlock()
	msg.Msg = ev
	msg.IsForMe = isForMe
	close(msg.C)
}

func (this *PushMessageManager) GetMsg(idx uint16) *PushMessage {
	return this.msg[idx]
}

