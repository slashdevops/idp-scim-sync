package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
)

// Group represents a group entity.
type Group struct {
	IPID     string `json:"ipid,omitempty"`
	SCIMID   string `json:"scimid,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	HashCode string `json:"hashCode,omitempty"`
}

// MarshalBinary implements the gob.GobEncoder interface for Group entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (g Group) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(g.IPID); err != nil {
		return nil, err
	}
	if err := enc.Encode(g.Name); err != nil {
		return nil, err
	}
	if err := enc.Encode(g.Email); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for Group entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (g *Group) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&g.IPID); err != nil {
		return err
	}
	if err := dec.Decode(&g.Name); err != nil {
		return err
	}
	if err := dec.Decode(&g.Email); err != nil {
		return err
	}
	return nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (g *Group) SetHashCode() {
	g.HashCode = Hash(g)
}

// GroupsResult represents a group result list entity.
type GroupsResult struct {
	Items     int      `json:"items"`
	HashCode  string   `json:"hashCode"`
	Resources []*Group `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsResult entity.
func (gr *GroupsResult) MarshalJSON() ([]byte, error) {
	if gr.Resources == nil {
		gr.Resources = make([]*Group, 0)
	}
	return json.MarshalIndent(*gr, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gr *GroupsResult) SetHashCode() {
	copyResources := make([]*Group, len(gr.Resources))
	copy(copyResources, gr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &GroupsResult{
		Items:     gr.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].HashCode < copyStruct.Resources[j].HashCode
	})

	gr.HashCode = Hash(copyStruct)
}
