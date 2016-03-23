package errorinfo

import (
	"errors"
)

var (
	ClientErrBadComLineFormat = errors.New("CLIENT_ERROR bad command line format")
	ClientErrBadDataChunk     = errors.New("CLIENT_ERROR bad data chunk")
	ClientErrUnknownProtocol  = errors.New("unknown protocol")
	FlagErr                   = errors.New("Flag Error")
	ExpireErr                 = errors.New("Expire Error")
	NotFounded                = errors.New("Not Founded")
	Error                     = errors.New("ERROR")
)
