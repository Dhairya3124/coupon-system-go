package service

import (
	"context"
	"sync"
	"time"

	"github.com/Dhairya3124/coupon-system-go/internal/cache"
	"github.com/Dhairya3124/coupon-system-go/internal/model"
	"github.com/Dhairya3124/coupon-system-go/internal/repository"
)

type couponService struct {
	mu    sync.RWMutex
	repo  repository.Coupon
	cache cache.Cache
}
type CouponService interface {
	CreateCouponService(ctx context.Context, coupon *model.Coupon) error
	ValidateCouponService(ctx context.Context, code string, cart *model.Cart) (bool, error)
	GetApplicableCoupons(ctx context.Context, cart *model.Cart) ([]*model.Coupon, error)
}

func NewCouponService(repo repository.Coupon, cache cache.Cache) CouponService {
	return &couponService{repo: repo, cache: cache}
}
func (s *couponService) CreateCouponService(ctx context.Context, coupon *model.Coupon) error {
	s.cache.Delete(generateCacheKey("applicable", nil))
	return s.repo.Create(ctx, coupon)

}
func (s *couponService) ValidateCouponService(ctx context.Context, code string, cart *model.Cart) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cacheKey := generateCacheKey("validatekey", code, cart)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(bool), nil
	}

	coupon, err := s.repo.GetCouponByCode(ctx, code)
	if err != nil {
		return false, err
	}
	if coupon == nil {
		return false, nil
	}
	now := time.Now()
	if !coupon.IsActive {
		return false, nil
	}
	if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
		return false, nil
	}
	if cart.Total < coupon.MinOrderValue {
		return false, nil
	}
	if coupon.UsageType == "one_time" {
		if coupon.UsageCount >= coupon.UsageLimit {
			return false, nil
		}
	}
	hasApplicableItem := false
	for _, item := range cart.Items {
		for _, applicableItem := range coupon.ApplicableMedicineIDs {
			if item.ID == applicableItem {
				hasApplicableItem = true
				break
			}
		}
		if hasApplicableItem {
			break
		}
	}
	if !hasApplicableItem {
		for _, item := range cart.Items {
			for _, applicableItem := range coupon.ApplicableCategories {
				if item.ID == applicableItem {
					hasApplicableItem = true
					break
				}
			}
			if hasApplicableItem {
				break
			}
		}
	}

	if !hasApplicableItem {
		return false, nil
	}

	coupon.UsageCount++
	if err := s.repo.UpdateCoupon(ctx, coupon); err != nil {
		return false, err
	}

	s.cache.Set(cacheKey, true)

	return true, nil

}
func (s *couponService) GetApplicableCoupons(ctx context.Context, cart *model.Cart) ([]*model.Coupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cacheKey := generateCacheKey("applicable", cart)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]*model.Coupon), nil
	}

	coupons, err := s.repo.GetAllCoupons(ctx)
	if err != nil {
		return nil, err
	}

	applicableCoupons := make([]*model.Coupon, 0)
	now := time.Now()

	for _, coupon := range coupons {
		if !coupon.IsActive {
			continue
		}

		if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
			continue
		}

		if cart.Total < coupon.MinOrderValue {
			continue
		}

		if coupon.UsageCount >= coupon.UsageLimit {
			continue
		}

		hasApplicableItem := false
		for _, item := range cart.Items {
			for _, applicableItem := range coupon.ApplicableMedicineIDs {
				if item.ID == applicableItem {
					hasApplicableItem = true
					break
				}
			}
			if hasApplicableItem {
				break
			}
		}

		if hasApplicableItem {
			applicableCoupons = append(applicableCoupons, coupon)
		}
		if !hasApplicableItem {
			for _, item := range cart.Items {
				for _, applicableItem := range coupon.ApplicableCategories {
					if item.ID == applicableItem {
						hasApplicableItem = true
						break
					}
				}
				if hasApplicableItem {
					break
				}
			}
		}
		if hasApplicableItem {
			applicableCoupons = append(applicableCoupons, coupon)
		}

	}

	s.cache.Set(cacheKey, applicableCoupons)

	return applicableCoupons, nil
}
func generateCacheKey(prefix string, params ...interface{}) string {
	return prefix
}
