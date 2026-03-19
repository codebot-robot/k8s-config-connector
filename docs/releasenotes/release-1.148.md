# Run [mdformat](go/mdformat) before publishing this release notes.

** This version is not yet released; this document is gathering release notes
for the future release **

*   Special shout-outs to @anhdle-sso, @barney-s, @cheftako, @gurusai-voleti, @justinsb, @katrielt, @maqiuyujoyce, and @xiaoweim for their contributions to this release.

## New Alpha Resources (Direct Reconciler)

*   `AccessContextManagerServicePerimeter`
    *   Create direct CRD Go type for AccessContextManagerServicePerimeter. ([PR 6970](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6970))

## Reconciliation Improvements

*   [PR 6774](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6774)
    Add structured reporting diff to BigQueryAnalyticsHubDataExchange
*   [PR 7083](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7083)
    Support alternative controller comparison in preview mode

## New features

*   [PR 6671](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6671)
    Add documentation for enabling VerticalPodAutoscaler in Config Connector
*   [PR 7075](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7075)
    Introduce skip-name-validation flag and consolidate tests

## Bug Fixes

*   [PR 6943](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6943)
    Handle ALREADY_EXISTS in TagKey and TagValue controllers
*   [PR 7106](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7106)
    Fix typo in container.yaml service mapping file
*   [PR 7115](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7115)
    Restore missing descriptions in CloudBuildTrigger CRD
