package gamedef

import (
	"partyframe/thirdparty/cellnet"
	_ "partyframe/thirdparty/cellnet/codec/json"
	"partyframe/thirdparty/cellnet/util"
	"partyframe/thirdparty/goobjfmt"
	"reflect"
)

type TestEchoJsonACK struct {
	Content string
}

func (m *TestEchoJsonACK) String() string { return goobjfmt.CompactTextString(m) }

func init() {

	// coredef.proto
	cellnet.RegisterMessageMeta("json", "gamedef.TestEchoJsonACK", reflect.TypeOf((*TestEchoJsonACK)(nil)).Elem(), util.StringHash("gamedef.TestEchoJsonACK"))
}
