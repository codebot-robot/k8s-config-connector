for file in $(grep -l 'cnrm.cloud.google.com/dcl2crd: "true"' config/crds/resources/*.yaml); do
  group=$(yq '.spec.group' $file)
  kind=$(yq '.spec.names.kind' $file)
  version=$(yq '.spec.versions[0].name' $file)
  
  if [ "$version" == "v1beta1" ]; then
    group_dir=$(echo $group | awk -F'.' '{print $1}')
    types_file="apis/$group_dir/v1beta1/$(echo $kind | tr '[:upper:]' '[:lower:]')_types.go"
    
    if ! ls apis/$group_dir/v1beta1/*_types.go >/dev/null 2>&1; then
      echo "$group_dir $kind"
    fi
  fi
done
