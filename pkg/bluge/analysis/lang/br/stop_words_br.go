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

package br

import (
	"github.com/blugelabs/bluge/analysis"
)

// this content was obtained from:
// lucene-4.7.2/analysis/common/src/resources/org/apache/lucene/analysis/
// ` was changed to ' to allow for literal string

var StopWordsBytes = []byte(`# This file was created by Jacques Savoy and is distributed under the BSD license.
a
ainda
alem
ambas
ambos
antes
ao
aonde
aos
apos
aquele
aqueles
as
assim
com
como
contra
contudo
cuja
cujas
cujo
cujos
da
das
de
dela
dele
deles
demais
depois
desde
desta
deste
dispoe
dispoem
diversa
diversas
diversos
do
dos
durante
e
ela
elas
ele
eles
em
entao
entre
essa
essas
esse
esses
esta
estas
este
estes
ha
isso
isto
logo
mais
mas
mediante
menos
mesma
mesmas
mesmo
mesmos
na
nas
nao
nas
nem
nesse
neste
nos
o
os
ou
outra
outras
outro
outros
pelas
pelas
pelo
pelos
perante
pois
por
porque
portanto
proprio
propios
quais
qual
qualquer
quando
quanto
que
quem
quer
se
seja
sem
sendo
seu
seus
sob
sobre
sua
suas
tal
tambem
teu
teus
toda
todas
todo
todos
tua
tuas
tudo
um
uma
umas
uns
`)

func StopWords() analysis.TokenMap {
	rv := analysis.NewTokenMap()
	rv.LoadBytes(StopWordsBytes)
	return rv
}
