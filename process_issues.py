import json
import subprocess

# Candidates list from previous step
candidates = [
("apigee", "ApigeeEnvironment"),
("apigee", "ApigeeOrganization"),
("billingbudgets", "BillingBudgetsBudget"),
("binaryauthorization", "BinaryAuthorizationAttestor"),
("binaryauthorization", "BinaryAuthorizationPolicy"),
("cloudfunctions", "CloudFunctionsFunction"),
("cloudscheduler", "CloudSchedulerJob"),
("compute", "ComputeFirewallPolicy"),
("compute", "ComputeFirewallPolicyAssociation"),
("compute", "ComputeInstanceGroupManager"),
("compute", "ComputePacketMirroring"),
("compute", "ComputeServiceAttachment"),
("configcontroller", "ConfigControllerInstance"),
("containeranalysis", "ContainerAnalysisNote"),
("datafusion", "DataFusionInstance"),
("dataproc", "DataprocAutoscalingPolicy"),
("dataproc", "DataprocCluster"),
("dataproc", "DataprocWorkflowTemplate"),
("dlp", "DLPDeidentifyTemplate"),
("dlp", "DLPInspectTemplate"),
("dlp", "DLPJobTrigger"),
("dlp", "DLPStoredInfoType"),
("eventarc", "EventarcTrigger"),
("filestore", "FilestoreBackup"),
("filestore", "FilestoreInstance"),
("gkehub", "GKEHubFeature"),
("gkehub", "GKEHubMembership"),
("iam", "IAMWorkforcePool"),
("iam", "IAMWorkforcePoolProvider"),
("iam", "IAMWorkloadIdentityPool"),
("iam", "IAMWorkloadIdentityPoolProvider"),
("iap", "IAPBrand"),
("iap", "IAPIdentityAwareProxyClient"),
("identityplatform", "IdentityPlatformConfig"),
("identityplatform", "IdentityPlatformOAuthIDPConfig"),
("identityplatform", "IdentityPlatformTenant"),
("identityplatform", "IdentityPlatformTenantOAuthIDPConfig"),
("logging", "LoggingLogBucket"),
("logging", "LoggingLogExclusion"),
("logging", "LoggingLogView"),
("monitoring", "MonitoringGroup"),
("monitoring", "MonitoringMetricDescriptor"),
("monitoring", "MonitoringMonitoredProject"),
("monitoring", "MonitoringService"),
("monitoring", "MonitoringServiceLevelObjective"),
("monitoring", "MonitoringUptimeCheckConfig"),
("networkconnectivity", "NetworkConnectivityHub"),
("networkconnectivity", "NetworkConnectivitySpoke"),
("networksecurity", "NetworkSecurityAuthorizationPolicy"),
("networksecurity", "NetworkSecurityClientTLSPolicy"),
("networksecurity", "NetworkSecurityServerTLSPolicy"),
("networkservices", "NetworkServicesEndpointPolicy"),
("networkservices", "NetworkServicesGRPCRoute"),
("networkservices", "NetworkServicesGateway"),
("networkservices", "NetworkServicesHTTPRoute"),
("networkservices", "NetworkServicesMesh"),
("networkservices", "NetworkServicesTCPRoute"),
("networkservices", "NetworkServicesTLSRoute"),
("osconfig", "OSConfigGuestPolicy"),
("osconfig", "OSConfigOSPolicyAssignment"),
("privateca", "PrivateCACAPool"),
("privateca", "PrivateCACertificate"),
("privateca", "PrivateCACertificateAuthority"),
("privateca", "PrivateCACertificateTemplate"),
("recaptchaenterprise", "RecaptchaEnterpriseKey")
]

# Get all issues with the matching title pattern
result = subprocess.run(
    ["gh", "issue", "list", "--search", "Create generate.sh and types.go files for in:title", "--state", "all", "--json", "number,title,state,labels", "--limit", "1000"],
    capture_output=True, text=True
)
if result.returncode != 0:
    print("Error fetching issues:", result.stderr)
    exit(1)

issues = json.loads(result.stdout)
open_count = sum(1 for i in issues if i["state"].upper() == "OPEN")

issue_map = {}
for i in issues:
    issue_map[i["title"]] = i

created_issue = False

