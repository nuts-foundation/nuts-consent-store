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
)

// PatientConsent defines struct for patient_consent table.
// ID refers to the HMAC id for a custodian(subject-actor)
type PatientConsent struct {
	ID         string `gorm:"primary_key"`
	Actor      string `gorm:"not null"`
	Custodian  string `gorm:"not null"`
	Records    []ConsentRecord
	Subject    string `gorm:"not null"`
}

func (PatientConsent) TableName() string {
	return "patient_consent"
}

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
	Hash       		 string    `gorm:"not null"`
	Resources        []Resource
}

func (ConsentRecord) TableName() string {
	return "consent_record"
}

// Resource defines struct for resource table
type Resource struct {
	ConsentRecordID uint
	ResourceType    string `gorm:"not null"`
}

func (Resource) TableName() string {
	return "resource"
}

func (se *PatientConsent) String() string {
	return fmt.Sprintf("%s@%s for %s", se.Subject, se.Custodian, se.Actor)
}

// SameTriple compares this PatientConsent with another one on just Actor, Custiodian and Subject
func (se *PatientConsent) SameTriple(other *PatientConsent) bool {
	return se.Subject == other.Subject && se.Custodian == other.Custodian && se.Actor == other.Actor
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
