import os
import json
import subprocess

# get candidates from find_resources.py
output = subprocess.check_output(['python3', 'find_resources.py'])
candidates_text = output.decode('utf-8').strip().split('\n')
candidates = [line.split(' ') for line in candidates_text if line]

# get all existing issues
output = subprocess.check_output(['gh', 'issue', 'list', '--search', 'Create generate.sh and types.go files for in:title', '--state', 'all', '--json', 'number,title,state,labels', '--limit', '200'])
issues = json.loads(output.decode('utf-8'))

# count open issues for this task
open_issues_count = len([i for i in issues if i['state'] == 'OPEN'])

issues_by_title = {i['title'].lower(): i for i in issues}

target_labels = {"overseer", "area/direct", "priority/medium"}

issue_created = False
printed_skip_msg = False

for group, kind in candidates:
    title = f"Create generate.sh and types.go files for {group} {kind}"
    lower_title = title.lower()
    
    if lower_title in issues_by_title:
        issue = issues_by_title[lower_title]
        existing_labels = {l['name'] for l in issue['labels']}
        missing_labels = target_labels - existing_labels
        if missing_labels:
            print(f"Injecting labels {missing_labels} into issue #{issue['number']} for {group} {kind}")
            subprocess.call(['gh', 'issue', 'edit', str(issue['number']), '--add-label', ','.join(missing_labels)])
        continue
    
    if not issue_created:
        if open_issues_count >= 10:
            if not printed_skip_msg:
                print(f"There are already 10 pending issues ({open_issues_count} actually). Skipping creating new ones until some of the existing issues are resolved.")
                printed_skip_msg = True
        else:
            print(f"Creating issue for {group} {kind}...")
            # Here we would normally create it
            issue_created = True
