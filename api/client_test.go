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
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
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
	t.Run("201", func(t *testing.T) {
		client := testClient(201, []byte{})

		c := []pkg.ConsentRule{{}}
		err := client.RecordConsent(context.TODO(), c)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}
	})

	t.Run("too many rules returns error", func(t *testing.T) {
		client := testClient(201, []byte{})

		c := []pkg.ConsentRule{}
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

		err := client.RecordConsent(context.TODO(), []pkg.ConsentRule{consentRule()})

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

func TestHttpClient_ConsentAuth(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		tr := "true"
		resp, _ := json.Marshal(ConsentCheckResponse{ConsentGiven: &tr})
		client := testClient(200, resp)

		cr := pkg.ConsentRule{}
		res, err := client.ConsentAuth(context.TODO(), cr, "test")

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

		_, err := client.ConsentAuth(context.TODO(), consentRule(), "resource")

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

		_, err := client.ConsentAuth(context.TODO(), consentRule(), "resource")

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
		resp, _ := json.Marshal(ConsentQueryResponse{Results: []SimplifiedConsent{
			{
				Resources: []string{"test"},
				Actors:    []Identifier{"actor"},
				Subject:   Identifier("subject"),
				Custodian: Identifier("custodian"),
			},
		}})
		client := testClient(200, resp)

		res, err := client.QueryConsentForActor(context.TODO(), "actor", "test")

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if len(res) != 1 {
			t.Errorf("Expected 1 consentRule, got %d", len(res))
		}
	})
}

func TestHttpClient_QueryConsentForActorAndSubject(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		resp, _ := json.Marshal(ConsentQueryResponse{Results: []SimplifiedConsent{
			{
				Resources: []string{"test"},
				Actors:    []Identifier{"actor"},
				Subject:   Identifier("subject"),
				Custodian: Identifier("custodian"),
			},
		}})
		client := testClient(200, resp)

		res, err := client.QueryConsentForActorAndSubject(context.TODO(), "actor", "urn:subject")

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if len(res) != 1 {
			t.Errorf("Expected 1 consentRule, got %d", len(res))
		}
	})

	t.Run("client returns error", func(t *testing.T) {
		client := testClient(500, []byte("error"))

		_, err := client.QueryConsentForActorAndSubject(context.TODO(), "actor", "urn:subject")

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

		_, err := client.QueryConsentForActorAndSubject(context.TODO(), "actor", "urn:subject")

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

		_, err := client.QueryConsentForActorAndSubject(context.TODO(), "actor", "urn:subject")

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