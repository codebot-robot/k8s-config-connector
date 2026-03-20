import json
import subprocess

def run_cmd(cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    return result.stdout.strip()

issues_json = run_cmd('gh issue list --search "Create generate.sh and types.go files for" --state all --json number,title,labels --limit 100')
if issues_json:
    issues = json.loads(issues_json)

    for i in issues:
        title = i["title"]
        labels = {l["name"] for l in i["labels"]}
        required = {"overseer", "area/direct", "priority/medium"}
        missing = required - labels
        if missing and title.startswith("Create generate.sh and types.go files for"):
            num = i["number"]
            print(f"Issue {num} missing labels: {missing}")
            run_cmd(f'gh issue edit {num} --add-label "{",".join(missing)}"')
    print("Label update complete.")
else:
    print("No issues found or error fetching issues.")
