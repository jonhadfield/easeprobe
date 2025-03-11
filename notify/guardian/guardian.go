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
	APIURL             string `yaml:"apiUrl" json:"apiUrl" jsonschema:"required,format=uri,title=API URL,description=The Guardian Notification API URL"`
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

// SendGuardian is the wrapper for SendGuardianNotification
func (c *NotifyConfig) SendGuardian(title, text string) error {
	err := c.SendGuardianNotification(text)
	if err != nil {
		return err
	}
	return nil
}

type NotificationFormat struct {
	Name      string            `json:"name"`
	Endpoint  string            `json:"endpoint"`
	Time      time.Time         `json:"time"`
	Timestamp int64             `json:"timestamp"`
	Rtt       int               `json:"rtt"`
	Status    string            `json:"status"`
	Prestatus string            `json:"prestatus"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels"`
}

// SendGuardianNotification will send the notification to guardian.
func (c *NotifyConfig) SendGuardianNotification(text string) error {
	var notification NotificationFormat
	if err := json.Unmarshal([]byte(text), &notification); err != nil {
		return err
	}

	type notificationSendFormat struct {
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

	n := notificationSendFormat{
		Name:      notification.Name,
		Endpoint:  notification.Endpoint,
		Time:      notification.Time,
		Timestamp: notification.Timestamp,
		Rtt:       notification.Rtt,
		Status:    notification.Status,
		Prestatus: notification.Prestatus,
		Message:   notification.Message,
		Tags:      notification.Labels,
	}

	updated, err := json.Marshal(n)
	apiUrl := c.APIURL
	log.Debugf("[%s] - API %s", c.Kind(), apiUrl)
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBufferString(string(updated)))
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
