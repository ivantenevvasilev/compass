BEGIN;

-- all views below should be dropped because they depend on the 'short_description' column that will be later modified in the respective tables
DROP VIEW tenants_packages;
DROP VIEW tenants_products;
DROP VIEW tenants_bundles;
DROP VIEW tenants_specifications;
DROP VIEW tenants_apis;
DROP VIEW tenants_events;

DROP VIEW packages_tenants;
DROP VIEW products_tenants;
DROP VIEW bundles_tenants;
DROP VIEW api_definitions_tenants;
DROP VIEW event_api_definitions_tenants;

ALTER TABLE packages
ALTER COLUMN short_description TYPE VARCHAR(255);

ALTER TABLE bundles
ALTER COLUMN short_description TYPE VARCHAR(255);

ALTER TABLE products
ALTER COLUMN short_description TYPE VARCHAR(255);

ALTER TABLE api_definitions
ALTER COLUMN short_description TYPE VARCHAR(255);

ALTER TABLE event_api_definitions
ALTER COLUMN short_description TYPE VARCHAR(255);


-- while re-creating this view, the 'support_info' column had to be explicitly omitted (thus selecting the rest columns one by one). This is needed because when it was introduced in
-- "add-support-info-to-packages.up" migration, this view was not adapted to include it. In the current "up" migration, the 'support_info' column is included
-- because all package columns are selected there (p.*). If we do not exclude the "support_info" column here, the "add-support-info-to-packages.down"
-- migration will fail because when trying to drop the column this view will already rely on it.
CREATE OR REPLACE VIEW packages_tenants
    (id, ord_id, title, short_description, description, version, package_links,
     links, licence_type, tags, countries, labels, policy_level, app_id, custom_policy_level, vendor,
     part_of_products, line_of_business, industry, resource_hash, tenant_id, owner)
AS
SELECT  p.id,
        p.ord_id,
        p.title,
        p.short_description,
        p.description,
        p.version,
        p.package_links,
        p.links,
        p.licence_type,
        p.tags,
        p.countries,
        p.labels,
        p.policy_level,
        p.app_id,
        p.custom_policy_level,
        p.vendor,
        p.part_of_products,
        p.line_of_business,
        p.industry,
        p.resource_hash,
        ta.tenant_id, ta.owner FROM packages AS p
                                            INNER JOIN tenant_applications AS ta ON ta.id = p.app_id;

-- while re-creating this view, the 'documentation_labels' column had to be explicitly omitted (thus selecting the rest columns one by one). This is needed because when it was introduced in
-- "ord_documentation_labels.up" migration, this view was not adapted to include it. In the current "up" migration, the 'documentation_labels' column is included
-- because all product columns are selected there (p.*). If we do not exclude the "documentation_labels" column here, the "ord_documentation_labels.down"
-- migration will fail because when trying to drop the column this view will already rely on it.
CREATE OR REPLACE VIEW products_tenants (ord_id, app_id, title, short_description, vendor, parent, labels,
                                         correlation_ids, id, tenant_id, owner)
AS
SELECT p.ord_id,
       p.app_id,
       p.title,
       p.short_description,
       p.vendor,
       p.parent,
       p.labels,
       p.correlation_ids,
       p.id, ta.tenant_id, ta.owner FROM products AS p
                                            INNER JOIN tenant_applications AS ta ON ta.id = p.app_id;


-- while re-creating this view, the 'documentation_labels' column had to be explicitly omitted (thus selecting the rest columns one by one). This is needed because when it was introduced in
-- "ord_documentation_labels.up" migration, this view was not adapted to include it. In the current "up" migration, the 'documentation_labels' column is included
-- because all bundle columns are selected there (b.*). If we do not exclude the "documentation_labels" column here, the "ord_documentation_labels.down"
-- migration will fail because when trying to drop the column this view will already rely on it.
CREATE OR REPLACE VIEW bundles_tenants
        (id, app_id, name, description, instance_auth_request_json_schema,
         default_instance_auth, ord_id, short_description, links, labels, credential_exchange_strategies, ready,
         created_at, updated_at, deleted_at, error, correlation_ids, tenant_id, owner)
AS
SELECT  b.id,
        b.app_id,
        b.name,
        b.description,
        b.instance_auth_request_json_schema,
        b.default_instance_auth,
        b.ord_id,
        b.short_description,
        b.links,
        b.labels,
        b.credential_exchange_strategies,
        b.ready,
        b.created_at,
        b.updated_at,
        b.deleted_at,
        b.error,
        b.correlation_ids, ta.tenant_id, ta.owner FROM bundles AS b
                                            INNER JOIN tenant_applications ta ON ta.id = b.app_id;

