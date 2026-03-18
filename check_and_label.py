import json
import subprocess

with open('existing_issues.json') as f:
    issues = json.load(f)

# we found these candidates before
candidates = [
    ("dataproc", "DataprocCluster"),
    ("osconfig", "OSConfigOSPolicyAssignment"),
    ("networkconnectivity", "NetworkConnectivitySpoke"),
    ("dlp", "DLPJobTrigger"),
    ("cloudscheduler", "CloudSchedulerJob"),
    ("networkservices", "NetworkServicesTLSRoute"),
    ("dlp", "DLPDeidentifyTemplate"),
    ("identityplatform", "IdentityPlatformTenantOAuthIDPConfig"),
    ("dataproc", "DataprocAutoscalingPolicy"),
    ("networkconnectivity", "NetworkConnectivityHub"),
    ("dlp", "DLPInspectTemplate"),
    ("cloudfunctions", "CloudFunctionsFunction"),
    ("dlp", "DLPStoredInfoType"),
    ("filestore", "FilestoreBackup"),
    ("recaptchaenterprise", "RecaptchaEnterpriseKey"),
    ("identityplatform", "IdentityPlatformOAuthIDPConfig"),
    ("binaryauthorization", "BinaryAuthorizationAttestor"),
    ("osconfig", "OSConfigGuestPolicy"),
    ("binaryauthorization", "BinaryAuthorizationPolicy"),
    ("eventarc", "EventarcTrigger"),
    ("identityplatform", "IdentityPlatformTenant"),
    ("datafusion", "DataFusionInstance"),
    ("networkservices", "NetworkServicesEndpointPolicy"),
    ("billingbudgets", "BillingBudgetsBudget"),
    ("networkservices", "NetworkServicesHTTPRoute"),
    ("identityplatform", "IdentityPlatformConfig"),
    ("containeranalysis", "ContainerAnalysisNote"),
    ("networkservices", "NetworkServicesMesh"),
    ("networkservices", "NetworkServicesTCPRoute"),
    ("configcontroller", "ConfigControllerInstance"),
    ("filestore", "FilestoreInstance"),
    ("dataproc", "DataprocWorkflowTemplate"),
    ("networkservices", "NetworkServicesGateway"),
    ("networkservices", "NetworkServicesGRPCRoute")
]

for group, kind in candidates:
    # check if issue exists in the list
    for issue in issues:
        title = issue['title'].lower()
        if group.lower() in title and kind.lower() in title:
            # issue exists! check labels
            labels = [l['name'] for l in issue['labels']]
            missing = []
            for req in ['overseer', 'area/direct', 'priority/medium']:
                if req not in labels:
                    missing.append(req)
            if missing:
                print(f"Adding labels {missing} to issue {issue['number']}")
                subprocess.run(["gh", "issue", "edit", str(issue['number']), "--add-label", ",".join(missing)])
            
