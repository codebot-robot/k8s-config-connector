import os
import yaml
import glob

crd_dir = "config/crds/resources"
candidates = []
for filename in glob.glob(os.path.join(crd_dir, "*.yaml")):
    with open(filename, 'r') as f:
        try:
            docs = list(yaml.safe_load_all(f))
            for doc in docs:
                if not doc: continue
                labels = doc.get('metadata', {}).get('labels', {})
                if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                    versions = doc.get('spec', {}).get('versions', [])
                    if any(v.get('name') == 'v1beta1' for v in versions):
                        group = doc.get('spec', {}).get('group', '').split('.')[0]
                        kind = doc.get('spec', {}).get('names', {}).get('kind', '')
                        
                        types_pattern = f"apis/{group}/v1beta1/*_types.go"
                        if not glob.glob(types_pattern):
                            candidates.append((group, kind))
        except Exception as e:
            pass

for group, kind in sorted(candidates):
    print(f"{group} {kind}")
