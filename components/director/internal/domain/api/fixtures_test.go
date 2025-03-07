package api_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/domain/version"

	"github.com/kyma-incubator/compass/components/director/internal/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal/repo"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

const (
	apiDefID         = "ddddddddd-dddd-dddd-dddd-dddddddddddd"
	specID           = "sssssssss-ssss-ssss-ssss-ssssssssssss"
	tenantID         = "b91b59f7-2563-40b2-aba9-fef726037aa3"
	externalTenantID = "eeeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	bundleID         = "bbbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	packageID        = "ppppppppp-pppp-pppp-pppp-pppppppppppp"
	ordID            = "com.compass.ord.v1"
	localTenantID    = "localTenantID"
	extensible       = `{"supported":"automatic","description":"Please find the extensibility documentation"}`
	successors       = `["sap.s4:apiResource:API_BILL_OF_MATERIAL_SRV:v2"]`
	resourceHash     = "123456"
	publicVisibility = "public"
	targetURL        = "https://test-url.com/api"

	firstTargetURL  = "https://test-url.com"
	secondTargetURL = "https://test2-url.com"
)

var (
	fixedTimestamp       = time.Now()
	firstBundleID        = "id1"
	secondBundleID       = "id2"
	thirdBundleID        = "id3"
	frURL                = "foo.bar"
	spec                 = "spec"
	appID                = "aaaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	appTemplateVersionID = "fffffffff-ffff-aaaa-ffff-aaaaaaaaaaaa"
)

func fixAPIDefinitionModel(id string, name, targetURL string) *model.APIDefinition {
	return &model.APIDefinition{
		Name:       name,
		TargetURLs: api.ConvertTargetURLToJSONArray(targetURL),
		BaseEntity: &model.BaseEntity{ID: id},
		Visibility: str.Ptr(publicVisibility),
	}
}

func fixAPIDefinitionWithPackageModel(id, name string) *model.APIDefinition {
	return &model.APIDefinition{
		PackageID:  str.Ptr(packageID),
		Name:       name,
		TargetURLs: api.ConvertTargetURLToJSONArray(targetURL),
		Version:    &model.Version{},
		BaseEntity: &model.BaseEntity{
			ID:    id,
			Ready: true,
		},
	}
}

func fixFullAPIDefinitionModelWithAppID(placeholder string) (model.APIDefinition, model.Spec, model.BundleReference) {
	apiDef, spec, bundleRef := fixFullAPIDefinitionModel(apiDefID, placeholder)
	apiDef.ApplicationID = str.Ptr(appID)

	return apiDef, spec, bundleRef
}

func fixFullAPIDefinitionModelWithAppTemplateVersionID(placeholder string) (model.APIDefinition, model.Spec, model.BundleReference) {
	apiDef, spec, bundleRef := fixFullAPIDefinitionModel(apiDefID, placeholder)
	apiDef.ApplicationTemplateVersionID = str.Ptr(appTemplateVersionID)

	return apiDef, spec, bundleRef
}

