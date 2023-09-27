package notification

type Notification interface {
	Notify(of string) string
}