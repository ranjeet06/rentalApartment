// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// GORSK - Go(lang) restful starter kit
//
// API Docs for GORSK v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 2.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package api

import (
	"context"
	"crypto/sha1"
	_ "database/sql"
	"github.com/ribice/gorsk/pkg/apartmentUserUtl/apartmentUserJwt"
	"github.com/ribice/gorsk/pkg/api/apartment/apartment_service"
	"github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_service"
	"github.com/ribice/gorsk/pkg/api/auth"
	al "github.com/ribice/gorsk/pkg/api/auth/logging"
	at "github.com/ribice/gorsk/pkg/api/auth/transport"
	"github.com/ribice/gorsk/pkg/api/cache/redis"
	"github.com/ribice/gorsk/pkg/api/car/service"
	"github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_service"
	"github.com/ribice/gorsk/pkg/utl/jwt"
	authMw "github.com/ribice/gorsk/pkg/utl/middleware/auth"
	"gorm.io/gorm"
	"log"
	"os"
	"time"

	"github.com/ribice/gorsk/pkg/utl/zlog"

	"github.com/ribice/gorsk/pkg/api/password"
	pl "github.com/ribice/gorsk/pkg/api/password/logging"
	pt "github.com/ribice/gorsk/pkg/api/password/transport"
	"github.com/ribice/gorsk/pkg/api/user"
	ul "github.com/ribice/gorsk/pkg/api/user/logging"
	ut "github.com/ribice/gorsk/pkg/api/user/transport"

	cl "github.com/ribice/gorsk/pkg/api/car/logging"
	ct "github.com/ribice/gorsk/pkg/api/car/transport"

	apartmentTransport "github.com/ribice/gorsk/pkg/api/apartment/apartment_controller"
	apartmentUserTransport "github.com/ribice/gorsk/pkg/api/apartment_user/apartment_user_controller"
	favApartment "github.com/ribice/gorsk/pkg/api/user_favorite_apartment/user_favorite_apartment_controller"

	apartmentUserAuthMw "github.com/ribice/gorsk/pkg/apartmentUserUtl/middleware/auth"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/postgres"
	"github.com/ribice/gorsk/pkg/utl/rbac"
	"github.com/ribice/gorsk/pkg/utl/secure"
	pgtransaction "gorm.io/driver/postgres"

	"github.com/ribice/gorsk/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.New(os.Getenv("DATABASE_URL"), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	dbURL := "postgresql://ranjeet:3298@localhost:5432/rental?sslmode=disable"

	newDb, err := gorm.Open(pgtransaction.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.Service{}

	apartmentUserJwt, err := apartmentUserJwt.JwtNew(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}
	dummyJwt, err := apartmentUserJwt.GetDummyToken()
	if err != nil {
		return err
	}
	_ = dummyJwt

	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/", cfg.App.SwaggerUIPath)

	authMiddleware := authMw.Middleware(jwt)

	apartmentUserMiddleware := apartmentUserAuthMw.Middleware(apartmentUserJwt)

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec, rbac), log), e, authMiddleware)

	v1 := e.Group("/v1")
	v1.Use(authMiddleware)

	ut.NewHTTP(ul.New(user.Initialize(db, rbac, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, rbac, sec), log), v1)
	v2 := e.Group("/v2")
	v2.Use(apartmentUserMiddleware)
	ct.NewHTTP(cl.New(service.Initialize(db)), v2)

	apartmentService := apartment_service.Initialize(db, newDb)
	apartmentUserService := apartment_user_service.InitializeUser(db)
	favApartmentService := user_favorite_apartment_service.Initialize(db, apartmentService, apartmentUserService)
	ctx := context.Background()
	cacheService := redis.InitializeRedisCache("localhost:6379", 900*time.Second, ctx)
	apartmentTransport.NewHTTP(apartmentService, cacheService, v2)
	apartmentUserTransport.NewHTTP(apartmentUserService, cacheService, v2, apartmentUserMiddleware)
	favApartment.NewHTTP(favApartmentService, cacheService, v2)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