func fixFullAPIDefinitionModel(apiDefID, placeholder string) (model.APIDefinition, model.Spec, model.BundleReference) {
	apiType := model.APISpecTypeOpenAPI
	spec := model.Spec{
		ID:         specID,
		Data:       str.Ptr("spec_data_" + placeholder),
		Format:     model.SpecFormatYaml,
		ObjectType: model.APISpecReference,
		ObjectID:   apiDefID,
		APIType:    &apiType,
	}

	deprecated := false
	forRemoval := false

	v := &model.Version{
		Value:           "v1.1",
		Deprecated:      &deprecated,
		DeprecatedSince: str.Ptr("v1.0"),
		ForRemoval:      &forRemoval,
	}

	apiBundleReference := model.BundleReference{
		BundleID:            str.Ptr(bundleID),
		ObjectType:          model.BundleAPIReference,
		ObjectID:            str.Ptr(apiDefID),
		APIDefaultTargetURL: str.Ptr(fmt.Sprintf("https://%s.com", placeholder)),
	}

	boolVar := false
	return model.APIDefinition{
		PackageID:                               str.Ptr(packageID),
		Name:                                    placeholder,
		Description:                             str.Ptr("desc_" + placeholder),
		TargetURLs:                              api.ConvertTargetURLToJSONArray(fmt.Sprintf("https://%s.com", placeholder)),
		Group:                                   str.Ptr("group_" + placeholder),
		OrdID:                                   str.Ptr(ordID),
		LocalTenantID:                           str.Ptr(localTenantID),
		ShortDescription:                        str.Ptr("shortDescription"),
		SystemInstanceAware:                     &boolVar,
		PolicyLevel:                             nil,
		CustomPolicyLevel:                       nil,
		APIProtocol:                             str.Ptr("apiProtocol"),
		Tags:                                    json.RawMessage("[]"),
		Countries:                               json.RawMessage("[]"),
		Links:                                   json.RawMessage("[]"),
		APIResourceLinks:                        json.RawMessage("[]"),
		ReleaseStatus:                           str.Ptr("releaseStatus"),
		SunsetDate:                              str.Ptr("sunsetDate"),
		Successors:                              json.RawMessage(successors),
		ChangeLogEntries:                        json.RawMessage("[]"),
		Labels:                                  json.RawMessage("[]"),
		Visibility:                              str.Ptr(publicVisibility),
		Disabled:                                &boolVar,
		PartOfProducts:                          json.RawMessage("[]"),
		LineOfBusiness:                          json.RawMessage("[]"),
		Industry:                                json.RawMessage("[]"),
		ImplementationStandard:                  str.Ptr("implementationStandard"),
		CustomImplementationStandard:            str.Ptr("customImplementationStandard"),
		CustomImplementationStandardDescription: str.Ptr("customImplementationStandardDescription"),
		Version:                                 v,
		Extensible:                              json.RawMessage(extensible),
		ResourceHash:                            str.Ptr(resourceHash),
		SupportedUseCases:                       json.RawMessage("[]"),
		DocumentationLabels:                     json.RawMessage("[]"),
		CorrelationIDs:                          json.RawMessage("[]"),
		Direction:                               str.Ptr("mixed"),
		BaseEntity: &model.BaseEntity{
			ID:        apiDefID,
			Ready:     true,
			CreatedAt: &fixedTimestamp,
			UpdatedAt: &time.Time{},
			DeletedAt: &time.Time{},
			Error:     nil,
		},
	}, spec, apiBundleReference
}

func fixFullGQLAPIDefinition(placeholder string) *graphql.APIDefinition {
	data := graphql.CLOB("spec_data_" + placeholder)

	spec := &graphql.APISpec{
		Data:         &data,
		Format:       graphql.SpecFormatYaml,
		Type:         graphql.APISpecTypeOpenAPI,
		DefinitionID: apiDefID,
	}

	deprecated := false
	forRemoval := false

	v := &graphql.Version{
		Value:           "v1.1",
		Deprecated:      &deprecated,
		DeprecatedSince: str.Ptr("v1.0"),
		ForRemoval:      &forRemoval,
	}

	return &graphql.APIDefinition{
		BundleID:    bundleID,
		Name:        placeholder,
		Description: str.Ptr("desc_" + placeholder),
		Spec:        spec,
		TargetURL:   fmt.Sprintf("https://%s.com", placeholder),
		Group:       str.Ptr("group_" + placeholder),
		Version:     v,
		BaseEntity: &graphql.BaseEntity{
			ID:        apiDefID,
			Ready:     true,
			Error:     nil,
			CreatedAt: timeToTimestampPtr(fixedTimestamp),
			UpdatedAt: timeToTimestampPtr(time.Time{}),
			DeletedAt: timeToTimestampPtr(time.Time{}),
		},
	}
}

func fixModelAPIDefinitionInput(name, description string, group string) (*model.APIDefinitionInput, *model.SpecInput) {
	data := "data"
	apiType := model.APISpecTypeOpenAPI

	spec := &model.SpecInput{
		Data:         &data,
		APIType:      &apiType,
		Format:       model.SpecFormatYaml,
		FetchRequest: &model.FetchRequestInput{},
	}

	deprecated := false
	deprecatedSince := ""
	forRemoval := false

	v := &model.VersionInput{
		Value:           "1.0.0",
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}

	return &model.APIDefinitionInput{
		Name:         name,
		Description:  &description,
		TargetURLs:   api.ConvertTargetURLToJSONArray(targetURL),
		Group:        &group,
		VersionInput: v,
	}, spec
}

