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

package consent

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"sync"
)

type consentStoreConfig struct{}

type DefaultConsentStore struct {
	connectionString string
	db              *gorm.DB

	configOnce sync.Once
	Config consentStoreConfig
}

var instance *DefaultConsentStore
var oneEngine sync.Once

func ConsentStore() *DefaultConsentStore {
	oneEngine.Do(func() {
		instance = &DefaultConsentStore{
			connectionString: "file:test.db?cache=shared",
		}
	})

	return instance
}

func (cs *DefaultConsentStore) ConsentAuth(context context.Context, consentRule pkg.ConsentRule, resourceType string) (bool, error) {
	target := &pkg.ConsentRule{}
	copy := pkg.ConsentRule{
		Actor: consentRule.Actor,
		Custodian: consentRule.Custodian,
		Subject: consentRule.Subject,
	}

	// this will always fill target, but if a record does not exist, resources will be empty
	if err := cs.db.Table("consent_rule").Where(copy).Preload("Resources").FirstOrInit(&target).Error; err != nil {
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
		tx.Rollback()
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

func (cs *DefaultConsentStore) QueryConsentForActor(context context.Context, actor string, query string) ([]pkg.ConsentRule, error) {
	var rules []pkg.ConsentRule

	if err := cs.db.Where("Actor = ?", actor).Preload("Resources").Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}

func (cs *DefaultConsentStore) QueryConsentForActorAndSubject(context context.Context, actor string, subject string) ([]pkg.ConsentRule, error) {
	var rules []pkg.ConsentRule

	if err := cs.db.Where("Actor = ? AND Subject = ?", actor, subject).Preload("Resources").Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}
