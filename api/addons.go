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
	"errors"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"time"
)

// ToPatientConsent converts the SimplifiedConsent object to an internal PatientConsent
func (sc CreateConsentRequest) ToPatientConsent() (pkg.PatientConsent, error) {
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
		Records: records,
	}, nil
}

// ToConsentRecord converts the API consent record object to the internal DB object
func (cr ConsentRecord) ToConsentRecord() (pkg.ConsentRecord, error) {
	var resources []pkg.Resource

	for _, a := range cr.Resources {
		resources = append(resources, pkg.Resource{ResourceType: a})
	}

	validFrom, err := time.Parse("2006-01-02", string(cr.ValidFrom))
	if err != nil {
		return pkg.ConsentRecord{}, err
	}
	validTo, err := time.Parse("2006-01-02", string(cr.ValidTo))
	if err != nil {
		return pkg.ConsentRecord{}, err
	}

	return pkg.ConsentRecord{
		ValidFrom: validFrom,
		ValidTo:   validTo,
		Hash:      *cr.RecordHash,
		Resources: resources,
	}, nil
}

// FromPatientConsent converts a slice of pkg.PatientConsent to a slice of SimplifiedConsent
// it cannot convert when multiple actors are involved
func FromPatientConsent(patientConsent []pkg.PatientConsent) ([]SimplifiedConsent, error) {
	var (
		firstActor string
		consent    []SimplifiedConsent
	)

	for _, c := range patientConsent {
		if firstActor == "" {
			firstActor = c.Actor
		} else {
			if firstActor != c.Actor {
				return nil, errors.New("Can not convert consent rules with multiple actors")
			}
		}
		for _, r := range c.Records {
			var resources []string
			for _, r2 := range r.Resources {
				resources = append(resources, r2.ResourceType)
			}
			consent = append(consent, SimplifiedConsent{
				Id:        r.Hash,
				Subject:   Identifier(c.Subject),
				Custodian: Identifier(c.Custodian),
				Actor:     Identifier(c.Actor),
				Resources: resources,
				ValidFrom: ValidFrom(r.ValidFrom.Format("2006-01-02")),
				ValidTo:   ValidTo(r.ValidTo.Format("2006-01-02")),
			})
		}
	}

	return consent, nil
}
