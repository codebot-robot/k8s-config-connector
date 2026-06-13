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

// krm.group: compute.cnrm.cloud.google.com
// krm.version: v1beta1
// proto.service: google.cloud.compute.v1

package compute

import (
	pb "cloud.google.com/go/compute/apiv1/computepb"
	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/compute/v1beta1"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
)

func ComputeServiceAttachmentSpec_v1beta1_FromProto(mapCtx *direct.MapContext, in *pb.ServiceAttachment) *krm.ComputeServiceAttachmentSpec {
	if in == nil {
		return nil
	}
	out := &krm.ComputeServiceAttachmentSpec{}
	out.ConnectionPreference = direct.ValueOf(in.ConnectionPreference)
	out.ConsumerAcceptLists = direct.Slice_FromProto(mapCtx, in.ConsumerAcceptLists, ServiceattachmentConsumerAcceptLists_v1beta1_FromProto)
	out.ConsumerRejectLists = ComputeServiceAttachmentSpec_ConsumerRejectLists_FromProto(mapCtx, in.ConsumerRejectLists)
	out.Description = in.Description
	out.EnableProxyProtocol = in.EnableProxyProtocol

	if len(in.NatSubnets) > 0 {
		out.NatSubnets = make([]krm.ComputeSubnetworkRef, len(in.NatSubnets))
		for i, x := range in.NatSubnets {
			out.NatSubnets[i] = krm.ComputeSubnetworkRef{External: x}
		}
	}

	if in.ProducerForwardingRule != nil {
		out.TargetServiceRef = &refsv1beta1.ComputeForwardingRuleRef{External: *in.ProducerForwardingRule}
	}

	return out
}

func ComputeServiceAttachmentSpec_v1beta1_ToProto(mapCtx *direct.MapContext, in *krm.ComputeServiceAttachmentSpec) *pb.ServiceAttachment {
	if in == nil {
		return nil
	}
	out := &pb.ServiceAttachment{}
	out.ConnectionPreference = direct.LazyPtr(in.ConnectionPreference)
	out.ConsumerAcceptLists = direct.Slice_ToProto(mapCtx, in.ConsumerAcceptLists, ServiceattachmentConsumerAcceptLists_v1beta1_ToProto)
	out.ConsumerRejectLists = ComputeServiceAttachmentSpec_ConsumerRejectLists_ToProto(mapCtx, in.ConsumerRejectLists)
	out.Description = in.Description
	out.EnableProxyProtocol = in.EnableProxyProtocol

	if len(in.NatSubnets) > 0 {
		out.NatSubnets = make([]string, len(in.NatSubnets))
		for i, x := range in.NatSubnets {
			out.NatSubnets[i] = x.External
		}
	}

	if in.TargetServiceRef != nil {
		out.ProducerForwardingRule = &in.TargetServiceRef.External
	}

	return out
}

func ComputeServiceAttachmentStatus_v1beta1_FromProto(mapCtx *direct.MapContext, in *pb.ServiceAttachment) *krm.ComputeServiceAttachmentStatus {
	if in == nil {
		return nil
	}
	out := &krm.ComputeServiceAttachmentStatus{}
	out.ConnectedEndpoints = direct.Slice_FromProto(mapCtx, in.ConnectedEndpoints, ServiceattachmentConnectedEndpointsStatus_v1beta1_FromProto)
	out.Fingerprint = in.Fingerprint
	if in.Id != nil {
		val := int64(*in.Id)
		out.Id = &val
	}
	if in.PscServiceAttachmentId != nil {
		out.PscServiceAttachmentId = &krm.ServiceattachmentPscServiceAttachmentIdStatus{}
		if in.PscServiceAttachmentId.High != nil {
			val := int64(*in.PscServiceAttachmentId.High)
			out.PscServiceAttachmentId.High = &val
		}
		if in.PscServiceAttachmentId.Low != nil {
			val := int64(*in.PscServiceAttachmentId.Low)
			out.PscServiceAttachmentId.Low = &val
		}
	}
	out.Region = in.Region
	out.SelfLink = in.SelfLink
	return out
}

