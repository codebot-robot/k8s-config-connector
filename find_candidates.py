import yaml
import glob
import os

candidates = []
for file in glob.glob("config/crds/resources/*.yaml"):
    try:
        with open(file, 'r') as f:
            docs = list(yaml.safe_load_all(f))
            for doc in docs:
                if not doc: continue
                labels = doc.get('metadata', {}).get('labels', {})
                if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                    versions = doc.get('spec', {}).get('versions', [])
                    for v in versions:
                        if v.get('name') == 'v1beta1':
                            group = doc.get('spec', {}).get('group', '').split('.')[0]
                            kind = doc.get('spec', {}).get('names', {}).get('kind', '')
                            
                            types_pattern = f"apis/{group}/v1beta1/*_types.go"
                            if not glob.glob(types_pattern):
                                candidates.append((group, kind))
    except Exception as e:
        pass

for group, kind in sorted(list(set(candidates))):
    print(f"{group} {kind}")
