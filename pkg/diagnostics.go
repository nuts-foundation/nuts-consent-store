/*
 * Nuts consent store
 * Copyright (C) 2020. Nuts community
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

package pkg

import (
	"fmt"

	core "github.com/nuts-foundation/nuts-go-core"
)

type dbDiagnosticResult struct {
	connectionString string
	pingError        error
}

// Name returns the name of the dbDiagnosticResult
func (ddr dbDiagnosticResult) Name() string {
	return "DB"
}

// String returns the outcome of the dbDiagnosticResult
func (ddr dbDiagnosticResult) String() string {
	if ddr.pingError == nil {
		return fmt.Sprintf("connection string: %s, ping: true", ddr.connectionString)
	}

	return fmt.Sprintf("connection string: %s, ping: false, error: %v", ddr.connectionString, ddr.pingError)
}

// Diagnostics returns the slice of DiagnosticResults indicating the state of this engine
func (cs *ConsentStore) Diagnostics() []core.DiagnosticResult {
	dbState := dbDiagnosticResult{
		connectionString: cs.Config.Connectionstring,
		pingError:        cs.sqlDb.Ping(),
	}

	return []core.DiagnosticResult{
		dbState,
	}
}
