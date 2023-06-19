package api

type State string

const (
	StateInit           State = ""
	StateDeployReqName  State = "DPLY_REQ_NAME"
	StateDeployReqImage State = "DPLY_REQ_IMAGE"
	StateDeployConfirm  State = "DPLY_CNFRM"
)

type StateData struct {
	State     State
	ImageText string
	Image     string
	Scale     int
	Name      string
}

func (s *StateData) GetState() State {
	if s == nil {
		return StateInit
	}
	return s.State
}
