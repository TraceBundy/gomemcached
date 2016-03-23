package request

import (
	"bufio"
	"constarg"
	"errorinfo"
	"io"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Protocol uint16
	Ctype    constarg.CommandType
	Key      string
	Expire   time.Duration
	Flag     uint16
	Cas      uint64
	Data     []byte
	Num      uint32
}

func Handle(r io.Reader) (*Request, error) {
	reader := bufio.NewReader(r)
	head, err := reader.Peek(1)
	if err == nil {
		if head[0] == byte(constarg.BINARY) {
			return ReceiveBinary(reader)
		} else {
			return ReceiveAscii(reader)
		}
		return nil, errorinfo.ClientErrUnknownProtocol
	} else {
		return nil, err
	}
}
func ReceiveBinary(reader *bufio.Reader) (*Request, error) {
	return nil, nil
}
func ReceiveCRLF(reader *bufio.Reader) ([]byte, error) {
	buf := make([]byte, 100)
	var data []byte
	for {
		len, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf[0:len]...)
		if len > 2 && buf[len-2] == '\r' && buf[len-1] == '\n' {
			break
		}
	}
	return data[0 : len(data)-2], nil
}
func ReceiveAscii(reader *bufio.Reader) (*Request, error) {
	data, err := ReceiveCRLF(reader)
	if err != nil {
		return nil, err
	}
	command := strings.Split(string(data), " ")
	req := new(Request)
	c, ok := constarg.CommandTypeMap[strings.ToUpper(command[0])]
	if !ok {
		return nil, errorinfo.Error
	}
	req.Ctype = c
	if c == constarg.SET || c == constarg.ADD || c == constarg.REPLACE {
		if len(command) != 5 {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Key = command[1]
		flag, err := strconv.ParseUint(command[2], 10, 16)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		expire, err := strconv.ParseInt(command[3], 10, 64)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		total, err := strconv.ParseInt(command[4], 10, 64)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Flag = uint16(flag)
		req.Expire = time.Duration(expire) * time.Second
		data, err := ReceiveCRLF(reader)
		req.Data = data
		if len(req.Data) != int(total) {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
	} else if c == constarg.GET || c == constarg.GETS {
		if len(command) != 2 {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Key = command[1]
	} else if c == constarg.DELETE {
		if len(command) != 3 {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Key = command[1]
		expire, err := strconv.ParseInt(command[2], 10, 64)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Expire = time.Duration(expire) * time.Second
	} else if c == constarg.INCR || c == constarg.DECR {
		if len(command) != 3 {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Key = command[1]
		value, err := strconv.ParseUint(command[2], 10, 32)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Num = uint32(value)
	} else if c == constarg.VERSION {
		if len(command) != 1 {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
	} else if c == constarg.CAS {
		if len(command) != 6 {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Key = command[1]
		flag, err := strconv.ParseUint(command[2], 10, 16)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		expire, err := strconv.ParseInt(command[3], 10, 64)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		cas, err := strconv.ParseUint(command[5], 10, 64)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		total, err := strconv.ParseInt(command[4], 10, 64)
		if err != nil {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
		req.Flag = uint16(flag)
		req.Expire = time.Duration(expire) * time.Second
		req.Cas = cas
		data, err := ReceiveCRLF(reader)
		req.Data = data
		if len(req.Data) != int(total) {
			return nil, errorinfo.ClientErrBadComLineFormat
		}
	}
	return req, nil
}
