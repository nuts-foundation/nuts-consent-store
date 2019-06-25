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
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"reflect"
	"sync"
	"testing"
)

func TestNewConsentStoreClient(t *testing.T) {
	t.Run("returns ConsentStore by default", func(t *testing.T) {
		i := pkg.ConsentStoreInstance()
		i.Config.Mode = "server"
		i.ConfigOnce = sync.Once{}
		cc := NewConsentStoreClient()

		if reflect.TypeOf(cc).String() != "*pkg.ConsentStore" {
			t.Errorf("Expected Client to be of type *consent.ConsentStore, got %s", reflect.TypeOf(cc))
		}
	})

	t.Run("invalid configuration panics", func(t *testing.T) {
		i := pkg.ConsentStoreInstance()
		i.Config.Connectionstring = "file:data.db?mode=readonly"
		i.Config.Mode = "server"
		i.ConfigOnce = sync.Once{}

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic")
			}
			i.ConfigOnce = sync.Once{}
		}()

		NewConsentStoreClient()
	})

	t.Run("returns APIClient in client mode", func(t *testing.T) {
		i := pkg.ConsentStoreInstance()
		i.Config.Mode = "client"
		defer func() {
			i.ConfigOnce = sync.Once{}
		}()
		cc := NewConsentStoreClient()

		expected := "api.HttpClient"
		if reflect.TypeOf(cc).String() != expected {
			t.Errorf("Expected Client to be of type %s, got %s", expected, reflect.TypeOf(cc))
		}
	})
}
