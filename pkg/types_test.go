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
	"testing"
)

func TestConsentRule_String(t *testing.T) {
	t.Run("outputs string representation of model", func(t *testing.T) {
		model := ConsentRule{
			Actor: "actor",
			Custodian: "custodian",
			Subject: "subject",
			Resources: ResourcesFromStrings([]string{"resources"}),
		}

		out := model.String()

		if out != "subject@custodian for actor: resources" {
			t.Errorf("Expected subject@custodian for actor: resources, Got [%s]", out)
		}
	})
}

func TestResource_String(t *testing.T) {
	t.Run("outputs string representation of model", func(t *testing.T) {
		model := Resource{
			ResourceType: "resource",
		}

		out := model.String()

		if out != "resource" {
			t.Errorf("Expected resources, Got [%s]", out)
		}
	})
}

func TestConsentRule_SameTriple(t *testing.T) {
	t.Run("returns true for same Actor, Custodian and Subject", func(t *testing.T) {
		if !testConsent().SameTriple(testConsent()) {
			t.Errorf("Expected structs to be the same")
		}
	})

	t.Run("ignores resources", func(t *testing.T) {
		other := testConsent()
		other.Resources = nil

		if !testConsent().SameTriple(other) {
			t.Errorf("Expected structs to be the same")
		}
	})

	t.Run("returns false for different actor", func(t *testing.T) {
		other := testConsent()
		other.Actor = ""

		if testConsent().SameTriple(other) {
			t.Errorf("Expected structs to be different")
		}
	})

	t.Run("returns false for different Custodian", func(t *testing.T) {
		other := testConsent()
		other.Custodian = ""

		if testConsent().SameTriple(other) {
			t.Errorf("Expected structs to be different")
		}
	})

	t.Run("returns false for different Subject", func(t *testing.T) {
		other := testConsent()
		other.Subject = ""

		if testConsent().SameTriple(other) {
			t.Errorf("Expected structs to be different")
		}
	})
}

func testConsent() *ConsentRule {
	return &ConsentRule{
		Actor: "actor",
		Custodian: "custodian",
		Subject: "subject",
		Resources: []Resource{{ResourceType:"resource"}},
	}
}
