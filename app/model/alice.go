package model

type AliceSessionID string

type AliceSessionStateID string

const (
	AliceStateNew     AliceSessionStateID = "NEW"
	AliceStateItemAdd AliceSessionStateID = "ITEM_ADD"
)

type AliceSessionState struct {
	State AliceSessionStateID

	ListName      string
	ListConfirmed bool
	ListAsked     bool

	ItemText      string
	ItemConfirmed bool
	ItemAsked     bool
}

type AliceSession struct {
	ID    AliceSessionID
	State *AliceSessionState
}
