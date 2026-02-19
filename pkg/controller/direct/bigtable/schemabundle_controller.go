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

package bigtable

import (
	"context"
	"fmt"

	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/bigtable/v1alpha1"
	bigtablev1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/bigtable/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/config"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/directbase"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/registry"

	gcp "cloud.google.com/go/bigtable"
	bigtablepb "cloud.google.com/go/bigtable/admin/apiv2/adminpb"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

func init() {
	registry.RegisterModel(krm.BigtableSchemaBundleGVK, NewBigtableSchemaBundleModel)
}

func NewBigtableSchemaBundleModel(ctx context.Context, config *config.ControllerConfig) (directbase.Model, error) {
	return &modelBigtableSchemaBundle{config: *config}, nil
}

var _ directbase.Model = &modelBigtableSchemaBundle{}

type modelBigtableSchemaBundle struct {
	config config.ControllerConfig
}

func (m *modelBigtableSchemaBundle) client(ctx context.Context, projectID, instanceID string) (*gcp.AdminClient, error) {
	opts, err := m.config.GRPCClientOptions()
	if err != nil {
		return nil, err
	}
	gcpClient, err := gcp.NewAdminClient(ctx, projectID, instanceID, opts...)
	if err != nil {
		return nil, fmt.Errorf("building Bigtable SchemaBundle client: %w", err)
	}
	return gcpClient, nil
}

func (m *modelBigtableSchemaBundle) AdapterForObject(ctx context.Context, op *directbase.AdapterForObjectOperation) (directbase.Adapter, error) {
	u := op.GetUnstructured()
	reader := op.Reader
	obj := &krm.BigtableSchemaBundle{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &obj); err != nil {
		return nil, fmt.Errorf("error converting to %T: %w", obj, err)
	}

	id, err := krm.NewSchemaBundleIdentity(ctx, reader, obj)
	if err != nil {
		return nil, err
	}

	// Get bigtable admin GCP client.
	tableRef := id.Parent()
	instanceRef := tableRef.Parent
	projectRef := instanceRef.Parent

	adminClient, err := m.client(ctx, projectRef.ProjectID, instanceRef.Id)
	if err != nil {
		return nil, fmt.Errorf("error creating admin client: %w", err)
	}
	return &BigtableSchemaBundleAdapter{
		id:        id,
		gcpClient: adminClient,
		desired:   obj,
	}, nil
}

func (m *modelBigtableSchemaBundle) AdapterForURL(ctx context.Context, url string) (directbase.Adapter, error) {
	// TODO: Support URLs
	return nil, nil
}

type BigtableSchemaBundleAdapter struct {
	id        *krm.SchemaBundleIdentity
	gcpClient *gcp.AdminClient
	desired   *krm.BigtableSchemaBundle
	actual    *bigtablepb.SchemaBundle
}

var _ directbase.Adapter = &BigtableSchemaBundleAdapter{}

// Find retrieves the GCP resource.
func (a *BigtableSchemaBundleAdapter) Find(ctx context.Context) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("getting BigtableSchemaBundle", "name", a.id)

	info, err := a.gcpClient.GetSchemaBundle(ctx, a.id.Parent().ID(), a.id.ID())
	if err != nil {
		if direct.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("getting BigtableSchemaBundle %q: %w", a.id, err)
	}

	bundle := &bigtablepb.SchemaBundle{}
	if err := proto.Unmarshal(info.SchemaBundle, bundle); err != nil {
		return false, fmt.Errorf("unmarshalling SchemaBundle from info: %w", err)
	}

	a.actual = bundle
	return true, nil
}

