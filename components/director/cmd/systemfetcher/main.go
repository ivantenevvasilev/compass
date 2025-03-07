package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/domain/certsubjectmapping"

	"github.com/kyma-incubator/compass/components/director/internal/domain/bundleinstanceauth"

	"github.com/kyma-incubator/compass/components/director/internal/domain/formationconstraint/operators"

	"github.com/kyma-incubator/compass/components/director/internal/model"

	"github.com/kyma-incubator/compass/components/director/internal/domain/systemssync"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenantbusinesstype"

	"github.com/kyma-incubator/compass/components/director/internal/domain/formationconstraint"
	"github.com/kyma-incubator/compass/components/director/internal/domain/formationtemplateconstraintreferences"

	databuilder "github.com/kyma-incubator/compass/components/director/internal/domain/webhook/datainputbuilder"

	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"

	"github.com/kyma-incubator/compass/components/director/internal/domain/formationassignment"

	webhookclient "github.com/kyma-incubator/compass/components/director/pkg/webhook_client"

	"github.com/kyma-incubator/compass/components/director/internal/domain/formationtemplate"

	"github.com/kyma-incubator/compass/components/director/internal/domain/formation"

	"github.com/kyma-incubator/compass/components/director/internal/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal/domain/apptemplate"
	"github.com/kyma-incubator/compass/components/director/internal/domain/auth"
	bundleutil "github.com/kyma-incubator/compass/components/director/internal/domain/bundle"
	"github.com/kyma-incubator/compass/components/director/internal/domain/bundlereferences"
	"github.com/kyma-incubator/compass/components/director/internal/domain/document"
	"github.com/kyma-incubator/compass/components/director/internal/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal/domain/fetchrequest"
	"github.com/kyma-incubator/compass/components/director/internal/domain/integrationsystem"
	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/domain/labeldef"
	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	runtimectx "github.com/kyma-incubator/compass/components/director/internal/domain/runtime_context"
	"github.com/kyma-incubator/compass/components/director/internal/domain/scenarioassignment"
	"github.com/kyma-incubator/compass/components/director/internal/domain/spec"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/domain/version"
	"github.com/kyma-incubator/compass/components/director/internal/domain/webhook"
	"github.com/kyma-incubator/compass/components/director/internal/features"
	"github.com/kyma-incubator/compass/components/director/internal/systemfetcher"
	"github.com/kyma-incubator/compass/components/director/internal/uid"
	"github.com/kyma-incubator/compass/components/director/pkg/accessstrategy"
	pkgAuth "github.com/kyma-incubator/compass/components/director/pkg/auth"
	"github.com/kyma-incubator/compass/components/director/pkg/certloader"
	configprovider "github.com/kyma-incubator/compass/components/director/pkg/config"
	"github.com/kyma-incubator/compass/components/director/pkg/executor"
	httputil "github.com/kyma-incubator/compass/components/director/pkg/http"
	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/normalizer"
	oauth "github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	gcli "github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

const discoverSystemsOpMode = "DISCOVER_SYSTEMS"

type config struct {
	APIConfig      systemfetcher.APIConfig
	OAuth2Config   oauth.Config
	SystemFetcher  systemfetcher.Config
	Database       persistence.DatabaseConfig
	TemplateConfig appTemplateConfig

	Log log.Config

	Features features.Config

	ConfigurationFile string

	ConfigurationFileReload time.Duration `envconfig:"default=1m"`
	ClientTimeout           time.Duration `envconfig:"default=60s"`

	CertLoaderConfig certloader.Config

	SelfRegisterDistinguishLabelKey string `envconfig:"APP_SELF_REGISTER_DISTINGUISH_LABEL_KEY"`

	ORDWebhookMappings string `envconfig:"APP_ORD_WEBHOOK_MAPPINGS"`

	ExternalClientCertSecretName string `envconfig:"APP_EXTERNAL_CLIENT_CERT_SECRET_NAME"`
	ExtSvcClientCertSecretName   string `envconfig:"APP_EXT_SVC_CLIENT_CERT_SECRET_NAME"`
}

