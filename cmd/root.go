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
	goflag "flag"
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-consent-store/pkg/consent"
	"github.com/nuts-foundation/nuts-consent-store/pkg/generated"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var engine = consent.NewConsentStoreEngine()
var rootCmd = engine.Cmd

func Execute() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	if err := engine.Configure(); err != nil {
		panic(err)
	}

	// todo: as standalone, for now do it here
	if err := engine.Start(); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "run the store as standalone web server",
		Run: func(cmd *cobra.Command, args []string) {

			// start webserver
			e := echo.New()
			generated.RegisterHandlers(e, consent.ConsentStore())
			e.Logger.Fatal(e.Start(":1323"))
		},
	})



	rootCmd.Execute()

	// todo: as standalone, for now do it here
	if err := engine.Shutdown(); err != nil {
		panic(err)
	}
}
