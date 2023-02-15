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
	"sort"
	"sync"

	"github.com/zinclabs/zincsearch/pkg/zutils/hash"
	"github.com/zinclabs/zincsearch/pkg/zutils/hash/fnv64"
)

type Rendezvous struct {
	nodes map[string]int
	nstr  []string
	nhash []uint64
	hash  hash.Hasher
	lock  sync.RWMutex
}

type scoreNode struct {
	name  string
	score uint64
}

func New() *Rendezvous {
	return NewWithHash(fnv64.NewDefaultHasher())
}

func NewWithHash(hash hash.Hasher) *Rendezvous {
	return &Rendezvous{
		nodes: make(map[string]int, 0),
		nstr:  make([]string, 0),
		nhash: make([]uint64, 0),
		hash:  hash,
	}
}

func (r *Rendezvous) Lookup(k string) string {
	khash := r.Hash(k)

	r.lock.RLock()
	defer r.lock.RUnlock()

	var midx int
	var mhash uint64
	for i, nhash := range r.nhash {
		if h := xorshiftMult64(khash ^ nhash); h > mhash {
			midx = i
			mhash = h
		}
	}

	return r.nstr[midx]
}

func (r *Rendezvous) LookupTopN(k string, n int) []string {
	khash := r.Hash(k)

	r.lock.RLock()
	scored := make([]scoreNode, len(r.nstr))
	for i, nhash := range r.nhash {
		h := xorshiftMult64(khash ^ nhash)
		scored[i] = scoreNode{name: r.nstr[i], score: h}
	}
	r.lock.RUnlock()

	sort.Slice(scored, func(i, j int) bool { return scored[i].score > scored[j].score })

	names := make([]string, 0, n)
	for i := 0; i < n && i < len(r.nstr); i++ {
		names = append(names, scored[i].name)
	}
	return names
}

func (r *Rendezvous) Contains(node string) bool {
	r.lock.RLock()
	_, ok := r.nodes[node]
	r.lock.RUnlock()
	return ok
}

func (r *Rendezvous) Add(node string) {
	if r.Contains(node) {
		return
	}

	r.lock.Lock()
	r.nodes[node] = len(r.nstr)
	r.nstr = append(r.nstr, node)
	r.nhash = append(r.nhash, r.Hash(node))
	r.lock.Unlock()
}

func (r *Rendezvous) Remove(node string) {
	r.lock.Lock()

	// find index of node to remove
	nidx := r.nodes[node]

	// remove from the slices
	l := len(r.nstr)
	r.nstr[nidx] = r.nstr[l]
	r.nstr = r.nstr[:l]

	r.nhash[nidx] = r.nhash[l]
	r.nhash = r.nhash[:l]

	// update the map
	delete(r.nodes, node)
	moved := r.nstr[nidx]
	r.nodes[moved] = nidx

	r.lock.Unlock()
}

func (r *Rendezvous) List() []string {
	r.lock.RLock()
	ns := make([]string, len(r.nstr))
	copy(ns, r.nstr)
	r.lock.RUnlock()
	return ns
}

func (r *Rendezvous) Len() int {
	r.lock.RLock()
	n := len(r.nstr)
	r.lock.RUnlock()
	return n
}

func (r *Rendezvous) Hash(name string) uint64 {
	return r.hash.Sum64(name)
}

func xorshiftMult64(x uint64) uint64 {
	// uses the "xorshift*" mix function which is simple and effective
	// see: https://en.wikipedia.org/wiki/Xorshift#xorshift*
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	return x * 0x2545F4914F6CDD1D
}
