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

package bn

import (
	"github.com/blugelabs/bluge/analysis"
)

// this content was obtained from:
// lucene-4.7.2/analysis/common/src/resources/org/apache/lucene/analysis/
// ` was changed to ' to allow for literal string

var StopWordsBytes = []byte(`# This file was created by Jacques Savoy and is distributed under the BSD license.
# See http://members.unine.ch/jacques.savoy/clef/index.html.
# This file was created by Jacques Savoy and is distributed under the BSD license
এই
ও
থেকে
করে
এ
না
ওই
এক্
নিয়ে
করা
বলেন
সঙ্গে
যে
এব
তা
আর
কোনো
বলে
সেই
দিন
হয়
কি
দু
পরে
সব
দেওয়া
মধ্যে
এর
সি
শুরু
কাজ
কিছু
কাছে
সে
তবে
বা
বন
আগে
জ্নজন
পি
পর
তো
ছিল
এখন
আমরা
প্রায়
দুই
আমাদের
তাই
অন্য
গিয়ে
প্রযন্ত
মনে
নতুন
মতো
কেখা
প্রথম
আজ
টি
ধামার
অনেক
বিভিন্ন
র
হাজার
জানা
নয়
অবশ্য
বেশি
এস
করে
কে
হতে
বি
কয়েক
সহ
বেশ
এমন
এমনি
কেন
কেউ
নেওয়া
চেষ্টা
লক্ষ
বলা
কারণ
আছে
শুধু
তখন
যা
এসে
চার
ছিল
যদি
আবার
কোটি
উত্তর
সামনে
উপর
বক্তব্য
এত
প্রাথমিক
উপরে
আছে
প্রতি
কাজে
যখন
খুব
বহু
গেল
পেয়্র্
চালু
ই
নাগাদ
থাকা
পাচ
যাওয়া
রকম
সাধারণ
কমনে
`)

func StopWords() analysis.TokenMap {
	rv := analysis.NewTokenMap()
	rv.LoadBytes(StopWordsBytes)
	return rv
}
