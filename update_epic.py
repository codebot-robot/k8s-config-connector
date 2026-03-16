import json
import sys
import subprocess

out = subprocess.check_output(['gh', 'issue', 'view', '5954', '--json', 'body']).decode('utf-8')
body = json.loads(out)['body']

issue_link = "https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/7024"

if issue_link not in body:
    if "### Subtasks" not in body:
        body += "\n\n### Subtasks\n"
    body += f"\n- [ ] {issue_link}"

with open("new_epic_body.md", "w") as f:
    f.write(body)
