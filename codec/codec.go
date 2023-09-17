package codec

type MessageType int

type Message struct {
}

type Codec interface {
	Reader
	Writer
	Close() error
	String() string
}

type Reader interface {
	ReadHeader(*Message, MessageType) error
	ReadBody(interface{}) error
}

type Writer interface {
	Write(*Message, interface{}) error
}

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
