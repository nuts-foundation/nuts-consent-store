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
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Iso8601DateTime is the date format used in the API for denoting a zoned date time
const Iso8601DateTime = "2006-01-02T15:04:05-07:00"

// PatientConsent defines struct for patient_consent table.
// ID refers to the HMAC id for a custodian(subject-actor)
type PatientConsent struct {
	ID        string `gorm:"primary_key"`
	Actor     string `gorm:"not null"`
	Custodian string `gorm:"not null"`
	Records   []ConsentRecord
	Subject   string `gorm:"not null"`
}

// TableName returns the SQL table for this type
func (PatientConsent) TableName() string {
	return "patient_consent"
}

// BeforeDelete makes sure the ConsentRecords of a PatientConsent gets deleted too
func (pc *PatientConsent) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Delete(ConsentRecord{}, "patient_consent_id = ?", pc.ID).Error
}

// DataClasses combines all consent data classes from all records
func (pc PatientConsent) DataClasses() []DataClass {
	var dataClasses []DataClass
	for _, r := range pc.Records {
		dataClasses = append(dataClasses, r.DataClasses...)
	}
	return dataClasses
}

// ConsentRecord represents the individual records/attachments for a PatientConsent
// Changes to ConsentRecords are chained by PreviousHash pointing to Hash. All member of the chain can be found by the UUID
// The UUID remains internal
type ConsentRecord struct {
	ID               uint `gorm:"AUTO_INCREMENT"`
	PatientConsentID string
	ValidFrom        time.Time `gorm:"not null"`
	ValidTo          *time.Time
	Hash             string `gorm:"not null"`
	PreviousHash     *string
	Version          uint   `gorm:"DEFAULT:1"`
	UUID             string `gorm:"column:uuid;not null"`
	DataClasses      []DataClass
}

// TableName returns the SQL table for this type
func (ConsentRecord) TableName() string {
	return "consent_record"
}

// BeforeDelete makes sure the DataClasses of a ConsentRecords gets deleted too
func (cr *ConsentRecord) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Delete(DataClass{}, "consent_record_id = ?", cr.ID).Error
}

// DataClass defines struct for data_class table
type DataClass struct {
	ConsentRecordID uint
	Code            string `gorm:"not null"`
}

// TableName returns the SQL table for this type
func (DataClass) TableName() string {
	return "data_class"
}

func (pc *PatientConsent) String() string {
	return fmt.Sprintf("%s@%s for %s", pc.Subject, pc.Custodian, pc.Actor)
}

// SameTriple compares this PatientConsent with another one on just Actor, Custiodian and Subject
func (pc *PatientConsent) SameTriple(other *PatientConsent) bool {
	return pc.Subject == other.Subject && pc.Custodian == other.Custodian && pc.Actor == other.Actor
}

func (r *DataClass) String() string {
	return r.Code
}

// DataClassesFromStrings converts a slice of strings to a slice of Recources
func DataClassesFromStrings(list []string) []DataClass {
	a := make([]DataClass, len(list))
	for i, l := range list {
		a[i] = DataClass{Code: l}
	}
	return a
}
