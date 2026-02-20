package worker

import (
	"context"
	"time"

	"easy-arbitra/backend/internal/repository"
)

type FeatureBuilder struct {
	repo *repository.FeatureRepository
}

func NewFeatureBuilder(repo *repository.FeatureRepository) *FeatureBuilder {
	return &FeatureBuilder{repo: repo}
}

func (b *FeatureBuilder) Name() string { return "feature_builder" }

func (b *FeatureBuilder) RunOnce(ctx context.Context) error {
	return b.repo.BuildDaily(ctx, time.Now().UTC())
}
