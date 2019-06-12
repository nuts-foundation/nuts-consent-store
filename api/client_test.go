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
	}
}

func TestHttpClient_RecordConsent(t *testing.T) {
	t.Run("201", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			return &http.Response{
				StatusCode: 201,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
				Header:     http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		})

		c := []pkg.ConsentRule{{}}
		err := client.RecordConsent(context.TODO(), c)

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}
	})
}

func TestHttpClient_ConsentAuth(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		tr := "true"
		resp, _ := json.Marshal(ConsentCheckResponse{ConsentGiven: &tr})
		client := newTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(resp)),
				Header:     http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		})

		cr := pkg.ConsentRule{}
		res, err := client.ConsentAuth(context.TODO(), cr, "test")

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if !res {
			t.Errorf("Expected auth to be true, got false")
		}
	})
}

func TestHttpClient_QueryConsentForActor(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		resp, _ := json.Marshal(ConsentQueryResponse{Results: []SimplifiedConsent{
			{
				Resources: []string{"test"},
				Actors: []Identifier{"actor"},
				Subject:Identifier("subject"),
				Custodian:Identifier("custodian"),
			},
		}})
		client := newTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(resp)),
				Header:     http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		})

		res, err := client.QueryConsentForActor(context.TODO(), "actor", "test")

		if err != nil {
			t.Errorf("Expected no error, got [%s]", err.Error())
		}

		if len(res) != 1 {
			t.Errorf("Expected 1 consentRule, got %d", len(res))
		}
	})
}