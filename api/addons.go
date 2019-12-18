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

package api

import (
	"time"

	"github.com/nuts-foundation/nuts-consent-store/pkg"
)

// ToPatientConsent converts the api PatientConsent struct to an internal PatientConsent
func (sc PatientConsent) ToPatientConsent() (pkg.PatientConsent, error) {
	var records []pkg.ConsentRecord

	for _, r := range sc.Records {
		cr, err := r.ToConsentRecord()
		if err != nil {
			return pkg.PatientConsent{}, err
		}
		records = append(records, cr)
	}

	return pkg.PatientConsent{
		ID:        sc.Id,
		Subject:   string(sc.Subject),
		Custodian: string(sc.Custodian),
		Actor:     string(sc.Actor),
		Records:   records,
	}, nil
}

// ToConsentRecord converts the API consent record object to the internal DB object
func (cr ConsentRecord) ToConsentRecord() (pkg.ConsentRecord, error) {
	var resources []pkg.DataClass

	for _, a := range cr.DataClasses {
		resources = append(resources, pkg.DataClass{Code: a})
	}

	validFrom, err := time.Parse(pkg.Iso8601DateTime, string(cr.ValidFrom))
	if err != nil {
		return pkg.ConsentRecord{}, err
	}
	var validTo time.Time
	if cr.ValidTo != nil {
		validTo, err = time.Parse(pkg.Iso8601DateTime, string(*cr.ValidTo))
		if err != nil {
			return pkg.ConsentRecord{}, err
		}
	}

	return pkg.ConsentRecord{
		ValidFrom:    validFrom,
		ValidTo:      &validTo,
		Hash:         cr.RecordHash,
		PreviousHash: cr.PreviousRecordHash,
		DataClasses:  resources,
	}, nil
}

// FromPatientConsent converts a slice of pkg.PatientConsent to a slice of api.PatientConsent
func FromPatientConsents(pc []pkg.PatientConsent) []PatientConsent {
	var consents []PatientConsent

	for _, c := range pc {
		consents = append(consents, FromPatientConsent(c))
	}

	return consents
}

// FromPatientConsent converts a pkg.PatientConsent to a PatientConsent
func FromPatientConsent(pc pkg.PatientConsent) PatientConsent {
	var records []ConsentRecord

	for _, r := range pc.Records {
		cr := FromConsentRecord(r)
		records = append(records, cr)
	}

	return PatientConsent{
		Id:        pc.ID,
		Subject:   Identifier(pc.Subject),
		Custodian: Identifier(pc.Custodian),
		Actor:     Identifier(pc.Actor),
		Records:   records,
	}
}

// FromConsentRecord converts the DB type to api type
func FromConsentRecord(consentRecord pkg.ConsentRecord) ConsentRecord {
	var resources []string
	for _, r2 := range consentRecord.DataClasses {
		resources = append(resources, r2.Code)
	}

	version := int(consentRecord.Version)

	cr := ConsentRecord{
		PreviousRecordHash: consentRecord.PreviousHash,
		RecordHash:         consentRecord.Hash,
		DataClasses:        resources,
		ValidFrom:          ValidFrom(consentRecord.ValidFrom.Format(pkg.Iso8601DateTime)),
		Version:            &version,
	}

	if consentRecord.ValidTo != nil {
		validTo := ValidTo(consentRecord.ValidTo.Format(pkg.Iso8601DateTime))
		cr.ValidTo = &validTo
	}

	return cr
}
