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
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiWrapper struct {
	Cs *pkg.ConsentStore
}

func (w *ApiWrapper) CreateConsent(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	var createRequest = &SimplifiedConsent{}
	err = json.Unmarshal(buf, createRequest)

	if len(createRequest.Subject) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing subject in createRequest")
	}

	if len(createRequest.Custodian) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing custodian in createRequest")
	}

	if len(createRequest.Actors) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actors in createRequest")
	}

	if len(createRequest.Resources) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing resources in createRequest")
	}

	err = w.Cs.RecordConsent(ctx.Request().Context(), createRequest.ToConsentRule())

	if err != nil {
		return err
	}

	return ctx.NoContent(201)
}

func (w *ApiWrapper) CheckConsent(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
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

	if len(checkRequest.ResourceType) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing resourceType in checkRequest")
	}

	cr := checkRequest.ToConsentRule()
	cr.Resources = nil
	auth, err := w.Cs.ConsentAuth(ctx.Request().Context(), cr, checkRequest.ResourceType)

	if err != nil {
		return err
	}

	authValue := "false"
	if auth {
		authValue = "true"
	}

	checkResponse := ConsentCheckResponse{
		ConsentGiven: &authValue,
	}

	return ctx.JSON(200, checkResponse)
}

func (w *ApiWrapper) QueryConsent(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	var checkRequest = &ConsentQueryRequest{}
	err = json.Unmarshal(buf, checkRequest)

	if len(checkRequest.Actor) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing actor in queryRequest")
	}

	query := checkRequest.Query.(string)

	if len(query) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query in queryRequest")
	}

	var rules []pkg.ConsentRule

	if strings.Index(query, "urn") == 0 {
		rules, err = w.Cs.QueryConsentForActorAndSubject(ctx.Request().Context(), string(checkRequest.Actor), query)
	} else {
		rules, err = w.Cs.QueryConsentForActor(ctx.Request().Context(), string(checkRequest.Actor), query)
	}

	if err != nil {
		return err
	}

	logrus.Debugf("Found %d results", len(rules))

	results, err := FromSimplifiedConsentRule(rules)

	if err != nil {
		return err
	}

	return ctx.JSON(200,
		ConsentQueryResponse{
			Results: results,
			TotalResults: int32(len(results)),
		})
}
