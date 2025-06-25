package domain

type PositionType int

const (
	POS_HOLD PositionType = 0
	POS_BUY  PositionType = 1
	POS_SELL PositionType = 2
)

type Position struct {
	PosType    PositionType
	EntryPrice float64
}
