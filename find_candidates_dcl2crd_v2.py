import os
import yaml
import glob

crd_dir = "config/crds/resources"
candidates = []

for file in glob.glob(f"{crd_dir}/*.yaml"):
    with open(file, 'r') as f:
        try:
            docs = list(yaml.safe_load_all(f))
            for doc in docs:
                if not doc: continue
                labels = doc.get('metadata', {}).get('labels', {})
                if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                    # Check if v1beta1 is in versions
                    versions = doc.get('spec', {}).get('versions', [])
                    is_beta = any(v.get('name') == 'v1beta1' for v in versions)
                    if is_beta:
                        group = doc.get('spec', {}).get('group', '').split('.')[0]
                        kind = doc.get('spec', {}).get('names', {}).get('kind', '')
                        
                        # Check if apis/<group>/v1beta1/*_types.go exists
                        types_pattern = f"apis/{group}/v1beta1/*_types.go"
                        if not glob.glob(types_pattern):
                            candidates.append(f"{group} {kind}")
        except Exception as e:
            pass

for c in sorted(set(candidates)):
    print(c)