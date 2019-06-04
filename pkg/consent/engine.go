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

package consent

import (
	"context"
	"database/sql"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-consent-store/migrations"
	"github.com/nuts-foundation/nuts-consent-store/pkg"
	"github.com/nuts-foundation/nuts-consent-store/pkg/generated"
	engine "github.com/nuts-foundation/nuts-go/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go/types"
	"strings"
)

func NewConsentStoreEngine() *engine.Engine {
	cs := ConsentStore()

	return &engine.Engine{
		Name: "ConsentStore",
		Cmd: Cmd(),
		Configure: cs.Configure,
		Config: &cs.Config,
		ConfigKey: "cstore",
		Routes: func(router runtime.EchoRouter) {
			generated.RegisterHandlers(router, cs)
		},
		Start: cs.Start,
		Shutdown: cs.Shutdown,
	}
}

func Cmd() *cobra.Command {
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
				return types.Error{Msg: "requires an actor argument"}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := NewConsentStoreClient()

			var (
				consentList []pkg.ConsentRule
				err         error
			)

			if len(args) > 1 {
				consentList, err = csc.QueryConsentForActorAndSubject(context.TODO(), args[0], args[1])
			} else {
				consentList, err = csc.QueryConsentForActor(context.TODO(), args[0], "*")
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
		Use:     "record [subject] [custodian] [actor] [resources]",
		Example: "record urn:oid:2.16.840.1.113883.2.4.6.3:999999990 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 Observation,Patient",
		Short:   "record a new consent in store, resources is comma-separated",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return types.Error{Msg: "requires 4 arguments"}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := NewConsentStoreClient()

			resources := pkg.ResourcesFromStrings(strings.Split(args[3], ","))

			err := csc.RecordConsent(context.TODO(), []pkg.ConsentRule{
				{
					Subject:   args[0],
					Custodian: args[1],
					Actor:     args[2],
					Resources: resources,
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
		Use:     "check [subject] [custodian] [actor] [resource]",
		Example: "check urn:oid:2.16.840.1.113883.2.4.6.3:999999990 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 urn:oid:2.16.840.1.113883.2.4.6.1:00000007 Observation",
		Short:   "check if consent is given for the given combination",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return types.Error{Msg: "requires 4 arguments"}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			csc := NewConsentStoreClient()

			auth, err := csc.ConsentAuth(context.TODO(), pkg.ConsentRule{
				Subject:   args[0],
				Custodian: args[1],
				Actor:     args[2],
			}, args[3])

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

func (cs *DefaultConsentStore) Configure() error {
	var err error

	cs.configOnce.Do(func() {
		db, err := sql.Open("sqlite3", cs.connectionString)
		if err != nil {
			return
		}
		defer db.Close()

		// 1 ping
		err = db.Ping()
		if err != nil {
			return
		}

		// migrate
		err = runMigrations(db)
		if err != nil {
			return
		}

	})

	return err
}

func runMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	// wrap assets into Resource
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)

	if err != nil {
		return err
	}

	// run migrations
	m, err := migrate.NewWithInstance("go-bindata", d, "test", driver)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		return err
	}

	return nil
}

func (cs *DefaultConsentStore) Shutdown() error {
	return cs.db.Close()
}

func (cs *DefaultConsentStore) Start() error {
	var err error

	// gorm db connection
	cs.db, err = gorm.Open("sqlite3", cs.connectionString)

	return err
}

