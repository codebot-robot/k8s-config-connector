import yaml
import os
import glob

crd_dir = 'config/crds/resources'
crd_files = glob.glob(os.path.join(crd_dir, '*.yaml'))

candidates = []

for fpath in crd_files:
    with open(fpath, 'r') as f:
        try:
            docs = list(yaml.safe_load_all(f))
            for doc in docs:
                if not doc: continue
                labels = doc.get('metadata', {}).get('labels', {})
                if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                    group = doc.get('spec', {}).get('group', '')
                    short_group = group.split('.')[0]
                    kind = doc.get('spec', {}).get('names', {}).get('kind', '')
                    
                    versions = doc.get('spec', {}).get('versions', [])
                    has_v1beta1 = False
                    for v in versions:
                        if v.get('name') == 'v1beta1':
                            has_v1beta1 = True
                    
                    if has_v1beta1:
                        # Check if types.go exists
                        # Wait, the prompt says: check if *any* _types.go for this specific kind?
                        # No, the prompt says: "And the resource should not have a types.go file generated yet: ls -al apis/<GROUP>/v1beta1/*_types.go"
                        # Wait, usually it's apis/<GROUP>/v1beta1/<kind.lower()>_types.go.
                        # Let's check for <kind.lower()>_types.go
                        types_file = f"apis/{short_group}/v1beta1/{kind.lower()}_types.go"
                        if not os.path.exists(types_file):
                            candidates.append((short_group, kind))
        except Exception as e:
            pass

candidates.sort()
for c in candidates:
    print(f"{c[0]} {c[1]}")
