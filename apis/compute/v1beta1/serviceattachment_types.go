// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ComputeServiceAttachmentGVK = GroupVersion.WithKind("ComputeServiceAttachment")

type ServiceattachmentConsumerAcceptLists struct {
	/* The value of the limit to set. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachmentConsumerProjectLimit.connection_limit
	ConnectionLimit *int64 `json:"connectionLimit,omitempty"`

	// +required
	ProjectRef refsv1beta1.ProjectRef `json:"projectRef"`
}

// ComputeServiceAttachmentSpec defines the desired state of ComputeServiceAttachment
// +kcc:spec:proto=google.cloud.compute.v1.ServiceAttachment
type ComputeServiceAttachmentSpec struct {
	// The project that this resource belongs to.
	// +required
	ProjectRef *refsv1beta1.ProjectRef `json:"projectRef"`

	// Immutable. The location of this resource.
	// +required
	Location string `json:"location"`

	// The ComputeServiceAttachment name. If not given, the metadata.name will be used.
	ResourceID *string `json:"resourceID,omitempty"`

	/* The connection preference of service attachment. The value can be set to `ACCEPT_AUTOMATIC`. An `ACCEPT_AUTOMATIC` service attachment is one that always accepts the connection from consumer forwarding rules. Possible values: CONNECTION_PREFERENCE_UNSPECIFIED, ACCEPT_AUTOMATIC, ACCEPT_MANUAL */
	// +required
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.connection_preference
	ConnectionPreference string `json:"connectionPreference"`

	/* Projects that are allowed to connect to this service attachment. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.consumer_accept_lists
	ConsumerAcceptLists []ServiceattachmentConsumerAcceptLists `json:"consumerAcceptLists,omitempty"`

	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.consumer_reject_lists
	ConsumerRejectLists []refsv1beta1.ProjectRef `json:"consumerRejectLists,omitempty"`

	/* An optional description of this resource. Provide this property when you create the resource. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.description
	Description *string `json:"description,omitempty"`

	/* Immutable. If true, enable the proxy protocol which is for supplying client TCP/IP address data in TCP connections that traverse proxies on their way to destination servers. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.enable_proxy_protocol
	EnableProxyProtocol *bool `json:"enableProxyProtocol,omitempty"`

	// +required
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.nat_subnets
	NatSubnets []ComputeSubnetworkRef `json:"natSubnets"`

	/* Immutable. */
	// +required
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.producer_forwarding_rule
	TargetServiceRef *refsv1beta1.ComputeForwardingRuleRef `json:"targetServiceRef"`
}

type ServiceattachmentConnectedEndpointsStatus struct {
	/* The url of a connected endpoint. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachmentConnectedEndpoint.endpoint
	Endpoint *string `json:"endpoint,omitempty"`

	/* The PSC connection id of the connected endpoint. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachmentConnectedEndpoint.psc_connection_id
	PscConnectionId *int64 `json:"pscConnectionId,omitempty"`

	/* The status of a connected endpoint to this service attachment. Possible values: PENDING, RUNNING, DONE */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachmentConnectedEndpoint.status
	Status *string `json:"status,omitempty"`
}

type ServiceattachmentPscServiceAttachmentIdStatus struct {
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.Uint128.high
	High *int64 `json:"high,omitempty"`

	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.Uint128.low
	Low *int64 `json:"low,omitempty"`
}

// ComputeServiceAttachmentStatus defines the config connector machine state of ComputeServiceAttachment
// +kcc:status:proto=google.cloud.compute.v1.ServiceAttachment
type ComputeServiceAttachmentStatus struct {
	/* Conditions represent the latest available observations of the
	   object's current state. */
	Conditions []v1alpha1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the generation of the resource that was most recently observed by the Config Connector controller. If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource.
	ObservedGeneration *int64 `json:"observedGeneration,omitempty"`

	// A unique specifier for the ComputeServiceAttachment resource in GCP.
	ExternalRef *string `json:"externalRef,omitempty"`

	/* An array of connections for all the consumers connected to this service attachment. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.connected_endpoints
	ConnectedEndpoints []ServiceattachmentConnectedEndpointsStatus `json:"connectedEndpoints,omitempty"`

	/* Fingerprint of this resource. This field is used internally during updates of this resource. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.fingerprint
	Fingerprint *string `json:"fingerprint,omitempty"`

	/* The unique identifier for the resource type. The server generates this identifier. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.id
	Id *int64 `json:"id,omitempty"`

	/* An 128-bit global unique ID of the PSC service attachment. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.psc_service_attachment_id
	PscServiceAttachmentId *ServiceattachmentPscServiceAttachmentIdStatus `json:"pscServiceAttachmentId,omitempty"`

	/* URL of the region where the service attachment resides. This field applies only to the region resource. You must specify this field as part of the HTTP request URL. It is not settable as a field in the request body. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.region
	Region *string `json:"region,omitempty"`

	/* Server-defined URL for the resource. */
	// +optional
	// +kcc:proto:field=google.cloud.compute.v1.ServiceAttachment.self_link
	SelfLink *string `json:"selfLink,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories=gcp,shortName=gcpcomputeserviceattachment;gcpcomputeserviceattachments
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/dcl2crd=true"
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/managed-by-kcc=true"
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/stability-level=stable"
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/system=true"
// +kubebuilder:printcolumn:name="Age",JSONPath=".metadata.creationTimestamp",type="date"
// +kubebuilder:printcolumn:name="Ready",JSONPath=".status.conditions[?(@.type=='Ready')].status",type="string",description="When 'True', the most recent reconcile of the resource succeeded"
// +kubebuilder:printcolumn:name="Status",JSONPath=".status.conditions[?(@.type=='Ready')].reason",type="string",description="The reason for the value in 'Ready'"
// +kubebuilder:printcolumn:name="Status Age",JSONPath=".status.conditions[?(@.type=='Ready')].lastTransitionTime",type="date",description="The last transition time for the value in 'Status'"

// ComputeServiceAttachment is the Schema for the ComputeServiceAttachment API
// +k8s:openapi-gen=true
type ComputeServiceAttachment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec   ComputeServiceAttachmentSpec   `json:"spec,omitempty"`
	Status ComputeServiceAttachmentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ComputeServiceAttachmentList contains a list of ComputeServiceAttachment
type ComputeServiceAttachmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComputeServiceAttachment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComputeServiceAttachment{}, &ComputeServiceAttachmentList{})
}