-- while re-creating this view, the 'documentation_labels' column had to be explicitly omitted (thus selecting the rest columns one by one). This is needed because when it was introduced in
-- "ord_documentation_labels.up" migration, this view was not adapted to include it. In the current "up" migration, the 'documentation_labels' column is included
-- because all api columns are selected there (ad.*). If we do not exclude the "documentation_labels" column here, the "ord_documentation_labels.down"
-- migration will fail because when trying to drop the column this view will already rely on it.
CREATE OR REPLACE VIEW api_definitions_tenants
    (id, app_id, name, description, group_name, default_auth, version_value,
     version_deprecated, version_deprecated_since, version_for_removal, ord_id, short_description,
     system_instance_aware, api_protocol, tags, countries, links, api_resource_links, release_status,
     sunset_date, changelog_entries, labels, package_id, visibility, disabled, part_of_products,
     line_of_business, industry, ready, created_at, updated_at, deleted_at, error, implementation_standard,
     custom_implementation_standard, custom_implementation_standard_description, target_urls, extensible,
     successors, resource_hash, tenant_id, owner)
AS
SELECT ad.id,
       ad.app_id,
       ad.name,
       ad.description,
       ad.group_name,
       ad.default_auth,
       ad.version_value,
       ad.version_deprecated,
       ad.version_deprecated_since,
       ad.version_for_removal,
       ad.ord_id,
       ad.short_description,
       ad.system_instance_aware,
       ad.api_protocol,
       ad.tags,
       ad.countries,
       ad.links,
       ad.api_resource_links,
       ad.release_status,
       ad.sunset_date,
       ad.changelog_entries,
       ad.labels,
       ad.package_id,
       ad.visibility,
       ad.disabled,
       ad.part_of_products,
       ad.line_of_business,
       ad.industry,
       ad.ready,
       ad.created_at,
       ad.updated_at,
       ad.deleted_at,
       ad.error,
       ad.implementation_standard,
       ad.custom_implementation_standard,
       ad.custom_implementation_standard_description,
       ad.target_urls,
       ad.extensible,
       ad.successors,
       ad.resource_hash, ta.tenant_id, ta.owner FROM api_definitions AS ad
                                             INNER JOIN tenant_applications ta ON ta.id = ad.app_id;

-- while re-creating this view, the 'documentation_labels' column had to be explicitly omitted (thus selecting the rest columns one by one). This is needed because when it was introduced in
-- "ord_documentation_labels.up" migration, this view was not adapted to include it. In the current "up" migration, the 'documentation_labels' column is included
-- because all event columns are selected there (e.*). If we do not exclude the "documentation_labels" column here, the "ord_documentation_labels.down"
-- migration will fail because when trying to drop the column this view will already rely on it.
CREATE OR REPLACE VIEW event_api_definitions_tenants
    (id, app_id, name, description, group_name, version_value,
     version_deprecated, version_deprecated_since, version_for_removal, ord_id, short_description,
     system_instance_aware, changelog_entries, links, tags, countries, release_status, sunset_date, labels,
     package_id, visibility, disabled, part_of_products, line_of_business, industry, ready, created_at,
     updated_at, deleted_at, error, extensible, successors, resource_hash, tenant_id, owner)
AS
SELECT  e.id,
        e.app_id,
        e.name,
        e.description,
        e.group_name,
        e.version_value,
        e.version_deprecated,
        e.version_deprecated_since,
        e.version_for_removal,
        e.ord_id,
        e.short_description,
        e.system_instance_aware,
        e.changelog_entries,
        e.links,
        e.tags,
        e.countries,
        e.release_status,
        e.sunset_date,
        e.labels,
        e.package_id,
        e.visibility,
        e.disabled,
        e.part_of_products,
        e.line_of_business,
        e.industry,
        e.ready,
        e.created_at,
        e.updated_at,
        e.deleted_at,
        e.error,
        e.extensible,
        e.successors,
        e.resource_hash, ta.tenant_id, ta.owner FROM event_api_definitions AS e
                                            INNER JOIN tenant_applications ta ON ta.id = e.app_id;

-----

-- no changes to the rest of the views below - just re-creating them
CREATE OR REPLACE VIEW tenants_packages
            (tenant_id, provider_tenant_id, id, ord_id, title, short_description, description, version, package_links,
             links, licence_type, tags, countries, labels, policy_level, app_id, custom_policy_level, vendor,
             part_of_products, line_of_business, industry, resource_hash, support_info)
