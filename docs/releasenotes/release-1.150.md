*   Special shout-outs to anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, gurusai-voleti, justinsb, katrielt, maqiuyujoyce, xiaoweim for their contributions to this release.

## Reconciliation Improvements

We have added support for direct reconciliation to more resources, with opt-in behaviour. The API is unchanged. To use the direct reconciler, add the `alpha.cnrm.cloud.google.com/reconciler: direct` annotation to the corresponding Config Connector object. The following resources now have direct reconciliation support:

*   `AccessContextManagerServicePerimeter`
*   `AlloyDBUser`
*   `ComputeResourcePolicy`

Other reconciliation improvements include:

*   Allow specific integer type changes in CRD equivalence checks to prevent unnecessary updates.
*   Add structured reporting diff to BigQueryAnalyticsHubDataExchange.
*   Introduce `skip-name-validation` flag for resource name validation.
*   Support alternative controller comparison in preview mode.

## New features:

*   Add documentation and support for enabling VerticalPodAutoscaler in Config Connector.

## Bug Fixes:

*   Restore missing descriptions in CloudBuildTrigger CRD.
*   Fix typo in container.yaml service mapping file.
*   Fix `ALREADY_EXISTS` error handling in TagKey and TagValue controllers.
