package enum

type Item int

const (
	PENDING Item = iota
	DOING
	DONE
)

func (i Item) Find() ItemValue {
	switch i {
	case PENDING:
		return ItemValue{0, "pending"}
	case DOING:
		return ItemValue{1, "doing"}
	case DONE:
		return ItemValue{2, "done"}
	default:
		return ItemValue{0, "pending"}
	}
}

type ItemValue struct {
	INDEX int
	Name  string
}
