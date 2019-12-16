package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Map contains configuration
	Map map[string]interface{}

	// Slice contains a list of configuration options
	Slice []interface{}
)

// Add to the map (deep merge)
func (m Map) Add(cfg Map) error {
	// so we can not deep merge if we have `.` in our own keys, we need to ensure our keys are clean first
	for k, v := range m {
		var toClean Map
		if strings.Contains(k, ".") {
			if toClean == nil {
				toClean = make(Map)
			}
			toClean[k] = v
			delete(m, k)
		}
		if toClean != nil {
			if err := m.Add(toClean); err != nil {
				return err
			}
		}
	}

	for k, v := range cfg {
		if vv, ok := v.(map[string]interface{}); ok {
			v = Map(vv)
		} else if vv, ok := v.([]interface{}); ok {
			v = Slice(vv)
		}

		if strings.Contains(k, ".") {
			k, sub := strings.SplitN(k, ".", 2)[0], strings.SplitN(k, ".", 2)[1]
			if mm, ok := m[k]; !ok {
				m[k] = make(Map)
				if err := m[k].(Map).Add(Map{sub: v}); err != nil {
					return err
				}
			} else if mm, ok := mm.(Map); ok {
				if err := mm.Add(Map{sub: v}); err != nil {
					return err
				}
			} else {
				return errors.Errorf("config conflict: %q.%q: %v into %v", k, sub, v, m[k])
			}
		} else {
			_, mapleft := m[k].(Map)
			_, mapright := v.(Map)
			// if left side already is a map and will be assigned to nil in config_dev
			if mapleft && v == nil {
				m[k] = nil
			} else if mapleft && mapright {
				if err := m[k].(Map).Add(v.(Map)); err != nil {
					return err
				}
			} else if mapleft && !mapright {
				return errors.Errorf("config conflict: %q:%v into %v", k, v, m[k])
			} else if mapright {
				m[k] = make(Map)
				if err := m[k].(Map).Add(v.(Map)); err != nil {
					return err
				}
			} else {
				// convert non-float64 to float64
				switch vv := v.(type) {
				case int:
					v = float64(vv)
				case int8:
					v = float64(vv)
				case int16:
					v = float64(vv)
				case int32:
					v = float64(vv)
				case int64:
					v = float64(vv)
				case uint:
					v = float64(vv)
				case uint8:
					v = float64(vv)
				case uint16:
					v = float64(vv)
				case uint32:
					v = float64(vv)
				case uint64:
					v = float64(vv)
				case float32:
					v = float64(vv)
				}
				m[k] = v
			}
		}
	}
	return nil
}

// Flat map
func (m Map) Flat() Map {
	res := make(Map)

	for k, v := range m {
		res[k] = v
		if v, ok := v.(Map); ok {
			for sk, sv := range v.Flat() {
				res[k+"."+sk] = sv
			}
		}
	}

	return res
}

// MapInto tries to map the configuration map into a given interface
func (m Map) MapInto(out interface{}) error {
	jsonBytes, err := json.Marshal(m)

	if err != nil {
		return errors.Wrap(err, "Problem with marshaling map")
	}

	err = json.Unmarshal(jsonBytes, out)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Problem with unmarshaling into given structure %T", out))
	}

	return nil
}

// Get a value by it's path
func (m Map) Get(key string) (interface{}, bool) {
	keyParts := strings.SplitN(key, ".", 2)
	val, ok := m[keyParts[0]]
	if len(keyParts) == 2 {
		mm, ok := val.(Map)
		if ok {
			return mm.Get(keyParts[1])
		}
		return mm, false
	}
	return val, ok
}

// MapInto tries to map the configuration map into a given interface
func (s Slice) MapInto(out interface{}) error {
	jsonBytes, err := json.Marshal(s)

	if err != nil {
		return errors.Wrap(err, "Problem with marshaling map")
	}

	err = json.Unmarshal(jsonBytes, &out)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Problem with unmarshaling into given structure %T", out))
	}

	return nil
}
