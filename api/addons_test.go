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
	"testing"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/stretchr/testify/assert"
)

func TestFromSimplifiedConsentRule(t *testing.T) {
	t.Run("single patientConsent converted", func(t *testing.T) {
		scs, _ := FromPatientConsent([]pkg.PatientConsent{patientConsent()})

		if len(scs) != 1 {
			t.Error("Expected rules to have 1 item")
			return
		}

		sc := scs[0]

		assert.Equal(t, "patientConsentId", sc.Id)
		assert.Equal(t, Identifier("subject"), sc.Subject)
		assert.Equal(t, Identifier("custodian"), sc.Custodian)
		assert.Equal(t, Identifier("actor"), sc.Actor)
		assert.Equal(t, "", *sc.RecordHash)
		assert.Len(t, sc.Resources, 1)
		assert.Equal(t, "resource", sc.Resources[0])
	})

	t.Run("multiple actors gives error", func(t *testing.T) {
		crs := []pkg.PatientConsent{patientConsent(), patientConsent()}
		crs[1].Actor = "actor2"

		_, err := FromPatientConsent(crs)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "Can not convert consent rules with multiple actors"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func TestCreateConsentRequest_ToPatientConsent(t *testing.T) {
	version := 1
	sc := CreateConsentRequest{
		Actor:     "actor",
		Custodian: "custodian",
		Subject:   "subject",
		Records: []ConsentRecord{
			{
				RecordHash: random.String(8),
				Resources:  []string{"resource"},
				ValidFrom:  "2019-01-01T12:00:00+01:00",
				ValidTo:    "2020-01-01T12:00:00+01:00",
				Version:    &version,
			},
		},
	}

	t.Run("correct transform", func(t *testing.T) {
		pc, err := sc.ToPatientConsent()
		if err != nil {
			t.Fatal(err.Error())
		}

		assert.Equal(t, string(sc.Subject), pc.Subject)
		assert.Equal(t, string(sc.Custodian), pc.Custodian)
		assert.Equal(t, string(sc.Actor), pc.Actor)
		assert.Len(t, pc.Records, 1)
		assert.Equal(t, sc.Records[0].RecordHash, pc.Records[0].Hash)
		assert.Equal(t, sc.Records[0].Resources[0], pc.Resources()[0].ResourceType)
		assert.Equal(t, string(sc.Records[0].ValidFrom), pc.Records[0].ValidFrom.Format(pkg.Iso6801DateTime))
		assert.Equal(t, string(sc.Records[0].ValidTo), pc.Records[0].ValidTo.Format(pkg.Iso6801DateTime))
	})

	t.Run("Incorrect validTo returns error", func(t *testing.T) {
		sc.Records[0].ValidTo = "202-01-01"
		_, err := sc.ToPatientConsent()

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "parsing time \"202-01-01\" as \"2006-01-02T15:04:05-07:00\": cannot parse \"01-01\" as \"2006\""
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("Incorrect validFrom returns error", func(t *testing.T) {
		sc.Records[0].ValidFrom = "202-01-01"
		_, err := sc.ToPatientConsent()

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "parsing time \"202-01-01\" as \"2006-01-02T15:04:05-07:00\": cannot parse \"01-01\" as \"2006\""
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func TestFromConsentRecord(t *testing.T) {
	t.Run("correct transform", func(t *testing.T) {
		pc := FromConsentRecord(consentRecord())

		assert.Equal(t, "PreviousHash", *pc.PreviousRecordHash)
		assert.Equal(t, "Hash", pc.RecordHash)
		assert.Equal(t, 2, *pc.Version)
		assert.Equal(t, ValidTo("2001-09-12T12:00:00+02:00"), pc.ValidTo)
		assert.Equal(t, ValidFrom("2001-09-11T12:00:00+02:00"), pc.ValidFrom)
	})
}

func patientConsent() pkg.PatientConsent {
	return pkg.PatientConsent{
		ID:        "patientConsentId",
		Subject:   "subject",
		Custodian: "custodian",
		Actor:     "actor",
		Records: []pkg.ConsentRecord{
			{
				Resources: []pkg.Resource{
					{ResourceType: "resource"},
				},
			},
		},
	}
}

func consentRecord() pkg.ConsentRecord {
	t1, _ := time.Parse(pkg.Iso6801DateTime, "2001-09-11T12:00:00+02:00")
	t2, _ := time.Parse(pkg.Iso6801DateTime, "2001-09-12T12:00:00+02:00")

	prevH := "PreviousHash"

	return pkg.ConsentRecord{
		ID:               1,
		PatientConsentID: "PatientConsentID",
		ValidFrom:        t1,
		ValidTo:          t2,
		Hash:             "Hash",
		PreviousHash:     &prevH,
		Version:          2,
		UUID:             "UUID",
		Resources: []pkg.Resource{
			{ResourceType: "resource"},
		},
	}
}
