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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/gommon/random"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/nuts-foundation/nuts-go-core/mock"
)

func TestDefaultConsentStore_CheckConsent(t *testing.T) {
	client := defaultConsentStore()
	client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{consentRuleForQuery()})
	defer client.Cs.Shutdown()

	t.Run("API call returns 200 for no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		ccr := consentCheckRequest()
		ccr.Subject = "subject2"
		json, _ := json.Marshal(ccr)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		authValue := "false"
		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentCheckResponse{
			ConsentGiven: &authValue,
		})

		err := client.CheckConsent(echo)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})

	t.Run("API call returns 200 for auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		json, _ := json.Marshal(consentCheckRequest())
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		authValue := "true"
		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentCheckResponse{
			ConsentGiven: &authValue,
		})

		err := client.CheckConsent(echo)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})

	t.Run("API call returns 200 for auth with ValidAt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		ccr := consentCheckRequest()
		validAt := time.Now().Format("2006-01-02")
		ccr.ValidAt = &validAt
		json, _ := json.Marshal(ccr)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		authValue := "true"
		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentCheckResponse{
			ConsentGiven: &authValue,
		})

		err := client.CheckConsent(echo)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})

	t.Run("Missing body gives 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request)

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=missing body in request"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("Invalid ValidAt gives 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		ccr := consentCheckRequest()
		validAt := "202-01-01"
		ccr.ValidAt = &validAt
		json, _ := json.Marshal(ccr)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=invalid value for validAt: 202-01-01"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("Reading error gives 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{
			Body: errorCloser{},
		}

		echo.EXPECT().Request().Return(request)

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=error reading request: error"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("Invalid body gives 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request)

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=missing body in request"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing actor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentCheckRequest()
		consent.Actor = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing actor in checkRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})

	t.Run("API call returns 400 for missing subject", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentCheckRequest()
		consent.Subject = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing subject in checkRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})

	t.Run("API call returns 400 for missing custodian", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentCheckRequest()
		consent.Custodian = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=missing custodian in checkRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})

	t.Run("API call returns 400 for missing resource", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentCheckRequest()
		consent.ResourceType = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CheckConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing resourceType in checkRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})
}

func TestDefaultConsentStore_CreateConsent(t *testing.T) {
	client := defaultConsentStore()
	defer client.Cs.Shutdown()

	t.Run("API call returns 201 Created", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		json, _ := json.Marshal(testConsent())
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().NoContent(http.StatusCreated)

		if err := client.CreateConsent(echo); err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}
	})

	t.Run("Missing body gives 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request)

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=missing body in request"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing actor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Actor = Identifier("")

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing actor in createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error %s, got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Id = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing ID in createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error %s, got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing proofHash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.RecordHash = nil

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing recordHash in createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error %s, got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing subject", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Subject = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing subject in createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error %s, got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing custodian", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Custodian = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing custodian in createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error %s, got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing resources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Resources = []string{}

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing resources in createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})
}

