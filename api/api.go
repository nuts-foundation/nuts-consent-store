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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/sirupsen/logrus"
)

// Wrapper implements the ServerInterface for the base ConsentStore
type Wrapper struct {
	Cs *pkg.ConsentStore
}

// CreateConsent creates or updates a PatientConsent in the consent store
func (w *Wrapper) CreateConsent(ctx echo.Context) error {
	buf, err := readBody(ctx)
	if err != nil {
		return err
	}

	var createRequest = &PatientConsent{}
	err = json.Unmarshal(buf, createRequest)

	if len(createRequest.Id) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing ID in createRequest")
	}

	if len(createRequest.Subject) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing subject in createRequest")
	}

	if len(createRequest.Custodian) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing custodian in createRequest")
	}

	if len(createRequest.Actor) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actor in createRequest")
	}

	if len(createRequest.Records) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing records in createRequest")
	}

	for _, r := range createRequest.Records {
		if len(r.DataClasses) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "missing resources in one or more records within createRequest")
		}

		if len(r.RecordHash) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "missing recordHash in one or more records within createRequest")
		}
	}

	c, err := createRequest.ToPatientConsent()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = w.Cs.RecordConsent(ctx.Request().Context(), []pkg.PatientConsent{c})

	if err != nil {
		return err
	}

	return ctx.NoContent(201)
}

// CheckConsent checks if a given resource is allowed for a given actor, subject, custodian triple
func (w *Wrapper) CheckConsent(ctx echo.Context) error {
	buf, err := readBody(ctx)
	if err != nil {
		return err
	}

	var checkRequest = &ConsentCheckRequest{}
	err = json.Unmarshal(buf, checkRequest)

	if len(checkRequest.Subject) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing subject in checkRequest")
	}

	if len(checkRequest.Custodian) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing custodian in checkRequest")
	}

	if len(checkRequest.Actor) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actor in checkRequest")
	}

	if len(checkRequest.DataClass) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing dataClass in checkRequest")
	}

	var checkpoint *time.Time
	if checkRequest.ValidAt != nil {
		cp, err := time.Parse(time.RFC3339, *checkRequest.ValidAt)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid value for validAt: %s", *checkRequest.ValidAt))
		}
		checkpoint = &cp
	}

	auth, err := w.Cs.ConsentAuth(
		ctx.Request().Context(),
		string(checkRequest.Custodian),
		string(checkRequest.Subject),
		string(checkRequest.Actor),
		checkRequest.DataClass,
		checkpoint)

	if err != nil {
		return err
	}

	authValue := "no"
	if auth {
		authValue = "yes"
	}

	checkResponse := ConsentCheckResponse{
		ConsentGiven: &authValue,
	}

	return ctx.JSON(200, checkResponse)
}

// ErrorMissingHash is returned when the consentRecordHash parameter is missing
var ErrorMissingHash = errors.New("missing consentRecordHash")

// DeleteConsent deletes the consentRecord for a given consentRecordHash
func (w *Wrapper) DeleteConsent(ctx echo.Context, consentRecordHash string) error {
	if len(consentRecordHash) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ErrorMissingHash)
	}

	// delete record, if it doesn't exist an error is returned
	if f, err := w.Cs.DeleteConsentRecordByHash(ctx.Request().Context(), consentRecordHash); err != nil || !f {
		if !f {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.NoContent(202)
}

// FindConsentRecord returns a ConsentRecord based on a hash. A latest flag can be added to indicate a record may only be returned if it's the latest in the chain.
func (w *Wrapper) FindConsentRecord(ctx echo.Context, consentRecordHash string, params FindConsentRecordParams) error {
	if len(consentRecordHash) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ErrorMissingHash)
	}

	var (
		latest bool
		record pkg.ConsentRecord
		err    error
	)

	if params.Latest != nil {
		latest = *params.Latest
	}

	if record, err = w.Cs.FindConsentRecordByHash(ctx.Request().Context(), consentRecordHash, latest); err != nil {
		if errors.Is(err, pkg.ErrorNotFound) || errors.Is(err, pkg.ErrorConsentRecordNotLatest) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, FromConsentRecord(record))
}

// QueryConsent finds given consent for a combination of actor, subject and/or custodian
func (w *Wrapper) QueryConsent(ctx echo.Context) error {
	buf, err := readBody(ctx)
	if err != nil {
		return err
	}

	var checkRequest = &ConsentQueryRequest{}
	err = json.Unmarshal(buf, checkRequest)
	var (
		actor, custodian, subject *string
	)
	va := time.Now()

	if checkRequest.Actor != nil && len(*checkRequest.Actor) > 0 {
		actorString := string(*checkRequest.Actor)
		actor = &actorString
	}

	if checkRequest.Custodian != nil && len(*checkRequest.Custodian) > 0 {
		custodianString := string(*checkRequest.Custodian)
		custodian = &custodianString
	}

	var rules []pkg.PatientConsent

	if actor == nil && custodian == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actor or custodian in queryRequest")
	}

	if checkRequest.Subject != nil {
		s := string(*checkRequest.Subject)
		subject = &s
	}

	if checkRequest.ValidAt != nil {
		va, err = time.Parse(time.RFC3339, *checkRequest.ValidAt)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid format for validAt, required: %s", time.RFC3339))
		}
	}

	rules, err = w.Cs.QueryConsent(ctx.Request().Context(), actor, custodian, subject, &va)

	if err != nil {
		return err
	}

	logrus.Debugf("Found %d results", len(rules))

	results := FromPatientConsents(rules)

	if err != nil {
		return err
	}

	return ctx.JSON(200,
		ConsentQueryResponse{
			Results:      results,
			TotalResults: len(results),
		})
}

func readBody(ctx echo.Context) ([]byte, error) {
	req := ctx.Request()
	if req.Body == nil {
		msg := "missing body in request"
		logrus.Error(msg)
		return nil, echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := fmt.Sprintf("error reading request: %v", err)
		logrus.Error(msg)
		return nil, echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	return buf, nil
}
