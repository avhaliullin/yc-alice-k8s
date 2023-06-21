package api

import (
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/model"
)

type Request struct {
	Version                string    `json:"version,omitempty"`
	Session                Session   `json:"session,omitempty"`
	Request                *Req      `json:"request,omitempty"`
	AccountLinkingComplete *EmptyObj `json:"account_linking_complete_event,omitempty"`
	State                  ReqState  `json:"state,omitempty"`
}

type ReqState struct {
	Session StateData `json:"session,omitempty"`
}

type Session struct {
	MessageID int                  `json:"message_id,omitempty"`
	SessionID model.AliceSessionID `json:"session_id,omitempty"`
	User      *User                `json:"user,omitempty"`
	New       bool                 `json:"new"`
}

type User struct {
	ID    string `json:"user_id,omitempty"`
	Token string `json:"access_token,omitempty"`
}

type RequestType string

const (
	RequestTypeSimple RequestType = "SimpleUtterance"
	RequestTypeButton RequestType = "ButtonPressed"
)

type Req struct {
	Command           string         `json:"command,omitempty"`
	OriginalUtterance string         `json:"original_utterance,omitempty"`
	NLU               NLU            `json:"nlu,omitempty"`
	Type              RequestType    `json:"type,omitempty"`
	Payload           *ButtonPayload `json:"payload,omitempty"`
}

type NLU struct {
	Tokens   []string `json:"tokens,omitempty"`
	Intents  Intents  `json:"intents,omitempty"`
	Entities []Slot   `json:"entities,omitempty"`
}

type TokensRef struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type Resp struct {
	Text       string    `json:"text,omitempty"`
	TTS        string    `json:"tts,omitempty"`
	Buttons    []*Button `json:"buttons,omitempty"`
	EndSession bool      `json:"end_session"`
}

type Button struct {
	Title   string         `json:"title,omitempty"`
	Payload *ButtonPayload `json:"payload,omitempty"`
	URL     string         `json:"url,omitempty"`
	Hide    bool           `json:"hide"`
}

type Response struct {
	Version             string     `json:"version,omitempty"`
	Response            *Resp      `json:"response,omitempty"`
	StartAccountLinking *EmptyObj  `json:"start_account_linking,omitempty"`
	State               *StateData `json:"session_state,omitempty"`
}

func (r *Response) WithState(s *StateData) *Response {
	r.State = s
	return r
}

type ButtonPayload struct {
}

type EmptyObj struct {
}
