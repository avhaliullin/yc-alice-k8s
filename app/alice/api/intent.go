package api

type Intents struct {
	ListNamespaces    *EmptyObj           `json:"list_namespaces,omitempty"`
	CountPods         *IntentCountPods    `json:"count_pods,omitempty"`
	BrokenPods        *IntentBrokenPods   `json:"broken_pods,omitempty"`
	ServiceList       *IntentServiceList  `json:"service_list,omitempty"`
	IngressList       *IntentIngressList  `json:"ingress_list,omitempty"`
	DiscoverScenarios *EmptyObj           `json:"discover_scenarios,omitempty"`
	DeployStatus      *IntentDeployStatus `json:"deploy_status,omitempty"`
	Deploy            *IntentDeploy       `json:"deploy,omitempty"`
	ScaleDeploy       *IntentScaleDeploy  `json:"scale_deploy,omitempty"`
	DeleteDeploy      *IntentDeleteDeploy `json:"delete_deploy,omitempty"`

	EasterDBLaunch  *EmptyObj `json:"easter_db_launch,omitempty"`
	EasterWhatIsK8s *EmptyObj `json:"easter_what_is_k8s,omitempty"`
	EasterHowTo     *EmptyObj `json:"easter_how_to,omitempty"`
	LetsPlayK8S     *EmptyObj `json:"lets_play_k8s,omitempty"`

	Confirm *EmptyObj `json:"YANDEX.CONFIRM,omitempty"`
	Cancel  *EmptyObj `json:"cancel,omitempty"`
	Reject  *EmptyObj `json:"YANDEX.REJECT,omitempty"`
}

type IntentCountPods struct {
	Slots NamespaceSlots `json:"slots,omitempty"`
}

type IntentBrokenPods struct {
	Slots NamespaceSlots `json:"slots,omitempty"`
}

type IntentServiceList struct {
	Slots NamespaceSlots `json:"slots,omitempty"`
}

type IntentIngressList struct {
	Slots NamespaceSlots `json:"slots,omitempty"`
}

type IntentDeploy struct {
	Slots IntentDeploySlots `json:"slots,omitempty"`
}

type IntentDeploySlots struct {
	Image *Slot `json:"image,omitempty"`
	Scale *Slot `json:"scale,omitempty"`
	Name  *Slot `json:"name,omitempty"`
}

type IntentDeployStatus struct {
	Slots DeployInNamespaceSlots `json:"slots,omitempty"`
}

type IntentScaleDeploy struct {
	Slots IntentScaleDeploySlots `json:"slots,omitempty"`
}

type IntentScaleDeploySlots struct {
	Name  *Slot `json:"name,omitempty"`
	Scale *Slot `json:"scale,omitempty"`
}

type IntentDeleteDeploy struct {
	Slots DeployInDefaultNSSLots `json:"slots,omitempty"`
}

type DeployInDefaultNSSLots struct {
	Name *Slot `json:"name,omitempty"`
}

type DeployInNamespaceSlots struct {
	Name      *Slot `json:"name,omitempty"`
	Namespace *Slot `json:"namespace,omitempty"`
}

type NamespaceSlots struct {
	Namespace *Slot `json:"namespace,omitempty"`
}
