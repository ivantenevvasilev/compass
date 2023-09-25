package eventdef

import (
	"database/sql"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/domain/version"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
)

// Entity is a representation of a single EventDefinition in the database.
type Entity struct {
	ApplicationID                           sql.NullString `db:"app_id"`
	ApplicationTemplateVersionID            sql.NullString `db:"app_template_version_id"`
	PackageID                               sql.NullString `db:"package_id"`
	Name                                    string         `db:"name"`
	Description                             sql.NullString `db:"description"`
	GroupName                               sql.NullString `db:"group_name"`
	OrdID                                   sql.NullString `db:"ord_id"`
	LocalTenantID                           sql.NullString `db:"local_tenant_id"`
	ShortDescription                        sql.NullString `db:"short_description"`
	SystemInstanceAware                     sql.NullBool   `db:"system_instance_aware"`
	PolicyLevel                             sql.NullString `db:"policy_level"`
	CustomPolicyLevel                       sql.NullString `db:"custom_policy_level"`
	ChangeLogEntries                        sql.NullString `db:"changelog_entries"`
	Links                                   sql.NullString `db:"links"`
	Tags                                    sql.NullString `db:"tags"`
	Countries                               sql.NullString `db:"countries"`
	ReleaseStatus                           sql.NullString `db:"release_status"`
	SunsetDate                              sql.NullString `db:"sunset_date"`
	Successors                              sql.NullString `db:"successors"`
	Labels                                  sql.NullString `db:"labels"`
	Visibility                              string         `db:"visibility"`
	Disabled                                sql.NullBool   `db:"disabled"`
	PartOfProducts                          sql.NullString `db:"part_of_products"`
	LineOfBusiness                          sql.NullString `db:"line_of_business"`
	Industry                                sql.NullString `db:"industry"`
	ImplementationStandard                  sql.NullString `db:"implementation_standard"`
	CustomImplementationStandard            sql.NullString `db:"custom_implementation_standard"`
	CustomImplementationStandardDescription sql.NullString `db:"custom_implementation_standard_description"`
	Extensible                              sql.NullString `db:"extensible"`
	ResourceHash                            sql.NullString `db:"resource_hash"`
	Hierarchy                               sql.NullString `db:"hierarchy"`
	DocumentationLabels                     sql.NullString `db:"documentation_labels"`
	CorrelationIDs                          sql.NullString `db:"correlation_ids"`
	version.Version

	*repo.BaseEntity
}

// GetParent returns the parent type and the parent ID of the entity.
func (e *Entity) GetParent(_ resource.Type) (resource.Type, string) {
	if e.ApplicationID.Valid {
		return resource.Application, e.ApplicationID.String
	} else if e.ApplicationTemplateVersionID.Valid {
		return resource.ApplicationTemplateVersion, e.ApplicationTemplateVersionID.String
	}

	return "", ""
}

// DecorateWithTenantID decorates the entity with the given tenant ID.
func (e *Entity) DecorateWithTenantID(tenant string) interface{} {
	return struct {
		*Entity
		TenantID string `db:"tenant_id"`
	}{
		Entity:   e,
		TenantID: tenant,
	}
}
