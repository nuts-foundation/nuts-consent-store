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
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jinzhu/gorm"
	// import needed to enable the sqlite dialect for gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-consent-store/migrations"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// ConsentStoreConfig holds the config for the consent store
type ConsentStoreConfig struct {
	Connectionstring string
	Mode             string
	Address          string
}

// ConfigConnectionString is the config name for the connection string
const ConfigConnectionString = "connectionstring"
// ConfigMode is the config name for the mode of the store (server, client)
const ConfigMode = "mode"
// ConfigAddress is the config name for the api address when running in client mode
const ConfigAddress = "address"
// ConfigConnectionStringDefault is the default db connection string
const ConfigConnectionStringDefault = ":memory:"

// ConsentStore is the main data struct holding the config and references to the DB
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
	// QueryConsent can be used to query consent from a custodian/actor point of view.
	QueryConsent(context context.Context, actor *string, custodian *string, subject *string) ([]PatientConsent, error)
	// DeleteConsentRecordByHash removes a ConsentRecord from the db. Returns true if the record was found and deleted.
	DeleteConsentRecordByHash(context context.Context, proofHash string) (bool, error)
	// FindConsentRecordByHash find a consent record given its hash, the latest flag indicates the requirement if the record is the latest in the chain.
	FindConsentRecordByHash(context context.Context, proofHash string, latest bool) (ConsentRecord, error)
}

// ConsentStoreInstance returns a singleton consent store
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

// Logger returns the standard logger with a module field
func Logger() *logrus.Entry {
	return logrus.StandardLogger().WithField("module", "consent-store")
}

// Configure opens a DB connection and runs migrations
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

// ErrorInvalidValidTo is returned when the ValidTo from a ConsentRecord comes before the ValidFrom
var ErrorInvalidValidTo = errors.New("ConsentRecord validation failed: ValidTo must come after ValidFrom")

// RecordConsent records a list of PatientConsents, their records and their resources.
// For consent records that are updates, this function finds the version number and UUID from the previous record
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
				UUID:             uuid.NewV4().String(),
				Version:          1,
			}

			// if this is an update to an existing entry, find UUID and version
			if cr.PreviousHash != nil {
				var pcr ConsentRecord
				if err := tx.Where("hash = ?", *cr.PreviousHash).First(&pcr).Error; err != nil {
					tx.Rollback()
					if gorm.IsRecordNotFoundError(err) {
						return ErrorNotFound
					}
					return fmt.Errorf("error when finding existing consent record for hash %s: %w", *cr.PreviousHash, err)
				}
				tcr.PreviousHash = cr.PreviousHash
				tcr.Version = pcr.Version + 1
				tcr.UUID = pcr.UUID
			}

			if !tcr.ValidTo.After(tcr.ValidFrom) {
				tx.Rollback()
				return ErrorInvalidValidTo
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

// QueryConsent accepts actor, custodian and subject, if these are nil, it's not used in the query.
func (cs *ConsentStore) QueryConsent(context context.Context, _actor *string, _custodian *string, _subject *string) ([]PatientConsent, error) {
	var pc PatientConsent

	if _actor != nil {
		pc.Actor = *_actor
	}

	if _custodian != nil {
		pc.Custodian = *_custodian
	}

	if _subject != nil {
		pc.Subject = *_subject
	}

	var records []uint

	rows, err := cs.Db.Debug().Where(pc).
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
		var rID uint

		err = rows.Scan(&rID)
		if err != nil {
			return nil, err
		}

		records = append(records, rID)
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
			return false, ErrorNotFound
		}
		return false, err
	}

	if err := cs.Db.Debug().Delete(&record).Error; err != nil {
		return false, err
	}

	return true, nil
}

// FindConsentRecordByHash find a consent record given its hash, the latest flag indicates the requirement if the record is the latest in the chain.
func (cs *ConsentStore) FindConsentRecordByHash(context context.Context, proofHash string, latest bool) (ConsentRecord, error) {
	var (
		record ConsentRecord
		err    error
	)

	if latest {
		err = cs.findConsentRecordByHashGrouped(context, proofHash, &record)
	} else {
		err = cs.findConsentRecordByHashExact(context, proofHash, &record)
	}

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return record, ErrorNotFound
		}
		return record, err
	}

	return record, nil
}

// ErrorConsentRecordNotLatest is returned when the latest consent record for a chain is requested but given hash is not the latest
var ErrorConsentRecordNotLatest = errors.New("consent record for given hash is not the latest in the chain")

// ErrorNotFound is the same as Gorm.IsRecordNotFound
var ErrorNotFound = errors.New("record not found")

func (cs *ConsentStore) findConsentRecordByHashGrouped(context context.Context, proofHash string, record *ConsentRecord) error {
	var id uint

	// sub query broken
	var cr ConsentRecord
	if err := cs.Db.Debug().Where("hash = ?", proofHash).First(&cr).Error; err != nil {
		return err
	}

	rows, err := cs.Db.Debug().Where("uuid = ?", cr.UUID).
		Table("consent_record").
		Select("id, hash").
		Group("uuid").Having("max(version)").
		Rows()

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var h string
		err = rows.Scan(&id, &h)
		if err != nil {
			return err
		}

		if h != proofHash {
			return ErrorConsentRecordNotLatest
		}

		if rows.NextResultSet() {
			// for future safety...
			return errors.New("BUG in findConsentRecordByHashGrouped, unique result should have been given")
		}
	}

	// new queries can only be done after rows has been closed....
	rows.Close()

	return cs.Db.Debug().Where("id = ?", id).First(record).Error
}

func (cs *ConsentStore) findConsentRecordByHashExact(context context.Context, proofHash string, record *ConsentRecord) error {
	return cs.Db.Debug().Where("hash = ?", proofHash).First(record).Error
}
