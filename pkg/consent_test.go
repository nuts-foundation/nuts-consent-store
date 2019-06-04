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
	"context"
	"testing"
)

func TestDefaultConsentStore_RecordConsent_AuthConsent(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	t.Run("Recorded consent can be authorized against", func(t *testing.T) {

		rules := []ConsentRule{
			{
				Actor:     "actor",
				Custodian: "custodian",
				Subject:   "subject",
				Resources: []Resource{
					{
						ResourceType: "resource",
					},
				},
			},
		}

		err := client.RecordConsent(context.TODO(), rules)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		auth, err := client.ConsentAuth(context.TODO(), rules[0], "resource")

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if !auth {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Authorize non-existing consent returns false", func(t *testing.T) {

		rule := ConsentRule{
			Actor:     "actor2",
			Custodian: "custodian",
			Subject:   "subject",
			Resources: []Resource{
				{
					ResourceType: "resource",
				},
			},
		}

		auth, err := client.ConsentAuth(context.TODO(), rule, "resource")

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if auth {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestDefaultConsentStore_QueryConsentForActor(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	rules := []ConsentRule{
		{
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Resources: []Resource{
				{
					ResourceType: "resource",
				},
			},
		},
		{
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject2",
			Resources: []Resource{
				{
					ResourceType: "resource2",
				},
			},
		},
	}

	client.RecordConsent(context.TODO(), rules)

	t.Run("Recorded consent can be found", func(t *testing.T) {
		consent, err := client.QueryConsentForActor(context.TODO(), "actor", "*")

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 1 {
			t.Errorf("Expected 1 result, got [%d]", len(consent))
		}
	})

	t.Run("Non-recorded is not found", func(t *testing.T) {
		consent, err := client.QueryConsentForActor(context.TODO(), "actor3", "*")

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 0 {
			t.Errorf("Expected 0 results, got [%d]", len(consent))
		}
	})
}

func TestDefaultConsentStore_QueryConsentForActorAndSubject(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	rules := []ConsentRule{
		{
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Resources: []Resource{
				{
					ResourceType: "resource",
				},
			},
		},
		{
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject2",
			Resources: []Resource{
				{
					ResourceType: "resource2",
				},
			},
		},
	}

	client.RecordConsent(context.TODO(), rules)

	t.Run("Recorded consent can be found", func(t *testing.T) {
		consent, err := client.QueryConsentForActorAndSubject(context.TODO(), "actor", "subject")

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 1 {
			t.Errorf("Expected 1 result, got [%d]", len(consent))
		}
	})

	t.Run("Non-recorded is not found", func(t *testing.T) {
		consent, err := client.QueryConsentForActorAndSubject(context.TODO(), "actor", "subject2")

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 0 {
			t.Errorf("Expected 0 results, got [%d]", len(consent))
		}
	})
}

func defaultConsentStore() ConsentStore {
	client := ConsentStore{
		Config: ConsentStoreConfig{
			Connectionstring: ":memory:",
		},
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	client.RunMigrations(client.Db.DB())

	return client
}
