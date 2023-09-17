package conn

import (
	"fmt"
	"partyframe/codec"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

func GetTimeDay(uid int64) int64 {
	return time.Now().Unix()
}

// conn管理器
type ConnManage struct {
	// connMap
	connMap cmap.ConcurrentMap

	listLock   sync.Mutex
	ConnList   *Conn
	listCount  int64
	loginCheck int //app 小游戏登录检测开关
}

var (
	once     sync.Once
	instance *ConnManage
)

func GetManage() *ConnManage {
	once.Do(func() {
		instance = &ConnManage{
			connMap:    cmap.New(),
			ConnList:   nil,
			listCount:  0,
			loginCheck: 0,
		}
	})

	return instance
}

func (m *ConnManage) SetLoginCheck(loginCheck int) {
	m.loginCheck = loginCheck
}

// 添加一个连接
func (m *ConnManage) Set(uid int64, conn *Conn) bool {
	key := strconv.FormatInt(uid, 10)
	m.connMap.Set(key, conn)
	m.listAddTail(conn)
	return true
}

func (m *ConnManage) Get(uid int64) *Conn {
	key := strconv.FormatInt(uid, 10)
	value, ok := m.connMap.Get(key)

	if !ok {
		return nil
	}

	return value.(*Conn)
}

// 移除一个连接
func (m *ConnManage) Remove(uid int64) {
	key := strconv.FormatInt(uid, 10)
	value, ok := m.connMap.Get(key)
	if !ok {
		return
	}

	// wc := value.(*Conn).SocketSession.RawConn()
	// if c, ok := wc.(*websocket.Conn); ok {
	// 	f := c.CloseHandler()
	// 	f(3002, "服务器断开链接")
	// }

	value.(*Conn).SocketSession.Close()
	m.connMap.Remove(key)
	m.listRemove(value.(*Conn))
}

// 发送消息到指定连接
func (m *ConnManage) Send(uid int64, msg *codec.Message) bool {
	key := strconv.FormatInt(uid, 10)
	conn, ok := m.connMap.Get(key)
	if !ok {
		return false
	}
	//if msg.Res_Gate_Tick != nil {
	//	msg.Res_Gate_Tick.ResServerTime = GetTimeDay(uid)
	//}
	//msg.ResServerTime = GetTimeDay(uid) // time.Now().UnixNano() / 1000000
	conn.(*Conn).Send(msg)

	return true
}

func (m *ConnManage) AddResponsesCounter(uid int64, msg *codec.Message) {
	if msg == nil {
		return
	}
	var typeInfo = reflect.TypeOf(*msg)
	var valInfo = reflect.ValueOf(*msg)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		k := typeInfo.Field(i).Name
		v := valInfo.Field(i)

		kArr := strings.Split(k, "_")
		if len(kArr) > 1 {
			if !v.IsNil() {
				break
			}
		}
	}
}

// 发送消息到一批连接
func (m *ConnManage) SendToRange(uids []int64, msg *codec.Message) {
	//msg.ResServerTime = time.Now().UnixNano() / 1000000
	for _, uid := range uids {
		m.Send(uid, msg)
	}
}

// 发送消息到所有连接
func (m *ConnManage) SendToAll(msg *codec.Message) {
	//msg.ResServerTime = time.Now().UnixNano() / 1000000
	go func() {
		//defer tools.TryCatch()

		m.connMap.IterCb(func(key string, v interface{}) {
			conn := v.(*Conn)
			conn.Send(msg)
		})
	}()
}

func (m *ConnManage) doListAddTail(conn *Conn) {
	if m.loginCheck == 0 {
		return
	}
	if conn.Prev != nil || conn.Next != nil {
		//util.Log.Debug(fmt.Sprintf("list not add: %v", conn.Uid))
		return
	}
	if m.ConnList == nil {
		conn.Next = conn
		conn.Prev = conn
		m.ConnList = conn
	} else {
		last := m.ConnList.Prev
		last.Next = conn
		conn.Prev = last
		conn.Next = m.ConnList
		m.ConnList.Prev = conn
	}
	m.listCount += 1
	//util.Log.Debug(fmt.Sprintf("list add: %v", conn.Uid))
}
func (m *ConnManage) listAddTail(conn *Conn) {
	if m.loginCheck == 0 {
		return
	}
	m.listLock.Lock()
	defer m.listLock.Unlock()
	if conn != nil {
		m.doListAddTail(conn)
		m.listPrint()
	}
}

