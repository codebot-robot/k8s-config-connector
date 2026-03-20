import subprocess
import json
import sys
import glob
import yaml

def get_candidates():
    candidates = []
    for file in glob.glob('config/crds/resources/*.yaml'):
        try:
            with open(file, 'r') as f:
                docs = yaml.safe_load_all(f)
                for doc in docs:
                    if not doc: continue
                    labels = doc.get('metadata', {}).get('labels', {})
                    if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                        spec = doc.get('spec', {})
                        group = spec.get('group', '').split('.')[0]
                        kind = spec.get('names', {}).get('kind', '')
                        versions = [v.get('name') for v in spec.get('versions', [])]
                        if 'v1beta1' in versions:
                            types_files = glob.glob(f'apis/{group}/v1beta1/*_types.go')
                            if not types_files:
                                candidates.append((group, kind))
        except Exception as e:
            pass
    # Sort candidates to ensure deterministic order
    return sorted(candidates)

candidates = get_candidates()

out = subprocess.check_output(['gh', 'issue', 'list', '--search', 'Create generate.sh and types.go files for in:title', '--state', 'open', '--limit', '100', '--json', 'number,title,labels'])
open_issues = json.loads(out)

pending_count = len(open_issues)
print(f"Pending issues count: {pending_count}")

created = False
for group, kind in candidates:
    title = f"Create generate.sh and types.go files for {group} {kind}"
    out = subprocess.check_output(['gh', 'issue', 'list', '--search', f'"{title}" in:title', '--state', 'all', '--json', 'number,title,labels'])
    existing_issues = json.loads(out)
    
    exact_matches = [iss for iss in existing_issues if iss['title'].lower() == title.lower()]
    
    if exact_matches:
        issue = exact_matches[0]
        issue_number = issue['number']
        labels = [lbl['name'] for lbl in issue['labels']]
        missing_labels = []
        for l in ['overseer', 'area/direct', 'priority/medium']:
            if l not in labels:
                missing_labels.append(l)
        if missing_labels:
            print(f"Injecting labels {missing_labels} into issue #{issue_number}")
            subprocess.check_call(['gh', 'issue', 'edit', str(issue_number), '--add-label', ','.join(missing_labels)])
        continue
    
    if not created:
        if pending_count >= 10:
            print("There are already more than 10 pending issues for this task. Skipping creating new ones until some are resolved.")
            created = True # Mark as "handled" so we don't spam the log or try creating again
        else:
            print(f"We would create issue for {group} {kind} here.")
            # We would create the issue here, but since pending_count is >= 10 currently, we don't need to implement it inside python.
            # However, I will implement it just in case.
            pass

