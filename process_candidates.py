import json
import subprocess
import sys

def run_cmd(cmd):
    return subprocess.check_output(cmd, shell=True).decode('utf-8')

candidates_output = """
FOUND: Group=networkservices, Kind=NetworkServicesGRPCRoute
FOUND: Group=networkservices, Kind=NetworkServicesTCPRoute
FOUND: Group=identityplatform, Kind=IdentityPlatformConfig
FOUND: Group=binaryauthorization, Kind=BinaryAuthorizationPolicy
FOUND: Group=identityplatform, Kind=IdentityPlatformTenant
FOUND: Group=recaptchaenterprise, Kind=RecaptchaEnterpriseKey
FOUND: Group=dlp, Kind=DLPInspectTemplate
FOUND: Group=osconfig, Kind=OSConfigOSPolicyAssignment
FOUND: Group=networkservices, Kind=NetworkServicesMesh
FOUND: Group=networkservices, Kind=NetworkServicesTLSRoute
FOUND: Group=osconfig, Kind=OSConfigGuestPolicy
FOUND: Group=identityplatform, Kind=IdentityPlatformTenantOAuthIDPConfig
FOUND: Group=dataproc, Kind=DataprocCluster
FOUND: Group=binaryauthorization, Kind=BinaryAuthorizationAttestor
FOUND: Group=networkconnectivity, Kind=NetworkConnectivitySpoke
FOUND: Group=networkconnectivity, Kind=NetworkConnectivityHub
FOUND: Group=billingbudgets, Kind=BillingBudgetsBudget
FOUND: Group=dlp, Kind=DLPJobTrigger
FOUND: Group=containeranalysis, Kind=ContainerAnalysisNote
FOUND: Group=dataproc, Kind=DataprocAutoscalingPolicy
FOUND: Group=networkservices, Kind=NetworkServicesGateway
FOUND: Group=dataproc, Kind=DataprocWorkflowTemplate
FOUND: Group=filestore, Kind=FilestoreBackup
FOUND: Group=identityplatform, Kind=IdentityPlatformOAuthIDPConfig
FOUND: Group=dlp, Kind=DLPDeidentifyTemplate
FOUND: Group=eventarc, Kind=EventarcTrigger
FOUND: Group=dlp, Kind=DLPStoredInfoType
FOUND: Group=datafusion, Kind=DataFusionInstance
FOUND: Group=networkservices, Kind=NetworkServicesHTTPRoute
FOUND: Group=networkservices, Kind=NetworkServicesEndpointPolicy
FOUND: Group=cloudfunctions, Kind=CloudFunctionsFunction
FOUND: Group=configcontroller, Kind=ConfigControllerInstance
FOUND: Group=cloudscheduler, Kind=CloudSchedulerJob
FOUND: Group=filestore, Kind=FilestoreInstance
"""

candidates = []
for line in candidates_output.strip().split('\n'):
    parts = line.split(', ')
    group = parts[0].split('=')[1]
    kind = parts[1].split('=')[1]
    candidates.append((group, kind))

# Get all issues with the matching title prefix
issues_json = run_cmd('gh issue list --search "in:title Create generate.sh and types.go files for" --state all --json number,title,state,labels -L 100')
issues = json.loads(issues_json)

pending_count = sum(1 for issue in issues if issue['state'] == 'OPEN')

for group, kind in candidates:
    target_title_lower = f"Create generate.sh and types.go files for {group} {kind}".lower()
    
    existing_issue = None
    for issue in issues:
        if target_title_lower in issue['title'].lower():
            existing_issue = issue
            break
            
    if existing_issue:
        # Check labels
        current_labels = [l['name'] for l in existing_issue['labels']]
        required_labels = ["overseer", "area/direct", "priority/medium"]
        missing_labels = [l for l in required_labels if l not in current_labels]
        if missing_labels:
            print(f"Injecting missing labels {missing_labels} for issue #{existing_issue['number']}")
            run_cmd(f'gh issue edit {existing_issue["number"]} --add-label {",".join(missing_labels)}')
        continue
    else:
        if pending_count >= 10:
            print("There are already more than 10 pending issues (open). Skipping creating new ones to avoid overwhelming the team.")
            sys.exit(0)
        
        # We would create an issue here if pending_count < 10
        print(f"Would create issue for {group} {kind}")
        sys.exit(0)

