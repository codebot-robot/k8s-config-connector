package v1alpha1

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refs "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ refs.Ref = &NetworkSecurityTLSInspectionPolicyRef{}

// NetworkSecurityTLSInspectionPolicyRef represents a reference to a NetworkSecurityTLSInspectionPolicy.
type NetworkSecurityTLSInspectionPolicyRef struct {
	// A reference to an externally managed NetworkSecurityTLSInspectionPolicy resource.
	// Should be in the format "projects/{projectID}/locations/{location}/tlsInspectionPolicies/{tls_inspection_policy}".
	External string `json:"external,omitempty"`

	// The name of a NetworkSecurityTLSInspectionPolicy resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkSecurityTLSInspectionPolicy resource.
	Namespace string `json:"namespace,omitempty"`
}

func init() {
	refs.Register(&NetworkSecurityTLSInspectionPolicyRef{})
}

func (r *NetworkSecurityTLSInspectionPolicyRef) GetGVK() schema.GroupVersionKind {
	return NetworkSecurityTLSInspectionPolicyGVK
}

func (r *NetworkSecurityTLSInspectionPolicyRef) GetNamespacedName() client.ObjectKey {
	return client.ObjectKey{Name: r.Name, Namespace: r.Namespace}
}

func (r *NetworkSecurityTLSInspectionPolicyRef) GetExternal() string {
	return r.External
}

func (r *NetworkSecurityTLSInspectionPolicyRef) SetExternal(external string) {
	r.External = external
}

func (r *NetworkSecurityTLSInspectionPolicyRef) ValidateExternal(external string) error {
	id := &NetworkSecurityTLSInspectionPolicyIdentity{}
	return id.FromExternal(external)
}

func (r *NetworkSecurityTLSInspectionPolicyRef) ParseExternalToIdentity() (identity.Identity, error) {
	id := &NetworkSecurityTLSInspectionPolicyIdentity{}
	err := id.FromExternal(r.External)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *NetworkSecurityTLSInspectionPolicyRef) Normalize(ctx context.Context, reader client.Reader, namespace string) error {
	fallback := func(u *unstructured.Unstructured) string {
		id, err := getIdentityFromNetworkSecurityTLSInspectionPolicySpec(ctx, reader, u)
		if err != nil {
			return ""
		}
		return id.String()
	}
	return refs.NormalizeWithFallback(ctx, reader, r, namespace, fallback)
}
