package api

type Intents struct {
	ListNamespaces *EmptyObj          `json:"list_namespaces"`
	CountPods      *IntentCountPods   `json:"count_pods"`
	BrokenPods     *IntentBrokenPods  `json:"broken_pods"`
	ServiceList    *IntentServiceList `json:"service_list"`

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

type NamespaceSlots struct {
	Namespace *Slot `json:"namespace"`
}
