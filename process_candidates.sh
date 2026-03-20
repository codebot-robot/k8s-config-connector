#!/bin/bash
candidates=(
"filestore FilestoreBackup"
"binaryauthorization BinaryAuthorizationAttestor"
"networkservices NetworkServicesEndpointPolicy"
"cloudfunctions CloudFunctionsFunction"
"networkservices NetworkServicesHTTPRoute"
"dlp DLPJobTrigger"
"dlp DLPInspectTemplate"
"dataproc DataprocCluster"
"identityplatform IdentityPlatformConfig"
"dlp DLPStoredInfoType"
"cloudscheduler CloudSchedulerJob"
"billingbudgets BillingBudgetsBudget"
"osconfig OSConfigGuestPolicy"
"networkservices NetworkServicesTLSRoute"
"configcontroller ConfigControllerInstance"
"identityplatform IdentityPlatformTenant"
"binaryauthorization BinaryAuthorizationPolicy"
"networkservices NetworkServicesGRPCRoute"
"networkservices NetworkServicesTCPRoute"
"eventarc EventarcTrigger"
"recaptchaenterprise RecaptchaEnterpriseKey"
"filestore FilestoreInstance"
"datafusion DataFusionInstance"
"identityplatform IdentityPlatformOAuthIDPConfig"
"networkservices NetworkServicesGateway"
"dataproc DataprocWorkflowTemplate"
"identityplatform IdentityPlatformTenantOAuthIDPConfig"
"containeranalysis ContainerAnalysisNote"
"networkservices NetworkServicesMesh"
"osconfig OSConfigOSPolicyAssignment"
"networkconnectivity NetworkConnectivitySpoke"
"dataproc DataprocAutoscalingPolicy"
"dlp DLPDeidentifyTemplate"
"networkconnectivity NetworkConnectivityHub"
)

OPEN_ISSUES=$(gh issue list --search "Create generate.sh and types.go files for in:title label:overseer label:area/direct" --state open --limit 50 --json number -q 'length')

echo "Current open issues: $OPEN_ISSUES"

for candidate in "${candidates[@]}"; do
  group=$(echo $candidate | awk '{print $1}')
  kind=$(echo $candidate | awk '{print $2}')
  
  # Search for issue
  ISSUE_JSON=$(gh issue list --search "Create generate.sh and types.go files for $group $kind in:title" --state all --json number,labels -q '.[0]')
  
  if [ "$ISSUE_JSON" != "null" ] && [ -n "$ISSUE_JSON" ]; then
    ISSUE_NUM=$(echo "$ISSUE_JSON" | jq -r '.number')
    
    MISSING_LABELS=""
    if ! echo "$ISSUE_JSON" | jq -e '.labels[].name | select(. == "overseer")' > /dev/null; then
      MISSING_LABELS="overseer"
    fi
    if ! echo "$ISSUE_JSON" | jq -e '.labels[].name | select(. == "area/direct")' > /dev/null; then
      MISSING_LABELS="$MISSING_LABELS area/direct"
    fi
    if ! echo "$ISSUE_JSON" | jq -e '.labels[].name | select(. == "priority/medium")' > /dev/null; then
      MISSING_LABELS="$MISSING_LABELS priority/medium"
    fi
    
    if [ -n "$MISSING_LABELS" ]; then
      echo "Adding labels ($MISSING_LABELS) to issue #$ISSUE_NUM for $group $kind"
      # gh issue edit $ISSUE_NUM --add-label "$(echo $MISSING_LABELS | tr ' ' ',')"
      # For now, just echoing to see what would happen
      # Let's actually add the labels
      LABELS_ARGS=""
      for L in $MISSING_LABELS; do
        LABELS_ARGS="$LABELS_ARGS --add-label $L"
      done
      gh issue edit $ISSUE_NUM $LABELS_ARGS
    else
      echo "Issue #$ISSUE_NUM for $group $kind already has correct labels."
    fi
  else
    if [ "$OPEN_ISSUES" -ge 10 ]; then
      echo "More than 10 pending issues already exist ($OPEN_ISSUES). Skipping creating new issue for $group $kind."
    else
      echo "Would create issue for $group $kind, but let's do that via Python or another script."
    fi
  fi
done