func fixGQLAPIDefinitionInput(name, description string, group string) *graphql.APIDefinitionInput {
	data := graphql.CLOB("data")

	spec := &graphql.APISpecInput{
		Data:         &data,
		Type:         graphql.APISpecTypeOpenAPI,
		Format:       graphql.SpecFormatYaml,
		FetchRequest: &graphql.FetchRequestInput{},
	}

	deprecated := false
	deprecatedSince := ""
	forRemoval := false

	v := &graphql.VersionInput{
		Value:           "1.0.0",
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}

	return &graphql.APIDefinitionInput{
		Name:        name,
		Description: &description,
		TargetURL:   targetURL,
		Group:       &group,
		Spec:        spec,
		Version:     v,
	}
}

func fixEntityAPIDefinition(id string, name, targetURL string) *api.Entity {
	return &api.Entity{
		Name:       name,
		TargetURLs: repo.NewValidNullableString(`["` + targetURL + `"]`),
		BaseEntity: &repo.BaseEntity{ID: id},
		Visibility: publicVisibility,
	}
}

func fixFullEntityAPIDefinitionWithAppID(apiDefID, placeholder string) api.Entity {
	entity := fixFullEntityAPIDefinition(apiDefID, placeholder)
	entity.ApplicationID = repo.NewValidNullableString(appID)

	return entity
}

func fixFullEntityAPIDefinitionWithAppTemplateVersionID(apiDefID, placeholder string) api.Entity {
	entity := fixFullEntityAPIDefinition(apiDefID, placeholder)
	entity.ApplicationTemplateVersionID = repo.NewValidNullableString(appTemplateVersionID)

	return entity
}

func fixFullEntityAPIDefinition(apiDefID, placeholder string) api.Entity {
	return api.Entity{
		PackageID:                               repo.NewValidNullableString(packageID),
		Name:                                    placeholder,
		Description:                             repo.NewValidNullableString("desc_" + placeholder),
		Group:                                   repo.NewValidNullableString("group_" + placeholder),
		TargetURLs:                              repo.NewValidNullableString(`["` + fmt.Sprintf("https://%s.com", placeholder) + `"]`),
		OrdID:                                   repo.NewValidNullableString(ordID),
		LocalTenantID:                           repo.NewValidNullableString(localTenantID),
		ShortDescription:                        repo.NewValidNullableString("shortDescription"),
		SystemInstanceAware:                     repo.NewValidNullableBool(false),
		PolicyLevel:                             sql.NullString{},
		CustomPolicyLevel:                       sql.NullString{},
		APIProtocol:                             repo.NewValidNullableString("apiProtocol"),
		Tags:                                    repo.NewValidNullableString("[]"),
		Countries:                               repo.NewValidNullableString("[]"),
		Links:                                   repo.NewValidNullableString("[]"),
		APIResourceLinks:                        repo.NewValidNullableString("[]"),
		ReleaseStatus:                           repo.NewValidNullableString("releaseStatus"),
		SunsetDate:                              repo.NewValidNullableString("sunsetDate"),
		Successors:                              repo.NewValidNullableString(successors),
		ChangeLogEntries:                        repo.NewValidNullableString("[]"),
		Labels:                                  repo.NewValidNullableString("[]"),
		Visibility:                              publicVisibility,
		Disabled:                                repo.NewValidNullableBool(false),
		PartOfProducts:                          repo.NewValidNullableString("[]"),
		LineOfBusiness:                          repo.NewValidNullableString("[]"),
		Industry:                                repo.NewValidNullableString("[]"),
		ImplementationStandard:                  repo.NewValidNullableString("implementationStandard"),
		CustomImplementationStandard:            repo.NewValidNullableString("customImplementationStandard"),
		CustomImplementationStandardDescription: repo.NewValidNullableString("customImplementationStandardDescription"),
		Extensible:                              repo.NewValidNullableString(extensible),
		Version: version.Version{
			Value:           repo.NewNullableString(str.Ptr("v1.1")),
			Deprecated:      repo.NewValidNullableBool(false),
			DeprecatedSince: repo.NewNullableString(str.Ptr("v1.0")),
			ForRemoval:      repo.NewValidNullableBool(false),
		},
		ResourceHash:        repo.NewValidNullableString(resourceHash),
		SupportedUseCases:   repo.NewValidNullableString("[]"),
		DocumentationLabels: repo.NewValidNullableString("[]"),
		CorrelationIDs:      repo.NewValidNullableString("[]"),
		Direction:           repo.NewValidNullableString("mixed"),
		BaseEntity: &repo.BaseEntity{
			ID:        apiDefID,
			Ready:     true,
			CreatedAt: &fixedTimestamp,
			UpdatedAt: &time.Time{},
			DeletedAt: &time.Time{},
			Error:     sql.NullString{},
		},
	}
}

