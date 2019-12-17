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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/labstack/gommon/random"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestConsentStoreInstance(t *testing.T) {
	t.Run("returns same instance every time", func(t *testing.T) {
		if ConsentStoreInstance() != ConsentStoreInstance() {
			t.Error("Expected instance to be the same")
		}
	})
}

func TestConsentStore_RecordConsent(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	t.Run("It sets the version number to 1", func(t *testing.T) {
		validTo := time.Now().Add(time.Hour * +12)

		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   &validTo,
						Hash:      "234caef",
						DataClasses: []DataClass{
							{
								Code: "resource",
							},
						},
						UUID: uuid.NewV4().String(),
					},
				},
			},
		}

		err := client.RecordConsent(context.TODO(), rules)

		if assert.NoError(t, err) {
			a := "actor"
			s := "subject"
			pcs, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

			if assert.NoError(t, err) {
				assert.Equal(t, uint(1), pcs[0].Records[0].Version)
			}
		}
	})
}

func TestConsentStore_FindConsentRecord(t *testing.T) {
	client := defaultConsentStore()
	pc := patientConsent()
	if err := client.RecordConsent(context.TODO(), pc); err != nil {
		t.Fatal(err)
	}
	defer client.Shutdown()

	t.Run("without latest flag", func(t *testing.T) {

		cr, err := client.FindConsentRecordByHash(context.TODO(), pc[0].Records[0].Hash, false)

		assert.NoError(t, err)
		assert.Equal(t, 1, int(cr.Version))
	})

	t.Run("with latest flag", func(t *testing.T) {

		cr, err := client.FindConsentRecordByHash(context.TODO(), pc[0].Records[0].Hash, true)

		assert.NoError(t, err)
		assert.Equal(t, 1, int(cr.Version))
	})

	t.Run("not found", func(t *testing.T) {

		_, err := client.FindConsentRecordByHash(context.TODO(), "unknown", false)

		if assert.Error(t, err) {
			assert.True(t, errors.Is(err, ErrorNotFound))
		}
	})

	t.Run("not found with latest", func(t *testing.T) {

		_, err := client.FindConsentRecordByHash(context.TODO(), "unknown", true)

		if assert.Error(t, err) {
			assert.True(t, errors.Is(err, ErrorNotFound))
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
						Hash:      "234caef",
						DataClasses: []DataClass{
							{
								Code: "resource",
							},
						},
						UUID: uuid.NewV4().String(),
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

	t.Run("Updating an existing consent", func(t *testing.T) {

		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor333",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						ValidFrom: time.Now().Add(time.Hour * -24),
						Hash:      "234caefg",
						DataClasses: []DataClass{
							{
								Code: "resource",
							},
						},
					},
				},
			},
		}
		err := client.RecordConsent(context.TODO(), rules)
		if err != nil {
			t.Fatalf("Expected no error, got [%v]", err)
		}

		// Update the validTo of the record.
		newValidTo := time.Now().Add(time.Hour)
		rules[0].Records[0].ValidTo = &newValidTo
		hcp := rules[0].Records[0].Hash
		rules[0].Records[0].PreviousHash = &hcp
		rules[0].Records[0].Hash = "234caefh_2"

		err = client.RecordConsent(context.TODO(), rules)
		if assert.NoError(t, err) {
			a := "actor333"
			consent, err := client.QueryConsent(context.TODO(), &a, nil, nil, nil)
			if assert.NoError(t, err) {

				assert.Len(t, consent, 1)
				assert.Len(t, consent[0].Records, 1)
				assert.Len(t, consent[0].Records[0].DataClasses, 1)
				assert.Equal(t, uint(2), consent[0].Records[0].Version)
			}
		}
	})

	t.Run("Updating consent record not latest in chain", func(t *testing.T) {
		r := random.String(8)
		validTo := time.Now().Add(time.Hour * +12)

		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor3333",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   &validTo,
						Hash:      r,
						DataClasses: []DataClass{
							{
								Code: "resource",
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
		validTo = time.Now()
		rules[0].Records[0].ValidTo = &validTo
		rules[0].Records[0].PreviousHash = &r
		rules[0].Records[0].Hash = fmt.Sprintf("%s_2", r)

		err = client.RecordConsent(context.TODO(), rules)
		if assert.NoError(t, err) {
			// update again
			err = client.RecordConsent(context.TODO(), rules)
			assert.Error(t, err) // unique constraint violation
		}
	})

	t.Run("Updating unknown consent", func(t *testing.T) {
		r := random.String(8)

		validTo := time.Now().Add(time.Hour * +12)
		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor123",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						ValidFrom:    time.Now().Add(time.Hour * -24),
						ValidTo:      &validTo,
						Hash:         r,
						PreviousHash: &r,
						DataClasses: []DataClass{
							{
								Code: "resource",
							},
						},
					},
				},
			},
		}
		err := client.RecordConsent(context.TODO(), rules)

		if assert.Error(t, err) {
			assert.True(t, errors.Is(err, ErrorNotFound))
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

		validTo := time.Now().Add(time.Hour * -12)
		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor2",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						Hash:      random.String(8),
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   &validTo,
						DataClasses: []DataClass{
							{
								Code: "resource",
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

		validTo := time.Now().Add(time.Hour * +12)
		rules := []PatientConsent{
			{
				ID:        random.String(8),
				Actor:     "actor3",
				Custodian: "custodian",
				Subject:   "subject",

				Records: []ConsentRecord{
					{
						Hash:      random.String(8),
						ValidFrom: time.Now().Add(time.Hour * -24),
						ValidTo:   &validTo,
						DataClasses: []DataClass{
							{
								Code: "resource",
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

	validTo := time.Now().Add(time.Hour * +12)

	rules := []PatientConsent{
		{
			ID:        random.String(8),
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      random.String(8),
					DataClasses: []DataClass{
						{
							Code: "resource",
						},
					},
				},
			},
		},
		{
			ID:        random.String(8),
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject2",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      random.String(8),
					DataClasses: []DataClass{
						{
							Code: "resource2",
						},
					},
				},
			},
		},
	}

	client.RecordConsent(context.TODO(), rules)

	t.Run("Recorded consent can be found", func(t *testing.T) {
		a := "actor"
		consent, err := client.QueryConsent(context.TODO(), &a, nil, nil, nil)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 1 {
			t.Errorf("Expected 1 result, got [%d]", len(consent))
		}
	})

	t.Run("Recorded consent is not found outside time frame", func(t *testing.T) {
		a := "actor"
		tt := time.Now().Add(time.Hour * 13)
		consent, err := client.QueryConsent(context.TODO(), &a, nil, nil, &tt)

		if assert.NoError(t, err) {
			assert.Len(t, consent, 0)
		}
	})

	t.Run("Non-recorded is not found", func(t *testing.T) {
		a := "actor3"
		consent, err := client.QueryConsent(context.TODO(), &a, nil, nil, nil)

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

	validTo := time.Now().Add(time.Hour * +12)

	rules := []PatientConsent{
		{
			ID:        "123",
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      "1",
					DataClasses: []DataClass{
						{
							Code: "resource",
						},
					},
				},
			},
		},
		{
			ID:        "223",
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject2",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      "2",
					DataClasses: []DataClass{
						{
							Code: "resource2",
						},
					},
				},
			},
		},
	}

	client.RecordConsent(context.TODO(), rules)

	t.Run("Recorded consent can be found", func(t *testing.T) {
		a := "actor"
		s := "subject"
		consent, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}

		if len(consent) != 1 {
			t.Errorf("Expected 1 result, got [%d]", len(consent))
		}
	})

	t.Run("Non-recorded is not found", func(t *testing.T) {
		a := "actor"
		s := "subject2"
		consent, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

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
				//Connectionstring: ":memory:",
				Connectionstring: "file",
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
			Mode:             "server",
		},
	}

	if err := client.Configure(); err != nil {
		panic(err)
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	return client
}

func TestConsentStore_QueryConsent(t *testing.T) {
	client := defaultConsentStore()
	defer client.Shutdown()

	validTo := time.Now().Add(time.Hour * +12)

	rules := []PatientConsent{
		{
			ID:        random.String(8),
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      random.String(8),
					UUID:      "1",
					DataClasses: []DataClass{
						{
							Code: "resource",
						},
					},
				},
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      random.String(8),
					UUID:      "2",
					DataClasses: []DataClass{
						{
							Code: "resource",
						},
					},
				},
			},
		},
		{
			ID:        random.String(8),
			Actor:     "actor2",
			Custodian: "custodian2",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      random.String(8),
					UUID:      "3",
					DataClasses: []DataClass{
						{
							Code: "resource2",
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

		consent, err := client.QueryConsent(context.TODO(), nil, nil, &subject, nil)

		if assert.NoError(t, err) {
			assert.Len(t, consent, 2)
			assert.Equal(t, 3, len(consent[0].Records)+len(consent[1].Records))
			assert.Len(t, consent[0].Records[0].DataClasses, 1)
		}

	})

	t.Run("Recorded consent can be found by custodian", func(t *testing.T) {
		custodian := "custodian2"

		consent, err := client.QueryConsent(context.TODO(), nil, &custodian, nil, nil)

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
	validTo := time.Now().Add(time.Hour * 12)

	rules := []PatientConsent{
		{
			ID:        random.String(8),
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",
			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      hash,
					DataClasses: []DataClass{
						{
							Code: "resource",
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

func patientConsent() []PatientConsent {
	validTo := time.Now().Add(time.Hour * 12)
	return []PatientConsent{
		{
			ID:        random.String(8),
			Actor:     "actor",
			Custodian: "custodian",
			Subject:   "subject",

			Records: []ConsentRecord{
				{
					ValidFrom: time.Now().Add(time.Hour * -24),
					ValidTo:   &validTo,
					Hash:      random.String(8),
					DataClasses: []DataClass{
						{
							Code: "resource",
						},
					},
					UUID: uuid.NewV4().String(),
				},
			},
		},
	}
}
