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

package errors

import "errors"

var ErrCancelSignal = errors.New("cancelled") // just for cancel notice
var ErrNotFound = errors.New("not found")

var ErrIDNotFound = errors.New("id not found")
var ErrIDIsEmpty = errors.New("id is empty")

var ErrKeyNotFound = errors.New("key not found")
var ErrKeyIsEmpty = errors.New("key is empty")

var ErrIndexNotExists = errors.New("index not exists")
var ErrIndexIsExists = errors.New("index already exists")
var ErrIndexIsEmpty = errors.New("index is empty")

var ErrShardNotExists = errors.New("shard not exists")
var ErrShardIsExists = errors.New("shard already exists")

var ErrClusterModeNotSupported = errors.New("cluster mode not supported")
var ErrClusterTimeout = errors.New("cluster timeout")
