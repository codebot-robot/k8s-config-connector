import os
import glob
import re

crd_dir = 'config/crds/resources'
crd_files = glob.glob(os.path.join(crd_dir, '*.yaml'))

for f in crd_files:
    with open(f, 'r') as fp:
        content = fp.read()
    
    # Check for dcl2crd
    if 'cnrm.cloud.google.com/dcl2crd: "true"' not in content:
        continue
        
    group_match = re.search(r'\bgroup:\s*([^\n]+)', content)
    
    # Kind is in spec.names.kind
    kind_match = re.search(r'names:\n(?:.*\n)*?\s*kind:\s*([^\n]+)', content)
    if not kind_match:
        # try another regex
        kind_match = re.search(r'kind:\s*([^\n]+)', content[content.find('names:'):])
    
    if not group_match or not kind_match:
        continue
    
    group = group_match.group(1).strip().split('.')[0]
    kind = kind_match.group(1).strip()
    
    if 'name: v1beta1' not in content and '- name: v1beta1' not in content:
        continue
        
    type_files = glob.glob(f'apis/{group}/v1beta1/*_types.go')
    if type_files:
        continue
        
    print(f"FOUND: Group={group}, Kind={kind}")

