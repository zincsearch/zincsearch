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
	"encoding/base64"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/ider"
	"github.com/zinclabs/zincsearch/pkg/meta"
)

func CreateUser(id, name, plaintextPassword, role string) (*meta.User, error) {
	id = strings.ToLower(id)
	var newUser *meta.User
	existingUser, userExists, err := GetUser(id)
	if err != nil && !errors.Is(err, errors.ErrKeyNotFound) {
		return nil, err
	}

	if userExists {
		newUser = existingUser
		if plaintextPassword != "" {
			newUser.Salt = GenerateSalt()
			newUser.Password = GeneratePassword(plaintextPassword, newUser.Salt)
		}
		newUser.Name = name
		newUser.Role = role
		newUser.UpdatedAt = time.Now()
	} else {
		newUser = &meta.User{
			ID:        id,
			Name:      name,
			Role:      role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		newUser.Salt = GenerateSalt()
		newUser.Password = GeneratePassword(plaintextPassword, newUser.Salt)
	}

	err = SetUser(newUser.ID, *newUser)
	if err != nil {
		return nil, err
	}

	// cache user
	ZINC_CACHED_USERS.Set(newUser.ID, newUser)

	return newUser, nil
}

// Cache for password generation
var _cacheGeneratePassword = sync.Map{}

func GeneratePassword(password, salt string) string {
	key := password + ":" + salt
	if v, ok := _cacheGeneratePassword.Load(key); ok {
		return v.(string)
	}
	params := &Argon2Params{
		Memory:      2 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  128,
		KeyLength:   32,
		Time:        1,
		Threads:     1,
	}
	hash := argon2.IDKey([]byte(password), []byte(salt), params.Time, params.Memory, params.Threads, params.KeyLength)
	val := base64.StdEncoding.EncodeToString(hash)
	_cacheGeneratePassword.Store(key, val)
	return val
}

func GenerateSalt() string {
	return ider.Generate()
}

type Argon2Params struct {
	Time        uint32
	Memory      uint32
	Threads     uint8
	KeyLength   uint32
	SaltLength  uint32
	Parallelism uint8
	Iterations  uint32
}
