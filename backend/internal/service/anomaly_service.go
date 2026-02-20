package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"gorm.io/gorm"
)

type AnomalyService struct {
	repo       *repository.AnomalyRepository
	walletRepo *repository.WalletRepository
	tradeRepo  *repository.TradeRepository
	infoEdge   *InfoEdgeService
}

type AnomalyListQuery struct {
	Page         int
	PageSize     int
	Severity     *int16
	AlertType    string
	Acknowledged *bool
}

type AnomalyListResult struct {
	Items      []model.AnomalyAlert `json:"items"`
	Pagination Pagination           `json:"pagination"`
}

func NewAnomalyService(repo *repository.AnomalyRepository, walletRepo *repository.WalletRepository, tradeRepo *repository.TradeRepository, infoEdge *InfoEdgeService) *AnomalyService {
	return &AnomalyService{repo: repo, walletRepo: walletRepo, tradeRepo: tradeRepo, infoEdge: infoEdge}
}

func (s *AnomalyService) List(ctx context.Context, q AnomalyListQuery) (*AnomalyListResult, error) {
	rows, total, err := s.repo.List(ctx, repository.AnomalyListFilter{
		Severity:     q.Severity,
		AlertType:    q.AlertType,
		Acknowledged: q.Acknowledged,
		Limit:        q.PageSize,
		Offset:       (q.Page - 1) * q.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &AnomalyListResult{
		Items: rows,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *AnomalyService) Acknowledge(ctx context.Context, id int64) error {
	return s.repo.MarkAcknowledged(ctx, id)
}

func (s *AnomalyService) GetByID(ctx context.Context, id int64) (*model.AnomalyAlert, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (s *AnomalyService) Scan(ctx context.Context) error {
	walletIDs, err := s.walletRepo.ListIDs(ctx)
	if err != nil {
		return err
	}
	for _, walletID := range walletIDs {
		summary, err := s.tradeRepo.AggregateByWalletID(ctx, walletID)
		if err != nil {
			continue
		}
		if summary.Volume30D > 50000 {
			exists, err := s.repo.ExistsRecent(ctx, walletID, "large_position", 6)
			if err == nil && !exists {
				evidence, _ := json.Marshal(map[string]any{"volume_30d": summary.Volume30D, "threshold": 50000})
				_ = s.repo.Create(ctx, &model.AnomalyAlert{
					WalletID:    walletID,
					AlertType:   "large_position",
					Severity:    2,
					Evidence:    evidence,
					Description: fmt.Sprintf("wallet %d has unusually high 30d volume", walletID),
				})
			}
		}
		if summary.TotalTrades > 150 {
			exists, err := s.repo.ExistsRecent(ctx, walletID, "high_frequency", 6)
			if err == nil && !exists {
				evidence, _ := json.Marshal(map[string]any{"trade_count": summary.TotalTrades, "threshold": 150})
				_ = s.repo.Create(ctx, &model.AnomalyAlert{
					WalletID:    walletID,
					AlertType:   "high_frequency",
					Severity:    1,
					Evidence:    evidence,
					Description: fmt.Sprintf("wallet %d high trade frequency detected", walletID),
				})
			}
		}
		if summary.TradingPnL > 10000 {
			exists, err := s.repo.ExistsRecent(ctx, walletID, "pnl_spike", 6)
			if err == nil && !exists {
				evidence, _ := json.Marshal(map[string]any{"trading_pnl": summary.TradingPnL, "threshold": 10000})
				_ = s.repo.Create(ctx, &model.AnomalyAlert{
					WalletID:    walletID,
					AlertType:   "pnl_spike",
					Severity:    2,
					Evidence:    evidence,
					Description: fmt.Sprintf("wallet %d pnl spike observed", walletID),
				})
			}
		}
		info, err := s.infoEdge.Evaluate(ctx, walletID)
		if err == nil && info.Samples >= 10 && info.PValue < 0.05 && info.MeanDeltaMinutes <= -30 {
			exists, err := s.repo.ExistsRecent(ctx, walletID, "pre_event_timing", 6)
			if err == nil && !exists {
				evidence, _ := json.Marshal(map[string]any{
					"mean_delta_minutes": info.MeanDeltaMinutes,
					"p_value":            info.PValue,
					"samples":            info.Samples,
				})
				_ = s.repo.Create(ctx, &model.AnomalyAlert{
					WalletID:    walletID,
					AlertType:   "pre_event_timing",
					Severity:    3,
					Evidence:    evidence,
					Description: fmt.Sprintf("wallet %d shows significant pre-event timing edge", walletID),
				})
			}
		}
	}
	return nil
}

func ParseOptionalInt16(input string) (*int16, error) {
	if input == "" {
		return nil, nil
	}
	v, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		return nil, err
	}
	value := int16(v)
	return &value, nil
}

func ParseOptionalBool(input string) (*bool, error) {
	if input == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(input)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
