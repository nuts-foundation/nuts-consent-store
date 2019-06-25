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

package pkg

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-consent-store/migrations"
	"github.com/sirupsen/logrus"
	"sync"
)

type ConsentStoreConfig struct {
	Connectionstring string
	Mode             string
	Address          string
}

const ConfigConnectionString = "connectionstring"
const ConfigMode = "mode"
const ConfigAddress = "address"
const ConfigConnectionStringDefault = ":memory:"

type ConsentStore struct {
	Db *gorm.DB

	ConfigOnce sync.Once
	Config     ConsentStoreConfig
}

var instance *ConsentStore
var oneEngine sync.Once

// ConsentStoreClient defines all actions possible through a direct connection, command-line and REST api
type ConsentStoreClient interface {
	// ConsentAuth checks if a record exists in the Db for the given combination and returns a bool.
	ConsentAuth(context context.Context, consentRule ConsentRule, resourceType string) (bool, error)
	// RecordConsent records a record in the Db, this is not to be used to create a new distributed consent record. It's only valid for the local node.
	// It should only be called by the consent logic component (or for development purposes)
	RecordConsent(context context.Context, consent []ConsentRule) error
	// QueryConsentForActor can be used to perform full text searches on the backend. Scoped on actor only.
	QueryConsentForActor(context context.Context, actor string, query string) ([]ConsentRule, error)
	// QueryConsentForActorAndSubject can be used to list the custodians and resources for a given Actor and Subject.
	QueryConsentForActorAndSubject(context context.Context, actor string, subject string) ([]ConsentRule, error)
}

func ConsentStoreInstance() *ConsentStore {
	oneEngine.Do(func() {
		instance = &ConsentStore{
			Config: ConsentStoreConfig{
				Connectionstring: ConfigConnectionStringDefault,
			},
		}
	})

	return instance
}

func (cs *ConsentStore) Configure() error {
	var (
		err error
		db  *sql.DB
	)

	cs.ConfigOnce.Do(func() {
		if cs.Config.Mode == "server" {
			db, err = sql.Open("sqlite3", cs.Config.Connectionstring)
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
			err = cs.RunMigrations(db)
			if err != nil {
				return
			}
		}

	})

	return err
}

//Shutdown closes the db connections
func (cs *ConsentStore) Shutdown() error {
	return cs.Db.Close()
}

// Start opens the db connections
func (cs *ConsentStore) Start() error {
	var err error

	// gorm db connection
	cs.Db, err = gorm.Open("sqlite3", cs.Config.Connectionstring)

	// logging
	cs.Db.SetLogger(logrus.StandardLogger())

	return err
}

// RunMigrations runs all new migrations in order
func (cs *ConsentStore) RunMigrations(db *sql.DB) error {
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

func (cs *ConsentStore) ConsentAuth(context context.Context, consentRule ConsentRule, resourceType string) (bool, error) {
	target := &ConsentRule{}
	copy := ConsentRule{
		Actor:     consentRule.Actor,
		Custodian: consentRule.Custodian,
		Subject:   consentRule.Subject,
	}

	// this will always fill target, but if a record does not exist, resources will be empty
	if err := cs.Db.Table("consent_rule").Where(copy).Preload("Resources").FirstOrInit(&target).Error; err != nil {
		return false, err
	}

	var resources []Resource
	cs.Db.Find(&resources)

	for _, n := range target.Resources {
		if resourceType == n.ResourceType {
			return true, nil
		}
	}

	return false, nil
}

func (cs *ConsentStore) RecordConsent(context context.Context, consent []ConsentRule) error {

	// start transaction
	tx := cs.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, cr := range consent {
		tcr := ConsentRule{
			Actor:     cr.Actor,
			Custodian: cr.Custodian,
			Subject:   cr.Subject,
		}

		// first check if a consent record exists for subject, custodian and actor, if not create
		if err := tx.Where(tcr).FirstOrCreate(&tcr).Error; err != nil {
			tx.Rollback()
			return err
		}

		tcr.Resources = cr.Resources
		if err := tx.Save(&tcr).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (cs *ConsentStore) QueryConsentForActor(context context.Context, actor string, query string) ([]ConsentRule, error) {
	var rules []ConsentRule

	if err := cs.Db.Debug().Where("Actor = ?", actor).Preload("Resources").Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}

func (cs *ConsentStore) QueryConsentForActorAndSubject(context context.Context, actor string, subject string) ([]ConsentRule, error) {
	var rules []ConsentRule

	if err := cs.Db.Debug().Where("Actor = ? AND Subject = ?", actor, subject).Preload("Resources").Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}
