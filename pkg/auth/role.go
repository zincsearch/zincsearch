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
	"strings"
	"sync"
	"time"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/metadata"
)

var ZINC_CACHED_PERMISSIONS = cachedPermissions{pm: map[string]map[string]struct{}{}}

type cachedPermissions struct {
	pm   map[string]map[string]struct{}
	lock sync.RWMutex
}

func (t *cachedPermissions) Get(id string) (map[string]struct{}, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	role, ok := t.pm[id]
	return role, ok
}

func (t *cachedPermissions) Set(id string, pm map[string]struct{}) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.pm[id] = pm
}

func (t *cachedPermissions) Delete(id string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.pm, id)
}

func strArrayToMap(ss []string) map[string]struct{} {
	m := map[string]struct{}{}
	for _, v := range ss {
		m[v] = struct{}{}
	}
	return m
}

func CreateRole(id, name string, permissions []string) (*meta.Role, error) {
	id = strings.ToLower(id)
	if id == "admin" {
		return nil, errors.New(errors.ErrorTypeInvalidArgument, "role id admin not allowed")
	}
	var newRole *meta.Role
	existingRole, roleExists, err := GetRole(id)
	if err != nil && !errors.Is(err, errors.ErrKeyNotFound) {
		return nil, err
	}

	if roleExists {
		newRole = existingRole
		newRole.Name = name
		newRole.Permission = permissions
		newRole.UpdatedAt = time.Now()
	} else {
		newRole = &meta.Role{
			ID:         id,
			Name:       name,
			Permission: permissions,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	err = SetRole(newRole.ID, *newRole)
	if err != nil {
		return nil, err
	}

	ZINC_CACHED_PERMISSIONS.Set(newRole.ID, strArrayToMap(permissions))

	return newRole, nil
}

func GetRoles() ([]*meta.Role, error) {
	return metadata.Role.List(0, 0)
}

func GetRole(id string) (*meta.Role, bool, error) {
	if id == "" {
		return nil, false, errors.New(errors.ErrorTypeInvalidArgument, "role id is required")
	}
	role, err := metadata.Role.Get(id)
	if err != nil {
		return nil, false, err
	}
	return role, true, nil
}

func SetRole(id string, role meta.Role) error {
	return metadata.Role.Set(id, role)
}

func DeleteRole(id string) error {
	id = strings.ToLower(id)
	ZINC_CACHED_PERMISSIONS.Delete(id)
	return metadata.Role.Delete(id)
}
