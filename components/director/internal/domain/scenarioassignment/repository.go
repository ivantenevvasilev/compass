package scenarioassignment

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
)

const tableName string = `public.automatic_scenario_assignments`

var columns = []string{scenarioColumn, tenantColumn, selectorKeyColumn, selectorValueColumn}

var (
	tenantColumn        = "tenant_id"
	selectorKeyColumn   = "selector_key"
	selectorValueColumn = "selector_value"
	scenarioColumn      = "scenario"
)

// NewRepository missing godoc
func NewRepository(conv EntityConverter) *repository {
	return &repository{
		creator:         repo.NewCreatorGlobal(resource.AutomaticScenarioAssigment, tableName, columns),
		lister:          repo.NewListerWithEmbeddedTenant(resource.AutomaticScenarioAssigment, tableName, tenantColumn, columns),
		singleGetter:    repo.NewSingleGetterWithEmbeddedTenant(resource.AutomaticScenarioAssigment, tableName, tenantColumn, columns),
		pageableQuerier: repo.NewPageableQuerierWithEmbeddedTenant(resource.AutomaticScenarioAssigment, tableName, tenantColumn, columns),
		deleter:         repo.NewDeleterWithEmbeddedTenant(resource.AutomaticScenarioAssigment, tableName, tenantColumn),
		conv:            conv,
	}
}

type repository struct {
	creator         repo.CreatorGlobal
	singleGetter    repo.SingleGetter
	lister          repo.Lister
	pageableQuerier repo.PageableQuerier
	deleter         repo.Deleter
	conv            EntityConverter
}

// EntityConverter missing godoc
//go:generate mockery --name=EntityConverter --output=automock --outpkg=automock --case=underscore
type EntityConverter interface {
	ToEntity(assignment model.AutomaticScenarioAssignment) Entity
	FromEntity(assignment Entity) model.AutomaticScenarioAssignment
}

// Create missing godoc
func (r *repository) Create(ctx context.Context, model model.AutomaticScenarioAssignment) error {
	entity := r.conv.ToEntity(model)
	return r.creator.Create(ctx, entity)
}

// ListForSelector missing godoc
func (r *repository) ListForSelector(ctx context.Context, in model.LabelSelector, tenantID string) ([]*model.AutomaticScenarioAssignment, error) {
	var out EntityCollection

	conditions := repo.Conditions{
		repo.NewEqualCondition(selectorKeyColumn, in.Key),
		repo.NewEqualCondition(selectorValueColumn, in.Value),
	}

	if err := r.lister.List(ctx, tenantID, &out, conditions...); err != nil {
		return nil, errors.Wrap(err, "while getting automatic scenario assignments from db")
	}

	items := make([]*model.AutomaticScenarioAssignment, 0, len(out))

	for _, v := range out {
		item := r.conv.FromEntity(v)
		items = append(items, &item)
	}

	return items, nil
}

// GetForScenarioName missing godoc
func (r *repository) GetForScenarioName(ctx context.Context, tenantID, scenarioName string) (model.AutomaticScenarioAssignment, error) {
	var ent Entity

	conditions := repo.Conditions{
		repo.NewEqualCondition(scenarioColumn, scenarioName),
	}

	if err := r.singleGetter.Get(ctx, tenantID, conditions, repo.NoOrderBy, &ent); err != nil {
		return model.AutomaticScenarioAssignment{}, err
	}

	assignmentModel := r.conv.FromEntity(ent)

	return assignmentModel, nil
}

// List missing godoc
func (r *repository) List(ctx context.Context, tenantID string, pageSize int, cursor string) (*model.AutomaticScenarioAssignmentPage, error) {
	var collection EntityCollection
	page, totalCount, err := r.pageableQuerier.List(ctx, tenantID, pageSize, cursor, scenarioColumn, &collection)
	if err != nil {
		return nil, err
	}

	items := make([]*model.AutomaticScenarioAssignment, 0, len(collection))

	for _, ent := range collection {
		m := r.conv.FromEntity(ent)
		items = append(items, &m)
	}

	return &model.AutomaticScenarioAssignmentPage{
		Data:       items,
		TotalCount: totalCount,
		PageInfo:   page,
	}, nil
}

// DeleteForSelector missing godoc
func (r *repository) DeleteForSelector(ctx context.Context, tenantID string, selector model.LabelSelector) error {
	conditions := repo.Conditions{
		repo.NewEqualCondition(selectorKeyColumn, selector.Key),
		repo.NewEqualCondition(selectorValueColumn, selector.Value),
	}

	return r.deleter.DeleteMany(ctx, tenantID, conditions)
}

// DeleteForScenarioName missing godoc
func (r *repository) DeleteForScenarioName(ctx context.Context, tenantID string, scenarioName string) error {
	conditions := repo.Conditions{
		repo.NewEqualCondition(scenarioColumn, scenarioName),
	}

	return r.deleter.DeleteOne(ctx, tenantID, conditions)
}
