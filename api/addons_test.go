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
	"github.com/labstack/gommon/random"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"testing"
)

func TestFromSimplifiedConsentRule(t *testing.T) {
	t.Run("single consentRule converted", func(t *testing.T) {
		scs, _ := FromSimplifiedConsentRule([]pkg.PatientConsent{consentRule()})

		if len(scs) != 1 {
			t.Error("Expected rules to have 1 item")
			return
		}

		sc := scs[0]

		if sc.Subject != "subject" {
			t.Error("Expected Subject to equal [subject]")
		}

		if sc.Custodian != "custodian" {
			t.Error("Expected Custodian to equal [custodian]")
		}

		if sc.Actor != "actor" {
			t.Error("Expected Actor to equal actor")
		}

		if len(sc.Resources) != 1 {
			t.Error("Expected resources to have 1 item")
			return
		}

		if sc.Resources[0] != "resource" {
			t.Error("Expected Resource to equal [resource]")
		}
	})

	t.Run("multiple actors gives error", func(t *testing.T) {
		crs := []pkg.PatientConsent{consentRule(), consentRule()}
		crs[1].Actor = "actor2"

		_, err := FromSimplifiedConsentRule(crs)

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

func TestSimplifiedConsent_ToPatientConsent(t *testing.T) {
	hash := random.String(8)
	sc := SimplifiedConsent{
		Actor:     "actor",
		Custodian: "custodian",
		RecordHash: &hash,
		Resources: []string{"resource"},
		Subject:   "subject",
		ValidFrom: "2019-01-01",
		ValidTo:   "2020-01-01",
	}

	t.Run("correct transform", func(t *testing.T) {
		pc, _ := sc.ToPatientConsent()

		if pc.Subject != string(sc.Subject) {
			t.Error("Expected Subject to match")
			return
		}

		if pc.Custodian != string(sc.Custodian) {
			t.Error("Expected Custodian to match")
			return
		}

		if pc.Actor != string(sc.Actor) {
			t.Error("Expected Actor to match")
			return
		}

		if len(pc.Records) != 1 {
			t.Error("Expected 1 record in PatientConsent")
			return
		}

		if pc.Records[0].Hash != *sc.RecordHash {
			t.Error("Expected Hash to match")
			return
		}

		if pc.Resources()[0].ResourceType != sc.Resources[0] {
			t.Error("Expected Resources to match")
			return
		}

		if pc.Records[0].ValidFrom.Format("2006-01-02") != string(sc.ValidFrom) {
			t.Error("Expected ValidFrom to match")
			return
		}

		if pc.Records[0].ValidTo.Format("2006-01-02") != string(sc.ValidTo) {
			t.Error("Expected ValidTo to match")
			return
		}
	})

	t.Run("Incorrect validTo returns error", func(t *testing.T) {
		sc.ValidTo = "202-01-01"
		_, err := sc.ToPatientConsent()

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "parsing time \"202-01-01\" as \"2006-01-02\": cannot parse \"01-01\" as \"2006\""
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("Incorrect validFrom returns error", func(t *testing.T) {
		sc.ValidFrom = "202-01-01"
		_, err := sc.ToPatientConsent()

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "parsing time \"202-01-01\" as \"2006-01-02\": cannot parse \"01-01\" as \"2006\""
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func consentRule() pkg.PatientConsent {
	return pkg.PatientConsent{
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
