// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"context"
	"fmt"

	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ refsv1beta1.ExternalNormalizer = &EdgeCacheServiceRef{}

// EdgeCacheServiceRef defines the resource reference to NetworkServicesEdgeCacheService, which "External" field
// holds the GCP identifier for the KRM object.
type EdgeCacheServiceRef struct {
	// A reference to an externally managed NetworkServicesEdgeCacheService resource.
	// Should be in the format "projects/{{projectID}}/locations/global/edgeCacheServices/{{edgeCacheServiceID}}".
	External string `json:"external,omitempty"`

	// The name of a NetworkServicesEdgeCacheService resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkServicesEdgeCacheService resource.
	Namespace string `json:"namespace,omitempty"`
}

// NormalizedExternal provision the "External" value for other resource that depends on NetworkServicesEdgeCacheService.
// If the "External" is given in the other resource's spec.NetworkServicesEdgeCacheServiceRef, the given value will be used.
// Otherwise, the "Name" and "Namespace" will be used to query the actual NetworkServicesEdgeCacheService object from the cluster.
func (r *EdgeCacheServiceRef) NormalizedExternal(ctx context.Context, reader client.Reader, otherNamespace string) (string, error) {
	if r.External != "" && r.Name != "" {
		return "", fmt.Errorf("cannot specify both name and external on %s reference", NetworkServicesEdgeCacheServiceGVK.Kind)
	}
	// From given External
	if r.External != "" {
		id := &EdgeCacheServiceIdentity{}
		if err := id.FromExternal(r.External); err != nil {
			return "", err
		}
		return r.External, nil
	}

	// From the Config Connector object
	if r.Namespace == "" {
		r.Namespace = otherNamespace
	}
	key := types.NamespacedName{Name: r.Name, Namespace: r.Namespace}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(NetworkServicesEdgeCacheServiceGVK)
	if err := reader.Get(ctx, key, u); err != nil {
		if apierrors.IsNotFound(err) {
			return "", k8s.NewReferenceNotFoundError(u.GroupVersionKind(), key)
		}
		return "", fmt.Errorf("reading referenced %s %s: %w", NetworkServicesEdgeCacheServiceGVK, key, err)
	}
	// Get external from status.externalRef. This is the most trustworthy place.
	actualExternalRef, _, err := unstructured.NestedString(u.Object, "status", "externalRef")
	if err != nil {
		return "", fmt.Errorf("reading status.externalRef: %w", err)
	}
	if actualExternalRef == "" {
		return "", k8s.NewReferenceNotReadyError(u.GroupVersionKind(), key)
	}
	r.External = actualExternalRef
	return r.External, nil
}

var _ refsv1beta1.ExternalNormalizer = &EdgeCacheKeysetRef{}

// EdgeCacheKeysetRef defines the resource reference to NetworkServicesEdgeCacheKeyset, which "External" field
// holds the GCP identifier for the KRM object.
type EdgeCacheKeysetRef struct {
	// A reference to an externally managed NetworkServicesEdgeCacheKeyset resource.
	// Should be in the format "projects/{{projectID}}/locations/global/edgeCacheKeysets/{{edgeCacheKeysetID}}".
	External string `json:"external,omitempty"`

	// The name of a NetworkServicesEdgeCacheKeyset resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkServicesEdgeCacheKeyset resource.
	Namespace string `json:"namespace,omitempty"`
}

// NormalizedExternal provision the "External" value for other resource that depends on NetworkServicesEdgeCacheKeyset.
func (r *EdgeCacheKeysetRef) NormalizedExternal(ctx context.Context, reader client.Reader, otherNamespace string) (string, error) {
	if r.External != "" && r.Name != "" {
		return "", fmt.Errorf("cannot specify both name and external on EdgeCacheKeyset reference")
	}
	// From given External
	if r.External != "" {
		return r.External, nil
	}

	// From the Config Connector object
	if r.Namespace == "" {
		r.Namespace = otherNamespace
	}
	key := types.NamespacedName{Name: r.Name, Namespace: r.Namespace}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networkservices.cnrm.cloud.google.com",
		Version: "v1alpha1",
		Kind:    "NetworkServicesEdgeCacheKeyset",
	})
	if err := reader.Get(ctx, key, u); err != nil {
		if apierrors.IsNotFound(err) {
			return "", k8s.NewReferenceNotFoundError(u.GroupVersionKind(), key)
		}
		return "", fmt.Errorf("reading referenced NetworkServicesEdgeCacheKeyset %s: %w", key, err)
	}
	// Get external from status.externalRef.
	actualExternalRef, _, err := unstructured.NestedString(u.Object, "status", "externalRef")
	if err != nil {
		return "", fmt.Errorf("reading status.externalRef: %w", err)
	}
	if actualExternalRef == "" {
		return "", k8s.NewReferenceNotReadyError(u.GroupVersionKind(), key)
	}
	r.External = actualExternalRef
	return r.External, nil
}

