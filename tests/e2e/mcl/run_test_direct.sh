#!/bin/bash
# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

export CLUSTER_A="kcc-cluster-a"
export CLUSTER_B="kcc-cluster-b"
export CTX_A="kind-${CLUSTER_A}"
export CTX_B="kind-${CLUSTER_B}"
export LEASE_NAMESPACE="cnrm-system"
export WATCH_NAMESPACE="tenant-a"
export LEASE_NAME="kcc-leader-lease"
export SYNCER_NAME="${LEASE_NAME}-${WATCH_NAMESPACE}"

retry() {
  local count=$1
  local delay=$2
  shift 2
  local cmd="$@"
  for i in $(seq 1 $count); do
    if eval "$cmd"; then
      return 0
    fi
    sleep $delay
  done
  echo "FAILED: $cmd"
  return 1
}

retry 15 2 "kubectl --context='${CTX_A}' get multiclusterlease '${LEASE_NAME}' -n '${LEASE_NAMESPACE}'"
retry 15 2 "kubectl --context='${CTX_B}' get multiclusterlease '${LEASE_NAME}' -n '${LEASE_NAMESPACE}'"
retry 15 2 "kubectl --context='${CTX_A}' wait --for=condition=available deployment/syncer-controller-manager -n '${LEASE_NAMESPACE}' --timeout=5s"
retry 15 2 "kubectl --context='${CTX_B}' wait --for=condition=available deployment/syncer-controller-manager -n '${LEASE_NAMESPACE}' --timeout=5s"

retry 15 2 "kubectl --context='${CTX_A}' wait --for=condition=available deployment/multiclusterlease-controller-manager -n '${LEASE_NAMESPACE}' --timeout=5s"
retry 15 2 "kubectl --context='${CTX_B}' wait --for=condition=available deployment/multiclusterlease-controller-manager -n '${LEASE_NAMESPACE}' --timeout=5s"

echo "Waiting for REAL MCL Controller to elect a leader..."
retry 30 5 "kubectl --context='${CTX_A}' get multiclusterlease '${LEASE_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.status.globalHolderIdentity != null and .status.globalHolderIdentity != \"\"'"

leader=$(kubectl --context="${CTX_A}" get multiclusterlease "${LEASE_NAME}" -n "${LEASE_NAMESPACE}" -o json | jq -r '.status.globalHolderIdentity')
if [ "$leader" = "${CLUSTER_A}" ]; then
   follower="${CLUSTER_B}"
   leader_ctx="${CTX_A}"
   follower_ctx="${CTX_B}"
else
   follower="${CLUSTER_A}"
   leader_ctx="${CTX_B}"
   follower_ctx="${CTX_A}"
fi

echo "Elected Leader: $leader"

retry 15 2 "kubectl --context='${follower_ctx}' get krmsyncer '${SYNCER_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.spec.suspend == false'"
retry 15 2 "kubectl --context='${leader_ctx}' get krmsyncer '${SYNCER_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.spec.suspend == true'"

kubectl --context="${CTX_A}" apply -f tests/e2e/mcl/storagebucket.yaml
kubectl --context="${CTX_B}" apply -f tests/e2e/mcl/storagebucket.yaml

retry 15 2 "kubectl --context='${leader_ctx}' get storagebucket test-ha-bucket -n '${WATCH_NAMESPACE}' -o json | grep -E 'UpdateFailed|UpToDate'"
retry 40 3 "kubectl --context='${follower_ctx}' get storagebucket test-ha-bucket -n '${WATCH_NAMESPACE}' -o json | grep -E 'UpdateFailed|UpToDate'"

echo "Failing over to $follower..."
kubectl --context="${leader_ctx}" scale deployment cnrm-controller-manager -n "${WATCH_NAMESPACE}" --replicas=0
kubectl --context="${leader_ctx}" scale deployment multiclusterlease-controller-manager -n "${LEASE_NAMESPACE}" --replicas=0

retry 20 2 "kubectl --context='${follower_ctx}' get multiclusterlease '${LEASE_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.status.globalHolderIdentity == \"${follower}\"'"
retry 20 2 "kubectl --context='${follower_ctx}' get krmsyncer '${SYNCER_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.spec.suspend == true'"

kubectl --context="${follower_ctx}" patch storagebucket test-ha-bucket -n "${WATCH_NAMESPACE}" --type='merge' -p '{"spec":{"storageClass":"STANDARD"}}'
kubectl --context="${leader_ctx}" patch storagebucket test-ha-bucket -n "${WATCH_NAMESPACE}" --type='merge' -p '{"spec":{"storageClass":"STANDARD"}}'

retry 20 2 "kubectl --context='${follower_ctx}' get storagebucket test-ha-bucket -n '${WATCH_NAMESPACE}' -o json | jq -e '.status.observedGeneration == 2'"

echo "Failing back to $leader..."
kubectl --context="${leader_ctx}" scale deployment cnrm-controller-manager -n "${WATCH_NAMESPACE}" --replicas=1
retry 20 2 "kubectl --context='${leader_ctx}' get krmsyncer '${SYNCER_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.spec.suspend == false'"
kubectl --context="${leader_ctx}" patch storagebucket test-ha-bucket -n "${WATCH_NAMESPACE}" --type='merge' -p '{"spec":{"storageClass":"STANDARD"}}'

retry 40 3 "kubectl --context='${leader_ctx}' get storagebucket test-ha-bucket -n '${WATCH_NAMESPACE}' -o json | jq -e '.status.observedGeneration == 2'"

kubectl --context="${leader_ctx}" scale deployment multiclusterlease-controller-manager -n "${LEASE_NAMESPACE}" --replicas=1
kubectl --context="${follower_ctx}" scale deployment multiclusterlease-controller-manager -n "${LEASE_NAMESPACE}" --replicas=0

retry 20 2 "kubectl --context='${leader_ctx}' get multiclusterlease '${LEASE_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.status.globalHolderIdentity == \"${leader}\"'"
retry 20 2 "kubectl --context='${leader_ctx}' get krmsyncer '${SYNCER_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.spec.suspend == true'"
retry 20 2 "kubectl --context='${follower_ctx}' get krmsyncer '${SYNCER_NAME}' -n '${LEASE_NAMESPACE}' -o json | jq -e '.spec.suspend == false'"

kubectl --context="${leader_ctx}" patch storagebucket test-ha-bucket -n "${WATCH_NAMESPACE}" --type='merge' -p '{"spec":{"storageClass":"NEARLINE"}}'
kubectl --context="${follower_ctx}" patch storagebucket test-ha-bucket -n "${WATCH_NAMESPACE}" --type='merge' -p '{"spec":{"storageClass":"NEARLINE"}}'

retry 20 2 "kubectl --context='${leader_ctx}' get storagebucket test-ha-bucket -n '${WATCH_NAMESPACE}' -o json | jq -e '.status.observedGeneration == 3'"
retry 40 3 "kubectl --context='${follower_ctx}' get storagebucket test-ha-bucket -n '${WATCH_NAMESPACE}' -o json | jq -e '.status.observedGeneration == 3'"

echo "SUCCESS!"
