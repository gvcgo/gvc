package utils

import "encoding/json"

/*
Parser for koanf
*/
type JSON struct{}

func NewJsonParser() *JSON {
	return &JSON{}
}

// Unmarshal parses the given JSON bytes.
func (p *JSON) Unmarshal(b []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Marshal marshals the given config map to JSON bytes.
func (p *JSON) Marshal(o map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(o, "", "  ")
}
