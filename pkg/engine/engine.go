/*
 * Nuts consent store
 * Copyright (C) 2019. Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-consent-store/migrations"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/nuts-foundation/nuts-consent-store/pkg/generated"
	types "github.com/nuts-foundation/nuts-crypto/pkg"
	engine "github.com/nuts-foundation/nuts-go/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type ConsentStoreEngine interface {
	ConsentStoreClient
	generated.ServerInterface
	engine.Engine
}

type DefaultConsentStore struct {
	connectionString string
	db              *gorm.DB
}

var instance *DefaultConsentStore
var oneEngine sync.Once

func NewConsentStoreEngine() ConsentStoreEngine {
	oneEngine.Do(func() {
		instance = &DefaultConsentStore{
			connectionString: "file:test.db?cache=shared",
		}
	})

	return instance
}

func (cs *DefaultConsentStore) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consent-store",
		Short: "consent store commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "list [actor] [subject]?",
		Example: "list urn:oid:2.16.840.1.113883.2.4.6.1:00000007",
		Short:   "lists all consent records for the given actor and optional subject",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return types.Error{Msg: "requires an actor argument"}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := NewConsentStoreClient()

			var (
				consentList []pkg.ConsentRule
				err         error
			)

			if len(args) > 1 {
				consentList, err = csc.QueryConsentForActorAndSubject(context.TODO(), args[0], args[1])
			} else {
				consentList, err = csc.QueryConsentForActor(context.TODO(), args[0], "*")
			}

			if err != nil {
				fmt.Printf("Error finding consent records: %s\n", err.Error())
				return
			}

			fmt.Printf("Found %d records", len(consentList))
			println()
			for _, c := range consentList {
				fmt.Println(c.String())
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "record [subject] [custodian] [actor] [resources]",
		Example: "record urn:oid:2.16.840.1.113883.2.4.6.3:999999990 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 Observation,Patient",
		Short:   "record a new consent in store, resources is comma-separated",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return types.Error{Msg: "requires 4 arguments"}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := NewConsentStoreClient()

			resources := pkg.ResourcesFromStrings(strings.Split(args[3], ","))

			err := csc.RecordConsent(context.TODO(), []pkg.ConsentRule{
				{
					Subject:   args[0],
					Custodian: args[1],
					Actor:     args[2],
					Resources: resources,
				},
			})

			if err != nil {
				fmt.Printf("Error recording consent: %s\n", err.Error())
				return
			}

			fmt.Println("Consent recorded")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "check [subject] [custodian] [actor] [resource]",
		Example: "check urn:oid:2.16.840.1.113883.2.4.6.3:999999990 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 Observation",
		Short:   "check if consent is given for the given combination",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return types.Error{Msg: "requires 4 arguments"}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := NewConsentStoreClient()

			auth, err := csc.ConsentAuth(context.TODO(), pkg.ConsentRule{
				Subject:   args[0],
				Custodian: args[1],
				Actor:     args[2],
			}, args[3])

			if err != nil {
				fmt.Printf("Error checking consent: %s", err.Error())
				return
			}

			if auth {
				fmt.Println("Consent given")
			} else {
				fmt.Println("No consent given")
			}
		},
	})

	return cmd
}

var configOnce sync.Once
var ConfigDone bool

func (cs *DefaultConsentStore) Configure() error {
	var err error

	configOnce.Do(func() {
		db, err := sql.Open("sqlite3", cs.connectionString)
		if err != nil {
			return
		}
		defer db.Close()

		// 1 ping
		err = db.Ping()
		if err != nil {
			return
		}

		// migrate
		err = runMigrations(db)
		if err != nil {
			return
		}

	})

	return err
}

func runMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	// wrap assets into Resource
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)

	if err != nil {
		return err
	}

	// run migrations
	m, err := migrate.NewWithInstance("go-bindata", d, "test", driver)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		return err
	}

	return nil
}

func (cs *DefaultConsentStore) FlagSet() *pflag.FlagSet {
	return pflag.NewFlagSet("consent-store", pflag.ContinueOnError)
}

func (cs *DefaultConsentStore) Routes(router runtime.EchoRouter) {
	generated.RegisterHandlers(router, cs)
}

func (cs *DefaultConsentStore) Shutdown() error {
	return cs.db.Close()
}

func (cs *DefaultConsentStore) Start() error {
	var err error

	// gorm db connection
	cs.db, err = gorm.Open("sqlite3", cs.connectionString)

	return err
}

func (cs *DefaultConsentStore) ConsentAuth(context context.Context, consentRule pkg.ConsentRule, resourceType string) (bool, error) {
	target := &pkg.ConsentRule{}

	if err := cs.db.Table("consent_rule").Where(consentRule).Preload("Resources").First(&target).Error; err != nil {
		return false, err
	}

	var resources []pkg.Resource
	cs.db.Find(&resources)

	for _, n := range target.Resources {
		if resourceType == n.ResourceType {
			return true, nil
		}
	}

	return false, nil
}

func (cs *DefaultConsentStore) RecordConsent(context context.Context, consent []pkg.ConsentRule) error {

	// start transaction
	tx := cs.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, cr := range consent {
		tcr :=  pkg.ConsentRule{
			Actor: cr.Actor,
			Custodian: cr.Custodian,
			Subject: cr.Subject,
		}

		// first check if a consent record exists for subject, custodian and actor, if not create
		if err := tx.Where(tcr).FirstOrCreate(&tcr).Error; err != nil {
			return err
		}

		tcr.Resources = cr.Resources
		if err := tx.Save(&tcr).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (cs *DefaultConsentStore) QueryConsentForActor(context context.Context, actor string, query string) ([]pkg.ConsentRule, error) {
	var rules []pkg.ConsentRule

	if err := cs.db.Where("Actor = ?", actor).Preload("Resource").Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}

func (cs *DefaultConsentStore) QueryConsentForActorAndSubject(context context.Context, subject string, actor string) ([]pkg.ConsentRule, error) {
	var rules []pkg.ConsentRule

	if err := cs.db.Where("Actor = ? AND Subject = ?", actor, subject).Preload("Resource").Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}

func (cs *DefaultConsentStore) CreateConsent(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	var createRequest = &generated.SimplifiedConsent{}
	err = json.Unmarshal(buf, createRequest)

	if len(createRequest.Subject) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing subject in createRequest")
	}

	if len(createRequest.Custodian) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing custodian in createRequest")
	}

	if len(createRequest.Actors) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actors in createRequest")
	}

	if len(createRequest.Resources) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing resources in createRequest")
	}

	err = cs.RecordConsent(ctx.Request().Context(), createRequest.ToConsentRule())

	if err != nil {
		return err
	}

	return ctx.NoContent(201)
}

func (cs *DefaultConsentStore) CheckConsent(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	var checkRequest = &generated.ConsentCheckRequest{}
	err = json.Unmarshal(buf, checkRequest)

	if len(checkRequest.Subject) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing subject in checkRequest")
	}

	if len(checkRequest.Custodian) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing custodian in checkRequest")
	}

	if len(checkRequest.Actor) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actor in checkRequest")
	}

	if len(checkRequest.ResourceType) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing resourceType in checkRequest")
	}

	cr := checkRequest.ToConsentRule()
	cr.Resources = nil
	auth, err := cs.ConsentAuth(ctx.Request().Context(), cr, checkRequest.ResourceType)

	if err != nil {
		return err
	}

	authValue := "no"
	if auth {
		authValue = "true"
	}

	checkResponse := generated.ConsentCheckResponse{
		ConsentGiven: &authValue,
	}

	return ctx.JSON(200, checkResponse)
}

func (cs *DefaultConsentStore) QueryConsent(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	var checkRequest = &generated.ConsentQueryRequest{}
	err = json.Unmarshal(buf, checkRequest)

	if len(checkRequest.Actor) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actor in queryRequest")
	}

	query := checkRequest.Query.(string)

	if len(query) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query in queryRequest")
	}

	var rules []pkg.ConsentRule

	if strings.Index(query, "urn") == 0 {
		rules, err = cs.QueryConsentForActorAndSubject(ctx.Request().Context(), query, string(checkRequest.Actor))
	} else {
		rules, err = cs.QueryConsentForActor(ctx.Request().Context(), string(checkRequest.Actor), query)
	}

	if err != nil {
		return err
	}

	results, err := generated.FromSimplifiedConsentRule(rules)

	if err != nil {
		return err
	}

	return ctx.JSON(200,
		generated.ConsentQueryResponse{
			Results: results,
		})
}
