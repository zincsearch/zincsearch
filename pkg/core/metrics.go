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

package core

import (
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zincsearch/go-gin-prometheus"
)

var ZINC_METRICS *ginprometheus.Metric

func init() {
	ZINC_METRICS = &ginprometheus.Metric{
		ID:          "indexStats",                 // Identifier
		Name:        "index_stats",                // Metric Name
		Description: "Summary index stats metric", // Help Description
		Type:        "gauge_vec",                  // type associated with prometheus collector
		// Type Options:
		//	counter, counter_vec, gauge, gauge_vec,
		//	histogram, histogram_vec, summary, summary_vec
		Args: []string{"index", "field"},
	}
}

func SetMetricStatsByIndex(index, field string, val float64) {
	if ZINC_METRICS.MetricCollector == nil {
		return
	}
	ZINC_METRICS.MetricCollector.(*prometheus.GaugeVec).WithLabelValues(index, field).Set(val)
}

func IncrMetricStatsByIndex(index, field string) {
	if ZINC_METRICS.MetricCollector == nil {
		return
	}
	ZINC_METRICS.MetricCollector.(*prometheus.GaugeVec).WithLabelValues(index, field).Inc()
}
