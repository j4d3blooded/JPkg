package jpkg_impl

type HasherFlag uint8

const (
	HASHER_NONE HasherFlag = iota
)

type HasherHandler interface {
	Flag() HasherFlag
}

type NullHasherHandler struct {
}

func (n *NullHasherHandler) Flag() HasherFlag {
	return HASHER_NONE
}