func fixAPIDefinitionColumns() []string {
	return []string{"id", "app_id", "app_template_version_id", "package_id", "name", "description", "group_name", "ord_id", "local_tenant_id",
		"short_description", "system_instance_aware", "policy_level", "custom_policy_level", "api_protocol",
		"tags", "countries", "links", "api_resource_links", "release_status",
		"sunset_date", "changelog_entries", "labels", "visibility", "disabled", "part_of_products", "line_of_business",
		"industry", "version_value", "version_deprecated", "version_deprecated_since", "version_for_removal", "ready",
		"created_at", "updated_at", "deleted_at", "error", "implementation_standard", "custom_implementation_standard",
		"custom_implementation_standard_description", "target_urls", "extensible", "successors", "resource_hash", "supported_use_cases", "documentation_labels", "correlation_ids", "direction"}
}

func fixAPIDefinitionRow(id, placeholder string) []driver.Value {
	boolVar := false
	return []driver.Value{id, appID, repo.NewValidNullableString(""), packageID, placeholder, "desc_" + placeholder, "group_" + placeholder,
		ordID, localTenantID, "shortDescription", &boolVar, nil, nil, "apiProtocol", repo.NewValidNullableString("[]"),
		repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), "releaseStatus", "sunsetDate", repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), publicVisibility, &boolVar,
		repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), "v1.1", false, "v1.0", false, true, fixedTimestamp, time.Time{}, time.Time{}, nil,
		"implementationStandard", "customImplementationStandard", "customImplementationStandardDescription", repo.NewValidNullableString(`["` + fmt.Sprintf("https://%s.com", placeholder) + `"]`),
		repo.NewValidNullableString(extensible), repo.NewValidNullableString(successors), repo.NewValidNullableString(resourceHash), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("mixed")}
}

func fixAPIDefinitionRowForAppTemplateVersion(id, placeholder string) []driver.Value {
	boolVar := false
	return []driver.Value{id, repo.NewValidNullableString(""), appTemplateVersionID, packageID, placeholder, "desc_" + placeholder, "group_" + placeholder,
		ordID, localTenantID, "shortDescription", &boolVar, nil, nil, "apiProtocol", repo.NewValidNullableString("[]"),
		repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), "releaseStatus", "sunsetDate", repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), publicVisibility, &boolVar,
		repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), "v1.1", false, "v1.0", false, true, fixedTimestamp, time.Time{}, time.Time{}, nil,
		"implementationStandard", "customImplementationStandard", "customImplementationStandardDescription", repo.NewValidNullableString(`["` + fmt.Sprintf("https://%s.com", placeholder) + `"]`),
		repo.NewValidNullableString(extensible), repo.NewValidNullableString(successors), repo.NewValidNullableString(resourceHash), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("mixed")}
}

func fixAPICreateArgs(id string, apiDef *model.APIDefinition) []driver.Value {
	return []driver.Value{id, appID, repo.NewValidNullableString(""), packageID, apiDef.Name, apiDef.Description, apiDef.Group,
		apiDef.OrdID, apiDef.LocalTenantID, apiDef.ShortDescription, apiDef.SystemInstanceAware, nil, nil, apiDef.APIProtocol, repo.NewNullableStringFromJSONRawMessage(apiDef.Tags), repo.NewNullableStringFromJSONRawMessage(apiDef.Countries),
		repo.NewNullableStringFromJSONRawMessage(apiDef.Links), repo.NewNullableStringFromJSONRawMessage(apiDef.APIResourceLinks),
		apiDef.ReleaseStatus, apiDef.SunsetDate, repo.NewNullableStringFromJSONRawMessage(apiDef.ChangeLogEntries), repo.NewNullableStringFromJSONRawMessage(apiDef.Labels), apiDef.Visibility,
		apiDef.Disabled, repo.NewNullableStringFromJSONRawMessage(apiDef.PartOfProducts), repo.NewNullableStringFromJSONRawMessage(apiDef.LineOfBusiness), repo.NewNullableStringFromJSONRawMessage(apiDef.Industry),
		apiDef.Version.Value, apiDef.Version.Deprecated, apiDef.Version.DeprecatedSince, apiDef.Version.ForRemoval, apiDef.Ready, apiDef.CreatedAt, apiDef.UpdatedAt, apiDef.DeletedAt, apiDef.Error,
		apiDef.ImplementationStandard, apiDef.CustomImplementationStandard, apiDef.CustomImplementationStandardDescription, repo.NewNullableStringFromJSONRawMessage(apiDef.TargetURLs), extensible,
		repo.NewNullableStringFromJSONRawMessage(apiDef.Successors), apiDef.ResourceHash, repo.NewNullableStringFromJSONRawMessage(apiDef.SupportedUseCases), repo.NewNullableStringFromJSONRawMessage(apiDef.DocumentationLabels),
		repo.NewNullableStringFromJSONRawMessage(apiDef.CorrelationIDs), apiDef.Direction}
}

