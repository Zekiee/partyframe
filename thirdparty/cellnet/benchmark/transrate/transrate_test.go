package transrate

import (
	"partyframe/thirdparty/cellnet"
	"partyframe/thirdparty/cellnet/benchmark"
	_ "partyframe/thirdparty/cellnet/codec/pb" // 启用pb编码
	"partyframe/thirdparty/cellnet/proto/pb/gamedef"
	"partyframe/thirdparty/cellnet/socket"
	"partyframe/thirdparty/cellnet/util"
	"testing"
	"time"
)

var signal *util.SignalTester

// 测试地址
const benchmarkAddress = "127.0.0.1:7201"

// 客户端并发数量
const clientCount = 100

// 测试时间(秒)
const benchmarkSeconds = 10

func server() {

	queue := cellnet.NewEventQueue()
	qpsm := benchmark.NewQPSMeter(queue, func(qps int) {

		logger.Infof("TransmitRate: %d KBps", qps)

	})

	evd := socket.NewAcceptor(queue).Start(benchmarkAddress)

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.Event) {

		if qpsm.Acc() > benchmarkSeconds {
			signal.Done(1)
			logger.Infof("Average QPS: %d", qpsm.Average())
		}

		ev.Send(&gamedef.TestEchoACK{})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewConnector(queue).Start(benchmarkAddress)

	data := make([]byte, 4096)

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.Event) {

		ev.Send(&gamedef.TestEchoACK{
			Bytes: data,
		})

	})

	cellnet.RegisterMessage(evd, "coredef.SessionConnected", func(ev *cellnet.Event) {

		ev.Send(&gamedef.TestEchoACK{})

	})

	queue.StartLoop()

}

func TestTransRate(t *testing.T) {
	signal = util.NewSignalTester(t)

	// 超时时间为测试时间延迟一会
	signal.SetTimeout((benchmarkSeconds + 5) * time.Second)

	server()

	for i := 0; i < clientCount; i++ {
		go client()
	}

	signal.WaitAndExpect("recv time out", 1)

}
