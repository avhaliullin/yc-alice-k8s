package api

type State string

const (
	StateInit           State = ""
	StateDeployReqName  State = "DPLY_REQ_NAME"
	StateDeployReqImage State = "DPLY_REQ_IMAGE"
	StateDeployConfirm  State = "DPLY_CNFRM"

	StateDeployStatusReqName      State = "DPLY_ST_REQ_NAME"
	StateDeployStatusReqNamespace State = "DPLY_ST_REQ_NMSPC"

	StateScaleDeployReqName    State = "SCL_DPLY_REQ_NAME"
	StateScaleDeployReqScale   State = "SCL_DPLY_REQ_SCALE"
	StateScaleDeployReqConfirm State = "SCL_DPLY_REQ_CNFRM"

	StateDeleteDeployReqName    State = "DEL_DPLY_REQ_NAME"
	StateDeleteDeployReqConfirm State = "DEL_DPLY_REQ_CNFRM"
)

type StateData struct {
	State       State
	Image       string
	ImageID     string
	Scale       int
	DeployName  string
	DeployID    string
	Namespace   string
	NamespaceID string
}

func (s *StateData) GetState() State {
	if s == nil {
		return StateInit
	}
	return s.State
}
