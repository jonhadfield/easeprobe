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

// Package guardian is the guardian notification package.
package guardian

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/megaease/easeprobe/global"
	"github.com/megaease/easeprobe/notify/base"
	"github.com/megaease/easeprobe/report"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	// MaxMessageLength is the max message length of a Guardian bot
	MaxMessageLength = 4096
)

// NotifyConfig is the guardian notification configuration
type NotifyConfig struct {
	base.DefaultNotify `yaml:",inline"`
}

// Config configures the guardian configuration
func (c *NotifyConfig) Config(gConf global.NotifySettings) error {
	c.NotifyKind = "guardian"
	c.NotifyFormat = report.JSON
	c.NotifySendFunc = c.SendGuardian
	c.DefaultNotify.Config(gConf)
	log.Debugf("Notification [%s] - [%s] configuration: %+v", c.NotifyKind, c.NotifyName, c)
	return nil
}

// splitMessage splits the message into parts
func splitMessage(message string) []string {
	var parts []string
	for len(message) > 0 {
		if len(message) > MaxMessageLength {
			parts = append(parts, message[:MaxMessageLength])
			message = message[MaxMessageLength:]
		} else {
			parts = append(parts, message)
			message = ""
		}
	}
	return parts
}

// SendGuardian is the wrapper for SendGuardianNotification
func (c *NotifyConfig) SendGuardian(title, text string) error {
	fmt.Println("TITLE:", title)
	fmt.Println("TEXT:", text)
	// parts := splitMessage(text)
	// for _, part := range parts {
	err := c.SendGuardianNotification(text)
	if err != nil {
		return err
	}
	// }
	return nil
}

type GuardianNotificationFormat struct {
	Name      string            `json:"name"`
	Endpoint  string            `json:"endpoint"`
	Time      time.Time         `json:"time"`
	Timestamp int64             `json:"timestamp"`
	Rtt       int               `json:"rtt"`
	Status    string            `json:"status"`
	Prestatus string            `json:"prestatus"`
	Message   string            `json:"message"`
	Tags      map[string]string `json:"tags"`
}

// SendGuardianNotification will send the notification to guardian.
func (c *NotifyConfig) SendGuardianNotification(text string) error {
	var notification GuardianNotificationFormat
	if err := json.Unmarshal([]byte(text), &notification); err != nil {
		return err
	}
	tags := make(map[string]string)
	tags["hosting"] = "Third Party"
	tags["category"] = "application"
	tags["platform"] = "github"
	tags["application"] = "github"
	notification.Tags = tags
	updated, err := json.Marshal(notification)
	api := "https://guardian.tukaz.co.uk/api/v1/monitoring/create"
	fmt.Println("sending-request-to", api)
	fmt.Println(string(updated))
	log.Debugf("[%s] - API %s", c.Kind(), api)
	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBufferString(string(updated)))
	if err != nil {
		return err
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error response from Guardian - code [%d] - msg [%s]", resp.StatusCode, string(buf))
	}
	return nil
}
