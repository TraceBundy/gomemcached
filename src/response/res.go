package response

import (
	"constarg"
	"io"
)

type Response struct {
	Protocol uint16
	Ctype    constarg.CommandType
	Key      interface{}
	Body     []byte
}

func Transmit(writer io.Writer, res *Response) error {
	bodylen := len(res.Body)
	len := 0
	for len < bodylen {
		cur, err := writer.Write(res.Body)
		if err != nil {
			return err
		}
		len += cur
	}
	return nil
}