type appTemplateConfig struct {
	LabelFilter                    string `envconfig:"APP_TEMPLATE_LABEL_FILTER"`
	OverrideApplicationInput       string `envconfig:"APP_TEMPLATE_OVERRIDE_APPLICATION_INPUT"`
	PlaceholderToSystemKeyMappings string `envconfig:"APP_TEMPLATE_PLACEHOLDER_TO_SYSTEM_KEY_MAPPINGS"`
}

func main() {
	cfg := config{}
	err := envconfig.InitWithPrefix(&cfg, "APP")
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to load config"))
	}

	ctx, err := log.Configure(context.Background(), &cfg.Log)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to configure logger"))
	}

	cfgProvider := createAndRunConfigProvider(ctx, cfg)

	transact, closeFunc, err := persistence.Configure(ctx, cfg.Database)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to connect to the database"))
	}
	defer func() {
		err := closeFunc()
		if err != nil {
			log.D().Fatal(errors.Wrap(err, "failed to close database connection"))
		}
	}()

	certCache, err := certloader.StartCertLoader(ctx, cfg.CertLoaderConfig)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to initialize certificate loader"))
	}

	httpClient := &http.Client{Timeout: cfg.ClientTimeout}
	securedHTTPClient := pkgAuth.PrepareHTTPClient(cfg.ClientTimeout)
	mtlsClient := pkgAuth.PrepareMTLSClient(cfg.ClientTimeout, certCache, cfg.ExternalClientCertSecretName)
	extSvcMtlsClient := pkgAuth.PrepareMTLSClient(cfg.ClientTimeout, certCache, cfg.ExtSvcClientCertSecretName)

	sf, err := createSystemFetcher(ctx, cfg, cfgProvider, transact, httpClient, securedHTTPClient, mtlsClient, extSvcMtlsClient, certCache)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to initialize System Fetcher"))
	}

	if cfg.SystemFetcher.OperationalMode != discoverSystemsOpMode {
		log.C(ctx).Infof("The operatioal mode is set to %q, skipping systems discovery.", cfg.SystemFetcher.OperationalMode)
		return
	}

	if err = sf.SyncSystems(ctx); err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to sync systems"))
	}

	if err = sf.UpsertSystemsSyncTimestamps(ctx, transact); err != nil {
		log.D().Fatal(errors.Wrap(err, "failed to upsert systems synchronization timestamps in database"))
	}
}

