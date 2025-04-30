/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContext_Render(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		obj          any
		acceptHeader string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "render nil",
			obj:          nil,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "render status code",
			obj:          http.StatusBadRequest,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "render plain text",
			obj:          "Hello, World!",
			acceptHeader: "text/plain",
			expectedCode: http.StatusOK,
			expectedBody: "Hello, World!",
		},
		{
			name:         "render error",
			obj:          fmt.Errorf("test error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"test error"}`,
		},
		{
			name:         "render json by default",
			obj:          map[string]string{"message": "hello"},
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"hello"}`,
		},
		{
			name: "render xml",
			obj: struct {
				XMLName xml.Name `xml:"response"`
				Message string   `xml:"message"`
			}{Message: "hello"},
			acceptHeader: "application/xml",
			expectedCode: http.StatusOK,
			expectedBody: `<response><message>hello</message></response>`,
		},
		{
			name:         "render yaml",
			obj:          map[string]string{"message": "hello"},
			acceptHeader: "application/yaml",
			expectedCode: http.StatusOK,
			expectedBody: "message: hello\n",
		},
		{
			name:         "render jsonp",
			obj:          map[string]string{"message": "hello"},
			acceptHeader: "application/javascript",
			expectedCode: http.StatusOK,
			expectedBody: "callback({\"message\":\"hello\"});",
		},
		{
			name:         "render image data",
			obj:          []byte("fake image data"),
			acceptHeader: "image/jpeg",
			expectedCode: http.StatusOK,
			expectedBody: "fake image data",
		},
		{
			name:         "render audio data",
			obj:          []byte("fake audio data"),
			acceptHeader: "audio/mpeg",
			expectedCode: http.StatusOK,
			expectedBody: "fake audio data",
		},
		{
			name:         "render binary data",
			obj:          []byte("binary data"),
			acceptHeader: "application/octet-stream",
			expectedCode: http.StatusOK,
			expectedBody: "binary data",
		},
		{
			name:         "render invalid status code",
			obj:          1000,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/?callback=callback", nil)
			if tc.acceptHeader != "" {
				c.Request.Header.Set("Accept", tc.acceptHeader)
			}

			ctx := &Context{c}
			ctx.RenderAny(tc.obj)

			assert.Equal(t, tc.expectedCode, c.Writer.Status())
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
