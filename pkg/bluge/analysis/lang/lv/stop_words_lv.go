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

package lv

import (
	"github.com/blugelabs/bluge/analysis"
)

// this content was obtained from:
// lucene-4.7.2/analysis/common/src/resources/org/apache/lucene/analysis/
// ` was changed to ' to allow for literal string

var StopWordsBytes = []byte(`# This file was created by Jacques Savoy and is distributed under the BSD license.
# Set of Latvian stopwords from A Stemming Algorithm for Latvian, Karlis Kreslins
# the original list of over 800 forms was refined: 
#   pronouns, adverbs, interjections were removed
# 
# prepositions
aiz
ap
ar
apakš
ārpus
augšpus
bez
caur
dēļ
gar
iekš
iz
kopš
labad
lejpus
līdz
no
otrpus
pa
par
pār
pēc
pie
pirms
pret
priekš
starp
šaipus
uz
viņpus
virs
virspus
zem
apakšpus
# Conjunctions
un
bet
jo
ja
ka
lai
tomēr
tikko
turpretī
arī
kaut
gan
tādēļ
tā
ne
tikvien
vien
kā
ir
te
vai
kamēr
# Particles
ar
diezin
droši
diemžēl
nebūt
ik
it
taču
nu
pat
tiklab
iekšpus
nedz
tik
nevis
turpretim
jeb
iekam
iekām
iekāms
kolīdz
līdzko
tiklīdz
jebšu
tālab
tāpēc
nekā
itin
jā
jau
jel
nē
nezin
tad
tikai
vis
tak
iekams
vien
# modal verbs
būt  
biju 
biji
bija
bijām
bijāt
esmu
esi
esam
esat 
būšu     
būsi
būs
būsim
būsiet
tikt
tiku
tiki
tika
tikām
tikāt
tieku
tiec
tiek
tiekam
tiekat
tikšu
tiks
tiksim
tiksiet
tapt
tapi
tapāt
topat
tapšu
tapsi
taps
tapsim
tapsiet
kļūt
kļuvu
kļuvi
kļuva
kļuvām
kļuvāt
kļūstu
kļūsti
kļūst
kļūstam
kļūstat
kļūšu
kļūsi
kļūs
kļūsim
kļūsiet
# verbs
varēt
varēju
varējām
varēšu
varēsim
var
varēji
varējāt
varēsi
varēsiet
varat
varēja
varēs
`)

func StopWords() analysis.TokenMap {
	rv := analysis.NewTokenMap()
	rv.LoadBytes(StopWordsBytes)
	return rv
}
