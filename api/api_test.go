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
	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/nuts-foundation/nuts-go/mock"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDefaultConsentStore_CheckConsent(t *testing.T) {
	client := defaultConsentStore()
	client.Cs.RecordConsent(context.Background(), []pkg.ConsentRule{consentRuleForQuery()})
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
		if err.Error() != expected {
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
		if err.Error() != expected {
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
		if err.Error() != expected {
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
		if err.Error() != expected {
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
		if err.Error() != expected {
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
		if err.Error() != expected {
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
		if err.Error() != expected {
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

		client.CreateConsent(echo)
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
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing actor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Actors = []Identifier{}

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		if err.Error() != "code=400, message=missing actors in createRequest" {
			t.Errorf("Expected error code=400, message=missing actors in createRequest, got: [%s]", err.Error())
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

		if err.Error() != "code=400, message=missing subject in createRequest" {
			t.Errorf("Expected error code=400, message=missing subject in createRequest, got: [%s]", err.Error())
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

		if err.Error() != "code=400, message=missing custodian in createRequest" {
			t.Errorf("Expected error code=400, message=missing custodian in createRequest, got: [%s]", err.Error())
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
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})
}

func TestDefaultConsentStore_QueryConsent(t *testing.T) {
	client := defaultConsentStore()
	client.Cs.RecordConsent(context.Background(), []pkg.ConsentRule{consentRuleForQuery()})
	defer client.Cs.Shutdown()

	t.Run("API call returns 200 for empty query", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		ccr := consentQuery()
		ccr.Actor = "actor2"
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

		json, _ := json.Marshal(consentQuery())
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
					Actors: []Identifier{
						"actor",
					},
					Resources: []string{
						"resource",
					},
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
					Actors: []Identifier{
						"actor",
					},
					Resources: []string{
						"resource",
					},
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
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got: [%s]", expected, err.Error())
		}
	})

	t.Run("API call returns 400 for missing actor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := consentQuery()
		consent.Actor = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.QueryConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing actor in queryRequest"
		if err.Error() != expected {
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
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})
}

func testConsent() SimplifiedConsent {
	return SimplifiedConsent{
		Actors: []Identifier{
			Identifier("actor"),
		},
		Custodian: Identifier("custodian"),
		Subject:   Identifier("urn:subject"),
		Resources: []string{"resource"},
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

func consentQuery() ConsentQueryRequest {
	return ConsentQueryRequest{
		Actor: Identifier("actor"),
		Query: "subject",
	}
}

func defaultConsentStore() ApiWrapper {
	client := pkg.ConsentStore{
		Config: pkg.ConsentStoreConfig{
			Connectionstring: ":memory:",
		},
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	client.RunMigrations(client.Db.DB())

	return ApiWrapper{Cs: &client}
}

func consentRuleForQuery() pkg.ConsentRule {
	return pkg.ConsentRule{
		Subject:   "urn:subject",
		Custodian: "custodian",
		Actor:     "actor",
		Resources: []pkg.Resource{
			{ResourceType: "resource"},
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
