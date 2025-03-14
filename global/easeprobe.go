/*
 * Copyright (c) 2022, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package global is the package for GuardianProbe
package global

import (
	"os"
	"strings"
	"time"
	_ "time/tzdata" // revive:disable

	log "github.com/sirupsen/logrus"
)

// GuardianProbe is the information of the program
type GuardianProbe struct {
	Name       string         `yaml:"name"`
	IconURL    string         `yaml:"icon"`
	Version    string         `yaml:"version"`
	Host       string         `yaml:"host"`
	TimeFormat string         `yaml:"time_format"`
	TimeZone   string         `yaml:"time_zone"`
	TimeLoc    *time.Location `yaml:"-"`
}

var easeProbe *GuardianProbe

func init() {
	InitGuardianProbe(DefaultProg, DefaultIconURL)
}

// InitGuardianProbe the GuardianProbe
func InitGuardianProbe(name, icon string) {
	InitGuardianProbeWithTime(name, icon, DefaultTimeFormat, DefaultTimeZone)
}

// InitGuardianProbeWithTime init the GuardianProbe with time
func InitGuardianProbeWithTime(name, icon, tf, tz string) {
	host, err := os.Hostname()
	if err != nil {
		log.Errorf("Get Hostname Failed: %s", err)
		host = "unknown"
	}
	easeProbe = &GuardianProbe{
		Name:    name,
		IconURL: icon,
		Version: Ver,
		Host:    host,
	}
	SetTimeZone(tz)
	SetTimeFormat(tf)
}

// GetGuardianProbe return the GuardianProbe
func GetGuardianProbe() *GuardianProbe {
	if easeProbe == nil {
		InitGuardianProbe(DefaultProg, DefaultIconURL)
	}
	return easeProbe
}

// GetTimeFormat return the time format
func GetTimeFormat() string {
	e := GetGuardianProbe()
	return e.TimeFormat
}

// SetTimeFormat set the time format
func SetTimeFormat(tf string) {
	if strings.TrimSpace(tf) == "" {
		tf = DefaultTimeFormat
	}
	e := GetGuardianProbe()
	e.TimeFormat = tf
}

// GetTimeLocation return the time zone
func GetTimeLocation() *time.Location {
	e := GetGuardianProbe()
	return e.TimeLoc
}

// SetTimeZone set the time zone
func SetTimeZone(tz string) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Errorf("Load TimeZone Failed: %s, use UTC time zone", err)
		tz = "UTC"
		loc = time.UTC
	}
	e := GetGuardianProbe()
	e.TimeZone = tz
	e.TimeLoc = loc
}

// FooterString return the footer string
// e.g. "GuardianProbe v1.0.0 @ localhost"
func FooterString() string {
	e := GetGuardianProbe()
	return e.Name + " " + e.Version + " @ " + e.Host
}
