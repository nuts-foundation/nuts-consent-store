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

package engine

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/nuts-consent-store/pkg/generated"
	"github.com/nuts-foundation/nuts-go/mock"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDefaultConsentStore_CheckConsent(t *testing.T) {

}

func TestDefaultConsentStore_CreateConsent(t *testing.T) {
	client := createTempEngine()
	defer client.Shutdown()

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

	t.Run("API call returns 400 for missing actor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		consent := testConsent()
		consent.Actors = []generated.Identifier{}

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
}

func testConsent() generated.SimplifiedConsent {
	return generated.SimplifiedConsent{
		Actors: []generated.Identifier{
			generated.Identifier("actor"),
		},
		Custodian:generated.Identifier("custodian"),
		Subject:generated.Identifier("subject"),
		Resources: []string{"resource"},
	}
}

func TestDefaultConsentStore_QueryConsent(t *testing.T) {

}
