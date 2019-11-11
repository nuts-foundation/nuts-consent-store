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
	"errors"
	"fmt"
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
	"time"
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
	Db    *gorm.DB
	sqlDb *sql.DB

	ConfigOnce sync.Once
	Config     ConsentStoreConfig
}

var instance *ConsentStore
var oneEngine sync.Once

// ConsentStoreClient defines all actions possible through a direct connection, command-line and REST api
type ConsentStoreClient interface {
	// ConsentAuth checks if a record exists in the Db for the given combination and returns a bool. Checkpoint is optional and default to time.Now()
	ConsentAuth(context context.Context, custodian string, subject string, actor string, resourceType string, checkpoint *time.Time) (bool, error)
	// RecordConsent records a record in the Db, this is not to be used to create a new distributed consent record. It's only valid for the local node.
	// It should only be called by the consent logic component (or for development purposes)
	RecordConsent(context context.Context, consent []PatientConsent) error
	// QueryConsentForActor can be used to perform full text searches on the backend. Scoped on actor only.
	QueryConsentForActor(context context.Context, actor string, query string) ([]PatientConsent, error)
	// QueryConsentForActorAndSubject can be used to list the custodians and resources for a given Actor and Subject.
	QueryConsentForActorAndSubject(context context.Context, actor string, subject string) ([]PatientConsent, error)
	// QueryConsent can be used to query consent from a custodian/actor point of view.
	QueryConsent(context context.Context, actor *string, custodian *string, subject *string) ([]PatientConsent, error)
	// DeleteConsentRecordByHash removes a ConsentRecord from the db. Returns true if the record was found and deleted.
	DeleteConsentRecordByHash(context context.Context, proofHash string) (bool, error)
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

func Logger() *logrus.Entry {
	return logrus.StandardLogger().WithField("module", "consent-store")
}

func (cs *ConsentStore) Configure() error {
	var (
		err error
	)

	cs.ConfigOnce.Do(func() {
		if cs.Config.Mode == "server" {
			cs.sqlDb, err = sql.Open("sqlite3", cs.Config.Connectionstring)
			if err != nil {
				return
			}

			// 1 ping
			err = cs.sqlDb.Ping()
			if err != nil {
				return
			}

			// migrate
			err = cs.RunMigrations(cs.sqlDb)
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
	cs.Db, err = gorm.Open("sqlite3", cs.sqlDb)

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

// ConsentAuth checks if there is a consent for a given custodian, subject and actor for a certain resource at a given moment in time (checkpoint)
func (cs *ConsentStore) ConsentAuth(context context.Context, custodian string, subject string, actor string, resourceType string, checkpoint *time.Time) (bool, error) {
	target := &PatientConsent{}

	cp := time.Now()

	if checkpoint != nil {
		cp = *checkpoint
	}

	// this will always fill target, but if a record does not exist, resources will be empty
	var tdb = cs.Db.Debug()
	tdb = tdb.Table("patient_consent")
	tdb = tdb.Joins("JOIN consent_record ON consent_record.patient_consent_id = patient_consent.id")
	tdb = tdb.Preload("Records.Resources")
	tdb = tdb.Where("custodian = ? AND subject = ? AND actor = ?", custodian, subject, actor)
	tdb = tdb.Where("consent_record.valid_from <= ?", cp)
	tdb = tdb.Where("consent_record.valid_to > ?", cp)

	if err := tdb.FirstOrInit(&target).Error; err != nil {
		return false, err
	}

	for _, n := range target.Resources() {
		if resourceType == n.ResourceType {
			return true, nil
		}
	}

	return false, nil
}

// RecordConsent records a list of PatientConsents, their records and their resources.
func (cs *ConsentStore) RecordConsent(context context.Context, consent []PatientConsent) error {

	// start transaction
	tx := cs.Db.Begin().Debug()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, pr := range consent {
		if pr.ID == "" {
			return fmt.Errorf("id of patient consent cannot be empty")
		}
		tpc := PatientConsent{
			ID:        pr.ID,
			Actor:     pr.Actor,
			Custodian: pr.Custodian,
			Subject:   pr.Subject,
		}

		// first check if a consent record exists for subject, custodian and actor, if not create
		if err := tx.Where(tpc).Preload("Records").FirstOrCreate(&tpc).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, cr := range pr.Records {
			tcr := ConsentRecord{
				PatientConsentID: tpc.ID,
				Hash:             cr.Hash,
				ValidFrom:        cr.ValidFrom,
				ValidTo:          cr.ValidTo,
				UUID: 			  cr.UUID,
				Version:          cr.Version,
			}

			if !tcr.ValidTo.After(tcr.ValidFrom) {
				tx.Rollback()
				return errors.New("ConsentRecord validation failed: ValidTo must come after ValidFrom")
			}

			// Save all current resources
			tcr.Resources = cr.Resources
			if err := tx.Save(&tcr).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}




// QueryConsentForActor returns all PatientConsents for a given actor
func (cs *ConsentStore) QueryConsentForActor(context context.Context, actor string, query string) ([]PatientConsent, error) {
	var records []uint

	rows, err := cs.Db.Debug().Where("actor = ?", actor).
		Table("patient_consent").
		//Preload("Records").Preload("Records.Resources").
		Select("consent_record.id").
		Joins("left join consent_record on consent_record.patient_consent_id = patient_consent.id").
		Group("consent_record.uuid").Having("max(consent_record.version)").
		Rows()

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rId uint

		err = rows.Scan(&rId)
		if err != nil {
			return nil, err
		}

		records = append(records, rId)
	}

	// new queries can only be done after rows has been closed....
	rows.Close()

	return cs.patientConsentByConsentRecord(context, records)
}

func (cs *ConsentStore) patientConsentByConsentRecord(context context.Context, records []uint) ([]PatientConsent, error) {
	var consentMap = make(map[string]*PatientConsent)

	for _, ri := range records {
		var cr ConsentRecord
		if err := cs.Db.Debug().Where("id = ?", ri).Preload("Resources").Find(&cr).Error; err != nil {
			return nil, err
		}

		cpc := consentMap[cr.PatientConsentID]
		if cpc == nil {
			var pc PatientConsent
			if err := cs.Db.Debug().Where("id = ?", cr.PatientConsentID).Find(&pc).Error; err != nil {
				return nil, err
			}
			cpc = &pc
			consentMap[cr.PatientConsentID] = &pc
		}
		cpc.Records = append(cpc.Records, cr)
	}

	var consentList []PatientConsent

	for _, v := range consentMap {
		consentList = append(consentList, *v)
	}

	return consentList, nil
}

// QueryConsentForActorAndSubject  returns all PatientConsents for a given actor and subject
func (cs *ConsentStore) QueryConsentForActorAndSubject(context context.Context, actor string, subject string) ([]PatientConsent, error) {
	var records []uint

	rows, err := cs.Db.Debug().Where("actor = ? AND subject = ?", actor, subject).
		Table("patient_consent").
		//Preload("Records").Preload("Records.Resources").
		Select("consent_record.id").
		Joins("left join consent_record on consent_record.patient_consent_id = patient_consent.id").
		Group("consent_record.uuid").Having("max(consent_record.version)").
		Rows()

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rId uint

		err = rows.Scan(&rId)
		if err != nil {
			return nil, err
		}

		records = append(records, rId)
	}

	// new queries can only be done after rows has been closed....
	rows.Close()

	return cs.patientConsentByConsentRecord(context, records)
}

// QueryConsent accepts actor, custodian and subject, if these are nil, it uses a wildcard to query.
func (cs *ConsentStore) QueryConsent(context context.Context, _actor *string, _custodian *string, _subject *string) ([]PatientConsent, error) {
	var (
		actor, custodian, subject string
	)

	if actor = "%"; _actor != nil {
		actor = *_actor
	}

	if custodian = "%"; _custodian != nil {
		custodian = *_custodian
	}

	if subject = "%"; _subject != nil {
		subject = *_subject
	}

	var records []uint

	rows, err := cs.Db.Debug().Where("actor LIKE ? AND subject LIKE ? AND custodian LIKE ?", actor, subject, custodian).
		Table("patient_consent").
		//Preload("Records").Preload("Records.Resources").
		Select("consent_record.id").
		Joins("left join consent_record on consent_record.patient_consent_id = patient_consent.id").
		Group("consent_record.uuid").Having("max(consent_record.version)").
		Rows()

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rId uint

		err = rows.Scan(&rId)
		if err != nil {
			return nil, err
		}

		records = append(records, rId)
	}

	// new queries can only be done after rows has been closed....
	rows.Close()

	return cs.patientConsentByConsentRecord(context, records)
}

// DeleteConsentRecordByHash deletes a consent record by its hash. Returns boolean to indicate the success of the operation
func (cs *ConsentStore) DeleteConsentRecordByHash(context context.Context, proofHash string) (bool, error) {
	record := ConsentRecord{}

	if err := cs.Db.Debug().Where("hash = ?", proofHash).First(&record).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	if err := cs.Db.Debug().Delete(&record).Error; err != nil {
		return false, err
	}

	return true, nil
}