AS
SELECT DISTINCT t_apps.tenant_id,
                t_apps.provider_tenant_id,
                p.id,
                p.ord_id,
                p.title,
                p.short_description,
                p.description,
                p.version,
                p.package_links,
                p.links,
                p.licence_type,
                p.tags,
                p.countries,
                p.labels,
                p.policy_level,
                p.app_id,
                p.custom_policy_level,
                p.vendor,
                p.part_of_products,
                p.line_of_business,
                p.industry,
                p.resource_hash,
                p.support_info
FROM packages p
         JOIN (SELECT a1.id,
                      a1.tenant_id::text AS tenant_id,
                      a1.tenant_id::text AS provider_tenant_id
               FROM tenant_applications a1
               UNION ALL
               SELECT apps_subaccounts_func.id,
                      apps_subaccounts_func.tenant_id::text,
                      apps_subaccounts_func.provider_tenant_id::text
               FROM apps_subaccounts_func() apps_subaccounts_func(id, tenant_id, provider_tenant_id)
               UNION ALL
               SELECT ta.id AS app_id, ta.tenant_id::text AS consumer_tenant, tenant_runtimes.tenant_id::text AS provider_tenant
               FROM (SELECT labels.runtime_id, v ->> 0 AS consumer_tenant
                     FROM labels
                              JOIN jsonb_array_elements(labels.value) AS v ON TRUE
                     WHERE key = 'consumer_subaccount_ids') AS t_rts -- Get runtime and external consumer IDs pairs
                        JOIN business_tenant_mappings t ON t_rts.consumer_tenant = t.external_tenant -- Get runtime and internal consumer IDs pairs
                        JOIN apps_subaccounts_func() ta ON t.id = ta.tenant_id -- Get applications for consumer tenants
                        JOIN tenant_runtimes ON t_rts.runtime_id = tenant_runtimes.id) t_apps
              ON p.app_id = t_apps.id;

-----

CREATE OR REPLACE VIEW tenants_products
            (tenant_id, provider_tenant_id, ord_id, app_id, title, short_description, vendor, parent, labels,
             correlation_ids, id, documentation_labels)
AS
SELECT DISTINCT t_apps.tenant_id,
                t_apps.provider_tenant_id,
                p.ord_id,
                p.app_id,
                p.title,
                p.short_description,
                p.vendor,
                p.parent,
                p.labels,
                p.correlation_ids,
                p.id,
                p.documentation_labels
FROM products p
         JOIN (SELECT a1.id,
                      a1.tenant_id::text AS tenant_id,
                      a1.tenant_id::text AS provider_tenant_id
               FROM tenant_applications a1
               UNION ALL
               SELECT apps_subaccounts_func.id,
                      apps_subaccounts_func.tenant_id::text,
                      apps_subaccounts_func.provider_tenant_id::text
               FROM apps_subaccounts_func() apps_subaccounts_func(id, tenant_id, provider_tenant_id)
               UNION ALL
               SELECT ta.id AS app_id, ta.tenant_id::text AS consumer_tenant, tenant_runtimes.tenant_id::text AS provider_tenant
               FROM (SELECT labels.runtime_id, v ->> 0 AS consumer_tenant
                     FROM labels
                              JOIN jsonb_array_elements(labels.value) AS v ON TRUE
                     WHERE key = 'consumer_subaccount_ids') AS t_rts -- Get runtime and external consumer IDs pairs
                        JOIN business_tenant_mappings t ON t_rts.consumer_tenant = t.external_tenant -- Get runtime and internal consumer IDs pairs
                        JOIN apps_subaccounts_func() ta ON t.id = ta.tenant_id -- Get applications for consumer tenants
                        JOIN tenant_runtimes ON t_rts.runtime_id = tenant_runtimes.id) t_apps
              ON p.app_id = t_apps.id OR p.app_id IS NULL;

-----

CREATE OR REPLACE VIEW tenants_bundles
            (tenant_id, provider_tenant_id, id, app_id, name, description, instance_auth_request_json_schema,
             default_instance_auth, ord_id, short_description, links, labels, credential_exchange_strategies, ready,
             created_at, updated_at, deleted_at, error, correlation_ids)
