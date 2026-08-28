package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"slices"
	"strings"

	"github.com/slashdevops/idp-scim-sync/internal/deepcopy"
)

// Group represents a group entity.
type Group struct {
	IPID     string `json:"ipid,omitempty"`
	SCIMID   string `json:"scimid,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	HashCode string `json:"hashCode,omitempty"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Group entity.
func (g Group) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

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

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Group entity.
func (g *Group) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

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
	HashCode  string   `json:"hashCode,omitempty"`
	Resources []*Group `json:"resources"`
	Items     int      `json:"items"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for GroupsResult entity.
func (gr GroupsResult) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(gr.Items); err != nil {
		return nil, err
	}

	for _, g := range gr.Resources {
		if err := enc.Encode(g); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for GroupsResult entity.
func (gr *GroupsResult) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&gr.Items); err != nil {
		return err
	}

	for range gr.Items {
		var g Group
		if err := dec.Decode(&g); err != nil {
			return err
		}
		gr.Resources = append(gr.Resources, &g)
	}

	return nil
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
	// this copy is necessary to avoid changing the original data
	// with the sort.Slice function and always be consistent
	// when calculating the hash code
	c := compactNilPointers(deepcopy.SliceOfPointers(gr.Resources))

	// only these fields are used in the hash calculation
	copiedStruct := &GroupsResult{
		Items:     gr.Items,
		Resources: c,
	}

	// Order by hash code so the result is consistent regardless of the order in
	// which resources arrived. SortStableFunc rather than SortFunc: hash codes
	// are not guaranteed unique across resources, and an unstable sort would
	// leave tied elements in an unspecified order, making the hash of a set
	// containing ties non-deterministic.
	slices.SortStableFunc(copiedStruct.Resources, func(a, b *Group) int {
		return strings.Compare(a.HashCode, b.HashCode)
	})

	gr.HashCode = Hash(copiedStruct)
}
