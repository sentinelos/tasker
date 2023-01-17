package labels

import (
	"fmt"
	"sort"
	"strings"
)

// Labels is a map of label:value. It implements Labels.
type Labels map[string]string

func (l Labels) Set(key, value string) {
	l[key] = value
}

// Has returns whether the provided label exists in the map
func (l Labels) Has(key string) bool {
	_, exists := l[key]
	return exists
}

// Get returns the value in the map for the provided label
func (l Labels) Get(key string) string {
	return l[key]
}

func (l Labels) Delete(key string) {
	delete(l, key)
}

// Range iterate over element in labels
func (l Labels) Range(fn func(key, value string) bool) {
	for k, v := range l {
		if !fn(k, v) {
			break
		}
	}
}

// Keys returns sorted keys
func (l Labels) Keys() (keys []string) {
	for k := range l {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// Values returns sorted values
func (l Labels) Values() (values []string) {
	for _, v := range l {
		values = append(values, v)
	}

	sort.Strings(values)

	return values
}

// Clone returns a deep copy of Labels
func (l Labels) Clone() Labels {
	fc := Labels{}
	for k, v := range l {
		fc[k] = v
	}

	return fc
}

// Merge returns a merged of Labels
func (l Labels) Merge(labels Labels) Labels {
	keys := l.Keys()
	fm := make(Labels, len(l)+len(labels))
	for _, k := range keys {
		fm[k] = l.Get(k)
	}
	for k, v := range labels {
		fm[k] = v
	}

	return fm
}

// String returns all labels listed as a human-readable string.
func (l Labels) String() string {
	keys := l.Keys()
	selector := make([]string, 0, len(l))

	for _, k := range keys {
		selector = append(selector, k+"=\""+l.Get(k)+"\"")
	}

	// Sort for determinism.
	sort.StringSlice(selector).Sort()
	return strings.Join(selector, ",")
}

// Conflicts takes 2 maps and returns true if there a key match between
// the maps but the value doesn't match, and returns false in other cases
func Conflicts(labels1, labels2 Labels) bool {
	small := labels1
	big := labels2
	if len(labels2) < len(labels1) {
		small = labels2
		big = labels1
	}

	for k, v := range small {
		if val, match := big[k]; match {
			if val != v {
				return true
			}
		}
	}

	return false
}

// Merge combines given maps, and does not check for any conflicts
// between the maps. In case of conflicts, second map (labels2) wins
func Merge(labels1, labels2 Labels) Labels {
	mergedMap := Labels{}

	for k, v := range labels1 {
		mergedMap[k] = v
	}
	for k, v := range labels2 {
		mergedMap[k] = v
	}
	return mergedMap
}

// Equals returns true if the given maps are equal
func Equals(labels1, labels2 Labels) bool {
	if len(labels1) != len(labels2) {
		return false
	}

	for k, v := range labels1 {
		value, ok := labels2[k]
		if !ok {
			return false
		}
		if value != v {
			return false
		}
	}
	return true
}

// SelectorToLabels converts selector string to labels map and validates keys and values
func SelectorToLabels(selector string) (Labels, error) {
	labels := Labels{}

	if len(selector) == 0 {
		return labels, nil
	}

	for _, label := range strings.Split(selector, ",") {
		l := strings.Split(label, "=")

		if len(l) != 2 {
			return labels, fmt.Errorf("invalid selector: %s", l)
		}

		key := strings.TrimSpace(l[0])
		value := strings.TrimSpace(l[1])
		labels[key] = value
	}
	return labels, nil
}
