package service

type NotificationType int

const (
	Email NotificationType = iota
	SMS
	HTTP
	// Notifications can be different types (can be expanded)
)

type Notification struct {
	ID       int64
	Type     NotificationType
	Address  string
	IsSended bool
}

func (n *Notification) getID() int64 {
	return n.ID
}

func (n *Notification) getType() NotificationType {
	return n.Type
}

func (n *Notification) getAddress() string {
	return n.Address
}

func (n *Notification) isSended() bool {
	return n.IsSended
}
