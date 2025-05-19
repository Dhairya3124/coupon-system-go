package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dhairya3124/coupon-system-go/internal/model"
	"github.com/Dhairya3124/coupon-system-go/internal/service"
)

type CouponHandler struct {
	service service.CouponService
}
type ValidationRequest struct {
	CouponCode string     `json:"code"`
	Cart       model.Cart `json:"cart"`
}
type ValidationResponse struct {
	IsValid bool `json:"is_valid"`
}
type GetApplicableCouponsRequest struct {
	Items []model.CartItem `json:"items"`
	Total float64          `json:"total"`
}

func NewCouponHandler(service service.CouponService) *CouponHandler {
	return &CouponHandler{
		service: service,
	}
}

// CreateCouponHandler handles requests to create a coupon
// @Summary Create coupon
// @Description Create a new coupon
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body CreateCouponRequest true "Coupon details"
// @Success 201 {object} model.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/ [post]
func (h *CouponHandler) CreateCouponHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	coupon := &model.Coupon{
		Code:                  req.Code,
		DiscountType:          req.DiscountType,
		DiscountValue:         req.DiscountValue,
		MinOrderValue:         req.MinOrderValue,
		MaxDiscount:           req.MaxDiscount,
		StartDate:             req.StartDate,
		EndDate:               req.EndDate,
		UsageLimit:            req.UsageLimit,
		IsActive:              req.IsActive,
		ApplicableCategories:  req.ApplicableCategories,
		ApplicableMedicineIDs: req.ApplicableMedicineIDs,
		UsageType:             req.UsageType,
	}

	err := h.service.CreateCouponService(ctx, coupon)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}

// ValidateCouponHandler handles requests to validate a coupon
// @Summary Validate coupon
// @Description Validate a coupon against cart items
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body ValidateCouponRequest true "Coupon code and cart"
// @Success 200 {object} ValidateCouponResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/validate [post]
func (h *CouponHandler) ValidateCouponHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	isValid, err := h.service.ValidateCouponService(ctx, req.CouponCode, &req.Cart)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	resp := ValidationResponse{IsValid: isValid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

// ApplicableCouponHandler handles requests to get applicable coupons
// @Summary Get applicable coupons
// @Description Get coupons applicable to the given cart
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body GetApplicableCouponsRequest true "Cart items and total"
// @Success 200 {array} model.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/applicable [get]
func (h *CouponHandler) ApplicableCouponHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req GetApplicableCouponsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	cart := &model.Cart{
		Items: req.Items,
		Total: req.Total,
	}
	coupons, err := h.service.GetApplicableCoupons(ctx, cart)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coupons)

}