AS
SELECT DISTINCT t_apps.tenant_id,
                t_apps.provider_tenant_id,
                b.id,
                b.app_id,
                b.name,
                b.description,
                b.instance_auth_request_json_schema,
                b.default_instance_auth,
                b.ord_id,
                b.short_description,
                b.links,
                b.labels,
                b.credential_exchange_strategies,
                b.ready,
                b.created_at,
                b.updated_at,
                b.deleted_at,
                b.error,
                b.correlation_ids
FROM bundles b
         JOIN (SELECT a1.id,
                      a1.tenant_id::text AS tenant_id,
                      a1.tenant_id::text AS provider_tenant_id
               FROM tenant_applications a1
               UNION ALL
               SELECT apps_subaccounts_func.id,
                      apps_subaccounts_func.tenant_id::text,
                      apps_subaccounts_func.provider_tenant_id::text
               FROM apps_subaccounts_func() apps_subaccounts_func(id, tenant_id, provider_tenant_id)
               UNION ALL
               SELECT ta.id AS app_id, ta.tenant_id::text AS consumer_tenant, tenant_runtimes.tenant_id::text AS provider_tenant
               FROM (SELECT labels.runtime_id, v ->> 0 AS consumer_tenant
                     FROM labels
                              JOIN jsonb_array_elements(labels.value) AS v ON TRUE
                     WHERE key = 'consumer_subaccount_ids') AS t_rts -- Get runtime and external consumer IDs pairs
                        JOIN business_tenant_mappings t ON t_rts.consumer_tenant = t.external_tenant -- Get runtime and internal consumer IDs pairs
                        JOIN apps_subaccounts_func() ta ON t.id = ta.tenant_id -- Get applications for consumer tenants
                        JOIN tenant_runtimes ON t_rts.runtime_id = tenant_runtimes.id) t_apps
              ON b.app_id = t_apps.id;

-----

CREATE OR REPLACE VIEW tenants_apis
            (tenant_id, provider_tenant_id, id, app_id, name, description, group_name, default_auth, version_value,
             version_deprecated, version_deprecated_since, version_for_removal, ord_id, short_description,
             system_instance_aware, api_protocol, tags, countries, links, api_resource_links, release_status,
             sunset_date, changelog_entries, labels, package_id, visibility, disabled, part_of_products,
             line_of_business, industry, ready, created_at, updated_at, deleted_at, error, implementation_standard,
             custom_implementation_standard, custom_implementation_standard_description, target_urls, extensible,
             successors, resource_hash, documentation_labels)
AS
SELECT DISTINCT t_apps.tenant_id,
                t_apps.provider_tenant_id,
                apis.id,
                apis.app_id,
                apis.name,
                apis.description,
                apis.group_name,
                apis.default_auth,
                apis.version_value,
                apis.version_deprecated,
                apis.version_deprecated_since,
                apis.version_for_removal,
                apis.ord_id,
                apis.short_description,
                apis.system_instance_aware,
                CASE
                    WHEN apis.api_protocol IS NULL AND specs.api_spec_type::text = 'ODATA'::text THEN 'odata-v2'::text
                    WHEN apis.api_protocol IS NULL AND specs.api_spec_type::text = 'OPEN_API'::text THEN 'rest'::text
                    ELSE apis.api_protocol::text
                    END AS api_protocol,
                apis.tags,
                apis.countries,
                apis.links,
                apis.api_resource_links,
                apis.release_status,
                apis.sunset_date,
                apis.changelog_entries,
                apis.labels,
                apis.package_id,
                apis.visibility::text as visibility,
                apis.disabled,
                apis.part_of_products,
                apis.line_of_business,
                apis.industry,
                apis.ready,
                apis.created_at,
                apis.updated_at,
                apis.deleted_at,
                apis.error,
                apis.implementation_standard,
                apis.custom_implementation_standard,
                apis.custom_implementation_standard_description,
                apis.target_urls,
                apis.extensible,
                apis.successors,
                apis.resource_hash,
                apis.documentation_labels
