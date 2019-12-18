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
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// RoundTripFunc
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn RoundTripFunc) HttpClient {
	return HttpClient{
		ServerAddress: "http://localhost:1323",
		customClient: &http.Client{
			Transport: RoundTripFunc(fn),
		},
		Logger: logrus.StandardLogger().WithField("component", "API-client"),
	}
}

func TestHttpClient_RecordConsent(t *testing.T) {
	t.Run("empty patient consent returns 201", func(t *testing.T) {
		client := testClient(201, []byte{})

		c := []pkg.PatientConsent{{}}
		err := client.RecordConsent(context.TODO(), c)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}
	})

	t.Run("full consent with validTo date returns 201", func(t *testing.T) {
		client := testClient(201, []byte{})

		c := []pkg.PatientConsent{patientConsent()}
		validTo := time.Now()
		c[0].Records[0].ValidTo = &validTo
		err := client.RecordConsent(context.TODO(), c)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}
	})

	t.Run("no rules error", func(t *testing.T) {
		client := testClient(201, []byte{})

		c := []pkg.PatientConsent{}
		err := client.RecordConsent(context.TODO(), c)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "at least one consent record is needed"
		if expected != err.Error() {
			t.Errorf("Expected error [%s], got [%v]", expected, err)
		}
	})

	t.Run("too many rules returns error", func(t *testing.T) {
		client := testClient(201, []byte{})

		c := []pkg.PatientConsent{{}, {}}
		err := client.RecordConsent(context.TODO(), c)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "creating multiple consent records currently not supported"
		if expected != err.Error() {
			t.Errorf("Expected error [%s], got [%v]", expected, err)
		}
	})

	t.Run("body read error returns error", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			return &http.Response{
				StatusCode: 500,
				Body:       errorCloser{},
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		})

		err := client.RecordConsent(context.TODO(), []pkg.PatientConsent{patientConsent()})

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "error while reading response body: error"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func TestHttpClient_DeleteConsentRecordByHash(t *testing.T) {
	t.Run("202", func(t *testing.T) {
		client := testClient(202, []byte{})

		res, err := client.DeleteConsentRecordByHash(context.TODO(), "hash")

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
			return
		}

		if !res {
			t.Errorf("Expected auth to be true, got false")
		}
	})

	t.Run("500", func(t *testing.T) {
		client := testClient(500, []byte("some error"))

		res, err := client.DeleteConsentRecordByHash(context.TODO(), "hash")

		if err == nil {
			t.Errorf("Expected error, got nothing")
			return
		}

		if res {
			t.Errorf("Expected delete to be false, got true")
		}

		expected := "consent store returned 500, reason: some error"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func TestHttpClient_ConsentAuth(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		tr := "true"
		resp, _ := json.Marshal(ConsentCheckResponse{ConsentGiven: &tr})
		client := testClient(200, resp)

		res, err := client.ConsentAuth(context.TODO(), "", "", "", "test", nil)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if !res {
			t.Errorf("Expected auth to be true, got false")
		}
	})

	t.Run("200 with checkpoint", func(t *testing.T) {
		tr := "true"
		resp, _ := json.Marshal(ConsentCheckResponse{ConsentGiven: &tr})
		client := testClient(200, resp)

		now := time.Now()
		res, err := client.ConsentAuth(context.TODO(), "", "", "", "test", &now)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if !res {
			t.Errorf("Expected auth to be true, got false")
		}
	})

	t.Run("body read error returns error", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			return &http.Response{
				StatusCode: 500,
				Body:       errorCloser{},
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		})

		_, err := client.ConsentAuth(context.TODO(), "custodian", "subject", "actor", "resource", nil)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "error while reading response body: error"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("client returns invalid json gives error", func(t *testing.T) {
		client := testClient(200, []byte("{"))

		_, err := client.ConsentAuth(context.TODO(), "custodian", "subject", "actor", "resource", nil)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "could not unmarshal response body, reason: unexpected end of JSON input"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func TestHttpClient_QueryConsentForActor(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		validTo := ValidTo("2029-01-01T12:00:00+01:00")
		resp, _ := json.Marshal(ConsentQueryResponse{
			Results: []PatientConsent{
				{
					Records: []ConsentRecord{
						{
							DataClasses: []string{"test"},
							ValidFrom:   "2019-01-01T12:00:00+01:00",
							ValidTo:     &validTo,
						},
					},
					Actor:     "actor",
					Subject:   "subject",
					Custodian: "custodian",
				},
			},
		})
		client := testClient(200, resp)
		a := "actor"
		res, err := client.QueryConsent(context.TODO(), &a, nil, nil, nil)

		if assert.NoError(t, err) {
			assert.Len(t, res, 1)
		}
	})
}

func TestHttpClient_FindConsentRecordByHash(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		resp, _ := json.Marshal(FromConsentRecord(consentRecord()))
		client := testClient(200, resp)
		res, err := client.FindConsentRecordByHash(context.TODO(), "hash", false)

		if assert.NoError(t, err) {
			assert.Equal(t, "Hash", res.Hash)
		}
	})
}

func TestHttpClient_QueryConsentForActorAndSubject(t *testing.T) {
	a := "actor"
	s := "urn:subject"

	t.Run("200", func(t *testing.T) {
		validTo := ValidTo("2029-01-01T12:00:00+01:00")
		resp, _ := json.Marshal(ConsentQueryResponse{
			Results: []PatientConsent{
				{
					Records: []ConsentRecord{
						{
							ValidFrom:   "2019-01-01T12:00:00+01:00",
							ValidTo:     &validTo,
							DataClasses: []string{"test"},
						},
					},
					Actor:     "actor",
					Subject:   "subject",
					Custodian: "custodian",
				},
			},
		})
		client := testClient(200, resp)
		res, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if len(res) != 1 {
			t.Errorf("Expected 1 patientConsent, got %d", len(res))
		}
	})

	t.Run("200 without results in time frame", func(t *testing.T) {
		resp, _ := json.Marshal(ConsentQueryResponse{Results: []PatientConsent{}})
		client := testClient(200, resp)
		tt := time.Now().Add(time.Hour)
		res, err := client.QueryConsent(context.TODO(), &a, nil, &s, &tt)

		if err != nil {
			t.Fatalf("Expected no error, got [%s]", err.Error())
		}

		assert.Len(t, res, 0)
	})

	t.Run("client returns error", func(t *testing.T) {
		client := testClient(500, []byte("error"))

		_, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "consent store returned 500, reason: error"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("client returns invalid json", func(t *testing.T) {
		client := testClient(200, []byte("{"))

		_, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "could not unmarshal response body, reason: unexpected end of JSON input"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})

	t.Run("body read error returns error", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			return &http.Response{
				StatusCode: 500,
				Body:       errorCloser{},
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		})

		_, err := client.QueryConsent(context.TODO(), &a, nil, &s, nil)

		if err == nil {
			t.Error("Expected error, got nothing")
			return
		}

		expected := "error while reading response body: error"
		if err.Error() != expected {
			t.Errorf("Expected error [%s], got [%v]", expected, err.Error())
		}
	})
}

func testClient(status int, body []byte) HttpClient {
	return newTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: status,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		}
	})
}
