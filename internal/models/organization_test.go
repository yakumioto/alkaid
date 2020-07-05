/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package models

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"

	"github.com/yakumioto/alkaid/internal/storage"
	"github.com/yakumioto/alkaid/internal/storage/mock"
)

func testDefaultOrganization() *Organization {
	return &Organization{
		ID: "test",
	}
}

func TestOrganization_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := mock.NewMockDBProvider(ctl)
	db = mockDB

	testCases := []struct {
		input    *Organization
		expected error
	}{
		{
			input:    testDefaultOrganization(),
			expected: nil,
		},
		{
			input:    testDefaultOrganization(),
			expected: storage.ErrAlreadyExist,
		},
	}

	for _, tc := range testCases {
		mockDB.EXPECT().Create(tc.input).Return(tc.expected)
		assert.Equal(t, tc.expected, testDefaultOrganization().Create())
	}
}

func TestOrganization_Update(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := mock.NewMockDBProvider(ctl)
	db = mockDB

	testCases := []struct {
		input    *Organization
		expected error
	}{
		{
			input:    testDefaultOrganization(),
			expected: nil,
		},
		{
			input:    testDefaultOrganization(),
			expected: storage.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		mockDB.EXPECT().Update(tc.input, tc.input).Return(tc.expected)
		assert.Equal(t, tc.expected, testDefaultOrganization().Update(testDefaultOrganization().ID))
	}
}
