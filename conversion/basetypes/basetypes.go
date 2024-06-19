package basetypes

import "encoding/json"

type Int struct {
	Value  int64
	HasSet bool
	Unit   string
}

type Float64 struct {
	Value  float64
	HasSet bool
	Unit   string
}

type Bool struct {
	Value  bool
	HasSet bool
}

type String struct {
	Value  string
	HasSet bool
}

func (i *Int) Set(value int64, unit string) {
	i.Value = value
	i.HasSet = true
	i.Unit = unit
}

func (i Int) MarshalJSON() ([]byte, error) {
	if i.HasSet {
		if i.Unit != "" {
			return json.Marshal(struct {
				Value int64  `json:"value"`
				Unit  string `json:"unit"`
			}{
				Value: i.Value,
				Unit:  i.Unit,
			})
		} else {
			return json.Marshal(i.Value)
		}
	}
	return json.Marshal(nil)
}

func (f *Float64) Set(value float64, unit string) {
	f.Value = value
	f.HasSet = true
	f.Unit = unit
}

func (f Float64) MarshalJSON() ([]byte, error) {
	if f.HasSet {
		if f.Unit != "" {
			return json.Marshal(struct {
				Value float64 `json:"value"`
				Unit  string  `json:"unit"`
			}{
				Value: f.Value,
				Unit:  f.Unit,
			})
		} else {
			return json.Marshal(f.Value)
		}
	}
	return json.Marshal(nil)
}

func (b *Bool) Set(value bool) {
	b.Value = value
	b.HasSet = true
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b.HasSet {
		return json.Marshal(b.Value)
	}
	return json.Marshal(nil)
}

func (b *String) Set(value string) {
	b.Value = value
	b.HasSet = true
}

func (b String) MarshalJSON() ([]byte, error) {
	if b.HasSet {
		return json.Marshal(b.Value)
	}
	return json.Marshal(nil)
}
