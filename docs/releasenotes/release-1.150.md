*   Special shout-outs to @anhdle-sso, @cheftako, @gurusai-voleti, @justinsb, @katrielt, @maqiuyujoyce, @xiaoweim for their contributions to this release.

## New Beta Resources (Direct Reconciler):

*   `ComputeResourcePolicy`
*   `AlloyDBUser`
*   `AccessContextManagerServicePerimeter`

## New features:

*   Added documentation for enabling VerticalPodAutoscaler in Config Connector.

## Reconciliation Improvements

*   `BigQueryAnalyticsHubDataExchange`: Added structured reporting diff to improve change visibility.

## Bug Fixes:

*   `CloudBuildTrigger`: Restored missing descriptions in the CRD.
*   `TagKey` and `TagValue`: Handled `ALREADY_EXISTS` gracefully during resource creation.
*   Fixed a typo in the `container` service mapping file.
