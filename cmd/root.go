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

package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nuts-foundation/nuts-consent-store/api"
	"github.com/nuts-foundation/nuts-consent-store/engine"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	cfg "github.com/nuts-foundation/nuts-go-core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var e = engine.NewConsentStoreEngine()
var rootCmd = e.Cmd

func Execute() {
	c := cfg.NutsConfig()
	c.IgnoredPrefixes = append(c.IgnoredPrefixes, e.ConfigKey)
	c.RegisterFlags(rootCmd, e)
	if err := c.Load(rootCmd); err != nil {
		panic(err)
	}

	c.PrintConfig(logrus.StandardLogger())

	if err := c.InjectIntoEngine(e); err != nil {
		panic(err)
	}

	if err := e.Configure(); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "run the store as standalone web server",
		Run: func(cmd *cobra.Command, args []string) {

			// start webserver
			e := echo.New()
			e.HideBanner = true
			e.Use(middleware.Logger())
			api.RegisterHandlers(e, &api.Wrapper{Cs: pkg.ConsentStoreInstance()})
			logrus.Fatal(e.Start(":1323"))
		},
	})

	if err := e.Start(); err != nil {
		panic(err)
	}

	rootCmd.Execute()

	if err := e.Shutdown(); err != nil {
		panic(err)
	}
}
