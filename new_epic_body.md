As part of moving resources from terraform controllers to direct controllers, we want to create a normal CRD using the controller-runtime framework.

These direct CRDs are defined with go types under apis/<service>/v1beta1/<kind>_types.go

When we initially create the direct CRD, we want to make sure that we keep the same schema as the old Terraform CRD, for compatibility and so we can roll back.  Ideally the generated CRD (under `config/crds`) does not change.

Some changes are acceptable:

* Field descriptions can change, particularly for metadata.
* status.observedGeneration will now be an `int64`

Changes to the schema itself, such as a field being added or removed, are not acceptable.


### Subtasks

- [ ] https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/7024