var _ refsv1beta1.ExternalNormalizer = &EdgeSecurityPolicyRef{}

type EdgeSecurityPolicyRef struct {
	// A reference to an externally managed ComputeSecurityPolicy resource.
	// Should be in the format "projects/{{projectID}}/global/securityPolicies/{{securityPolicyID}}".
	External string `json:"external,omitempty"`

	// The name of a ComputeSecurityPolicy resource.
	Name string `json:"name,omitempty"`

	// The namespace of a ComputeSecurityPolicy resource.
	Namespace string `json:"namespace,omitempty"`
}

func (r *EdgeSecurityPolicyRef) NormalizedExternal(ctx context.Context, reader client.Reader, otherNamespace string) (string, error) {
	if r.External != "" && r.Name != "" {
		return "", fmt.Errorf("cannot specify both name and external on EdgeSecurityPolicy reference")
	}
	if r.External != "" {
		return r.External, nil
	}
	if r.Namespace == "" {
		r.Namespace = otherNamespace
	}
	key := types.NamespacedName{Name: r.Name, Namespace: r.Namespace}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "compute.cnrm.cloud.google.com",
		Version: "v1beta1",
		Kind:    "ComputeSecurityPolicy",
	})
	if err := reader.Get(ctx, key, u); err != nil {
		if apierrors.IsNotFound(err) {
			return "", k8s.NewReferenceNotFoundError(u.GroupVersionKind(), key)
		}
		return "", fmt.Errorf("reading referenced ComputeSecurityPolicy %s: %w", key, err)
	}
	actualExternalRef, _, err := unstructured.NestedString(u.Object, "status", "externalRef")
	if err != nil {
		return "", fmt.Errorf("reading status.externalRef: %w", err)
	}
	if actualExternalRef == "" {
		// Fallback to selfLink for older TF resources
		actualExternalRef, _, err = unstructured.NestedString(u.Object, "status", "selfLink")
		if err != nil {
			return "", fmt.Errorf("reading status.selfLink: %w", err)
		}
	}
	if actualExternalRef == "" {
		return "", k8s.NewReferenceNotReadyError(u.GroupVersionKind(), key)
	}
	r.External = actualExternalRef
	return r.External, nil
}

var _ refsv1beta1.ExternalNormalizer = &EdgeSslCertificateRef{}

type EdgeSslCertificateRef struct {
	// A reference to an externally managed ComputeSSLCertificate resource.
	// Should be in the format "projects/{{projectID}}/global/sslCertificates/{{sslCertificateID}}".
	External string `json:"external,omitempty"`

	// The name of a ComputeSSLCertificate resource.
	Name string `json:"name,omitempty"`

	// The namespace of a ComputeSSLCertificate resource.
	Namespace string `json:"namespace,omitempty"`
}

func (r *EdgeSslCertificateRef) NormalizedExternal(ctx context.Context, reader client.Reader, otherNamespace string) (string, error) {
	if r.External != "" && r.Name != "" {
		return "", fmt.Errorf("cannot specify both name and external on EdgeSslCertificate reference")
	}
	if r.External != "" {
		return r.External, nil
	}
	if r.Namespace == "" {
		r.Namespace = otherNamespace
	}
	key := types.NamespacedName{Name: r.Name, Namespace: r.Namespace}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "compute.cnrm.cloud.google.com",
		Version: "v1beta1",
		Kind:    "ComputeSSLCertificate",
	})
	if err := reader.Get(ctx, key, u); err != nil {
		if apierrors.IsNotFound(err) {
			return "", k8s.NewReferenceNotFoundError(u.GroupVersionKind(), key)
		}
		return "", fmt.Errorf("reading referenced ComputeSSLCertificate %s: %w", key, err)
	}
	actualExternalRef, _, err := unstructured.NestedString(u.Object, "status", "externalRef")
	if err != nil {
		return "", fmt.Errorf("reading status.externalRef: %w", err)
	}
	if actualExternalRef == "" {
		actualExternalRef, _, err = unstructured.NestedString(u.Object, "status", "selfLink")
		if err != nil {
			return "", fmt.Errorf("reading status.selfLink: %w", err)
		}
	}
	if actualExternalRef == "" {
		return "", k8s.NewReferenceNotReadyError(u.GroupVersionKind(), key)
	}
	r.External = actualExternalRef
	return r.External, nil
}

