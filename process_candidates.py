import json
import subprocess
import re

candidates = [
    ("networkservices", "NetworkServicesGateway"),
    ("networkservices", "NetworkServicesGRPCRoute"),
    ("billingbudgets", "BillingBudgetsBudget"),
    ("dlp", "DLPJobTrigger"),
    ("recaptchaenterprise", "RecaptchaEnterpriseKey"),
    ("cloudfunctions", "CloudFunctionsFunction"),
    ("identityplatform", "IdentityPlatformTenant"),
    ("dlp", "DLPStoredInfoType"),
    ("cloudscheduler", "CloudSchedulerJob"),
    ("identityplatform", "IdentityPlatformConfig"),
    ("eventarc", "EventarcTrigger"),
    ("networkservices", "NetworkServicesTLSRoute"),
    ("binaryauthorization", "BinaryAuthorizationAttestor"),
    ("identityplatform", "IdentityPlatformTenantOAuthIDPConfig"),
    ("networkservices", "NetworkServicesHTTPRoute"),
    ("networkconnectivity", "NetworkConnectivitySpoke"),
    ("dataproc", "DataprocCluster"),
    ("configcontroller", "ConfigControllerInstance"),
    ("identityplatform", "IdentityPlatformOAuthIDPConfig"),
    ("filestore", "FilestoreInstance"),
    ("osconfig", "OSConfigOSPolicyAssignment"),
    ("filestore", "FilestoreBackup"),
    ("networkservices", "NetworkServicesEndpointPolicy"),
    ("dlp", "DLPDeidentifyTemplate"),
    ("networkservices", "NetworkServicesMesh"),
    ("networkservices", "NetworkServicesTCPRoute"),
    ("dataproc", "DataprocWorkflowTemplate"),
    ("datafusion", "DataFusionInstance"),
    ("osconfig", "OSConfigGuestPolicy"),
    ("dlp", "DLPInspectTemplate"),
    ("networkconnectivity", "NetworkConnectivityHub"),
    ("containeranalysis", "ContainerAnalysisNote"),
    ("binaryauthorization", "BinaryAuthorizationPolicy"),
    ("dataproc", "DataprocAutoscalingPolicy")
]

# Fetch all issues related to the task
out = subprocess.check_output(['gh', 'issue', 'list', '--search', 'Create generate.sh and types.go files for', '--state', 'all', '--json', 'number,title,labels,state', '--limit', '100'])
issues = json.loads(out)

required_labels = {"overseer", "area/direct", "priority/medium"}

# check existing issues
open_count = sum(1 for iss in issues if iss['state'] == 'OPEN')
print(f"Total open issues: {open_count}")

for group, kind in candidates:
    # Check if issue exists
    # case insensitive title matching because of variations like Networkconnectivity vs networkconnectivity
    pattern = re.compile(f"Create generate.sh and types.go files for {group} {kind}", re.IGNORECASE)
    matching_issues = [iss for iss in issues if pattern.search(iss['title'])]
    
    if matching_issues:
        for iss in matching_issues:
            labels = {label['name'] for label in iss['labels']}
            missing = required_labels - labels
            if missing:
                print(f"Injecting labels {missing} into issue #{iss['number']} for {group} {kind}")
                # Inject labels
                subprocess.check_call(['gh', 'issue', 'edit', str(iss['number']), '--add-label', ','.join(missing)])
            else:
                print(f"Issue #{iss['number']} for {group} {kind} already has required labels.")
    else:
        # We would create an issue, but there are already > 10 open issues
        pass

