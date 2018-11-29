package http

type Hooker interface {
	Hook(x *Group)
}

const (
	AcInit Action = iota
	AcRequest
	AcResponse
	AcDone

	AcPacker
	AcUnPacker
)

func (m Action) String() string {
	switch m {
	case AcInit:
		return "AcInit"
	case AcRequest:
		return "AcRequest"
	case AcResponse:
		return "AcResponse"
	case AcDone:
		return "AcDone"
	default:
		return "AcUnknown"
	}
}
