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

package slack

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/o2ip/guardianprobe/global"
	"github.com/o2ip/guardianprobe/report"
	"github.com/stretchr/testify/assert"
)

func assertError(t *testing.T, err error, msg string) {
	t.Helper()
	assert.Error(t, err)
	assert.Equal(t, msg, err.Error())
}

func TestSlack(t *testing.T) {
	conf := &NotifyConfig{}
	conf.NotifyName = "dummy"
	err := conf.Config(global.NotifySettings{})
	assert.NoError(t, err)
	assert.Equal(t, "slack", conf.Kind())
	assert.Equal(t, report.Slack, conf.NotifyFormat)

	var client http.Client
	monkey.PatchInstanceMethod(reflect.TypeOf(&client), "Do", func(_ *http.Client, req *http.Request) (*http.Response, error) {
		r := io.NopCloser(strings.NewReader(`ok`))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	})
	err = conf.SendSlack("title", "message")
	assert.NoError(t, err)

	monkey.PatchInstanceMethod(reflect.TypeOf(&client), "Do", func(_ *http.Client, req *http.Request) (*http.Response, error) {
		r := io.NopCloser(strings.NewReader(`not found`))
		return &http.Response{
			StatusCode: 404,
			Body:       r,
		}, nil
	})
	err = conf.SendSlack("title", "message")
	assertError(t, err, "Error response from Slack with request body <message> - code [404] - msg [not found]")

	monkey.Patch(io.ReadAll, func(_ io.Reader) ([]byte, error) {
		return nil, errors.New("read error")
	})
	err = conf.SendSlack("title", "message")
	assertError(t, err, "read error")

	monkey.PatchInstanceMethod(reflect.TypeOf(&client), "Do", func(_ *http.Client, req *http.Request) (*http.Response, error) {
		return nil, errors.New("http do error")
	})
	err = conf.SendSlack("title", "message")
	assertError(t, err, "http do error")

	monkey.Patch(http.NewRequest, func(method string, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("new request error")
	})
	err = conf.SendSlack("title", "message")
	assertError(t, err, "new request error")

	monkey.UnpatchAll()
}
