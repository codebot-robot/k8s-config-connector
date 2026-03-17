*   Special shout-outs to @acpana, @anhdle-sso, @app/dependabot, @barney-s, @cheftako, @codebot-robot, @dhavalbhensdadiya-crest, @fkc1e100, @gemmahou, @justinsb, @katrielt, @maqiuyujoyce, @xiaoweim for their contributions to this release.

## New Fields

*   [`CloudDeployTarget`](https://cloud.google.com/config-connector/docs/reference/resource-docs/clouddeploy/clouddeploytarget)
    *   Updated `TargetIDs` field to `Ref`.

## Reconciliation Improvements

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

## Bug Fixes

*   Fixed a typo in `StorageDefaultObjectAccessControl` documentation.
*   Fixed use of `latest` in generation scripts.