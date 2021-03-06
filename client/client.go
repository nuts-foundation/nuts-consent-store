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

package client

import (
	"github.com/nuts-foundation/nuts-consent-store/api"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	core "github.com/nuts-foundation/nuts-go-core"
	"github.com/sirupsen/logrus"
	"time"
)

// NewConsentStoreClient creates a new Local- or RemoteClient for the nuts consent-store
func NewConsentStoreClient() pkg.ConsentStoreClient {
	consentStore := pkg.ConsentStoreInstance()

	if consentStore.Config.Mode == core.ServerEngineMode {
		if err := consentStore.Configure(); err != nil {
			logrus.Panic(err)
		}

		return consentStore
	} else {
		return api.HttpClient{
			ServerAddress: consentStore.Config.Address,
			Timeout:       time.Second,
			Logger: logrus.WithFields(logrus.Fields{
				"engine":    "consent-store",
				"component": "API-client",
			}),
		}
	}
}
