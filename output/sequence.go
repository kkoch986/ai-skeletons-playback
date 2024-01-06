package output

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Value struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type Sequence map[float64][]Value

func (s Sequence) AddValue(t float64, v Value) {
	_, ok := s[t]
	if !ok {
		s[t] = []Value{}
	}
	s[t] = append(s[t], v)
}

func (s Sequence) MarshalJSON() ([]byte, error) {
	m := make(map[string][]Value, len(s))
	for k, v := range s {
		m[fmt.Sprintf("%f", k)] = v
	}
	return json.Marshal(m)
}

func (s Sequence) UnmarshalJSON(p []byte) error {
	m := map[string][]Value{}
	err := json.Unmarshal(p, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		fv, _ := strconv.ParseFloat(k, 64)
		s[fv] = v
	}
	return nil
}
