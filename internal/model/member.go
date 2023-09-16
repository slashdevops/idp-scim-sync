package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
)

// Member represents a member entity.
type Member struct {
	IPID     string `json:"ipid"`
	SCIMID   string `json:"scimid"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	HashCode string `json:"hashCode"`
}

// MarshalBinary implements the gob.GobEncoder interface for Member entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (m Member) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(m.IPID); err != nil {
		return nil, err
	}
	if err := enc.Encode(m.Email); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the gob.GobDecoder interface for Member entity.
// This is necessary to avoid include the value in the field SCIMID until
// the hashcode calculation is done.
// the Hash function use gob to calculate the hash code.
func (m *Member) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&m.IPID); err != nil {
		return err
	}
	if err := dec.Decode(&m.Email); err != nil {
		return err
	}
	return nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (m *Member) SetHashCode() {
	m.HashCode = Hash(m)
}

// MembersResult represents a member result list entity.
type MembersResult struct {
	Items     int       `json:"items"`
	HashCode  string    `json:"hashCode"`
	Resources []*Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (mr *MembersResult) SetHashCode() {
	copyResources := make([]*Member, len(mr.Resources))
	copy(copyResources, mr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &MembersResult{
		Items:     mr.Items,
		Resources: copyResources,
	}

	// order the resources by their hash code to be consistency always
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].IPID < copyStruct.Resources[j].IPID
	})

	mr.HashCode = Hash(copyStruct)
}

// GroupMembers represents a group members entity.
type GroupMembers struct {
	Items     int       `json:"items"`
	HashCode  string    `json:"hashCode,omitempty"`
	Group     *Group    `json:"group"`
	Resources []*Member `json:"resources"`
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gm *GroupMembers) SetHashCode() {
	copyResources := make([]*Member, len(gm.Resources))
	copy(copyResources, gm.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &GroupMembers{
		Items:     gm.Items,
		Group:     gm.Group,
		Resources: copyResources,
	}

	// to order the members of the group we used the email of the members
	// because this never could be empty and it is unique
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].Email < copyStruct.Resources[j].Email
	})

	gm.HashCode = Hash(copyStruct)
}

// GroupsMembersResult represents a group members result list entity.
type GroupsMembersResult struct {
	Items     int             `json:"items"`
	HashCode  string          `json:"hashCode"`
	Resources []*GroupMembers `json:"resources"`
}

// MarshalJSON implements the json.Marshaler interface for GroupsMembersResult entity.
func (gmr *GroupsMembersResult) MarshalJSON() ([]byte, error) {
	if gmr.Resources == nil {
		gmr.Resources = make([]*GroupMembers, 0)
	}
	return json.MarshalIndent(*gmr, "", "  ")
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gmr *GroupsMembersResult) SetHashCode() {
	copyResources := make([]*GroupMembers, len(gmr.Resources))
	copy(copyResources, gmr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := GroupsMembersResult{
		Items:     gmr.Items,
		Resources: copyResources,
	}

	// to order the members of the group we used the email of the members
	// because this never could be empty and it is unique
	sort.Slice(copyStruct.Resources, func(i, j int) bool {
		return copyStruct.Resources[i].HashCode < copyStruct.Resources[j].HashCode
	})

	gmr.HashCode = Hash(copyStruct)
}
