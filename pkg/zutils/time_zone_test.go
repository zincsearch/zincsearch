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

package zutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeZone(t *testing.T) {
	date := "2020-02-02 02:02:02"
	layout := "2006-01-02 15:04:05"
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Location
		wantErr bool
	}{
		{
			name: "UTC",
			args: args{
				name: "UTC",
			},
			want:    time.UTC,
			wantErr: false,
		},
		{
			name: "Local",
			args: args{
				name: "Local",
			},
			want:    time.Local,
			wantErr: false,
		},
		{
			name: "+01:00",
			args: args{
				name: "+01:00",
			},
			want:    time.FixedZone("UTC", 3600),
			wantErr: false,
		},
		{
			name: "+08:00",
			args: args{
				name: "+08:00",
			},
			want:    time.FixedZone("UTC", 3600*8),
			wantErr: false,
		},
		{
			name: "-08:00",
			args: args{
				name: "-08:00",
			},
			want:    time.FixedZone("UTC", -3600*8),
			wantErr: false,
		},
		{
			name: "+0800",
			args: args{
				name: "+0800",
			},
			want:    time.FixedZone("UTC", 3600*8),
			wantErr: false,
		},
		{
			name: "-0800",
			args: args{
				name: "-0800",
			},
			want:    time.FixedZone("UTC", -3600*8),
			wantErr: false,
		},
		{
			name: "error timezone",
			args: args{
				name: "xxxx",
			},
			want:    time.UTC,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimeZone(tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			wantTime, _ := time.ParseInLocation(layout, date, tt.want)
			gotTime, err := time.ParseInLocation(layout, date, got)
			assert.NoError(t, err)
			assert.True(t, wantTime.Equal(gotTime))
		})
	}
}