func TestDefaultConsentStore_QueryConsent(t *testing.T) {
	client := defaultConsentStore()
	if err := client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{consentRuleForQuery()}); err != nil {
		t.Fatal(err)
	}
	defer client.Cs.Shutdown()

	t.Run("API call returns 200 for empty query", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		ccr := consentQuery()
		actor := Identifier("actor2")
		ccr.Actor = &actor
		json, _ := json.Marshal(ccr)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentQueryResponse{
			TotalResults: 0,
			Page:         PageDefinition{},
		})

		err := client.QueryConsent(echo)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})

	t.Run("API call returns 200 for results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		query := consentQuery()
		query.Query = "%"
		json, _ := json.Marshal(query)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentQueryResponse{
			TotalResults: 1,
			Results: []SimplifiedConsent{
				{
					Subject:   Identifier("urn:subject"),
					Custodian: Identifier("custodian"),
					Actor:     Identifier("actor"),
					Resources: []string{
						"resource",
					},
					ValidFrom: ValidFrom(time.Now().Add(time.Hour * -24).Format("2006-01-02")),
					ValidTo:   ValidTo(time.Now().Add(time.Hour * 24).Format("2006-01-02")),
				},
			},
			Page: PageDefinition{},
		})

		err := client.QueryConsent(echo)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})

	t.Run("API call returns 200 with results for subject search", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		query := consentQuery()
		query.Query = "urn:subject"
		json, _ := json.Marshal(query)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentQueryResponse{
			TotalResults: 1,
			Results: []SimplifiedConsent{
				{
					Subject:   Identifier("urn:subject"),
					Custodian: Identifier("custodian"),
					Actor:     Identifier("actor"),
					Resources: []string{
						"resource",
					},
					ValidFrom: ValidFrom(time.Now().Add(time.Hour * -24).Format("2006-01-02")),
					ValidTo:   ValidTo(time.Now().Add(time.Hour * 24).Format("2006-01-02")),
				},
			},
			Page: PageDefinition{},
		})

		err := client.QueryConsent(echo)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})

	t.Run("Missing body gives 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request)

		err := client.QueryConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
			return
		}

		expected := "code=400, message=missing body in request"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing actor and custodian", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentQuery()
		consent.Actor = nil
		consent.Custodian = nil

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.QueryConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing actor or custodian in queryRequest"
		if !strings.HasPrefix(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})

	t.Run("API call returns 400 for missing query", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentQuery()
		consent.Query = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.QueryConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing query in queryRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})
}

func TestDefaultConsentStore_DeleteConsent(t *testing.T) {
	client := defaultConsentStore()
	crq := consentRuleForQuery()
	client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{crq})
	defer client.Cs.Shutdown()

	t.Run("missing proofHash returns 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.DeleteConsent(echo, "")

		if err == nil {
			t.Error("Expected error, got nothing", err)
		}

		expected := "code=400, message=missing proofHash"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("Unknown proofHash returns 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.DeleteConsent(echo, "a")

		if err == nil {
			t.Error("Expected error, got nothing", err)
		}

		expected := "code=404, message=no ConsentRecord found for given hash"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("Correct delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().NoContent(202)

		err := client.DeleteConsent(echo, crq.Records[0].Hash)

		if err != nil {
			t.Errorf("Expected no error, got [%v]", err)
		}
	})
}

func testConsent() SimplifiedConsent {
	hash := random.String(8)
	return SimplifiedConsent{
		Id:         random.String(8),
		Actor:      Identifier("actor"),
		Custodian:  Identifier("custodian"),
		Subject:    Identifier("urn:subject"),
		RecordHash: &hash,
		Resources:  []string{"resource"},
		ValidFrom:  ValidFrom("2019-01-01"),
		ValidTo:    ValidTo("2030-01-01"),
	}
}

func consentCheckRequest() ConsentCheckRequest {
	return ConsentCheckRequest{
		Subject:      Identifier("urn:subject"),
		Custodian:    Identifier("custodian"),
		Actor:        Identifier("actor"),
		ResourceType: "resource",
	}
}

func consentQuery() QueryConsentJSONRequestBody {
	actor := Identifier("actor")
	return QueryConsentJSONRequestBody{
		Actor: &actor,
		Query: "subject",
	}
}

func defaultConsentStore() ApiWrapper {
	client := pkg.ConsentStore{
		Config: pkg.ConsentStoreConfig{
			Connectionstring: ":memory:",
			Mode:             "server",
		},
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	client.RunMigrations(client.Db.DB())

	return ApiWrapper{Cs: &client}
}

func consentRuleForQuery() pkg.PatientConsent {
	return pkg.PatientConsent{
		ID:        random.String(8),
		Subject:   "urn:subject",
		Custodian: "custodian",
		Actor:     "actor",
		Records: []pkg.ConsentRecord{
			{
				ValidFrom: time.Now().Add(time.Hour * -24),
				ValidTo:   time.Now().Add(time.Hour * 24),
				Resources: []pkg.Resource{
					{ResourceType: "resource"},
				},
				Hash: random.String(8),
			},
		},
	}
}

type errorCloser struct{}

func (errorCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func (errorCloser) Close() error {
	return errors.New("error")
}
