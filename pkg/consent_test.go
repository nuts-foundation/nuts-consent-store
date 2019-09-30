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
	"github.com/labstack/gommon/random"
	"testing"
	"time"
)

func TestConsentStoreInstance(t *testing.T) {
	t.Run("returns same instance every time", func(t *testing.T) {
		if ConsentStoreInstance() != ConsentStoreInstance() {
			t.Error("Expected instance to be the same")
		}
	})
}

func TestConsentStore_RecordConsent_AuthConsent(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	t.Run("Recorded consent can be authorized against", func(t *testing.T) {

		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   time.Now().Add(time.Hour * +12),
						Hash:      "234caef",
						Resources: []Resource{
							{
								ResourceType: "resource",
							},
						},
					},
				},
			},
		}

		err := client.RecordConsent(context.TODO(), rules)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		auth, err := client.ConsentAuth(context.TODO(), "custodian", "subject", "actor", "resource", nil)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if !auth {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Updating a existing consent", func(t *testing.T) {

		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor333",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   time.Now().Add(time.Hour * +12),
						Hash:      "234caefg",
						Resources: []Resource{
							{
								ResourceType: "resource",
							},
						},
					},
					{
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   time.Now().Add(time.Hour * +12),
						Hash:      "334caefg",
						Resources: []Resource{
							{
								ResourceType: "resource",
							},
						},
					},
				},
			},
		}
		err := client.RecordConsent(context.TODO(), rules)
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		// Update the validTo of the record.
		rules[0].Records[0].ValidTo = time.Now()
		rules[0].Records[0].Hash = "234caefh"

		err = client.RecordConsent(context.TODO(), rules)
		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		consent, err := client.QueryConsentForActor(context.TODO(), "actor333", "*")
		if len(consent) != 1 {
			t.Errorf("Expected 1 patientConsent, got [%d]", len(consent))
		}
		if len(consent[0].Records) != 2 {
			t.Errorf("Expected 2 records, got: [%d]", len(consent[0].Records))
		}
		if len(consent[0].Records[0].Resources) != 1 {
			t.Errorf("Expected 1 rule, got: [%d]", len(consent[0].Records[0].Resources))
		}
	})

	t.Run("Authorize non-existing consent returns false", func(t *testing.T) {

		auth, err := client.ConsentAuth(context.TODO(), "custodian", "subject", "actor2", "resource", nil)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if auth {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("Authorize against expired consent returns false", func(t *testing.T) {

		rules := []PatientConsent{
			{
				Actor:     "actor2",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						Hash:      random.String(8),
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   time.Now().Add(time.Hour * -12),
						Resources: []Resource{
							{
								ResourceType: "resource",
							},
						},
					},
				},
			},
		}

		err := client.RecordConsent(context.TODO(), rules)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		auth, err := client.ConsentAuth(context.TODO(), "custodian", "subject", "actor2", "resource", nil)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if auth {
			t.Errorf("Expected false, got true")
		}
	})

	t.Run("Authorize against consent at a different point in time returns false", func(t *testing.T) {

		rules := []PatientConsent{
			{
				Actor:     "actor3",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						Hash:      random.String(8),
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   time.Now().Add(time.Hour * 12),
						Resources: []Resource{
							{
								ResourceType: "resource",
							},
						},
					},
				},
			},
		}

		err := client.RecordConsent(context.TODO(), rules)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		cp := time.Now().Add(time.Hour * -36)
		auth, err := client.ConsentAuth(context.TODO(), "custodian", "subject", "actor3", "resource", &cp)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if auth {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestConsentStore_QueryConsentForActor(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	rules := []PatientConsent{
		{
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      random.String(8),
					Resources: []Resource{
						{
							ResourceType: "resource",
						},
					},
				},
			},
		},
		{
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject2",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      random.String(8),
					Resources: []Resource{
						{
							ResourceType: "resource2",
						},
					},
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

func TestConsentStore_QueryConsentForActorAndSubject(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	rules := []PatientConsent{
		{
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      "1",
					Resources: []Resource{
						{
							ResourceType: "resource",
						},
					},
				},
			},
		},
		{
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject2",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      "2",
					Resources: []Resource{
						{
							ResourceType: "resource2",
						},
					},
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

func TestConsentStore_Configure(t *testing.T) {
	t.Run("works OK for in memory db", func(t *testing.T) {
		client := ConsentStore{
			Config: ConsentStoreConfig{
				Connectionstring: ":memory:",
				Mode:             "server",
			},
		}

		if err := client.Configure(); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("gives error if DB can't be found", func(t *testing.T) {
		client := ConsentStore{
			Config: ConsentStoreConfig{
				Connectionstring: "file:test.db?mode=ro",
				Mode:             "server",
			},
		}

		err := client.Configure()

		if err == nil {
			t.Errorf("Expected error, got nothing")
			return
		}

		expected := "unable to open database file"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
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

func TestConsentStore_QueryConsent(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	rules := []PatientConsent{
		{
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      random.String(8),
					Resources: []Resource{
						{
							ResourceType: "resource",
						},
					},
				},
			},
		},
		{
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      random.String(8),
					Resources: []Resource{
						{
							ResourceType: "resource2",
						},
					},
				},
			},
		},
	}

	if err := client.RecordConsent(context.TODO(), rules); err != nil {
		t.Fatal(err)
	}

	t.Run("Recorded consent can be found by subject", func(t *testing.T) {
		subject := "subject"

		consent, err := client.QueryConsent(context.TODO(), nil, nil, &subject)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 2 {
			t.Errorf("Expected 2 results, got [%d]", len(consent))
		}
	})

	t.Run("Recorded consent can be found by custodian", func(t *testing.T) {
		custodian := "custodian2"

		consent, err := client.QueryConsent(context.TODO(), nil, &custodian, nil)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 1 {
			t.Errorf("Expected 1 results, got [%d]", len(consent))
		}

		if consent[0].Custodian != custodian {
			t.Errorf("Expected custodian to be [%s] got [%s]", custodian, consent[0].Custodian)
		}
	})
}

func TestConsentStore_DeleteConsentRecordByHash(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()
	hash := random.String(8)

	rules := []PatientConsent{
		{
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   time.Now().Add(time.Hour * 12),
					Hash:      hash,
					Resources: []Resource{
						{
							ResourceType: "resource",
						},
					},
				},
			},
		},
	}

	if err := client.RecordConsent(context.TODO(), rules); err != nil {
		t.Fatal(err)
	}

	t.Run("Not found returns false", func(t *testing.T) {
		val, _ := client.DeleteConsentRecordByHash(context.TODO(), "unknown")

		if val {
			t.Error("Expected record to not be deleted")
		}
	})

	t.Run("Record is deleted", func(t *testing.T) {
		val, _ := client.DeleteConsentRecordByHash(context.TODO(), hash)

		if !val {
			t.Error("Expected record to be deleted")
		}

		client.ConsentAuth(context.TODO(), "custodian", "subject", "actor", "resource", nil)
	})
}
