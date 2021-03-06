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
	"context"
	"errors"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-consent-store/api"
	"github.com/nuts-foundation/nuts-consent-store/client"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	engine "github.com/nuts-foundation/nuts-go-core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewConsentStoreEngine() *engine.Engine {
	cs := pkg.ConsentStoreInstance()

	return &engine.Engine{
		Name:        "ConsentStore",
		Cmd:         cmd(),
		Configure:   cs.Configure,
		Config:      &cs.Config,
		ConfigKey:   "cstore",
		Diagnostics: cs.Diagnostics,
		FlagSet:     flagSet(),
		Routes: func(router engine.EchoRouter) {
			api.RegisterHandlers(router, &api.Wrapper{Cs: cs})
		},
		Start:    cs.Start,
		Shutdown: cs.Shutdown,
	}
}

func flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("cstore", pflag.ContinueOnError)

	flags.String(pkg.ConfigConnectionString, pkg.ConfigConnectionStringDefault, "Db connectionString")
	flags.String(pkg.ConfigAddress, "localhost:1323", "Address of the server when in client mode")
	flags.String(pkg.ConfigMode, "", "server or client, when client it uses the HttpClient")

	return flags
}

func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consent-store",
		Short: "consent store commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "list [actor] [subject]?",
		Example: "list urn:oid:2.16.840.1.113883.2.4.6.1:00000007",
		Short:   "lists all consent records for the given actor and optional subject",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires an actor argument")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := client.NewConsentStoreClient()

			var (
				consentList []pkg.PatientConsent
				err         error
			)

			if len(args) > 1 {
				consentList, err = csc.QueryConsent(context.TODO(), &args[0], nil, &args[1], nil)
			} else {
				consentList, err = csc.QueryConsent(context.TODO(), &args[0], nil, nil, nil)
			}

			if err != nil {
				logrus.Errorf("Error finding consent records: %s\n", err.Error())
				return
			}

			logrus.Errorf("Found %d records\n\n", len(consentList))

			for _, c := range consentList {
				logrus.Errorln(c.String())
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "record [subject] [custodian] [actor] [dataClasses]",
		Example: "record urn:oid:2.16.840.1.113883.2.4.6.3:999999990 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:1.3.6.1.4.1.54851:1:MEDICAL",
		Short:   "record a new consent in store, dataClasses is comma-separated",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return errors.New("requires 4 arguments")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := client.NewConsentStoreClient()

			resources := pkg.DataClassesFromStrings(strings.Split(args[3], ","))

			err := csc.RecordConsent(context.TODO(), []pkg.PatientConsent{
				{
					Subject:   args[0],
					Custodian: args[1],
					Actor:     args[2],
					Records: []pkg.ConsentRecord{
						{
							DataClasses: resources,
						},
					},
				},
			})

			if err != nil {
				logrus.Errorf("Error recording consent: %s\n", err.Error())
				return
			}

			logrus.Errorln("Consent recorded")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "check [subject] [custodian] [actor] [dataClass]",
		Example: "check urn:oid:2.16.840.1.113883.2.4.6.3:999999990 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:1.3.6.1.4.1.54851:1:MEDICAL",
		Short:   "check if there's an active consent record for the given combination",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return errors.New("requires 4 arguments")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := client.NewConsentStoreClient()

			auth, err := csc.ConsentAuth(
				context.TODO(),
				args[0],
				args[1],
				args[2],
				args[3],
				nil)

			if err != nil {
				logrus.Errorf("Error checking consent: %s", err.Error())
				return
			}

			if auth {
				logrus.Errorln("Consent given")
			} else {
				logrus.Errorln("No consent given")
			}
		},
	})

	return cmd
}
