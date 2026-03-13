// Copyright 2024 Google LLC
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

package mockcompute

import (
	"context"
	"fmt"
	"time"

	pb "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/compute/v1"
	pbv1beta "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/compute/v1beta"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/klog/v2"

	"github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/pkg/storage"
)

type computeOperations struct {
	storage storage.Storage
}

func newComputeOperationsService(storage storage.Storage) *computeOperations {
	return &computeOperations{
		storage: storage,
	}
}

func (s *computeOperations) globalOperationFQN(projectID string, name string) string {
	return "projects/" + projectID + "/global/operations/" + name
}

func (s *computeOperations) globalOrganizationOperationFQN(name string) string {
	return "locations/global/operations/" + name
}

func (s *computeOperations) regionalOperationFQN(projectID string, region string, name string) string {
	return "projects/" + projectID + "/regions/" + region + "/operations/" + name
}

func (s *computeOperations) zonalOperationFQN(projectID string, zone string, name string) string {
	return "projects/" + projectID + "/zones/" + zone + "/operations/" + name
}

// Deprecated: use startGlobalLRO
func (s *computeOperations) newLRO(ctx context.Context, projectID string) (*pb.Operation, error) {
	log := klog.FromContext(ctx)

	now := time.Now()
	millis := now.UnixMilli()
	id := string(uuid.NewUUID())
	nanos := now.UnixNano()

	op := &pb.Operation{}

	op.StartTime = PtrTo(formatTime(now))
	op.InsertTime = PtrTo(formatTime(now))
	op.Id = PtrTo(uint64(nanos))

	op.Progress = PtrTo(int32(0))

	name := fmt.Sprintf("operation-%d-%s", millis, id)
	op.Name = PtrTo(name)
	op.Kind = PtrTo("compute#operation")
	fqn := s.globalOperationFQN(projectID, name)

	op.SelfLink = PtrTo(buildComputeSelfLink(ctx, fqn))

	op.Status = PtrTo(pb.Operation_DONE)

	log.Info("storing operation", "fqn", fqn)
	if err := s.storage.Create(ctx, fqn, op); err != nil {
		return nil, err
	}
	return op, nil
}

func (s *computeOperations) startLRO0(ctx context.Context, op *pb.Operation, fqn string, callback func() (proto.Message, error)) (*pb.Operation, error) {
	log := klog.FromContext(ctx)

	now := time.Now()
	nanos := now.UnixNano()

	if op == nil {
		op = &pb.Operation{}
	}

	op.StartTime = PtrTo(formatTime(now))
	op.InsertTime = PtrTo(formatTime(now))
	op.Id = PtrTo(uint64(nanos))

	// Specific to ComputeFirewallPolicy
	// Remove targetId and targetLink when status is RUNNING to match realGCP operation
	// ref: https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/2800/commits/32fdacd53d59c36626fce16f2b0125a8a455f3d6#r1783642429
	targetId := op.TargetId
	targetLink := op.TargetLink
	if op.OperationType != nil && *op.OperationType == "createFirewallPolicy" {
		op.TargetId = nil
		op.TargetLink = nil
	}

	if op.Progress == nil {
		op.Progress = PtrTo(int32(0))
	}

	if op.Status == nil {
		op.Status = PtrTo(pb.Operation_RUNNING)
	}

	op.Kind = PtrTo("compute#operation")
	op.SelfLink = PtrTo(buildComputeSelfLink(ctx, fqn))

	log.Info("storing operation", "fqn", fqn)
	if err := s.storage.Create(ctx, fqn, op); err != nil {
		return nil, err
	}

	go func() {
		result, err := callback()
		finished := &pb.Operation{}
		if err2 := s.storage.Get(ctx, fqn, finished); err2 != nil {
			klog.Warningf("error getting LRO: %v", err2)
			return
		}

		finished.Progress = PtrTo(int32(100))
		finished.Status = PtrTo(pb.Operation_DONE)
		finished.EndTime = PtrTo(formatTime(time.Now()))

		// Specific to ComputeFirewallPolicy
		// Add targetId and targetLink back when status is DONE to match realGCP operation
		if op.OperationType != nil && *op.OperationType == "createFirewallPolicy" {
			finished.TargetId = targetId
			finished.TargetLink = targetLink
		}

		if err != nil {
			code := status.Code(err)
			message := err.Error()

			finished.Error = &pb.Error{
				Errors: []*pb.Errors{
					{
						Code:    PtrTo(code.String()),
						Message: PtrTo(message),
					},
				},
			}
			klog.Warningf("TODO: more fully handle LRO error %v", err)
		} else {
			// The LRO result does not appear to be returned in the operation
			klog.V(4).Infof("LRO result: %+v", result)
		}
		if err := s.storage.Update(ctx, fqn, finished); err != nil {
			klog.Warningf("error updating LRO: %v", err)
			return
		}
	}()

	return op, nil
}

