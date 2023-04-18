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
	"sync"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

var ZINC_CACHED_USERS = cachedUsers{users: map[string]*meta.User{}}

type cachedUsers struct {
	users map[string]*meta.User
	lock  sync.RWMutex
}

func (t *cachedUsers) Get(userID string) (*meta.User, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	user, ok := t.users[userID]
	return user, ok
}

func (t *cachedUsers) Set(userID string, user *meta.User) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.users[userID] = user
}

func (t *cachedUsers) Delete(userID string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.users, userID)
}
