package constarg

type CommandType uint32

var (
	GET     CommandType = 1
	GETS    CommandType = 2
	SET     CommandType = 3
	ADD     CommandType = 4
	REPLACE CommandType = 5
	DELETE  CommandType = 6
	INCR    CommandType = 7
	DECR    CommandType = 8
	VERSION CommandType = 9
	CAS     CommandType = 10
)
var (
	STORED           = "STORED\r\n"
	NotStored        = "NOT_STORED\r\n"
	NotFounded       = "Not Founded\r\n"
	CRLF             = "\r\n"
	BadDataChunk     = "bad data chunk"
	BadComLineFormat = "bad command line format"
	VALUE            = "VALUE"
	DELETED          = "DELETED\r\n"
	NOT_FOUND        = "NOT_FOUND\r\n"
	VERSION_INFO     = "VERSION 0.1\r\n"
	EXIST            = "EXIST\r\n"
)

const (
	BINARYREQ = iota
	ASCIIREQ
	BINARYRES
	ASCIIRES
)

var (
	BINARY = 0x80
)

type CommandStr string

var CommandStrMap = map[CommandType]string{
	GET:     "GET",
	GETS:    "GETS",
	SET:     "SET",
	ADD:     "ADD",
	REPLACE: "REPLACE",
	DELETE:  "DELETE",
	INCR:    "INCR",
	DECR:    "DECR",
	VERSION: "VERSION",
	CAS:     "CAS",
}
var CommandTypeMap = map[string]CommandType{
	"GET":     GET,
	"GETS":    GETS,
	"SET":     SET,
	"ADD":     ADD,
	"REPLACE": REPLACE,
	"DELETE":  DELETE,
	"INCR":    INCR,
	"VERSION": VERSION,
	"CAS":     CAS,
}
