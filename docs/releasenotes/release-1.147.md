*   Special shout-outs to @acpana, @anhdle-sso, @app/dependabot, @barney-s, @cheftako, @codebot-robot, @dhavalbhensdadiya-crest, @fkc1e100, @gemmahou, @justinsb, @katrielt, @maqiuyujoyce, @xiaoweim for their contributions to this release.

## Announcement

*   The `DeployCustomTargetType` (v1alpha1) resource is no longer supported. It has been replaced by the new `CloudDeployCustomTargetType` (v1alpha1) resource. Please remove any instances of `DeployCustomTargetType` resource.

## New Fields

*   [`CloudDeployTarget`](https://cloud.google.com/config-connector/docs/reference/resource-docs/clouddeploy/clouddeploytarget)
    *   Updated `TargetIDs` field to `Ref`.

## Reconciliation Improvements

We have added support for direct reconciliation to more resources, with opt-in behaviour. The API is unchanged. To use the direct reconciler, add the `alpha.cnrm.cloud.google.com/reconciler: direct` annotation to the corresponding Config Connector object. The following resources now have direct reconciliation support:

*   [`CloudBuildTrigger`](https://cloud.google.com/config-connector/docs/reference/resource-docs/cloudbuild/cloudbuildtrigger)

*   Added structured reporting diff to provide more reliable and transparent resource status updates during reconciliation for the following resources:
    *   `BigQueryDataset`
    *   `BigQueryReservationAssignment`
    *   `BigQueryTable`
    *   `CertificateManagerDNSAuthorization`
    *   `CloudIdentityGroup`
    *   `DataformRepository`
    *   `MetastoreService`
    *   `PrivilegedAccessManagerEntitlement`
    *   `VMwareEngineExternalAccessRule`
    *   `WorkflowsWorkflow`

## New features

*   Added `reconcilerOverride` support to the `preview` CLI. Users can configure the `ConfigConnectorContext` resources to be dynamically modified during the preview interception process with specific controller overrides (e.g., forcing `direct` or `dcl` mode).

## Bug Fixes

*   Fixed a typo in `StorageDefaultObjectAccessControl` documentation.
*   Fixed use of `latest` in generation scripts.
