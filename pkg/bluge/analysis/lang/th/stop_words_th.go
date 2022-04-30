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

package th

import (
	"github.com/blugelabs/bluge/analysis"
)

// this content was obtained from:
// lucene-4.7.2/analysis/common/src/resources/org/apache/lucene/analysis/
// ` was changed to ' to allow for literal string

var StopWordsBytes = []byte(`# This file was created by Jacques Savoy and is distributed under the BSD license.
# Thai stopwords from:
# "Opinion Detection in Thai Political News Columns
# Based on Subjectivity Analysis"
# Khampol Sukhum, Supot Nitsuwat, and Choochart Haruechaiyasak
ไว้
ไม่
ไป
ได้
ให้
ใน
โดย
แห่ง
แล้ว
และ
แรก
แบบ
แต่
เอง
เห็น
เลย
เริ่ม
เรา
เมื่อ
เพื่อ
เพราะ
เป็นการ
เป็น
เปิดเผย
เปิด
เนื่องจาก
เดียวกัน
เดียว
เช่น
เฉพาะ
เคย
เข้า
เขา
อีก
อาจ
อะไร
ออก
อย่าง
อยู่
อยาก
หาก
หลาย
หลังจาก
หลัง
หรือ
หนึ่ง
ส่วน
ส่ง
สุด
สําหรับ
ว่า
วัน
ลง
ร่วม
ราย
รับ
ระหว่าง
รวม
ยัง
มี
มาก
มา
พร้อม
พบ
ผ่าน
ผล
บาง
น่า
นี้
นํา
นั้น
นัก
นอกจาก
ทุก
ที่สุด
ที่
ทําให้
ทํา
ทาง
ทั้งนี้
ทั้ง
ถ้า
ถูก
ถึง
ต้อง
ต่างๆ
ต่าง
ต่อ
ตาม
ตั้งแต่
ตั้ง
ด้าน
ด้วย
ดัง
ซึ่ง
ช่วง
จึง
จาก
จัด
จะ
คือ
ความ
ครั้ง
คง
ขึ้น
ของ
ขอ
ขณะ
ก่อน
ก็
การ
กับ
กัน
กว่า
กล่าว
`)

func StopWords() analysis.TokenMap {
	rv := analysis.NewTokenMap()
	rv.LoadBytes(StopWordsBytes)
	return rv
}
