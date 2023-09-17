package tests

import (
	"partyframe/thirdparty/cellnet"
	jsongamedef "partyframe/thirdparty/cellnet/proto/json/gamedef" // json逻辑协议
	"partyframe/thirdparty/cellnet/util"
	"partyframe/thirdparty/cellnet/websocket"
	"testing"
)

var wsSignal *util.SignalTester

func wsServer() {

	queue := cellnet.NewEventQueue()

	// 注意, 如果http代理/VPN在运行时可能会导致无法连接, 请关闭

	p := websocket.NewAcceptor(queue).Start("http://127.0.0.1:8801/echo")
	p.SetName("server")

	cellnet.RegisterMessage(p, "coredef.SessionAccepted", func(ev *cellnet.Event) {

		logger.Debugln("client accepted")
	})

	cellnet.RegisterMessage(p, "gamedef.TestEchoJsonACK", func(ev *cellnet.Event) {

		msg := ev.Msg.(*jsongamedef.TestEchoJsonACK)

		logger.Debugln(msg.Content)

		ev.Send(&jsongamedef.TestEchoJsonACK{Content: "roger"})
	})

	queue.StartLoop()
}

func wsClient() {

	queue := cellnet.NewEventQueue()

	p := websocket.NewConnector(queue).Start("ws://127.0.0.1:8801/echo")
	p.SetName("client")

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		logger.Debugln("client connected")

		// 发送消息, 底层自动选择pb编码
		ev.Send(&jsongamedef.TestEchoJsonACK{
			Content: "hello",
		})

		wsSignal.Done(1)

	})

	cellnet.RegisterMessage(p, "gamedef.TestEchoJsonACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*jsongamedef.TestEchoJsonACK)

		logger.Debugln("client recv:", msg.Content)

		wsSignal.Done(2)
	})

	queue.StartLoop()

	wsSignal.WaitAndExpect("not recv data", 1, 2)

}

func TestWebsocketEcho(t *testing.T) {

	wsSignal = util.NewSignalTester(t)

	wsServer()

	wsClient()

}
