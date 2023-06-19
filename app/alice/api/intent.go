package api

type Intents struct {
	ListNamespaces    *EmptyObj           `json:"list_namespaces"`
	CountPods         *IntentCountPods    `json:"count_pods"`
	BrokenPods        *IntentBrokenPods   `json:"broken_pods"`
	ServiceList       *IntentServiceList  `json:"service_list"`
	IngressList       *IntentIngressList  `json:"ingress_list"`
	DiscoverScenarios *EmptyObj           `json:"discover_scenarios"`
	DeployStatus      *IntentDeployStatus `json:"deploy_status"`
	Deploy            *IntentDeploy       `json:"deploy"`
	ScaleDeploy       *IntentScaleDeploy  `json:"scale_deploy"`
	DeleteDeploy      *IntentDeleteDeploy `json:"delete_deploy"`

	EasterDBLaunch  *EmptyObj `json:"easter_db_launch"`
	EasterWhatIsK8s *EmptyObj `json:"easter_what_is_k8s"`
	EasterHowTo     *EmptyObj `json:"easter_how_to"`

	Confirm *EmptyObj `json:"YANDEX.CONFIRM"`
	Cancel  *EmptyObj `json:"cancel"`
	Reject  *EmptyObj `json:"YANDEX.REJECT"`
}

type IntentCountPods struct {
	Slots NamespaceSlots `json:"slots"`
}

type IntentBrokenPods struct {
	Slots NamespaceSlots `json:"slots"`
}

type IntentServiceList struct {
	Slots NamespaceSlots `json:"slots"`
}

type IntentIngressList struct {
	Slots NamespaceSlots `json:"slots"`
}

type IntentDeploy struct {
	Slots IntentDeploySlots `json:"slots"`
}

type IntentDeploySlots struct {
	Image *Slot `json:"image"`
	Scale *Slot `json:"scale"`
	Name  *Slot `json:"name"`
}

type IntentDeployStatus struct {
	Slots DeployInNamespaceSlots `json:"slots"`
}

type IntentScaleDeploy struct {
	Slots IntentScaleDeploySlots `json:"slots"`
}

type IntentScaleDeploySlots struct {
	Name  *Slot `json:"name"`
	Scale *Slot `json:"scale"`
}

type IntentDeleteDeploy struct {
	Slots DeployInDefaultNSSLots `json:"slots"`
}

type DeployInDefaultNSSLots struct {
	Name *Slot `json:"name"`
}

type DeployInNamespaceSlots struct {
	Name      *Slot `json:"name"`
	Namespace *Slot `json:"namespace"`
}

type NamespaceSlots struct {
	Namespace *Slot `json:"namespace"`
}
