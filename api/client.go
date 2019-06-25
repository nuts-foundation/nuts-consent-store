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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// HttpClient holds the server address and other basic settings for the http client
type HttpClient struct {
	ServerAddress string
	Timeout       time.Duration
	Logger 		  *logrus.Entry
	customClient  *http.Client
}

func (hb HttpClient) ConsentAuth(ctx context.Context, consentRule pkg.ConsentRule, resourceType string) (bool, error) {
	req := ConsentCheckRequest{
		Actor:        Identifier(consentRule.Actor),
		Custodian:    Identifier(consentRule.Custodian),
		Subject:      Identifier(consentRule.Subject),
		ResourceType: resourceType,
	}

	result, err := hb.client().CheckConsent(ctx, req)
	if err != nil {
		err := fmt.Errorf("error while checking for consent in consent-store: %v", err)
		hb.Logger.Error(err)
		return false, err
	}

	body, err := hb.checkResponse(result)
	if err != nil {
		return false, err
	}

	var ccr ConsentCheckResponse
	if err := json.Unmarshal(body, &ccr); err != nil {
		err := fmt.Errorf("could not unmarshal response body, reason: %v", err)
		return false, err
	}

	return *ccr.ConsentGiven == "true", nil
}

// RecordConsent currently only supports the creation of a single record
func (hb HttpClient) RecordConsent(ctx context.Context, consent []pkg.ConsentRule) error {
	var req SimplifiedConsent

	if len(consent) != 1 {
		err := errors.New("creating multiple consent records currently not supported")
		hb.Logger.Error(err)
		return err
	}

	req.Actors = []Identifier{Identifier(consent[0].Actor)}
	req.Custodian = Identifier(consent[0].Custodian)
	req.Subject = Identifier(consent[0].Subject)

	for _, r := range consent[0].Resources {
		req.Resources = append(req.Resources, r.ResourceType)
	}

	result, err := hb.client().CreateConsent(ctx, req)

	if err != nil {
		hb.Logger.Error("error while storing consent in consent-store", err)
		return err
	}

	_, err = hb.checkResponse(result)
	if err != nil {
		return err
	}

	return nil
}

func (hb HttpClient) QueryConsentForActor(ctx context.Context, actor string, query string) ([]pkg.ConsentRule, error) {
	var rules []pkg.ConsentRule

	req := ConsentQueryRequest{
		Actor: Identifier(actor),
		Query: query,
	}

	result, err := hb.client().QueryConsent(ctx, req)
	if err != nil {
		err = fmt.Errorf("error while querying for consent in consent-store: %v", err)
		hb.Logger.Error(err)
		return rules, err
	}

	body, err := hb.checkResponse(result)
	if err != nil {
		return nil, err
	}

	var cqr ConsentQueryResponse
	if err := json.Unmarshal(body, &cqr); err != nil {
		err = fmt.Errorf("could not unmarshal response body, reason: %v", err)
		hb.Logger.Error(err)
		return rules, err
	}

	for _, sr := range cqr.Results {
		rule := pkg.ConsentRule{
			Actor:     actor,
			Subject:   string(sr.Subject),
			Custodian: string(sr.Custodian),
		}

		for _, r := range sr.Resources {
			rule.Resources = append(rule.Resources, pkg.Resource{ResourceType: r})
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// QueryConsentForActorAndSubject does the same as QueryConsentForActor, the backend just checks if the query starts with urn:
func (hb HttpClient) QueryConsentForActorAndSubject(ctx context.Context, actor string, subject string) ([]pkg.ConsentRule, error) {
	return hb.QueryConsentForActor(ctx, actor, subject)
}

func (hb *HttpClient) checkResponse(result *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		err = fmt.Errorf("error while reading response body: %v", err)
		hb.Logger.Error(err)
		return nil, err
	}

	if result.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("consent store returned %d, reason: %s", result.StatusCode, body)
		hb.Logger.Error(err.Error())
		return nil, err
	}

	return body, nil
}

func (hb HttpClient) client() *Client {
	if hb.customClient != nil {
		return &Client{
			Server: fmt.Sprintf("http://%v", hb.ServerAddress),
			Client: *hb.customClient,
		}
	}

	return &Client{
		Server: fmt.Sprintf("http://%v", hb.ServerAddress),
	}
}
