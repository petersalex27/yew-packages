package status

type Status[T ~uint] interface {
	IsOk() bool
	NotOk() bool
	Is(T) bool
}