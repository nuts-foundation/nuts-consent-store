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
	"github.com/jinzhu/gorm"
	"time"
)

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

// Resources combines all resources from all records
func (pc PatientConsent) Resources() []Resource {
	var resources []Resource
	for _, r := range pc.Records {
		resources = append(resources, r.Resources...)
	}
	return resources
}

// ConsentRecord represents the individual records/attachments for a PatientConsent
type ConsentRecord struct {
	ID               uint `gorm:"AUTO_INCREMENT"`
	PatientConsentID string
	ValidFrom        time.Time `gorm:"not null"`
	ValidTo          time.Time `gorm:"not null"`
	Hash             string    `gorm:"not null"`
	Version 		 uint	   `gorm:"DEFAULT:1"`
	UUID 			 string    `gorm:"column:uuid"`
	Resources        []Resource
}

// TableName returns the SQL table for this type
func (ConsentRecord) TableName() string {
	return "consent_record"
}

// BeforeDelete makes sure the Resources of a ConsentRecords gets deleted too
func (cr *ConsentRecord) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Delete(Resource{}, "consent_record_id = ?", cr.ID).Error
}

// Resource defines struct for resource table
type Resource struct {
	ConsentRecordID uint
	ResourceType    string `gorm:"not null"`
}

// TableName returns the SQL table for this type
func (Resource) TableName() string {
	return "resource"
}

func (pc *PatientConsent) String() string {
	return fmt.Sprintf("%s@%s for %s", pc.Subject, pc.Custodian, pc.Actor)
}

// SameTriple compares this PatientConsent with another one on just Actor, Custiodian and Subject
func (pc *PatientConsent) SameTriple(other *PatientConsent) bool {
	return pc.Subject == other.Subject && pc.Custodian == other.Custodian && pc.Actor == other.Actor
}

func (r *Resource) String() string {
	return r.ResourceType
}

// ResourcesFromStrings converts a slice of strings to a slice of Recources
func ResourcesFromStrings(list []string) []Resource {
	a := make([]Resource, len(list))
	for i, l := range list {
		a[i] = Resource{ResourceType: l}
	}
	return a
}
