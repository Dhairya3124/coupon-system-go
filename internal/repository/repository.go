package repository

import (
	"context"
	"sync"

	"github.com/Dhairya3124/coupon-system-go/internal/model"
	"gorm.io/gorm"
)

type Coupon interface {
	Create(ctx context.Context, coupon *model.Coupon) error
	GetCouponByCode(ctx context.Context, code string) (*model.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon *model.Coupon) error
	GetAllCoupons(ctx context.Context) ([]*model.Coupon, error)
}

type couponStore struct {
	db *gorm.DB
	mu *sync.RWMutex
}

func NewCouponRepository(db *gorm.DB) Coupon {
	return &couponStore{db: db, mu: &sync.RWMutex{}}
}
func (c *couponStore) Create(ctx context.Context, coupon *model.Coupon) error {
	if err := c.db.WithContext(ctx).Create(coupon).Error; err != nil {
		return err

	}

	return nil
}

func (c *couponStore) GetCouponByCode(ctx context.Context, code string) (*model.Coupon, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var coupon model.Coupon
	if err := c.db.WithContext(ctx).Where("code = ?", code).First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &coupon, nil
}
func (c *couponStore) UpdateCoupon(ctx context.Context, coupon *model.Coupon) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.db.WithContext(ctx).Save(coupon).Error
}
func (c *couponStore) GetAllCoupons(ctx context.Context) ([]*model.Coupon, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var coupons []*model.Coupon
	if err := c.db.WithContext(ctx).Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}
