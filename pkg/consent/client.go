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

package consent

import (
	"context"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
)

// ConsentStoreClient defines all actions possible through a direct connection, command-line and REST api
type ConsentStoreClient interface {
	// ConsentAuth checks if a record exists in the db for the given combination and returns a bool.
	ConsentAuth(context context.Context, consentRule pkg.ConsentRule, resourceType string) (bool, error)
	// RecordConsent records a record in the db, this is not to be used to create a new distributed consent record. It's only valid for the local node.
	// It should only be called by the consent logic component (or for development purposes)
	RecordConsent(context context.Context, consent []pkg.ConsentRule) error
	// QueryConsentForActor can be used to perform full text searches on the backend. Scoped on actor only.
	QueryConsentForActor(context context.Context, actor string, query string) ([]pkg.ConsentRule, error)
	// QueryConsentForActorAndSubject can be used to list the custodians and resources for a given Actor and Subject.
	QueryConsentForActorAndSubject(context context.Context, actor string, subject string) ([]pkg.ConsentRule, error)
}

func NewConsentStoreClient() ConsentStoreClient {
	return ConsentStore()
}