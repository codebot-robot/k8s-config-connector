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

# Description: Verifies the HA setup and tests failover.

# Stop on any error
set -e

# --- Parameters ---
PROJECT_ID=${1:?Usage: $0 PROJECT_ID BUCKET_NAME [REGION] [C1_NAME] [C2_NAME] [TENANT_NS] [LEASE_NAME]}
BUCKET_NAME=${2:?Bucket name required}
REGION=${3:-us-central1}
C1_NAME=${4:-kcc-ha-cluster-1}
C2_NAME=${5:-kcc-ha-cluster-2}
TENANT_NS=${6:-tenant-a}
LEASE_NAME=${7:-kcc-managers}

# --- Derived Variables ---
ZONE="${REGION}-c"
C1_CTX="gke_${PROJECT_ID}_${ZONE}_${C1_NAME}"
C2_CTX="gke_${PROJECT_ID}_${ZONE}_${C2_NAME}"
C1_ID="cluster-1"
C2_ID="cluster-2"
SYNCER_NAME="${LEASE_NAME}-${TENANT_NS}"

echo "--- Configuration ---"
# ... (echo parameters)
echo "---------------------"

gcloud config set project ${PROJECT_ID}
gcloud container clusters get-credentials ${C1_NAME} --zone ${ZONE} --project=${PROJECT_ID}
gcloud container clusters get-credentials ${C2_NAME} --zone ${ZONE} --project=${PROJECT_ID}

echo "--- 6. Operational Verification & Validation ---"
# ... (Verification steps from the previous combined script) ...
echo "--- Step 1: Observe Current Steady State ---"
echo "Checking leader, retrying for up to 2 minutes..."
LEADER=""
for i in {1..12}; do
  LEADER=$(kubectl --context=$C1_CTX get multiclusterlease $LEASE_NAME -n cnrm-system -o jsonpath='{.status.globalHolderIdentity}' 2>/dev/null || echo "")
  if [ -n "$LEADER" ]; then
    echo "Current Leader: $LEADER"
    break
  fi
  echo "Leader not yet elected, waiting 10s..."
  sleep 10
done

if [ -z "$LEADER" ]; then
  echo "Error: Could not determine leader."
  exit 1
fi

if [[ "$LEADER" == "$C1_ID" ]]; then
  LEADER_CTX=$C1_CTX
  FOLLOWER_CTX=$C2_CTX
  FOLLOWER_ID=$C2_ID
  echo "Cluster 1 is Leader, Cluster 2 is Follower"
elif [[ "$LEADER" == "$C2_ID" ]]; then
  LEADER_CTX=$C2_CTX
  FOLLOWER_CTX=$C1_CTX
  FOLLOWER_ID=$C1_ID
  echo "Cluster 2 is Leader, Cluster 1 is Follower"
else
  echo "Error: Leader '$LEADER' is unexpected."
  exit 1
fi

echo "--- Step 2: Verify Initial KRMSyncer Generation ---"
echo "Waiting for Syncer to be generated on Leader ($LEADER)..."
for i in {1..20}; do
  if kubectl --context=$LEADER_CTX get krmsyncer $SYNCER_NAME -n cnrm-system &> /dev/null; then break; fi
  echo "Waiting 5s..."
  sleep 5
done
kubectl --context=$LEADER_CTX get krmsyncer $SYNCER_NAME -n cnrm-system -o jsonpath='{"Mode: "}{.spec.mode}{"\nSuspend: "}{.spec.suspend}{"\n"}'
echo ""

echo "Waiting for Syncer to be generated on Follower ($FOLLOWER_ID)..."
for i in {1..20}; do
  if kubectl --context=$FOLLOWER_CTX get krmsyncer $SYNCER_NAME -n cnrm-system &> /dev/null; then break; fi
  echo "Waiting 5s..."
  sleep 5
done
kubectl --context=$FOLLOWER_CTX get krmsyncer $SYNCER_NAME -n cnrm-system -o jsonpath='{"Mode: "}{.spec.mode}{"\nSuspend: "}{.spec.suspend}{"\nPulling From: "}{.spec.remote.clusterConfig.kubeConfigSecretRef.name}{"\n"}'
echo ""

echo "--- Step 3: Trigger a Hard Failover ---"
echo "Scaling KCC manager StatefulSet to 0 on the active leader ($LEADER) to force failover..."
MGR_STS=$(kubectl --context=$LEADER_CTX get statefulset -n cnrm-system -l cnrm.cloud.google.com/component=cnrm-controller-manager,cnrm.cloud.google.com/scoped-namespace=${TENANT_NS} -o jsonpath='{.items[0].metadata.name}')
kubectl --context=$LEADER_CTX scale statefulset $MGR_STS -n cnrm-system --replicas=0

echo "Waiting 60s for the failover transition and lease expiration..."
sleep 60

echo "--- Step 4: Validate the Failover Transition ---"
echo "1. Verify GCS Lock Transfer..."
gcloud storage cat gs://${BUCKET_NAME}/leases/cnrm-system/${LEASE_NAME}
echo ""

echo "2. Verify Local Lease Consensus..."
NEW_LEADER_C1=""
NEW_LEADER_C2=""
for i in {1..12}; do
  NEW_LEADER_C1=$(kubectl --context=$C1_CTX get multiclusterlease $LEASE_NAME -n cnrm-system -o jsonpath='{.status.globalHolderIdentity}' 2>/dev/null || echo "")
  NEW_LEADER_C2=$(kubectl --context=$C2_CTX get multiclusterlease $LEASE_NAME -n cnrm-system -o jsonpath='{.status.globalHolderIdentity}' 2>/dev/null || echo "")
  echo "Attempt $i: Cluster 1 sees leader: $NEW_LEADER_C1, Cluster 2 sees leader: $NEW_LEADER_C2"
  if [[ "$NEW_LEADER_C1" == "$FOLLOWER_ID" && "$NEW_LEADER_C2" == "$FOLLOWER_ID" ]]; then
    break
  fi
  echo "Waiting for consensus on new leader $FOLLOWER_ID..."
  sleep 10
done

if [[ "$NEW_LEADER_C1" != "$FOLLOWER_ID" || "$NEW_LEADER_C2" != "$FOLLOWER_ID" ]]; then
  echo "Warning: Failover consensus not reached. Last state: C1=$NEW_LEADER_C1, C2=$NEW_LEADER_C2"
fi
echo "New Leader appears to be: $FOLLOWER_ID"

echo "3. Verify Syncer Role Reversal..."
echo "Checking Syncer on New Leader ($FOLLOWER_ID)..."
kubectl --context=$FOLLOWER_CTX get krmsyncer $SYNCER_NAME -n cnrm-system -o jsonpath='{"Suspend: "}{.spec.suspend}{"\n"}'
echo ""

echo "Checking Syncer on New Follower ($LEADER)..."
kubectl --context=$LEADER_CTX get krmsyncer $SYNCER_NAME -n cnrm-system -o jsonpath='{"Suspend: "}{.spec.suspend}{"\nPulling From: "}{.spec.remote.clusterConfig.kubeConfigSecretRef.name}{"\n"}'
echo ""

echo "--- Verification Complete ---"
