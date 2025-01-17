package repositories

import (
	"context"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/G-Research/fasttrackml/pkg/api/aim2/dao/models"
	"github.com/G-Research/fasttrackml/pkg/common/dao/repositories"
)

// TagRepositoryProvider provides an interface to work with models.Tag entity.
type TagRepositoryProvider interface {
	repositories.BaseRepositoryProvider
	// GetTagsByNamespace returns the list of tags.
	GetTagsByNamespace(ctx context.Context, namespaceID uint) ([]models.Tag, error)
	// CreateExperimentTag creates new models.ExperimentTag entity connected to models.Experiment.
	CreateExperimentTag(ctx context.Context, experimentTag *models.ExperimentTag) error
	// GetTagKeysByParameters returns list of tag keys by requested parameters.
	GetTagKeysByParameters(ctx context.Context, namespaceID uint, experimentNames []string) ([]string, error)
}

// TagRepository repository to work with models.Tag entity.
type TagRepository struct {
	repositories.BaseRepositoryProvider
}

// NewTagRepository creates repository to work with models.Tag entity.
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		repositories.NewBaseRepository(db),
	}
}

// CreateExperimentTag creates new models.ExperimentTag entity connected to models.Experiment.
func (r TagRepository) CreateExperimentTag(ctx context.Context, experimentTag *models.ExperimentTag) error {
	if err := r.GetDB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(experimentTag).Error; err != nil {
		return eris.Wrapf(err, "error creating tag for experiment with id: %d", experimentTag.ExperimentID)
	}
	return nil
}

// GetTagsByNamespace returns the list of tags.
// TODO fix stub implementation
func (r TagRepository) GetTagsByNamespace(ctx context.Context, namespaceID uint) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.GetDB().WithContext(ctx).Find(&tags).Error; err != nil {
		return nil, err
	}
	return []models.Tag{}, nil
}

// GetTagKeysByParameters returns list of tag keys by requested parameters.
func (r TagRepository) GetTagKeysByParameters(
	ctx context.Context, namespaceID uint, experimentNames []string,
) ([]string, error) {
	// fetch and process tags.
	query := r.GetDB().WithContext(ctx).Model(
		&models.Tag{},
	).Joins(
		"JOIN runs USING(run_uuid)",
	).Joins(
		"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
		namespaceID,
	).Where(
		"runs.lifecycle_stage = ?", models.LifecycleStageActive,
	)
	if len(experimentNames) != 0 {
		query = query.Where("experiments.name IN ?", experimentNames)
	}

	var keys []string
	if err := query.Pluck("Key", &keys).Error; err != nil {
		return nil, eris.Wrap(err, "error getting tag keys by parameters")
	}
	return keys, nil
}
