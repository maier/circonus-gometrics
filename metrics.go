// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import "strings"

// SetMetricTags DEPRECATED sets the tags for the named metric and flags a check update is needed
func (m *CirconusMetrics) SetMetricTags(name string, tags []string) bool {
	if m.check.UsingMetricRules() {
		m.Log.Println("circonusgometrics.SetMetricTags DEPRECATED - use stream tags metric_name|ST[a=val,b=val,...] format")
		delete(m.metricTags, name)
		return m.AddMetricTags(name, tags)
	}
	return m.check.AddMetricTags(name, tags, false)
}

// AddMetricTags DEPRECATED appends tags to any existing tags for the named metric and flags a check update is needed
func (m *CirconusMetrics) AddMetricTags(name string, tags []string) bool {
	if m.check.UsingMetricRules() {
		m.Log.Println("circonusgometrics.AddMetricTags DEPRECATED - use stream tags in name: metric_name|ST[a=val,b=val,...]")
		if _, exists := m.metricTags[name]; !exists {
			m.metricTags[name] = map[string]string{}
		}
		for _, tag := range tags {
			parts := strings.Split(tag, ":")
			if len(parts) > 2 {
				return false
			}
			m.metricTags[name][parts[0]] = parts[1]
		}
		return true
	}
	return m.check.AddMetricTags(name, tags, true)
}
