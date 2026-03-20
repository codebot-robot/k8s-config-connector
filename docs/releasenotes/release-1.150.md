** This version is not yet released; this document is gathering release notes
for the future release **

*   Special shout-outs to @anhdle-sso, @barney-s, @cheftako, @gurusai-voleti, @justinsb, @katrielt, @maqiuyujoyce, and @xiaoweim for their contributions to this release.

## New Alpha Resources (Direct Reconciler):

*   `AccessContextManagerServicePerimeter`
    *   Manage the `AccessContextManagerServicePerimeter` resource using the direct reconciler.
*   `AlloyDBUser`
    *   Manage the `AlloyDBUser` resource using the direct reconciler.

## New features:

*   Support alternative controller comparison in preview mode.
*   Add structured reporting diff to `BigQueryAnalyticsHubDataExchange`.

## Bug Fixes:

*   Restore missing descriptions in `CloudBuildTrigger` CRD.
*   Fix typo in `container.yaml` service mapping file.
*   Handle `ALREADY_EXISTS` error appropriately in `TagKey` and `TagValue` controllers.