var _ refsv1beta1.ExternalNormalizer = &SslPolicyRef{}

type SslPolicyRef struct {
	// A reference to an externally managed ComputeSSLPolicy resource.
	// Should be in the format "projects/{{projectID}}/global/sslPolicies/{{sslPolicyID}}".
	External string `json:"external,omitempty"`

	// The name of a ComputeSSLPolicy resource.
	Name string `json:"name,omitempty"`

	// The namespace of a ComputeSSLPolicy resource.
	Namespace string `json:"namespace,omitempty"`
}

func (r *SslPolicyRef) NormalizedExternal(ctx context.Context, reader client.Reader, otherNamespace string) (string, error) {
	if r.External != "" && r.Name != "" {
		return "", fmt.Errorf("cannot specify both name and external on SslPolicy reference")
	}
	if r.External != "" {
		return r.External, nil
	}
	if r.Namespace == "" {
		r.Namespace = otherNamespace
	}
	key := types.NamespacedName{Name: r.Name, Namespace: r.Namespace}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "compute.cnrm.cloud.google.com",
		Version: "v1beta1",
		Kind:    "ComputeSSLPolicy",
	})
	if err := reader.Get(ctx, key, u); err != nil {
		if apierrors.IsNotFound(err) {
			return "", k8s.NewReferenceNotFoundError(u.GroupVersionKind(), key)
		}
		return "", fmt.Errorf("reading referenced ComputeSSLPolicy %s: %w", key, err)
	}
	actualExternalRef, _, err := unstructured.NestedString(u.Object, "status", "externalRef")
	if err != nil {
		return "", fmt.Errorf("reading status.externalRef: %w", err)
	}
	if actualExternalRef == "" {
		actualExternalRef, _, err = unstructured.NestedString(u.Object, "status", "selfLink")
		if err != nil {
			return "", fmt.Errorf("reading status.selfLink: %w", err)
		}
	}
	if actualExternalRef == "" {
		return "", k8s.NewReferenceNotReadyError(u.GroupVersionKind(), key)
	}
	r.External = actualExternalRef
	return r.External, nil
}

var _ refsv1beta1.ExternalNormalizer = &EdgeCacheOriginRef{}

type EdgeCacheOriginRef struct {
	// A reference to an externally managed NetworkServicesEdgeCacheOrigin resource.
	// Should be in the format "projects/{{projectID}}/locations/global/edgeCacheOrigins/{{edgeCacheOriginID}}".
	External string `json:"external,omitempty"`

	// The name of a NetworkServicesEdgeCacheOrigin resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkServicesEdgeCacheOrigin resource.
	Namespace string `json:"namespace,omitempty"`
}

func (r *EdgeCacheOriginRef) NormalizedExternal(ctx context.Context, reader client.Reader, otherNamespace string) (string, error) {
	if r.External != "" && r.Name != "" {
		return "", fmt.Errorf("cannot specify both name and external on EdgeCacheOrigin reference")
	}
	if r.External != "" {
		return r.External, nil
	}
	if r.Namespace == "" {
		r.Namespace = otherNamespace
	}
	key := types.NamespacedName{Name: r.Name, Namespace: r.Namespace}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networkservices.cnrm.cloud.google.com",
		Version: "v1alpha1",
		Kind:    "NetworkServicesEdgeCacheOrigin",
	})
	if err := reader.Get(ctx, key, u); err != nil {
		if apierrors.IsNotFound(err) {
			return "", k8s.NewReferenceNotFoundError(u.GroupVersionKind(), key)
		}
		return "", fmt.Errorf("reading referenced NetworkServicesEdgeCacheOrigin %s: %w", key, err)
	}
	actualExternalRef, _, err := unstructured.NestedString(u.Object, "status", "externalRef")
	if err != nil {
		return "", fmt.Errorf("reading status.externalRef: %w", err)
	}
	if actualExternalRef == "" {
		return "", k8s.NewReferenceNotReadyError(u.GroupVersionKind(), key)
	}
	r.External = actualExternalRef
	return r.External, nil
}
