package server

import (
	"cache"
	"constarg"
	"errorinfo"
	"net"
	"request"
	"response"
	"strconv"
	"sync/atomic"
)

var (
	table cache.Table
	cas   uint64
)

type HandleFunc func(*request.Request) (*response.Response, error)

func Init() {
	table = cache.GetTable("global")
}

var HandleFuncMap = map[constarg.CommandType]HandleFunc{
	constarg.SET:     HandleSet,
	constarg.GET:     HandleGet,
	constarg.GETS:    HandleGet,
	constarg.ADD:     HandleAdd,
	constarg.REPLACE: HandleReplace,
	constarg.DELETE:  HandleDelete,
	constarg.INCR:    HandleIncr,
	constarg.DECR:    HandleDecr,
	constarg.VERSION: HandleVersion,
	constarg.CAS:     HandleCas,
}

func ConnectionHandle(conn net.Conn) error {
	defer func() { conn.Close() }()
	return HandleIO(conn)
}

func HandleIO(conn net.Conn) error {
	var err error = nil
	for err == nil {
		err = HandleMessage(conn)
	}
	return err
}

func HandleMessage(conn net.Conn) error {
	req, err := request.Handle(conn)
	if err != nil {
		return response.Transmit(conn, &response.Response{
			Body: []byte(err.Error() + "\r\n"),
		})
	}
	handle := HandleFuncMap[req.Ctype]
	res, err := handle(req)
	return response.Transmit(conn, res)
}

func HandleSet(req *request.Request) (*response.Response, error) {
	atomic.AddUint64(&cas, 1)
	table.Insert(req.Key, req.Data, req.Flag, &cas, req.Expire, nil)
	res := new(response.Response)
	res.Body = []byte(constarg.STORED)
	return res, nil
}
func HandleGet(req *request.Request) (*response.Response, error) {
	res := new(response.Response)
	value, err := table.Get(req.Key)
	if err != nil {
		res.Body = []byte(constarg.NotFounded)
		return res, errorinfo.NotFounded
	}
	res.Body = append(res.Body, []byte(constarg.VALUE)...)
	res.Body = append(res.Body, []byte(" ")...)
	res.Body = append(res.Body, []byte(value.Key)...)
	res.Body = append(res.Body, []byte(" ")...)
	flag := strconv.Itoa(int(value.Flag))
	res.Body = append(res.Body, []byte(flag)...)
	res.Body = append(res.Body, []byte(" ")...)
	flag = strconv.Itoa(len(value.Value))
	res.Body = append(res.Body, []byte(flag)...)
	if req.Ctype == constarg.GETS {
		res.Body = append(res.Body, []byte(" ")...)
		cas := strconv.FormatUint(uint64(value.Cas), 10)
		res.Body = append(res.Body, []byte(cas)...)
	}
	res.Body = append(res.Body, []byte(" ")...)
	res.Body = append(res.Body, []byte(constarg.CRLF)...)
	res.Body = append(res.Body, value.Value...)
	res.Body = append(res.Body, []byte(constarg.CRLF)...)
	return res, nil
}
func HandleAdd(req *request.Request) (*response.Response, error) {
	atomic.AddUint64(&cas, 1)
	res := new(response.Response)
	err := table.Add(req.Key, req.Data, req.Flag, &cas, req.Expire, nil)
	if err != nil {
		res.Body = []byte(constarg.NotStored)
	}
	return res, nil
}
func HandleReplace(req *request.Request) (*response.Response, error) {
	atomic.AddUint64(&cas, 1)
	res := new(response.Response)
	err := table.Replace(req.Key, req.Data, req.Flag, &cas, req.Expire, nil)
	if err != nil {
		res.Body = []byte(constarg.NotStored)
	} else {
		res.Body = []byte(constarg.STORED)
	}
	return res, nil
}
func HandleDelete(req *request.Request) (*response.Response, error) {
	res := new(response.Response)
	err := table.Delete(req.Key)
	if err != nil {
		res.Body = []byte(constarg.NOT_FOUND)
	} else {
		res.Body = []byte(constarg.DELETED)
	}
	return res, nil
}
func HandleIncr(req *request.Request) (*response.Response, error) {
	res := new(response.Response)
	val, err := table.Incr(req.Key, req.Num)
	if err != nil {
		res.Body = []byte(constarg.NOT_FOUND)
	} else {
		flag := strconv.FormatUint(uint64(val), 10)
		res.Body = []byte(flag)
		res.Body = append(res.Body, []byte(constarg.CRLF)...)
	}
	return res, nil
}
func HandleDecr(req *request.Request) (*response.Response, error) {
	atomic.AddUint64(&cas, 1)
	res := new(response.Response)
	val, err := table.Decr(req.Key, req.Num)
	if err != nil {
		res.Body = []byte(constarg.NOT_FOUND)
	} else {
		flag := strconv.FormatUint(uint64(val), 10)
		res.Body = []byte(flag)
		res.Body = append(res.Body, []byte(constarg.CRLF)...)
	}
	return res, nil
}

func HandleCas(req *request.Request) (*response.Response, error) {
	err := table.Cas(req.Key, req.Data, req.Flag, &cas, req.Cas, req.Expire, nil)
	res := new(response.Response)
	if err != nil {
		res.Body = []byte(constarg.EXIST)
	} else {
		res.Body = []byte(constarg.STORED)
	}
	return res, nil
}

func HandleVersion(req *request.Request) (*response.Response, error) {
	res := new(response.Response)
	res.Body = []byte(constarg.VERSION_INFO)
	return res, nil
}
