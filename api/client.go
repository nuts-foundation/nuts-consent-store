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
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/sirupsen/logrus"
)

// HttpClient holds the server address and other basic settings for the http client
type HttpClient struct {
	ServerAddress string
	Timeout       time.Duration
	Logger        *logrus.Entry
	customClient  *http.Client
}

// FindConsentRecordByHash returns a ConsentRecord based on a hash. A latest flag can be added to indicate a record may only be returned if it's the latest in the chain.
func (hb HttpClient) FindConsentRecordByHash(context context.Context, consentRecordHash string, latest bool) (pkg.ConsentRecord, error) {

	var consentRecord pkg.ConsentRecord

	if len(consentRecordHash) == 0 {
		return consentRecord, ErrorMissingHash
	}

	result, err := hb.client().FindConsentRecord(context, consentRecordHash, &FindConsentRecordParams{Latest: &latest})
	if err != nil {
		err = fmt.Errorf("error while finding consent record in consent-store: %w", err)
		hb.Logger.Error(err)
		return consentRecord, err
	}

	body, err := hb.checkResponse(result)
	if err != nil {
		return consentRecord, err
	}

	var cr ConsentRecord
	if err := json.Unmarshal(body, &cr); err != nil {
		err = fmt.Errorf("could not unmarshal response body, reason: %w", err)
		hb.Logger.Error(err)
		return consentRecord, err
	}

	return cr.ToConsentRecord()
}

// QueryConsent returns PatientConsent records based on a combination of actor, custodian and subject. The only constraint is that either actor or custodian must not be empty.
func (hb HttpClient) QueryConsent(context context.Context, actor *string, custodian *string, subject *string, validAt *time.Time) ([]pkg.PatientConsent, error) {
	var (
		rules []pkg.PatientConsent
		req   QueryConsentJSONRequestBody
	)

	if validAt != nil {
		s := validAt.Format(time.RFC3339)
		req.ValidAt = &s
	}

	if actor != nil {
		a := Identifier(*actor)
		req.Actor = &a
	}

	if custodian != nil {
		c := Identifier(*custodian)
		req.Custodian = &c
	}

	if subject != nil {
		s := Identifier(*subject)
		req.Subject = &s
	}

	result, err := hb.client().QueryConsent(context, req)
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
		patientConsent, err := sr.ToPatientConsent()
		if err != nil {
			return rules, err
		}
		rules = append(rules, patientConsent)
	}

	return rules, nil
}

func (hb HttpClient) DeleteConsentRecordByHash(context context.Context, consentRecordHash string) (bool, error) {
	// delete record, if it doesn't exist an error is returned
	result, err := hb.client().DeleteConsent(context, consentRecordHash)
	if err != nil {
		err := fmt.Errorf("error while deleting consent in consent-store: %v", err)
		hb.Logger.Error(err)
		return false, err
	}

	_, err = hb.checkResponse(result)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ConsentAuth checks if there is an active consent for a given custodian, subject, actor, dataClass and an optional moment in time (checkpoint)
func (hb HttpClient) ConsentAuth(ctx context.Context, custodian string, subject string, actor string, dataClass string, checkpoint *time.Time) (bool, error) {
	req := CheckConsentJSONRequestBody{
		Actor:     Identifier(actor),
		Custodian: Identifier(custodian),
		Subject:   Identifier(subject),
		DataClass: dataClass,
	}

	if checkpoint != nil {
		s := checkpoint.Format(time.RFC3339)
		req.ValidAt = &s
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
func (hb HttpClient) RecordConsent(ctx context.Context, consent []pkg.PatientConsent) error {
	var req CreateConsentJSONRequestBody

	if len(consent) != 1 {
		var err error
		if len(consent) > 1 {
			err = errors.New("creating multiple consent records currently not supported")
		}
		if len(consent) == 0 {
			err = errors.New("at least one consent record is needed")
		}
		hb.Logger.Error(err)
		return err
	}

	req.Actor = Identifier(consent[0].Actor)
	req.Custodian = Identifier(consent[0].Custodian)
	req.Subject = Identifier(consent[0].Subject)

	for _, r := range consent[0].Records {
		version := int(r.Version)

		cr := ConsentRecord{
			RecordHash:         r.Hash,
			PreviousRecordHash: r.PreviousHash,
			ValidFrom:          ValidFrom(r.ValidFrom.Format(time.RFC3339)),
			Version:            &version,
		}
		if r.ValidTo != nil {
			validTo := ValidTo(r.ValidTo.Format(time.RFC3339))
			cr.ValidTo = &validTo

		}
		for _, sr := range r.DataClasses {
			cr.DataClasses = append(cr.DataClasses, sr.Code)
		}
		req.Records = append(req.Records, cr)
	}

	result, err := hb.client().CreateConsent(ctx, req)

	if err != nil {
		hb.Logger.Error("error while storing consent in consent-store", err)
		return err
	}

	_, err = hb.checkResponse(result)

	return err
}

// checkResponse analyzes response code and body. It returns the body.
// If body can not be read or status code >= 400 it returns the error.
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
			Client: hb.customClient,
		}
	}

	return &Client{
		Server: fmt.Sprintf("http://%v", hb.ServerAddress),
	}
}
