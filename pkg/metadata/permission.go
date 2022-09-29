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

package metadata

type permission struct {
	ps []string
}

var Permission = &permission{ps: []string{}}

func (p *permission) List() []string {
	return p.ps
}

func (p *permission) Add(v string) {
	has := false
	for _, o := range p.ps {
		if o == v {
			has = true
		}
	}
	if !has {
		p.ps = append(p.ps, v)
	}
}