// Create creates the resource in GCP based on `spec`.
func (a *BigtableSchemaBundleAdapter) Create(ctx context.Context, createOp *directbase.CreateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("creating BigtableSchemaBundle", "name", a.id)

	conf := &gcp.SchemaBundleConf{
		TableID:        a.id.Parent().ID(),
		SchemaBundleID: a.id.ID(),
		ProtoSchema:    &gcp.ProtoSchemaInfo{},
	}
	if a.desired.Spec.ProtoSchema != nil {
		conf.ProtoSchema.ProtoDescriptors = a.desired.Spec.ProtoSchema.ProtoDescriptors
	}

	if err := a.gcpClient.CreateSchemaBundle(ctx, conf); err != nil {
		return fmt.Errorf("creating BigtableSchemaBundle %s: %w", a.id, err)
	}
	log.V(2).Info("successfully created BigtableSchemaBundle", "name", a.id)

	status := &krm.BigtableSchemaBundleStatus{}
	status.ExternalRef = direct.LazyPtr(a.id.String())
	return createOp.UpdateStatus(ctx, status, nil)
}

// Update updates the resource in GCP based on `spec`.
func (a *BigtableSchemaBundleAdapter) Update(ctx context.Context, updateOp *directbase.UpdateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("updating BigtableSchemaBundle", "name", a.id)

	conf := gcp.UpdateSchemaBundleConf{
		SchemaBundleConf: gcp.SchemaBundleConf{
			TableID:        a.id.Parent().ID(),
			SchemaBundleID: a.id.ID(),
			ProtoSchema:    &gcp.ProtoSchemaInfo{},
		},
	}
	if a.desired.Spec.ProtoSchema != nil {
		conf.SchemaBundleConf.ProtoSchema.ProtoDescriptors = a.desired.Spec.ProtoSchema.ProtoDescriptors
	}

	if err := a.gcpClient.UpdateSchemaBundle(ctx, conf); err != nil {
		return fmt.Errorf("updating BigtableSchemaBundle %s: %w", a.id, err)
	}
	log.V(2).Info("successfully updated BigtableSchemaBundle", "name", a.id)

	status := &krm.BigtableSchemaBundleStatus{}
	status.ExternalRef = direct.LazyPtr(a.id.String())
	return updateOp.UpdateStatus(ctx, status, nil)
}

// Export maps the GCP object to a Config Connector resource `spec`.
func (a *BigtableSchemaBundleAdapter) Export(ctx context.Context) (*unstructured.Unstructured, error) {
	if a.actual == nil {
		return nil, fmt.Errorf("Find() not called")
	}
	u := &unstructured.Unstructured{}

	obj := &krm.BigtableSchemaBundle{}
	mapCtx := &direct.MapContext{}
	obj.Spec = direct.ValueOf(BigtableSchemaBundleSpec_FromProto(mapCtx, a.actual))
	if mapCtx.Err() != nil {
		return nil, mapCtx.Err()
	}

	// Set Parent Ref
	// obj.Spec.TableRef = ...
	// Since SchemaBundle is child of Table, we need to reconstruct TableRef from ID.
	// But the ID is already known.

	// a.id is SchemaBundleIdentity. a.id.Parent() is TableIdentity.
	tableID := a.id.Parent()
	obj.Spec.TableRef = krm.BigtableSchemaBundleParent{
		TableRef: bigtablev1beta1.TableRef{
			Name: tableID.ID(), // This assumes referencing by name in same namespace?
			// Or External?
			External: tableID.String(),
		},
	}.TableRef

	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	u.SetName(a.id.ID())
	u.SetGroupVersionKind(krm.BigtableSchemaBundleGVK)

	u.Object = uObj
	return u, nil
}

// Delete the resource from GCP service.
func (a *BigtableSchemaBundleAdapter) Delete(ctx context.Context, deleteOp *directbase.DeleteOperation) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("deleting BigtableSchemaBundle", "name", a.id)

	err := a.gcpClient.DeleteSchemaBundle(ctx, a.id.Parent().ID(), a.id.ID())
	if err != nil {
		if direct.IsNotFound(err) {
			log.V(2).Info("skipping delete for non-existent BigtableSchemaBundle", "name", a.id)
			return true, nil
		}
		return false, fmt.Errorf("deleting BigtableSchemaBundle %s: %w", a.id, err)
	}
	log.V(2).Info("successfully deleted BigtableSchemaBundle", "name", a.id)

	return true, nil
}
