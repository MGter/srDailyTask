package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
)

type LongTermItemService struct {
	repo *repository.LongTermItemRepository
}

func NewLongTermItemService() *LongTermItemService {
	return &LongTermItemService{repo: repository.NewLongTermItemRepository()}
}

func (s *LongTermItemService) Create(req *model.CreateLongTermItemRequest) (*model.LongTermItem, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if req.Price <= 0 {
		return nil, errors.New("price must be positive")
	}
	if req.PurchaseDate.IsZero() {
		return nil, errors.New("purchase_date is required")
	}

	now := time.Now()
	item := &model.LongTermItem{
		UserID:       req.UserID,
		Name:         name,
		Price:        roundMoney(req.Price),
		PurchaseDate: startOfDay(req.PurchaseDate),
		Status:       model.LongTermItemActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	s.fillCostFields(item, startOfDay(now))
	return item, nil
}

func (s *LongTermItemService) GetByUserID(userID uint64) ([]*model.LongTermItem, *model.LongTermSummary, error) {
	items, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	summary := &model.LongTermSummary{}
	today := startOfDay(time.Now())
	for _, item := range items {
		s.fillCostFields(item, today)
		if item.Status == model.LongTermItemActive {
			summary.ActiveCount++
			summary.ActiveDailyCost += item.DailyCost
		} else {
			summary.ScrappedCount++
		}
	}
	summary.ActiveDailyCost = roundMoney(summary.ActiveDailyCost)
	return items, summary, nil
}

func (s *LongTermItemService) Update(id uint64, req *model.UpdateLongTermItemRequest) error {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if item.UserID != req.UserID {
		return errors.New("item not found or not owned by user")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("name is required")
	}
	if req.Price <= 0 {
		return errors.New("price must be positive")
	}
	if req.PurchaseDate.IsZero() {
		return errors.New("purchase_date is required")
	}

	purchaseDate := startOfDay(req.PurchaseDate)
	item.Name = name
	item.Price = roundMoney(req.Price)
	item.PurchaseDate = purchaseDate
	item.FrozenDailyCost = nil
	if item.Status == model.LongTermItemScrapped && item.ScrapDate != nil {
		scrapDate := startOfDay(*item.ScrapDate)
		if scrapDate.Before(purchaseDate) {
			return errors.New("purchase_date cannot be after scrap_date")
		}
		ownedDays := daysBetweenInclusive(purchaseDate, scrapDate)
		frozenDailyCost := roundMoney(item.Price / float64(ownedDays))
		item.FrozenDailyCost = &frozenDailyCost
	}

	return s.repo.Update(item)
}

func (s *LongTermItemService) Scrap(id uint64, req *model.ScrapLongTermItemRequest) error {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if item.UserID != req.UserID {
		return errors.New("item not found or not owned by user")
	}
	if item.Status == model.LongTermItemScrapped {
		return errors.New("item already scrapped")
	}

	scrapDate := req.ScrapDate
	if scrapDate.IsZero() {
		scrapDate = time.Now()
	}
	scrapDate = startOfDay(scrapDate)
	purchaseDate := startOfDay(item.PurchaseDate)
	if scrapDate.Before(purchaseDate) {
		return errors.New("scrap_date cannot be before purchase_date")
	}

	ownedDays := daysBetweenInclusive(purchaseDate, scrapDate)
	frozenDailyCost := roundMoney(item.Price / float64(ownedDays))
	return s.repo.Scrap(id, req.UserID, scrapDate, frozenDailyCost)
}

func (s *LongTermItemService) Delete(id uint64, userID uint64) error {
	return s.repo.Delete(id, userID)
}

func (s *LongTermItemService) fillCostFields(item *model.LongTermItem, today time.Time) {
	purchaseDate := startOfDay(item.PurchaseDate)
	endDate := today
	if item.Status == model.LongTermItemScrapped && item.ScrapDate != nil {
		endDate = startOfDay(*item.ScrapDate)
	}
	item.OwnedDays = daysBetweenInclusive(purchaseDate, endDate)
	if item.Status == model.LongTermItemScrapped && item.FrozenDailyCost != nil {
		item.DailyCost = roundMoney(*item.FrozenDailyCost)
		return
	}
	item.DailyCost = roundMoney(item.Price / float64(item.OwnedDays))
}

func startOfDay(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
}

func daysBetweenInclusive(start time.Time, end time.Time) int {
	start = startOfDay(start)
	end = startOfDay(end)
	if end.Before(start) {
		return 1
	}
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		return 1
	}
	return days
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
