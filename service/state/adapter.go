package state

import "github.com/republicprotocol/renex-swapper-go/service/logger"

type Adapter interface {
	logger.Logger
	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error
	Delete([]byte) error
}