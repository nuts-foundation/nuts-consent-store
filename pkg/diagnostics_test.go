/*
 * Nuts consent store
 * Copyright (C) 2020. Nuts community
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

	"github.com/stretchr/testify/assert"
)

func TestEventOctopus_Diagnostics(t *testing.T) {
	client := defaultConsentStore()
	client.Configure()

	t.Run("Diagnostics returns 1 report", func(t *testing.T) {
		results := client.Diagnostics()

		assert.Len(t, results, 1)
	})

	t.Run("Diagnostics returns DB info", func(t *testing.T) {
		found := false
		results := client.Diagnostics()
		for _, r := range results {
			if r.Name() == "DB" {
				found = true
				assert.Equal(t, "connection string: :memory:, ping: true", r.String())
			}
		}

		assert.True(t, found)
	})

	client.Shutdown()

	t.Run("Diagnostics returns DB info when down", func(t *testing.T) {
		found := false
		results := client.Diagnostics()
		for _, r := range results {
			if r.Name() == "DB" {
				found = true
				assert.Equal(t, "connection string: :memory:, ping: false, error: sql: database is closed", r.String())
			}
		}

		assert.True(t, found)
	})
}