func fixAPICreateArgsForAppTemplateVersion(id string, apiDef *model.APIDefinition) []driver.Value {
	return []driver.Value{id, repo.NewValidNullableString(""), appTemplateVersionID, packageID, apiDef.Name, apiDef.Description, apiDef.Group,
		apiDef.OrdID, apiDef.LocalTenantID, apiDef.ShortDescription, apiDef.SystemInstanceAware, nil, nil, apiDef.APIProtocol, repo.NewNullableStringFromJSONRawMessage(apiDef.Tags), repo.NewNullableStringFromJSONRawMessage(apiDef.Countries),
		repo.NewNullableStringFromJSONRawMessage(apiDef.Links), repo.NewNullableStringFromJSONRawMessage(apiDef.APIResourceLinks),
		apiDef.ReleaseStatus, apiDef.SunsetDate, repo.NewNullableStringFromJSONRawMessage(apiDef.ChangeLogEntries), repo.NewNullableStringFromJSONRawMessage(apiDef.Labels), apiDef.Visibility,
		apiDef.Disabled, repo.NewNullableStringFromJSONRawMessage(apiDef.PartOfProducts), repo.NewNullableStringFromJSONRawMessage(apiDef.LineOfBusiness), repo.NewNullableStringFromJSONRawMessage(apiDef.Industry),
		apiDef.Version.Value, apiDef.Version.Deprecated, apiDef.Version.DeprecatedSince, apiDef.Version.ForRemoval, apiDef.Ready, apiDef.CreatedAt, apiDef.UpdatedAt, apiDef.DeletedAt, apiDef.Error,
		apiDef.ImplementationStandard, apiDef.CustomImplementationStandard, apiDef.CustomImplementationStandardDescription, repo.NewNullableStringFromJSONRawMessage(apiDef.TargetURLs), extensible,
		repo.NewNullableStringFromJSONRawMessage(apiDef.Successors), apiDef.ResourceHash, repo.NewNullableStringFromJSONRawMessage(apiDef.SupportedUseCases), repo.NewNullableStringFromJSONRawMessage(apiDef.DocumentationLabels),
		repo.NewNullableStringFromJSONRawMessage(apiDef.CorrelationIDs), apiDef.Direction}
}

func fixModelFetchRequest(id, url string, timestamp time.Time) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     id,
		URL:    url,
		Auth:   nil,
		Mode:   "SINGLE",
		Filter: nil,
		Status: &model.FetchRequestStatus{
			Condition: model.FetchRequestStatusConditionInitial,
			Timestamp: timestamp,
		},
		ObjectType: model.APISpecFetchRequestReference,
		ObjectID:   specID,
	}
}

func fixModelBundleReference(bundleID, apiID string) *model.BundleReference {
	return &model.BundleReference{
		BundleID:   str.Ptr(bundleID),
		ObjectType: model.BundleAPIReference,
		ObjectID:   str.Ptr(apiID),
	}
}

func fixGQLFetchRequest(url string, timestamp time.Time) *graphql.FetchRequest {
	return &graphql.FetchRequest{
		Filter: nil,
		Mode:   graphql.FetchModeSingle,
		Auth:   nil,
		URL:    url,
		Status: &graphql.FetchRequestStatus{
			Timestamp: graphql.Timestamp(timestamp),
			Condition: graphql.FetchRequestStatusConditionInitial,
		},
	}
}

func timeToTimestampPtr(time time.Time) *graphql.Timestamp {
	t := graphql.Timestamp(time)
	return &t
}
