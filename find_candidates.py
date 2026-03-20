import os
import re
import subprocess
import glob

dcl2crd_files = subprocess.check_output("grep -rl 'cnrm.cloud.google.com/dcl2crd: \"true\"' config/crds/resources", shell=True).decode('utf-8').strip().split('\n')

for file in dcl2crd_files:
    if not file: continue
    
    with open(file, 'r') as f:
        content = f.read()
        
    group_match = re.search(r'group:\s*(.+)', content)
    
    kind_spec_match = re.search(r'names:\n.*kind:\s*([^\s]+)', content, re.MULTILINE | re.DOTALL)
    if not kind_spec_match:
        # Fallback to general kind match if the regex misses
        kind_spec_match = re.search(r'kind:\s*([A-Za-z0-9]+)', content)
        
    if group_match and kind_spec_match:
        group = group_match.group(1).strip()
        kind = kind_spec_match.group(1).strip()
        if kind == "CustomResourceDefinition":
            # Search again specifically under names:
            m = re.search(r'names:\s*\n(?:\s+.*\n)*\s+kind:\s*([A-Za-z0-9]+)', content)
            if m:
                kind = m.group(1)

        short_group = group.split('.')[0]
        
        if 'v1beta1' in content: # rough check for beta version
            types_go_pattern = f"apis/{short_group}/v1beta1/*_types.go"
            types_go_files = glob.glob(types_go_pattern)
            
            if not types_go_files:
                print(f"Candidate: {short_group} {kind} from {file}")
