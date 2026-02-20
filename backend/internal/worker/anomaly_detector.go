package worker

import (
	"context"

	"easy-arbitra/backend/internal/service"
)

type AnomalyDetector struct {
	service *service.AnomalyService
}

func NewAnomalyDetector(service *service.AnomalyService) *AnomalyDetector {
	return &AnomalyDetector{service: service}
}

func (d *AnomalyDetector) Name() string { return "anomaly_detector" }

func (d *AnomalyDetector) RunOnce(ctx context.Context) error {
	return d.service.Scan(ctx)
}
