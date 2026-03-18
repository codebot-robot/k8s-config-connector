import yaml
import glob
import os

candidate_files = glob.glob("config/crds/resources/*.yaml")
for file in candidate_files:
    with open(file, 'r') as f:
        try:
            data = yaml.safe_load(f)
        except:
            continue
            
    labels = data.get('metadata', {}).get('labels', {})
    if labels.get('cnrm.cloud.google.com/dcl2crd') == "true":
        versions = data.get('spec', {}).get('versions', [])
        for v in versions:
            if v.get('name') == 'v1beta1':
                group = data.get('spec', {}).get('group', '').split('.')[0]
                kind = data.get('spec', {}).get('names', {}).get('kind', '')
                
                # Check if types.go exists
                types_files = glob.glob(f"apis/{group}/v1beta1/*_types.go")
                if not types_files:
                    print(f"Candidate: Group={group}, Kind={kind}")
                break

