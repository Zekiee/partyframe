package main

import (
	"bufio"
	"os"
	"partyframe/thirdparty/cellnet"
	"partyframe/thirdparty/cellnet/examples/chat/proto/chatproto"
	"partyframe/thirdparty/cellnet/socket"
	"strings"
)

func ReadConsole(callback func(string)) {

	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		text = strings.TrimRight(text, "\n\r ")

		text = strings.TrimLeft(text, " ")

		callback(text)
	}
}

func main() {
	queue := cellnet.NewEventQueue()

	peer := socket.NewConnector(queue).Start("127.0.0.1:8801")
	peer.SetName("client")

	cellnet.RegisterMessage(peer, "chatproto.ChatACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*chatproto.ChatACK)

		logger.Infof("sid%d say: %s", msg.Id, msg.Content)
	})

	queue.StartLoop()

	ReadConsole(func(str string) {

		peer.(socket.Connector).DefaultSession().Send(&chatproto.ChatREQ{
			Content: str,
		})
	})
}
