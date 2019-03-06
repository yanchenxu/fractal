package types

import "github.com/fractalplatform/fractal/common"

type DetailTx struct {
	TxHash      common.Hash
	InternalTxs []*InternalTx
}

type InternalTx struct {
	Action     *Action
	ActionType string
	GasUsed    uint64
	GasLimit   uint64
	Depth      uint64
	Error      error
}
