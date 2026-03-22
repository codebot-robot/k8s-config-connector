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

# Description: Creates two GKE clusters with Workload Identity enabled.

# Stop on any error
set -e

# --- Parameters ---
PROJECT_ID=${1:?Usage: $0 PROJECT_ID [REGION] [C1_NAME] [C2_NAME]}
REGION=${2:-us-central1}
C1_NAME=${3:-kcc-ha-cluster-1}
C2_NAME=${4:-kcc-ha-cluster-2}

# --- Derived Variables ---
ZONE="${REGION}-c"

echo "--- Configuration ---"
echo "PROJECT_ID: ${PROJECT_ID}"
echo "REGION: ${REGION}"
echo "ZONE: ${ZONE}"
echo "C1_NAME: ${C1_NAME}"
echo "C2_NAME: ${C2_NAME}"
echo "---------------------"

gcloud config set project ${PROJECT_ID}

echo "--- 0.2. Creating GKE Clusters ---"
if ! gcloud container clusters describe ${C1_NAME} --zone ${ZONE} --project=${PROJECT_ID} &> /dev/null; then
  echo "Creating Cluster 1: ${C1_NAME}..."
  gcloud container clusters create ${C1_NAME} --zone ${ZONE} --workload-pool="${PROJECT_ID}.svc.id.goog" --num-nodes=1 --machine-type=e2-standard-4 --project=${PROJECT_ID}
else
  echo "Cluster ${C1_NAME} already exists."
fi

if ! gcloud container clusters describe ${C2_NAME} --zone ${ZONE} --project=${PROJECT_ID} &> /dev/null; then
  echo "Creating Cluster 2: ${C2_NAME}..."
  gcloud container clusters create ${C2_NAME} --zone ${ZONE} --workload-pool="${PROJECT_ID}.svc.id.goog" --num-nodes=1 --machine-type=e2-standard-4 --project=${PROJECT_ID}
else
  echo "Cluster ${C2_NAME} already exists."
fi

echo "Getting cluster credentials..."
gcloud container clusters get-credentials ${C1_NAME} --zone ${ZONE} --project=${PROJECT_ID}
gcloud container clusters get-credentials ${C2_NAME} --zone ${ZONE} --project=${PROJECT_ID}

echo "--- Cluster Creation Complete ---"
