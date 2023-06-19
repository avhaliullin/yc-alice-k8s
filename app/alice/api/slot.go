package api

type Slot struct {
	Type   string      `json:"type"`
	Tokens *TokensRef  `json:"tokens"`
	Value  interface{} `json:"value"`
}

func (s *Slot) AsString() (string, bool) {
	if s == nil || s.Type != "YANDEX.STRING" {
		return "", false
	}
	value := s.Value.(string)
	return value, value != ""
}

func (s *Slot) AsInt() (int, bool) {
	if s == nil || s.Type != "YANDEX.NUMBER" {
		return 0, false
	}
	value := s.Value.(float64)
	return int(value), true
}
