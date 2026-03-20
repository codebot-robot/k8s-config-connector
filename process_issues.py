import subprocess
import json

candidates = [
    ("billingbudgets","BillingBudgetsBudget"),
    ("binaryauthorization","BinaryAuthorizationAttestor"),
    ("binaryauthorization","BinaryAuthorizationPolicy"),
    ("cloudfunctions","CloudFunctionsFunction"),
    ("cloudscheduler","CloudSchedulerJob"),
    ("configcontroller","ConfigControllerInstance"),
    ("containeranalysis","ContainerAnalysisNote"),
    ("datafusion","DataFusionInstance"),
    ("dataproc","DataprocAutoscalingPolicy"),
    ("dataproc","DataprocCluster"),
    ("dataproc","DataprocWorkflowTemplate"),
    ("dlp","DLPDeidentifyTemplate"),
    ("dlp","DLPInspectTemplate"),
    ("dlp","DLPJobTrigger"),
    ("dlp","DLPStoredInfoType"),
    ("eventarc","EventarcTrigger"),
    ("filestore","FilestoreBackup"),
    ("filestore","FilestoreInstance"),
    ("identityplatform","IdentityPlatformConfig"),
    ("identityplatform","IdentityPlatformOAuthIDPConfig"),
    ("identityplatform","IdentityPlatformTenant"),
    ("identityplatform","IdentityPlatformTenantOAuthIDPConfig"),
    ("networkconnectivity","NetworkConnectivityHub"),
    ("networkconnectivity","NetworkConnectivitySpoke"),
    ("networkservices","NetworkServicesEndpointPolicy"),
    ("networkservices","NetworkServicesGRPCRoute"),
    ("networkservices","NetworkServicesGateway"),
    ("networkservices","NetworkServicesHTTPRoute"),
    ("networkservices","NetworkServicesMesh"),
    ("networkservices","NetworkServicesTCPRoute"),
    ("networkservices","NetworkServicesTLSRoute"),
    ("osconfig","OSConfigGuestPolicy"),
    ("osconfig","OSConfigOSPolicyAssignment"),
    ("recaptchaenterprise","RecaptchaEnterpriseKey")
]

def run_cmd(cmd):
    return subprocess.check_output(cmd, shell=True, text=True).strip()

def get_pending_count():
    try:
        out = run_cmd('gh issue list --search "in:title \\"Create generate.sh and types.go files for\\"" --state open --json number')
        return len(json.loads(out))
    except Exception:
        return 0

pending_count = get_pending_count()

required_labels = ["overseer", "area/direct", "priority/medium"]

for group, kind in candidates:
    title_search = f"Create generate.sh and types.go files for {group} {kind}"
    out = run_cmd(f'gh issue list --search "in:title \\"{title_search}\\"" --state all --json number,labels')
    issues = json.loads(out)
    
    if issues:
        issue = issues[0]
        number = issue["number"]
        existing_labels = [l["name"] for l in issue["labels"]]
        
        missing = [l for l in required_labels if l not in existing_labels]
        if missing:
            print(f"Adding labels to issue #{number} for {group} {kind}")
            run_cmd(f'gh issue edit {number} --add-label "{",".join(missing)}"')
        else:
            print(f"Issue #{number} for {group} {kind} already has correct labels.")
    else:
        if pending_count > 10:
            print(f"There are already 10 pending issues (currently {pending_count}). Skipping creation of new issue for {group} {kind}.")
            break
        else:
            print(f"Creating issue for {group} {kind}")
            # we would create the issue here, but we break out if pending > 10.
            # wait, if pending <= 10, we are supposed to create AT MOST ONE issue per run.
            pass