FROM api_definitions apis
         JOIN (SELECT a1.id,
                      a1.tenant_id::text AS tenant_id,
                      a1.tenant_id::text AS provider_tenant_id
               FROM tenant_applications a1
               UNION ALL
               SELECT apps_subaccounts_func.id,
                      apps_subaccounts_func.tenant_id::text,
                      apps_subaccounts_func.provider_tenant_id::text
               FROM apps_subaccounts_func() apps_subaccounts_func(id, tenant_id, provider_tenant_id)
               UNION ALL
               SELECT ta.id AS app_id, ta.tenant_id::text AS consumer_tenant, tenant_runtimes.tenant_id::text AS provider_tenant
               FROM (SELECT labels.runtime_id, v ->> 0 AS consumer_tenant
                     FROM labels
                              JOIN jsonb_array_elements(labels.value) AS v ON TRUE
                     WHERE key = 'consumer_subaccount_ids') AS t_rts -- Get runtime and external consumer IDs pairs
                        JOIN business_tenant_mappings t ON t_rts.consumer_tenant = t.external_tenant -- Get runtime and internal consumer IDs pairs
                        JOIN apps_subaccounts_func() ta ON t.id = ta.tenant_id -- Get applications for consumer tenants
                        JOIN tenant_runtimes ON t_rts.runtime_id = tenant_runtimes.id) t_apps
              ON apis.app_id = t_apps.id
         LEFT JOIN specifications specs ON apis.id = specs.api_def_id;

-----

CREATE OR REPLACE VIEW tenants_events
            (tenant_id, provider_tenant_id, id, app_id, name, description, group_name, version_value,
             version_deprecated, version_deprecated_since, version_for_removal, ord_id, short_description,
             system_instance_aware, changelog_entries, links, tags, countries, release_status, sunset_date, labels,
             package_id, visibility, disabled, part_of_products, line_of_business, industry, ready, created_at,
             updated_at, deleted_at, error, extensible, successors, resource_hash)
AS
SELECT DISTINCT t_apps.tenant_id,
                t_apps.provider_tenant_id,
                events.id,
                events.app_id,
                events.name,
                events.description,
                events.group_name,
                events.version_value,
                events.version_deprecated,
                events.version_deprecated_since,
                events.version_for_removal,
                events.ord_id,
                events.short_description,
                events.system_instance_aware,
                events.changelog_entries,
                events.links,
                events.tags,
                events.countries,
                events.release_status,
                events.sunset_date,
                events.labels,
                events.package_id,
                events.visibility::text as visibility,
                events.disabled,
                events.part_of_products,
                events.line_of_business,
                events.industry,
                events.ready,
                events.created_at,
                events.updated_at,
                events.deleted_at,
                events.error,
                events.extensible,
                events.successors,
                events.resource_hash
FROM event_api_definitions events
         JOIN (SELECT a1.id,
                      a1.tenant_id::text AS tenant_id,
                      a1.tenant_id::text AS provider_tenant_id
               FROM tenant_applications a1
               UNION ALL
               SELECT apps_subaccounts_func.id,
                      apps_subaccounts_func.tenant_id::text,
                      apps_subaccounts_func.provider_tenant_id::text
               FROM apps_subaccounts_func() apps_subaccounts_func(id, tenant_id, provider_tenant_id)
               UNION ALL
               SELECT ta.id AS app_id, ta.tenant_id::text AS consumer_tenant, tenant_runtimes.tenant_id::text AS provider_tenant
               FROM (SELECT labels.runtime_id, v ->> 0 AS consumer_tenant
                     FROM labels
                              JOIN jsonb_array_elements(labels.value) AS v ON TRUE
                     WHERE key = 'consumer_subaccount_ids') AS t_rts -- Get runtime and external consumer IDs pairs
                        JOIN business_tenant_mappings t ON t_rts.consumer_tenant = t.external_tenant -- Get runtime and internal consumer IDs pairs
                        JOIN apps_subaccounts_func() ta ON t.id = ta.tenant_id -- Get applications for consumer tenants
                        JOIN tenant_runtimes ON t_rts.runtime_id = tenant_runtimes.id) t_apps
              ON events.app_id = t_apps.id;

-----

CREATE OR REPLACE VIEW tenants_specifications
            (tenant_id, provider_tenant_id, id, api_def_id, event_def_id, spec_data, api_spec_format, api_spec_type,
             event_spec_format, event_spec_type, custom_type, created_at)
AS
SELECT DISTINCT t_api_event_def.tenant_id,
                t_api_event_def.provider_tenant_id,
                spec.id,
                spec.api_def_id,
                spec.event_def_id,
                spec.spec_data,
                spec.api_spec_format,
                spec.api_spec_type,
                spec.event_spec_format,
                spec.event_spec_type,
                spec.custom_type,
                spec.created_at
FROM specifications spec
         JOIN (SELECT a.id,
                      a.tenant_id,
                      a.provider_tenant_id
               FROM tenants_apis a
               UNION ALL
               SELECT e.id,
                      e.tenant_id,
                      e.provider_tenant_id
               FROM tenants_events e) t_api_event_def
              ON spec.api_def_id = t_api_event_def.id OR spec.event_def_id = t_api_event_def.id;

COMMIT;