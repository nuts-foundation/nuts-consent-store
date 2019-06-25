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
)

// ToConsentRule converts the SimplifiedConsent object to an internal ConsentRule
func (sc SimplifiedConsent) ToConsentRule() []pkg.ConsentRule {

	var rules = make([]pkg.ConsentRule, len(sc.Actors))
	var resources = make([]pkg.Resource, len(sc.Resources))

	for _, a := range sc.Resources {
		resources = append(resources, pkg.Resource{ResourceType: a})
	}

	for _, a := range sc.Actors {

		rules = append(rules, pkg.ConsentRule{
			Subject:   string(sc.Subject),
			Custodian: string(sc.Custodian),
			Actor:     string(a),
			Resources: resources,
		})
	}

	return rules
}

// ToConsentRule converts the ConsentCheckRequest object to an internal ConsentRule
func (sc ConsentCheckRequest) ToConsentRule() pkg.ConsentRule {

	return pkg.ConsentRule{
		Subject:   string(sc.Subject),
		Custodian: string(sc.Custodian),
		Actor:     string(sc.Actor),
		Resources: []pkg.Resource{{ResourceType: sc.ResourceType}},
	}
}

// FromSimplifiedConsentRule converts a slice of pkg.ConsentRule to a slice of SimplifiedConsent
// it cannot convert when numtiple actors are involved
func FromSimplifiedConsentRule(rules []pkg.ConsentRule) ([]SimplifiedConsent, error) {
	var (
		firstActor string
		consent    []SimplifiedConsent
	)

	for _, r := range rules {
		if firstActor == "" {
			firstActor = r.Actor
		} else {
			if firstActor != r.Actor {
				return nil, errors.New("Can not convert consent rules with multiple actors")
			}
		}
		var resources []string
		for _, r2 := range r.Resources {
			resources = append(resources, r2.ResourceType)
		}

		consent = append(consent, SimplifiedConsent{
			Subject:   Identifier(r.Subject),
			Custodian: Identifier(r.Custodian),
			Actors:    []Identifier{Identifier(r.Actor)},
			Resources: resources,
		})
	}

	return consent, nil
}