func (m *ConnManage) doListRemove(conn *Conn) {
	if m.loginCheck == 0 {
		return
	}
	if conn.Prev == nil && conn.Next == nil {
		//util.Log.Debug(fmt.Sprintf("list not remove: %v", conn.Uid))
		return
	}
	if conn.Next == conn {
		m.ConnList = nil
	} else {
		prev := conn.Prev
		next := conn.Next
		prev.Next = next
		next.Prev = prev
		if conn == m.ConnList {
			m.ConnList = next
		}
	}
	conn.Next = nil
	conn.Prev = nil
	m.listCount -= 1
	//util.Log.Debug(fmt.Sprintf("list remove: %v", conn.Uid))
}
func (m *ConnManage) listRemove(conn *Conn) {
	if m.loginCheck == 0 {
		return
	}
	m.listLock.Lock()
	defer m.listLock.Unlock()
	if conn != nil {
		m.doListRemove(conn)
		m.listPrint()
	}
}

func (m *ConnManage) listPrint() {
	if m.loginCheck == 0 {
		return
	}
	//return
	//m.listLock.Lock()
	//defer m.listLock.Unlock()
	if m.ConnList == nil {
		return
	}
	str := ""
	head := m.ConnList
	node := head
	for {
		str += fmt.Sprintf("[%v] ", node.Uid)
		node = node.Next
		if node == head {
			break
		}
	}
	//util.Log.Debug(fmt.Sprintf("listPrint: %v, listCount:%d", str, m.listCount))
}

func (m *ConnManage) UpdateConnTime(conn *Conn) {
	if m.loginCheck == 0 {
		return
	}
	m.listLock.Lock()
	defer m.listLock.Unlock()
	m.doListRemove(conn)
	m.doListAddTail(conn)
	m.listPrint()
}

func (m *ConnManage) CheckConnTimeout() {
	if m.loginCheck == 0 {
		return
	}
	m.listLock.Lock()
	defer m.listLock.Unlock()
	if m.ConnList == nil {
		return
	}
	//util.Log.Debug(fmt.Sprintf("list CheckConnTimeout listCount[%v] mapCount[%v]", m.listCount, m.connMap.Count()))
	head := m.ConnList
	node := head
	now := time.Now().Unix()
	delnodes := make([]*Conn, 0)
	for {
		diff := now - node.LastUpdateTime
		if diff > 20 { // 目前client tick是10秒一个
			// 断线了 清在线状态
			//key := fmt.Sprintf("user:online:%d", node.Uid)
			//nredis.GetClient().Del(key)
			////models.UserModelDelOnline(node.Uid, util.RunFuncName())
			delnodes = append(delnodes, node)
			//util.Log.Info(fmt.Sprintf("list CheckConnTimeout user[%v] del", node.Uid))
		} else {
			// 后面的肯定都是没有超时的
			break
		}
		node = node.Next
		if node == head {
			break
		}
	}
	if len(delnodes) > 0 {
		//util.Log.Info(fmt.Sprintf("list CheckConnTimeout Delbefore listCount[%v] mapCount[%v] delcount[%v]", m.listCount, m.connMap.Count(), len(delnodes)))
		for _, conn := range delnodes {
			m.doListRemove(conn)
		}
		//util.Log.Info(fmt.Sprintf("list CheckConnTimeout Delafter listCount[%v] mapCount[%v]", m.listCount, m.connMap.Count()))
	}
	m.listPrint()
}

func (m *ConnManage) GetListCount() int64 {
	m.listLock.Lock()
	defer m.listLock.Unlock()
	return m.listCount
}

func (m *ConnManage) DebugPrintList() {
	m.listLock.Lock()
	defer m.listLock.Unlock()
	if m.ConnList == nil {
		return
	}
	str := ""
	head := m.ConnList
	node := head
	for {
		str += fmt.Sprintf("[%v] ", node.Uid)
		node = node.Next
		if node == head {
			break
		}
	}
	fmt.Printf("DebugPrintList len[%v] list: %v\n", m.listCount, str)
}