func ComputeServiceAttachmentStatus_v1beta1_ToProto(mapCtx *direct.MapContext, in *krm.ComputeServiceAttachmentStatus) *pb.ServiceAttachment {
	if in == nil {
		return nil
	}
	out := &pb.ServiceAttachment{}
	out.ConnectedEndpoints = direct.Slice_ToProto(mapCtx, in.ConnectedEndpoints, ServiceattachmentConnectedEndpointsStatus_v1beta1_ToProto)
	out.Fingerprint = in.Fingerprint
	if in.Id != nil {
		val := uint64(*in.Id)
		out.Id = &val
	}
	if in.PscServiceAttachmentId != nil {
		out.PscServiceAttachmentId = &pb.Uint128{}
		if in.PscServiceAttachmentId.High != nil {
			val := uint64(*in.PscServiceAttachmentId.High)
			out.PscServiceAttachmentId.High = &val
		}
		if in.PscServiceAttachmentId.Low != nil {
			val := uint64(*in.PscServiceAttachmentId.Low)
			out.PscServiceAttachmentId.Low = &val
		}
	}
	out.Region = in.Region
	out.SelfLink = in.SelfLink
	return out
}

func ServiceattachmentConsumerAcceptLists_v1beta1_FromProto(mapCtx *direct.MapContext, in *pb.ServiceAttachmentConsumerProjectLimit) *krm.ServiceattachmentConsumerAcceptLists {
	if in == nil {
		return nil
	}
	out := &krm.ServiceattachmentConsumerAcceptLists{}
	if in.ConnectionLimit != nil {
		val := int64(*in.ConnectionLimit)
		out.ConnectionLimit = &val
	}
	if in.ProjectIdOrNum != nil {
		out.ProjectRef.External = *in.ProjectIdOrNum
	}
	return out
}

func ServiceattachmentConsumerAcceptLists_v1beta1_ToProto(mapCtx *direct.MapContext, in *krm.ServiceattachmentConsumerAcceptLists) *pb.ServiceAttachmentConsumerProjectLimit {
	if in == nil {
		return nil
	}
	out := &pb.ServiceAttachmentConsumerProjectLimit{}
	if in.ConnectionLimit != nil {
		val := uint32(*in.ConnectionLimit)
		out.ConnectionLimit = &val
	}
	if in.ProjectRef.External != "" {
		out.ProjectIdOrNum = &in.ProjectRef.External
	}
	return out
}

func ComputeServiceAttachmentSpec_ConsumerRejectLists_FromProto(mapCtx *direct.MapContext, in []string) []refsv1beta1.ProjectRef {
	if in == nil {
		return nil
	}
	var out []refsv1beta1.ProjectRef
	for _, i := range in {
		out = append(out, refsv1beta1.ProjectRef{
			External: i,
		})
	}
	return out
}

func ComputeServiceAttachmentSpec_ConsumerRejectLists_ToProto(mapCtx *direct.MapContext, in []refsv1beta1.ProjectRef) []string {
	if in == nil {
		return nil
	}
	var out []string
	for _, i := range in {
		out = append(out, i.External)
	}
	return out
}

func ServiceattachmentConnectedEndpointsStatus_v1beta1_FromProto(mapCtx *direct.MapContext, in *pb.ServiceAttachmentConnectedEndpoint) *krm.ServiceattachmentConnectedEndpointsStatus {
	if in == nil {
		return nil
	}
	out := &krm.ServiceattachmentConnectedEndpointsStatus{}
	out.Endpoint = in.Endpoint
	if in.PscConnectionId != nil {
		val := int64(*in.PscConnectionId)
		out.PscConnectionId = &val
	}
	out.Status = in.Status
	return out
}

func ServiceattachmentConnectedEndpointsStatus_v1beta1_ToProto(mapCtx *direct.MapContext, in *krm.ServiceattachmentConnectedEndpointsStatus) *pb.ServiceAttachmentConnectedEndpoint {
	if in == nil {
		return nil
	}
	out := &pb.ServiceAttachmentConnectedEndpoint{}
	out.Endpoint = in.Endpoint
	if in.PscConnectionId != nil {
		val := uint64(*in.PscConnectionId)
		out.PscConnectionId = &val
	}
	out.Status = in.Status
	return out
}
