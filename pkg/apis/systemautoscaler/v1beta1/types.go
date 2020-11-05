package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: decide whether to keep these or use `register.go`
// const (
// 	GroupName string = "PodScale.polimi.it"
// 	Kind      string = "PodScale"
// 	Version   string = "v1beta1"
// 	Plural    string = "PodScales"
// 	Singluar  string = "PodScale"
// 	ShortName string = "sysaut"
// 	Name      string = Plural + "." + GroupName
// )

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceLevelAgreementList is a list of ServiceLevelAgreement resources
type ServiceLevelAgreementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ServiceLevelAgreement `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceLevelAgreement is a configuration for the autoscaling system.
// It sets a requirement on the services that matches the selector.
type ServiceLevelAgreement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceLevelAgreementSpec `json:"spec"`
}

// ServiceLevelAgreementSpec defines the agreement specifying the
// metric requirement to honor by System Autoscaler, a Selector used
// to match a service with the Service Level Agreement and the
// default resources assigned to pods in case the `requests` field
// is empty in the `PodSpec`.
type ServiceLevelAgreementSpec struct {
	// Specify the metric on which the requirement is set.
	Metric MetricRequirement `json:"metric"`
	// Specify the default resources assigned to pods in case `requests` field is empty in `PodSpec`.
	DefaultResources v1.ResourceList `json:"defaultResources,omitempty" protobuf:"bytes,3,rep,name=defaultResources,casttype=ResourceList,castkey=ResourceName"`
	// Specify the selector to match Services and Service Level Agreement
	Selector *metav1.LabelSelector `json:"selector"`
}

// MetricRequirement specifies a requirement for a metric.
// This means that System Autoscaler will try to honor the
// agreement, making the service metric coherent with it.
// Only one MetricRequirement per ServiceLevelAgreement resource
// must be set to avoid ambiguity.
// Currently supports only ResponseTime.
//
// i.e.: the metric type is the Response Time and the value
// is 4 units of time. This means that the system will try
// to keep the service response time below 4 on average.
type MetricRequirement struct {
	ResponseTime *int32 `json:"responseTime,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodScaleList is a list of PodScale resources
type PodScaleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PodScale `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodScale defines the mapping between a `ServiceLevelAgreement` and a
// `Pod` matching the selector. It also keeps track of the resource values
// computed by `Recommender` and adjusted by `Contention Manager`.
type PodScale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodScaleSpec   `json:"spec"`
	Status PodScaleStatus `json:"status"`
}

// PodScaleSpec is the spec for a PodScale resource
type PodScaleSpec struct {
	SLAAgreement     SLAReference    `json:"serviceLevelAgreement"`
	PodRef           PodReference    `json:"pod"`
	DesiredResources v1.ResourceList `json:"desired,omitempty" protobuf:"bytes,3,rep,name=desired,casttype=ResourceList,castkey=ResourceName"`
}

// SLAReference specify the Service Level Agreement associated with the PodScale resource.
type SLAReference struct {
	// TODO: Decide if use all the objectMeta or something lighter
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// PodReference specify the Pod associated with the PodScale resource.
type PodReference struct {
	// TODO: Decide if use all the objectMeta or something lighter
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// TODO: Decide if useful or not

// PodScaleStatus contains the resources patched by the
// `Contention Manager` according to the available node resources
// and other pods' SLA
type PodScaleStatus struct {
	ActualResources v1.ResourceList `json:"actual,omitempty" protobuf:"bytes,3,rep,name=actual,casttype=ResourceList,castkey=ResourceName"`
}
