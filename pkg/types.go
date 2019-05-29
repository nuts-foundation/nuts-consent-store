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

package pkg

import (
	"fmt"
	"strings"
)

// ConsentRule defines struct for consent_rule table.
type ConsentRule struct {
	ID        uint   `gorm:"AUTO_INCREMENT"`
	Actor     string `gorm:"not null"`
	Custodian string `gorm:"not null"`
	Resources []Resource
	Subject   string `gorm:"not null"`
}

func (ConsentRule) TableName() string {
	return "consent_rule"
}

// Resource defines struct for resource table
type Resource struct {
	ConsentRuleID uint
	ResourceType  string `gorm:"not null"`
}

func (Resource) TableName() string {
	return "resource"
}

func (se *ConsentRule) String() string {
	return fmt.Sprintf("%s@%s for %s: %s", se.Subject, se.Custodian, se.Actor, resourceJoin(se.Resources, ","))
}

// SameTriple compares this ConsentRule with another one on just Actor, Custiodian and Subject
func (se *ConsentRule) SameTriple(other *ConsentRule) bool {
	return se.Subject == other.Subject && se.Custodian == other.Custodian && se.Actor == other.Actor
}

func (r *Resource) String() string {
	return r.ResourceType
}

// ResourcesFromStrings converts a slice of strings to a slice of Recources
func ResourcesFromStrings(list []string) []Resource {
	a := make([]Resource, len(list))
	for i, l := range list {
		a[i] = Resource{ResourceType: l}
	}
	return a
}

func resourceJoin(slice []Resource, sep string) string {
	a := make([]string, len(slice))

	for i, r := range slice {
		a[i] = r.ResourceType
	}

	return strings.Join(a, sep)
}
