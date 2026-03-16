import os
import yaml

crd_dir = 'config/crds/resources'
for filename in os.listdir(crd_dir):
    if not filename.endswith('.yaml'): continue
    filepath = os.path.join(crd_dir, filename)
    with open(filepath, 'r') as f:
        try:
            doc = yaml.safe_load(f)
        except Exception:
            continue
        labels = doc.get('metadata', {}).get('labels', {})
        if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
            group = doc.get('spec', {}).get('group', '')
            kind = doc.get('spec', {}).get('names', {}).get('kind', '')
            
            # Find beta version
            versions = doc.get('spec', {}).get('versions', [])
            has_beta = False
            for v in versions:
                if 'beta' in v.get('name', ''):
                    has_beta = True
                    break
            
            if has_beta:
                group_short = group.split('.')[0]
                # Check if types.go exists
                types_dir = f"apis/{group_short}/v1beta1"
                if not os.path.exists(types_dir):
                    print(f"Candidate: Group={group_short}, Kind={kind}, File={filename} (Dir {types_dir} does not exist)")
                else:
                    types_files = [f for f in os.listdir(types_dir) if f.endswith('_types.go')]
                    # Since kind can have different casing, let's just check if there's any *_types.go for that kind, 
                    # actually the prompt says "the resource should not have a types.go file generated yet"
                    # We can check if `kind.lower() + "_types.go"` exists or similar
                    found = False
                    for f in types_files:
                        if kind.lower() in f.lower() or f.lower().startswith(kind.lower()):
                            found = True
                    if not found:
                        print(f"Candidate: Group={group_short}, Kind={kind}, File={filename} (No types.go found in {types_dir})")