func createSystemFetcher(ctx context.Context, cfg config, cfgProvider *configprovider.Provider, tx persistence.Transactioner, httpClient, securedHTTPClient, mtlsClient, extSvcMtlsClient *http.Client, certCache certloader.Cache) (*systemfetcher.SystemFetcher, error) {
	ordWebhookMapping, err := application.UnmarshalMappings(cfg.ORDWebhookMappings)
	if err != nil {
		return nil, errors.Wrap(err, "failed while unmarshalling ord webhook mappings")
	}

	tenantConverter := tenant.NewConverter()
	tenantBusinessTypeConverter := tenantbusinesstype.NewConverter()
	authConverter := auth.NewConverter()
	frConverter := fetchrequest.NewConverter(authConverter)
	versionConverter := version.NewConverter()
	docConverter := document.NewConverter(frConverter)
	webhookConverter := webhook.NewConverter(authConverter)
	specConverter := spec.NewConverter(frConverter)
	apiConverter := api.NewConverter(versionConverter, specConverter)
	eventAPIConverter := eventdef.NewConverter(versionConverter, specConverter)
	labelDefConverter := labeldef.NewConverter()
	labelConverter := label.NewConverter()
	intSysConverter := integrationsystem.NewConverter()
	bundleConverter := bundleutil.NewConverter(authConverter, apiConverter, eventAPIConverter, docConverter)
	appConverter := application.NewConverter(webhookConverter, bundleConverter)
	runtimeConverter := runtime.NewConverter(webhookConverter)
	bundleReferenceConverter := bundlereferences.NewConverter()
	runtimeContextConverter := runtimectx.NewConverter()
	formationConverter := formation.NewConverter()
	formationTemplateConverter := formationtemplate.NewConverter(webhookConverter)
	assignmentConverter := scenarioassignment.NewConverter()
	appTemplateConverter := apptemplate.NewConverter(appConverter, webhookConverter)
	formationAssignmentConverter := formationassignment.NewConverter()
	formationConstraintConverter := formationconstraint.NewConverter()
	formationTemplateConstraintReferencesConverter := formationtemplateconstraintreferences.NewConverter()
	systemsSyncConverter := systemssync.NewConverter()
	bundleInstanceAuthConv := bundleinstanceauth.NewConverter(authConverter)
	certSubjectMappingConv := certsubjectmapping.NewConverter()

	tenantRepo := tenant.NewRepository(tenantConverter)
	tenantBusinessTypeRepo := tenantbusinesstype.NewRepository(tenantBusinessTypeConverter)
	runtimeRepo := runtime.NewRepository(runtimeConverter)
	applicationRepo := application.NewRepository(appConverter)
	labelRepo := label.NewRepository(labelConverter)
	labelDefRepo := labeldef.NewRepository(labelDefConverter)
	webhookRepo := webhook.NewRepository(webhookConverter)
	apiRepo := api.NewRepository(apiConverter)
	eventAPIRepo := eventdef.NewRepository(eventAPIConverter)
	specRepo := spec.NewRepository(specConverter)
	docRepo := document.NewRepository(docConverter)
	fetchRequestRepo := fetchrequest.NewRepository(frConverter)
	intSysRepo := integrationsystem.NewRepository(intSysConverter)
	bundleRepo := bundleutil.NewRepository(bundleConverter)
	bundleReferenceRepo := bundlereferences.NewRepository(bundleReferenceConverter)
	runtimeContextRepo := runtimectx.NewRepository(runtimeContextConverter)
	formationRepo := formation.NewRepository(formationConverter)
	formationTemplateRepo := formationtemplate.NewRepository(formationTemplateConverter)
	scenarioAssignmentRepo := scenarioassignment.NewRepository(assignmentConverter)
	appTemplateRepo := apptemplate.NewRepository(appTemplateConverter)
	formationAssignmentRepo := formationassignment.NewRepository(formationAssignmentConverter)
	formationConstraintRepo := formationconstraint.NewRepository(formationConstraintConverter)
	formationTemplateConstraintReferencesRepo := formationtemplateconstraintreferences.NewRepository(formationTemplateConstraintReferencesConverter)
	systemsSyncRepo := systemssync.NewRepository(systemsSyncConverter)
	bundleInstanceAuthRepo := bundleinstanceauth.NewRepository(bundleInstanceAuthConv)
	certSubjectMappingRepo := certsubjectmapping.NewRepository(certSubjectMappingConv)

	uidSvc := uid.NewService()
	tenantSvc := tenant.NewService(tenantRepo, uidSvc, tenantConverter)
	tenantBusinessTypeSvc := tenantbusinesstype.NewService(tenantBusinessTypeRepo, uidSvc)
	labelSvc := label.NewLabelService(labelRepo, labelDefRepo, uidSvc)
	intSysSvc := integrationsystem.NewService(intSysRepo, uidSvc)
	scenariosSvc := labeldef.NewService(labelDefRepo, labelRepo, scenarioAssignmentRepo, tenantRepo, uidSvc)
	fetchRequestSvc := fetchrequest.NewService(fetchRequestRepo, httpClient, accessstrategy.NewDefaultExecutorProvider(certCache, cfg.ExternalClientCertSecretName, cfg.ExtSvcClientCertSecretName))
	specSvc := spec.NewService(specRepo, fetchRequestRepo, uidSvc, fetchRequestSvc)
	bundleReferenceSvc := bundlereferences.NewService(bundleReferenceRepo, uidSvc)
	apiSvc := api.NewService(apiRepo, uidSvc, specSvc, bundleReferenceSvc)
	eventAPISvc := eventdef.NewService(eventAPIRepo, uidSvc, specSvc, bundleReferenceSvc)
	docSvc := document.NewService(docRepo, fetchRequestRepo, uidSvc)
	bundleInstanceAuthSvc := bundleinstanceauth.NewService(bundleInstanceAuthRepo, uidSvc)
	bundleSvc := bundleutil.NewService(bundleRepo, apiSvc, eventAPISvc, docSvc, bundleInstanceAuthSvc, uidSvc)
	scenarioAssignmentSvc := scenarioassignment.NewService(scenarioAssignmentRepo, scenariosSvc)
	tntSvc := tenant.NewServiceWithLabels(tenantRepo, uidSvc, labelRepo, labelSvc, tenantConverter)
	webhookClient := webhookclient.NewClient(securedHTTPClient, mtlsClient, extSvcMtlsClient)
	appTemplateSvc := apptemplate.NewService(appTemplateRepo, webhookRepo, uidSvc, labelSvc, labelRepo, applicationRepo)
	webhookLabelBuilder := databuilder.NewWebhookLabelBuilder(labelRepo)
	webhookTenantBuilder := databuilder.NewWebhookTenantBuilder(webhookLabelBuilder, tenantRepo)
	certSubjectInputBuilder := databuilder.NewWebhookCertSubjectBuilder(certSubjectMappingRepo)
	webhookDataInputBuilder := databuilder.NewWebhookDataInputBuilder(applicationRepo, appTemplateRepo, runtimeRepo, runtimeContextRepo, webhookLabelBuilder, webhookTenantBuilder, certSubjectInputBuilder)
	formationConstraintSvc := formationconstraint.NewService(formationConstraintRepo, formationTemplateConstraintReferencesRepo, uidSvc, formationConstraintConverter)
	constraintEngine := operators.NewConstraintEngine(tx, formationConstraintSvc, tenantSvc, scenarioAssignmentSvc, nil, nil, formationRepo, labelRepo, labelSvc, applicationRepo, runtimeContextRepo, formationTemplateRepo, formationAssignmentRepo, cfg.Features.RuntimeTypeLabelKey, cfg.Features.ApplicationTypeLabelKey)
	notificationsBuilder := formation.NewNotificationsBuilder(webhookConverter, constraintEngine, cfg.Features.RuntimeTypeLabelKey, cfg.Features.ApplicationTypeLabelKey)
	notificationsGenerator := formation.NewNotificationsGenerator(applicationRepo, appTemplateRepo, runtimeRepo, runtimeContextRepo, labelRepo, webhookRepo, webhookDataInputBuilder, notificationsBuilder)
	notificationSvc := formation.NewNotificationService(tenantRepo, webhookClient, notificationsGenerator, constraintEngine, webhookConverter, formationTemplateRepo)
	faNotificationSvc := formationassignment.NewFormationAssignmentNotificationService(formationAssignmentRepo, webhookConverter, webhookRepo, tenantRepo, webhookDataInputBuilder, formationRepo, notificationsBuilder, runtimeContextRepo, labelSvc, cfg.Features.RuntimeTypeLabelKey, cfg.Features.ApplicationTypeLabelKey)
	formationAssignmentStatusSvc := formationassignment.NewFormationAssignmentStatusService(formationAssignmentRepo, constraintEngine, faNotificationSvc)
	formationAssignmentSvc := formationassignment.NewService(formationAssignmentRepo, uidSvc, applicationRepo, runtimeRepo, runtimeContextRepo, notificationSvc, faNotificationSvc, labelSvc, formationRepo, formationAssignmentStatusSvc, cfg.Features.RuntimeTypeLabelKey, cfg.Features.ApplicationTypeLabelKey)
	formationStatusSvc := formation.NewFormationStatusService(formationRepo, labelDefRepo, scenariosSvc, notificationSvc, constraintEngine)
	formationSvc := formation.NewService(tx, applicationRepo, labelDefRepo, labelRepo, formationRepo, formationTemplateRepo, labelSvc, uidSvc, scenariosSvc, scenarioAssignmentRepo, scenarioAssignmentSvc, tntSvc, runtimeRepo, runtimeContextRepo, formationAssignmentSvc, faNotificationSvc, notificationSvc, constraintEngine, webhookRepo, formationStatusSvc, cfg.Features.RuntimeTypeLabelKey, cfg.Features.ApplicationTypeLabelKey)
	appSvc := application.NewService(&normalizer.DefaultNormalizator{}, cfgProvider, applicationRepo, webhookRepo, runtimeRepo, labelRepo, intSysRepo, labelSvc, bundleSvc, uidSvc, formationSvc, cfg.SelfRegisterDistinguishLabelKey, ordWebhookMapping)
	systemsSyncSvc := systemssync.NewService(systemsSyncRepo)

	authProvider := pkgAuth.NewMtlsTokenAuthorizationProvider(cfg.OAuth2Config, cfg.ExternalClientCertSecretName, certCache, pkgAuth.DefaultMtlsClientCreator)
	client := &http.Client{
		Transport: httputil.NewSecuredTransport(httputil.NewHTTPTransportWrapper(http.DefaultTransport.(*http.Transport)), authProvider),
		Timeout:   cfg.APIConfig.Timeout,
	}
	oauthMtlsClient := systemfetcher.NewOauthMtlsClient(cfg.OAuth2Config, certCache, client)
	systemsAPIClient := systemfetcher.NewClient(cfg.APIConfig, oauthMtlsClient)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.SystemFetcher.DirectorSkipSSLValidation,
		},
	}

	httpTransport := httputil.NewCorrelationIDTransport(httputil.NewErrorHandlerTransport(httputil.NewHTTPTransportWrapper(tr)))

	securedClient := &http.Client{
		Transport: httpTransport,
		Timeout:   cfg.SystemFetcher.DirectorRequestTimeout,
	}

	graphqlClient := gcli.NewClient(cfg.SystemFetcher.DirectorGraphqlURL, gcli.WithHTTPClient(securedClient))
	directorClient := &systemfetcher.DirectorGraphClient{
		Client:        graphqlClient,
		Authenticator: pkgAuth.NewServiceAccountTokenAuthorizationProvider(),
	}

	dataLoader := systemfetcher.NewDataLoader(tx, cfg.SystemFetcher, appTemplateSvc, intSysSvc)
	if err := dataLoader.LoadData(ctx, os.ReadDir, os.ReadFile); err != nil {
		return nil, err
	}

	if err := loadSystemsSynchronizationTimestamps(ctx, tx, systemsSyncSvc); err != nil {
		return nil, errors.Wrap(err, "failed while loading systems synchronization timestamps")
	}

	var placeholdersMapping []systemfetcher.PlaceholderMapping
	if err := json.Unmarshal([]byte(cfg.TemplateConfig.PlaceholderToSystemKeyMappings), &placeholdersMapping); err != nil {
		return nil, errors.Wrapf(err, "while unmarshaling placeholders mapping")
	}

	if err := calculateTemplateMappings(ctx, cfg, tx, appTemplateSvc, placeholdersMapping); err != nil {
		return nil, errors.Wrap(err, "failed while calculating application templates mappings")
	}

	templateRenderer, err := systemfetcher.NewTemplateRenderer(appTemplateSvc, appConverter, cfg.TemplateConfig.OverrideApplicationInput, placeholdersMapping)
	if err != nil {
		return nil, errors.Wrapf(err, "while creating template renderer")
	}

	return systemfetcher.NewSystemFetcher(tx, tenantSvc, appSvc, systemsSyncSvc, tenantBusinessTypeSvc, templateRenderer, systemsAPIClient, directorClient, cfg.SystemFetcher), nil
}