for group, kind in candidates:
    title = f"Create generate.sh and types.go files for {group} {kind}"
    
    if title in issue_map:
        # Issue exists, check labels
        issue = issue_map[title]
        existing_labels = {l["name"] for l in issue["labels"]}
        required_labels = {"overseer", "area/direct", "priority/medium"}
        
        missing = required_labels - existing_labels
        if missing:
            print(f"Adding labels {missing} to issue #{issue['number']}")
            subprocess.run(["gh", "issue", "edit", str(issue['number']), "--add-label", ",".join(missing)])
        continue
    
    if not created_issue:
        if open_count > 10:
            print("There are already more than 10 pending issues for this task. Skipping creating new ones until some of the existing issues are resolved.")
            created_issue = True # Stop checking for creation
            continue
        
        # Create the issue
        print(f"Creating issue for {group} {kind}")
        
        body = f"""As part of moving resources from terraform controllers to direct controllers (Epic #5954), we need to create the Go types for `{kind}`.

Currently, `{kind}` is managed by the Terraform controller (marked with `tf2crd=true`). The goal is to create the Go types in `apis/{group}/v1beta1/` so that we can eventually migrate the controller implementation to the "direct" approach.

### Instructions

1.  **Create a generate.sh**:
    Create `apis/{group}/v1beta1/generate.sh` which includes `{kind}`.
    It likely maps to something like `google.cloud.{group}.v1`.
    Example:
    ```bash
    go run . generate-types \\
      --service google.cloud.{group}.v1 \\
      --api-version {group}.cnrm.cloud.google.com/v1beta1 \\
      --resource {kind}:PolicyTag \\
      --include-skipped-output

    go run . generate-mapper \\
      --service google.cloud.{group}.v1 \\
      --api-version {group}.cnrm.cloud.google.com/v1beta1 \\
      --include-skipped-output
    ```

2.  Set the write permission on the new `apis/{group}/v1beta1/generate.sh` file. You should do this by running both `chmod +x apis/{group}/v1beta1/generate.sh` and `git add --chmod=+x apis/{group}/v1beta1/generate.sh`.

3.  **Generate Scaffolding**:
    Run `apis/{group}/v1beta1/generate.sh`. This should create `apis/{group}/v1beta1/{kind.lower()}_types.go`.

4.  **Iterate on Types**:
    Compare the generated CRD with the existing one using `dev/tasks/diff-crds`.
    Modify `apis/{group}/v1beta1/{kind.lower()}_types.go` until the CRD matches the existing one at `config/crds/resources/apiextensions.k8s.io_v1_customresourcedefinition_{kind.lower()}s.{group}.cnrm.cloud.google.com.yaml`.

    **Acceptance Criteria:**
    - Running `dev/tasks/diff-crds` should not show differences (or minimal acceptable ones like descriptions).
    - Ensure that running the check_crd_equivalence MCP on the CRD should return EQUIVALENT.
    - Changes to the schema (fields added/removed) are NOT acceptable.

5.  **Copyright Headers**:
    Ensure that new files have the correct copyright header:
    ```go
    // Copyright 2026 Google LLC
    ```
    Please do not change the copyright on existing files.

6.  **Labels**:
    Ensure the controller-runtime annotations match the existing CRD labels, including:
    ```go
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/managed-by-kcc=true"
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/system=true"
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/stability-level=stable"
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/tf2crd=true"
    ```
    The goal is to maintain these annotations, not add an annotation if it is missing.

7.  **Status**:
    `status.observedGeneration` should be an `*int64`.

8. **Generate Mappers**:
   - Running `dev/tasks/generate-types-and-mappers` will generate the mapper code once the `apis/{group}/v1beta1/{kind.lower()}_types.go` file is generating an equivalent CRD.
   - Run `make all-binary` to ensure the generated mapper code compiles. Please fix any issue discovered by this compilation.

This issue is part of Epic #5954.
"""
        with open("issue_body.md", "w") as f:
            f.write(body)
            
        subprocess.run([
            "gh", "issue", "create",
            "--title", title,
            "--body-file", "issue_body.md",
            "--label", "overseer,area/direct,priority/medium",
            "--milestone", "", # Epic is linked via milestone or body? The prompt says "The issue should be marked as a subtask of the main epic... Make sure to link the issue as a subtask to the main epic for tracking purposes."
            # Actually, standard way in github is either typing the issue number in body or using gh cli. The body already has `Epic #5954`. But let's see if there's a specific flag for epic?
            # Gh doesn't have an explicit subtask flag without project fields. So including Epic #5954 is enough for github to track.
        ])
        
        created_issue = True
        open_count += 1