func (s *computeOperations) startRegionalLRO(ctx context.Context, projectID string, region string, op *pb.Operation, callback func() (proto.Message, error)) (*pb.Operation, error) {
	now := time.Now()
	millis := now.UnixMilli()
	id := uuid.NewUUID()

	name := fmt.Sprintf("operation-%d-%s", millis, id)
	fqn := s.regionalOperationFQN(projectID, region, name)

	op.Name = PtrTo(name)
	op.Region = PtrTo(buildComputeSelfLink(ctx, "projects/"+projectID+"/regions/"+region))
	return s.startLRO0(ctx, op, fqn, callback)
}

func (s *computeOperations) startGlobalLRO(ctx context.Context, projectID string, op *pb.Operation, callback func() (proto.Message, error)) (*pb.Operation, error) {
	now := time.Now()
	millis := now.UnixMilli()
	id := uuid.NewUUID()

	name := fmt.Sprintf("operation-%d-%s", millis, id)
	fqn := s.globalOperationFQN(projectID, name)

	op.Name = PtrTo(name)
	return s.startLRO0(ctx, op, fqn, callback)
}

func (s *computeOperations) startGlobalOrganizationLRO(ctx context.Context, op *pb.Operation, callback func() (proto.Message, error)) (*pb.Operation, error) {
	now := time.Now()
	millis := now.UnixMilli()
	id := uuid.NewUUID()

	name := fmt.Sprintf("operation-%d-%s", millis, id)
	fqn := s.globalOrganizationOperationFQN(name)

	op.Name = PtrTo(name)
	return s.startLRO0(ctx, op, fqn, callback)
}

// Gets the latest state of a long-running operation.  Clients can use this
// method to poll the operation result at intervals as recommended by the API
// service.
func (s *computeOperations) getOperation(ctx context.Context, fqn string) (*pb.Operation, error) {
	op := &pb.Operation{}
	if err := s.storage.Get(ctx, fqn, op); err != nil {
		return nil, err
	}

	return op, nil
}

func (s *computeOperations) getOperationV1Beta(ctx context.Context, fqn string) (*pbv1beta.Operation, error) {
	op := &pb.Operation{}
	if err := s.storage.Get(ctx, fqn, op); err != nil {
		return nil, err
	}

	return convertOperationV1ToV1Beta(op), nil
}

func (s *computeOperations) startZonalLRO(ctx context.Context, projectID string, zone string, op *pbv1beta.Operation, callback func() (proto.Message, error)) (*pbv1beta.Operation, error) {
	now := time.Now()
	millis := now.UnixMilli()
	id := uuid.NewUUID()

	name := fmt.Sprintf("operation-%d-%s", millis, id)
	fqn := s.zonalOperationFQN(projectID, zone, name)

	op.Name = PtrTo(name)
	op.Zone = PtrTo(buildComputeSelfLink(ctx, "projects/"+projectID+"/zones/"+zone))
	return s.startZonalLRO0(ctx, op, fqn, callback)
}

func (s *computeOperations) startZonalLRO0(ctx context.Context, opv1beta *pbv1beta.Operation, fqn string, callback func() (proto.Message, error)) (*pbv1beta.Operation, error) {
	op := convertOperationV1BetaToV1(opv1beta)

	res, err := s.startLRO0(ctx, op, fqn, callback)
	if err != nil {
		return nil, err
	}
	return convertOperationV1ToV1Beta(res), nil
}

func convertOperationV1ToV1Beta(in *pb.Operation) *pbv1beta.Operation {
	if in == nil {
		return nil
	}
	out := &pbv1beta.Operation{}
	// Use protojson to convert between types with same JSON tags
	b, err := protojson.Marshal(in)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal operation: %v", err))
	}
	if err := protojson.Unmarshal(b, out); err != nil {
		panic(fmt.Sprintf("failed to unmarshal operation: %v", err))
	}
	return out
}

func convertOperationV1BetaToV1(in *pbv1beta.Operation) *pb.Operation {
	if in == nil {
		return nil
	}
	out := &pb.Operation{}
	// Use protojson to convert between types with same JSON tags
	b, err := protojson.Marshal(in)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal operation: %v", err))
	}
	if err := protojson.Unmarshal(b, out); err != nil {
		panic(fmt.Sprintf("failed to unmarshal operation: %v", err))
	}
	return out
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
