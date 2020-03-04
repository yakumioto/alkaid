/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yakumioto/alkaid/internal/db"
	"github.com/yakumioto/alkaid/internal/scheduler"
)

func dockerNetworkPrune() {
	c, _ := client.NewEnvClient()
	_, _ = c.NetworksPrune(context.Background(), filters.NewArgs())
}

func testInit() *gin.Engine {
	dockerNetworkPrune()
	rand.Seed(time.Now().Unix())
	gin.SetMode(gin.ReleaseMode)
	scheduler.Init()

	_ = db.Init("file:test.sqlite3", "mode=memory")

	Init()

	return gin.New()
}

func newJSONRequest(method, url string, data interface{}) *http.Request {
	if data == nil {
		return newNoneBodyRequest(method, url)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(dataBytes))

	return req
}

func newNoneBodyRequest(method, url string) *http.Request {
	req := httptest.NewRequest(method, url, nil)

	return req
}

func testHTTPEqual(t *testing.T, r http.Handler, method, url string, data interface{}, expectCode int, expectBody interface{}) {
	w := httptest.NewRecorder()

	var req *http.Request
	switch method {
	case http.MethodPost:
		req = newJSONRequest(method, url, data)
	case http.MethodGet:
		req = newNoneBodyRequest(method, url)
	}

	r.ServeHTTP(w, req)

	assert.Equal(t, expectCode, w.Code)
	switch expectBody.(type) {
	case nil:
	case string:
		assert.Equal(t, expectBody, w.Body.String())
	default:
		expectBytes, _ := json.Marshal(expectBody)
		assert.Equal(t, fmt.Sprintln(string(expectBytes)), w.Body.String())
	}
}
