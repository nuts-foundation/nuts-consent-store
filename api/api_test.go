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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"

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
		validAt := time.Now().Format(pkg.Iso8601DateTime)
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

	t.Run("API call returns 400 for missing consentRecordHash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Records[0].RecordHash = ""

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing recordHash in one or more records within createRequest"
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
		consent.Records[0].Resources = []string{}

		json, _ := json.Marshal(consent)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.CreateConsent(echo)

		if err == nil {
			t.Error("Expected error got nothing")
		}

		expected := "code=400, message=missing resources in one or more records within createRequest"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error [%s], got: [%v]", expected, err)
		}
	})
}

// ConsentQueryResponseMatcher a gomock matcher for ConsentQueryResponse (contains pointers)
type ConsentQueryResponseMatcher struct {
	want ConsentQueryResponse
}

// Matches checks the json of want and got objects
func (c ConsentQueryResponseMatcher) Matches(x interface{}) bool {
	resp, ok := x.(ConsentQueryResponse)

	if !ok {
		return false
	}

	wantBytes, err := json.Marshal(c.want)
	if err != nil {
		return false
	}

	gotBytes, err := json.Marshal(resp)
	if err != nil {
		return false
	}

	return string(wantBytes) == string(gotBytes)
}

func (c ConsentQueryResponseMatcher) String() string {
	return fmt.Sprintf("%v", c.want)
}

func TestDefaultConsentStore_QueryConsent(t *testing.T) {
	client := defaultConsentStore()
	crq := consentRuleForQuery()
	if err := client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{crq}); err != nil {
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
		json, _ := json.Marshal(query)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		v := 1
		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentQueryResponseMatcher{want: ConsentQueryResponse{
			TotalResults: 1,
			Results: []PatientConsent{
				{
					Id:        crq.ID,
					Subject:   "subject",
					Custodian: "custodian",
					Actor:     "actor",
					Records: []ConsentRecord{
						{
							RecordHash: crq.Records[0].Hash,
							Resources: []string{
								"resource",
							},
							ValidFrom: ValidFrom(time.Now().Add(time.Hour * -24).Format(pkg.Iso8601DateTime)),
							ValidTo:   ValidTo(time.Now().Add(time.Hour * 24).Format(pkg.Iso8601DateTime)),
							Version:   &v,
						},
					},
				},
			},
			Page: PageDefinition{},
		}})

		err := client.QueryConsent(echo)

		assert.NoError(t, err)
	})

	t.Run("API call returns 200 without results in time frame", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		query := consentQuery()
		tt := time.Now().Add(time.Hour * 25)
		query.ValidAt = &tt
		json, _ := json.Marshal(query)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentQueryResponse{
			TotalResults: 0,
			Results:      nil,
			Page:         PageDefinition{},
		})

		err := client.QueryConsent(echo)

		assert.NoError(t, err)
	})

	t.Run("API call returns 200 with results for subject search", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		query := consentQuery()
		subj := Identifier("subject")
		query.Subject = &subj
		json, _ := json.Marshal(query)
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		v := 1
		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, ConsentQueryResponse{
			TotalResults: 1,
			Results: []PatientConsent{
				{
					Id:        crq.ID,
					Subject:   "subject",
					Custodian: "custodian",
					Actor:     "actor",
					Records: []ConsentRecord{
						{
							RecordHash: crq.Records[0].Hash,
							Resources: []string{
								"resource",
							},
							ValidFrom: ValidFrom(time.Now().Add(time.Hour * -24).Format(pkg.Iso8601DateTime)),
							ValidTo:   ValidTo(time.Now().Add(time.Hour * 24).Format(pkg.Iso8601DateTime)),
							Version:   &v,
						},
					},
				},
			},
			Page: PageDefinition{},
		})

		err := client.QueryConsent(echo)

		assert.NoError(t, err)
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
}

func TestDefaultConsentStore_DeleteConsent(t *testing.T) {
	client := defaultConsentStore()
	crq := consentRuleForQuery()
	client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{crq})
	defer client.Cs.Shutdown()

	t.Run("missing consentRecordHash returns 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.DeleteConsent(echo, "")

		if assert.Error(t, err) {
			expected := "code=400, message=missing consentRecordHash"
			assert.Contains(t, err.Error(), expected)
		}
	})

	t.Run("Unknown consentRecordHash returns 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.DeleteConsent(echo, "a")

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), pkg.ErrorNotFound.Error())
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

func TestDefaultConsentStore_FindConsentRecord(t *testing.T) {
	client := defaultConsentStore()
	crq := consentRuleForQuery()
	client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{crq})
	defer client.Cs.Shutdown()

	t.Run("missing consentRecordHash returns 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.FindConsentRecord(echo, "", FindConsentRecordParams{})

		if assert.Error(t, err) {
			expected := "code=400, message=missing consentRecordHash"
			assert.Contains(t, err.Error(), expected)
		}
	})

	t.Run("Unknown consentRecordHash returns 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		err := client.FindConsentRecord(echo, "a", FindConsentRecordParams{})

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), pkg.ErrorNotFound.Error())
		}
	})

	t.Run("Correct find", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()
		echo.EXPECT().JSON(200, gomock.Any())

		tt := true
		err := client.FindConsentRecord(echo, crq.Records[0].Hash, FindConsentRecordParams{Latest: &tt})

		assert.NoError(t, err)
	})

	t.Run("find previous with latest flag", func(t *testing.T) {
		crq2 := consentRuleForQuery()
		crq2.ID = crq.ID
		crq2.Records[0].PreviousHash = &crq.Records[0].Hash
		if err := client.Cs.RecordConsent(context.Background(), []pkg.PatientConsent{crq2}); err != nil {
			assert.Fail(t, err.Error())
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		request := &http.Request{}

		echo.EXPECT().Request().Return(request).AnyTimes()

		tt := true
		err := client.FindConsentRecord(echo, crq.Records[0].Hash, FindConsentRecordParams{Latest: &tt})

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), pkg.ErrorConsentRecordNotLatest.Error())
		}
	})
}

func testConsent() PatientConsent {
	return PatientConsent{
		Id:        random.String(8),
		Actor:     Identifier("actor"),
		Custodian: Identifier("custodian"),
		Subject:   Identifier("subject"),
		Records: []ConsentRecord{
			{
				RecordHash: random.String(8),
				Resources:  []string{"resource"},
				ValidFrom:  ValidFrom("2019-01-01T12:00:00+01:00"),
				ValidTo:    ValidTo("2030-01-01T12:00:00+01:00"),
			},
		},
	}
}

func consentCheckRequest() ConsentCheckRequest {
	return ConsentCheckRequest{
		Subject:      Identifier("subject"),
		Custodian:    Identifier("custodian"),
		Actor:        Identifier("actor"),
		ResourceType: "resource",
	}
}

func consentQuery() QueryConsentJSONRequestBody {
	actor := Identifier("actor")
	subject := Identifier("subject")
	return QueryConsentJSONRequestBody{
		Actor:   &actor,
		Subject: &subject,
	}
}

func defaultConsentStore() Wrapper {
	client := pkg.ConsentStore{
		Config: pkg.ConsentStoreConfig{
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

	client.RunMigrations(client.Db.DB())

	return Wrapper{Cs: &client}
}

func consentRuleForQuery() pkg.PatientConsent {
	return pkg.PatientConsent{
		ID:        random.String(8),
		Subject:   "subject",
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
