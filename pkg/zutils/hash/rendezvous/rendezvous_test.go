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

package rendezvous

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	nodes = []string{"node1", "node2", "node3", "node4", "node5"}
	key   = "1QvEL0YywgM"
)

func TestRendezvous_Lookup(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				k: key,
			},
			want: "node2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			for _, node := range nodes {
				r.Add(node)
			}
			got := r.Lookup(tt.args.k)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRendezvous_LookupTopN(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				k: key,
			},
			want: []string{"node2", "node1", "node4", "node5", "node3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			for _, node := range nodes {
				r.Add(node)
			}
			gots := r.LookupTopN(tt.args.k, 10)
			assert.Equal(t, tt.want, gots)
		})
	}
}

func TestRendezvous_List(t *testing.T) {
	r := New()
	for _, node := range nodes {
		r.Add(node)
		r.Add(node)
	}
	got := r.List()
	assert.Equal(t, nodes, got)
}

func TestRendezvous_Len(t *testing.T) {
	r := New()
	for _, node := range nodes {
		r.Add(node)
	}
	r.Add(nodes[1])
	assert.Equal(t, len(nodes), r.Len())
}
