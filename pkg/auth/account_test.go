/* Copyright 2022 Zinc Labs Inc. and Contributors
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

package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

func TestUpdateAccount(t *testing.T) {
	const id = "testupdateaccount"
	_, err := CreateUser(id, "Old Name", "oldpass1", "admin")
	assert.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		password    string
		newName     string
		newPassword string
		wantErr     bool
		wantName    string
		wantValid   string // password that must authenticate afterwards
	}{
		{name: "nothing to change", id: id, password: "oldpass1", wantErr: true, wantName: "Old Name", wantValid: "oldpass1"},
		{name: "wrong password", id: id, password: "wrong", newName: "X", newPassword: "newpass1", wantErr: true, wantName: "Old Name", wantValid: "oldpass1"},
		{name: "unknown user", id: "nobody", password: "oldpass1", newName: "X", wantErr: true, wantName: "Old Name", wantValid: "oldpass1"},
		{name: "name only", id: id, password: "oldpass1", newName: "New Name", wantName: "New Name", wantValid: "oldpass1"},
		{name: "password only", id: id, password: "oldpass1", newPassword: "newpass1", wantName: "New Name", wantValid: "newpass1"},
		{name: "both, case insensitive id", id: "TestUpdateAccount", password: "newpass1", newName: "Final", newPassword: "newpass2", wantName: "Final", wantValid: "newpass2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateAccount(tt.id, tt.password, tt.newName, tt.newPassword)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantName, got.Name)
			}
			_, ok := VerifyCredentials(id, tt.wantValid)
			assert.True(t, ok)
			// persisted, not only cached
			stored, exists, err := GetUser(id)
			assert.NoError(t, err)
			assert.True(t, exists)
			assert.Equal(t, tt.wantName, stored.Name)
			assert.Equal(t, GeneratePassword(tt.wantValid, stored.Salt), stored.Password)
		})
	}
	assert.NoError(t, DeleteUser(id))
}

type accountUserStoreMock struct {
	set func(string, meta.User) error
}

func (m accountUserStoreMock) Set(id string, user meta.User) error {
	return m.set(id, user)
}

func TestUpdateAccountCacheSnapshot(t *testing.T) {
	storageErr := errors.New("account storage failure")
	for _, tt := range []struct {
		name string
		err  error
	}{
		{name: "storage_failure", err: storageErr},
		{name: "success"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			id := "testupdateaccountcachesnapshot_" + tt.name
			_, err := CreateUser(id, "Old Name", "oldpass1", "admin")
			require.NoError(t, err)
			t.Cleanup(func() { assert.NoError(t, DeleteUser(id)) })

			snapshot, ok := VerifyCredentials(id, "oldpass1")
			require.True(t, ok)
			before := *snapshot
			storedBefore, exists, err := GetUser(id)
			require.NoError(t, err)
			require.True(t, exists)

			originalStore := accountUserStore
			t.Cleanup(func() { accountUserStore = originalStore })
			calls := 0
			accountUserStore = accountUserStoreMock{set: func(key string, user meta.User) error {
				calls++
				assert.Equal(t, id, key)
				cached, exists := ZINC_CACHED_USERS.Get(id)
				assert.True(t, exists)
				assert.Same(t, snapshot, cached)
				assert.Equal(t, before, *snapshot)
				if tt.err != nil {
					return tt.err
				}
				return SetUser(key, user)
			}}

			got, err := UpdateAccount(id, "oldpass1", "New Name", "newpass1")
			assert.Equal(t, 1, calls)
			assert.Equal(t, before, *snapshot)
			cached, exists := ZINC_CACHED_USERS.Get(id)
			require.True(t, exists)
			stored, exists, readErr := GetUser(id)
			require.NoError(t, readErr)
			require.True(t, exists)
			_, oldValid := VerifyCredentials(id, "oldpass1")
			_, newValid := VerifyCredentials(id, "newpass1")
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
				assert.Nil(t, got)
				assert.Same(t, snapshot, cached)
				assert.Equal(t, before, *cached)
				assert.Equal(t, storedBefore, stored)
				assert.True(t, oldValid)
				assert.False(t, newValid)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.NotSame(t, snapshot, cached)
				assert.Same(t, got, cached)
				assert.Equal(t, "New Name", cached.Name)
				assert.Equal(t, cached.Name, stored.Name)
				assert.Equal(t, cached.Salt, stored.Salt)
				assert.Equal(t, cached.Password, stored.Password)
				assert.Equal(t, GeneratePassword("newpass1", stored.Salt), stored.Password)
				assert.False(t, oldValid)
				assert.True(t, newValid)
			}
		})
	}
}
