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
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-go-core/mock"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	err error
}

var siws = []*ServerInterfaceWrapper{
	serverInterfaceWrapper(nil), serverInterfaceWrapper(errors.New("Server error")),
}

func (t *testServer) DeleteConsent(ctx echo.Context, proofHash string) error {
	return t.err
}

func (t *testServer) CreateConsent(ctx echo.Context) error {
	return t.err
}

func (t *testServer) CheckConsent(ctx echo.Context) error {
	return t.err
}

func (t *testServer) QueryConsent(ctx echo.Context) error {
	return t.err
}



func TestServerInterfaceWrapper_CheckConsent(t *testing.T) {
	for _, siw := range siws {
		t.Run("CheckConsent call returns expected error", func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/?", nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err := siw.CheckConsent(c)
			tsi := siw.Handler.(*testServer)
			if tsi.err != err {
				t.Errorf("Expected argument doesn't match given err %v <> %v", tsi.err, err)
			}
		})
	}
}

func TestServerInterfaceWrapper_CreateConsent(t *testing.T) {
	for _, siw := range siws {
		t.Run("CreateConsent call returns expected error", func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/?", nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err := siw.CreateConsent(c)
			tsi := siw.Handler.(*testServer)
			if tsi.err != err {
				t.Errorf("Expected argument doesn't match given err %v <> %v", tsi.err, err)
			}
		})
	}
}

func TestServerInterfaceWrapper_QueryConsent(t *testing.T) {
	for _, siw := range siws {
		t.Run("QueryConsent call returns expected error", func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/?", nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err := siw.QueryConsent(c)
			tsi := siw.Handler.(*testServer)
			if tsi.err != err {
				t.Errorf("Expected argument doesn't match given err %v <> %v", tsi.err, err)
			}
		})
	}
}

func TestRegisterHandlers(t *testing.T) {
	t.Run("Registers routes for crypto module", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockEchoRouter(ctrl)

		echo.EXPECT().POST("/consent", gomock.Any())
		echo.EXPECT().POST("/consent/check", gomock.Any())
		echo.EXPECT().POST("/consent/query", gomock.Any())
		echo.EXPECT().DELETE("/consent/:proofHash", gomock.Any())

		RegisterHandlers(echo, &testServer{})
	})
}

func serverInterfaceWrapper(err error) *ServerInterfaceWrapper {
	return &ServerInterfaceWrapper{
		Handler: &testServer{err: err},
	}
}
