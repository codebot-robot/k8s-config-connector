import yaml
import glob
import subprocess

# Files with dcl2crd: "true"
dcl2crd_files = subprocess.check_output("grep -rl 'cnrm.cloud.google.com/dcl2crd: \"true\"' config/crds/resources", shell=True).decode('utf-8').strip().split('\n')

candidates = []
for file in dcl2crd_files:
    if not file: continue
    
    try:
        with open(file, 'r') as f:
            docs = yaml.safe_load_all(f)
            for doc in docs:
                if not doc: continue
                
                # Check beta versions
                is_beta = False
                versions = doc.get('spec', {}).get('versions', [])
                for v in versions:
                    if 'beta' in v.get('name', ''):
                        is_beta = True
                        break
                        
                if not is_beta:
                    continue
                
                group = doc.get('spec', {}).get('group', '')
                kind = doc.get('spec', {}).get('names', {}).get('kind', '')
                
                if group and kind:
                    short_group = group.split('.')[0]
                    types_go_pattern = f"apis/{short_group}/v1beta1/*_types.go"
                    types_go_files = glob.glob(types_go_pattern)
                    
                    if not types_go_files:
                        candidates.append((group, kind))
    except Exception as e:
        print(f"Error parsing {file}: {e}")

print("--- CANDIDATES ---")
for g, k in sorted(candidates):
    print(f"{g} {k}")
