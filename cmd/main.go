package main

import (
	"log"
	"net/http"

	"github.com/Dhairya3124/coupon-system-go/internal/api"
	"github.com/Dhairya3124/coupon-system-go/internal/cache"
	"github.com/Dhairya3124/coupon-system-go/internal/db"
	"github.com/Dhairya3124/coupon-system-go/internal/repository"
	"github.com/Dhairya3124/coupon-system-go/internal/service"
	"github.com/go-chi/chi"

	_ "github.com/Dhairya3124/coupon-system-go/cmd/docs"
	// swaggerFiles "github.com/swaggo/files" // swagger embed files
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Coupon System API
// @version 1.0
// @description This is a coupon system API server.
// @host localhost:3000
// @BasePath /

func main() {
	dbConn := db.NewDB()
	repo := repository.NewCouponRepository(dbConn.DB)
	cache := cache.NewLRU(100)
	serv := service.NewCouponService(repo, cache)
	handler := api.NewCouponHandler(serv)
	r := chi.NewRouter()
	// docsURL := "/swagger/doc.json"

	r.Route("/v1", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:3000/swagger/doc.json"), //The url pointing to API definition
		))
		r.Route("/coupons", func(r chi.Router) {
			r.Post("/", handler.CreateCouponHandler)
			r.Get("/applicable", handler.ApplicableCouponHandler)
			r.Post("/validate", handler.ValidateCouponHandler)
		})

	})
	srv := http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	log.Fatal(srv.ListenAndServe())

}