func createAndRunConfigProvider(ctx context.Context, cfg config) *configprovider.Provider {
	provider := configprovider.NewProvider(cfg.ConfigurationFile)
	err := provider.Load()
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "error on loading configuration file"))
	}
	executor.NewPeriodic(cfg.ConfigurationFileReload, func(ctx context.Context) {
		if err = provider.Load(); err != nil {
			if err != nil {
				log.D().Fatal(errors.Wrap(err, "error from Reloader watch"))
			}
		}
		log.C(ctx).Infof("Successfully reloaded configuration file.")
	}).Run(ctx)

	return provider
}

func calculateTemplateMappings(ctx context.Context, cfg config, transact persistence.Transactioner, appTemplateSvc apptemplate.ApplicationTemplateService, placeholdersMapping []systemfetcher.PlaceholderMapping) error {
	applicationTemplates := make([]systemfetcher.TemplateMapping, 0)

	tx, err := transact.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	appTemplates, err := appTemplateSvc.ListByFilters(ctx, []*labelfilter.LabelFilter{labelfilter.NewForKey(cfg.TemplateConfig.LabelFilter)})
	if err != nil {
		return errors.Wrapf(err, "while listing application templates by label filter %q", cfg.TemplateConfig.LabelFilter)
	}

	selectFilterProperties := make(map[string]bool, 0)
	for _, appTemplate := range appTemplates {
		lbl, err := appTemplateSvc.ListLabels(ctx, appTemplate.ID)
		if err != nil {
			return errors.Wrapf(err, "while listing labels for application template with ID %q", appTemplate.ID)
		}

		applicationTemplates = append(applicationTemplates, systemfetcher.TemplateMapping{AppTemplate: appTemplate, Labels: lbl})

		addPropertiesFromAppTemplatePlaceholders(selectFilterProperties, appTemplate.Placeholders)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	systemfetcher.ApplicationTemplates = applicationTemplates
	systemfetcher.SelectFilter = createSelectFilter(selectFilterProperties, placeholdersMapping)
	systemfetcher.ApplicationTemplateLabelFilter = cfg.TemplateConfig.LabelFilter
	systemfetcher.SystemSourceKey = cfg.APIConfig.SystemSourceKey
	return nil
}

func loadSystemsSynchronizationTimestamps(ctx context.Context, transact persistence.Transactioner, systemSyncSvc systemfetcher.SystemsSyncService) error {
	systemSynchronizationTimestamps := make(map[string]map[string]systemfetcher.SystemSynchronizationTimestamp, 0)

	tx, err := transact.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	syncTimestamps, err := systemSyncSvc.List(ctx)
	if err != nil {
		return err
	}

	for _, s := range syncTimestamps {
		currentTimestamp := systemfetcher.SystemSynchronizationTimestamp{
			ID:                s.ID,
			LastSyncTimestamp: s.LastSyncTimestamp,
		}

		if _, ok := systemSynchronizationTimestamps[s.TenantID]; !ok {
			systemSynchronizationTimestamps[s.TenantID] = make(map[string]systemfetcher.SystemSynchronizationTimestamp, 0)
		}

		systemSynchronizationTimestamps[s.TenantID][s.ProductID] = currentTimestamp
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	systemfetcher.SystemSynchronizationTimestamps = systemSynchronizationTimestamps

	return nil
}

func getTopParentFromJSONPath(jsonPath string) string {
	prefix := "$."
	infix := "."

	topParent := strings.TrimPrefix(jsonPath, prefix)
	firstInfixIndex := strings.Index(topParent, infix)
	if firstInfixIndex == -1 {
		return topParent
	}

	return topParent[:firstInfixIndex]
}

func addPropertiesFromAppTemplatePlaceholders(selectFilterProperties map[string]bool, placeholders []model.ApplicationTemplatePlaceholder) {
	for _, placeholder := range placeholders {
		if placeholder.JSONPath != nil && len(*placeholder.JSONPath) > 0 {
			topParent := getTopParentFromJSONPath(*placeholder.JSONPath)
			if _, exists := selectFilterProperties[topParent]; !exists {
				selectFilterProperties[topParent] = true
			}
		}
	}
}

func createSelectFilter(selectFilterProperties map[string]bool, placeholdersMapping []systemfetcher.PlaceholderMapping) []string {
	selectFilter := make([]string, 0)

	for _, pm := range placeholdersMapping {
		topParent := getTopParentFromJSONPath(pm.SystemKey)
		if _, exists := selectFilterProperties[topParent]; !exists {
			selectFilterProperties[topParent] = true
		}
	}

	for property := range selectFilterProperties {
		selectFilter = append(selectFilter, property)
	}

	return selectFilter
}
