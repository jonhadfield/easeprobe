//go:build windows
// +build windows

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

package log

import (
	"github.com/o2ip/guardianprobe/report"
	log "github.com/sirupsen/logrus"
)

// ConfigLog is the config for log
// Windows platform only support log file notification
func (c *NotifyConfig) ConfigLog() error {
	c.NotifyKind = "log"
	c.NotifyFormat = report.Log
	c.NotifySendFunc = c.Log

	c.logger = log.New()
	return c.configLogFile()
}
