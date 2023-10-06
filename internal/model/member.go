package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"

	"github.com/slashdevops/idp-scim-sync/internal/deepcopy"
)

// Member represents a member entity.
type Member struct {
	IPID     string `json:"ipid,omitempty"`
	SCIMID   string `json:"scimid,omitempty"`
	Email    string `json:"email,omitempty"`
	Status   string `json:"status,omitempty"`
	HashCode string `json:"hashCode,omitempty"`
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

	if err := enc.Encode(m.Status); err != nil {
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

	if err := dec.Decode(&m.Status); err != nil {
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
	HashCode  string    `json:"hashCode,omitempty"`
	Resources []*Member `json:"resources"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for MembersResult entity.
func (mr MembersResult) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(mr.Items); err != nil {
		return nil, err
	}

	for _, member := range mr.Resources {
		if err := enc.Encode(member); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for MembersResult entity.
func (mr *MembersResult) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&mr.Items); err != nil {
		return err
	}

	for i := 0; i < mr.Items; i++ {
		var member Member
		if err := dec.Decode(&member); err != nil {
			return err
		}
		mr.Resources = append(mr.Resources, &member)
	}

	return nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (mr *MembersResult) SetHashCode() {
	// this copy is necessary to avoid changing the original data
	// with the sort.Slice function and always be consistent
	// when calculating the hash code
	c := deepcopy.SliceOfPointers(mr.Resources)

	// only these fields are used in the hash calculation
	copyStruct := &MembersResult{
		Items:     mr.Items,
		Resources: c,
	}

	// order the resources by their hash code to be consistency always
	// NOTE: review this, it may be a performance issue and may not be necessary
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

// MarshalBinary implements the encoding.BinaryMarshaler interface for GroupMembers entity.
func (gm GroupMembers) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(gm.Items); err != nil {
		return nil, err
	}

	if gm.Group != nil {
		if err := enc.Encode(gm.Group); err != nil {
			return nil, err
		}
	}

	for _, member := range gm.Resources {
		if err := enc.Encode(member); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for GroupMembers entity.
func (gm *GroupMembers) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&gm.Items); err != nil {
		return err
	}

	if err := dec.Decode(&gm.Group); err != nil {
		if err.Error() != "EOF" {
			return err
		}
	}

	for i := 0; i < gm.Items; i++ {
		var member Member
		if err := dec.Decode(&member); err != nil {
			return err
		}
		gm.Resources = append(gm.Resources, &member)
	}

	return nil
}

// SetHashCode is a helper function to avoid errors when calculating hash code.
// this method discards fields that are not used in the hash calculation.
// only fields coming from the Identity Provider are used.
func (gm *GroupMembers) SetHashCode() {
	// this copy is necessary to avoid changing the original data
	// with the sort.Slice function and always be consistent
	// when calculating the hash code
	c := deepcopy.SliceOfPointers(gm.Resources)

	// only these fields are used in the hash calculation
	copiedStruct := &GroupMembers{
		Items:     gm.Items,
		Group:     gm.Group,
		Resources: c,
	}

	// to order the members of the group we used the email of the members
	// because this never could be empty and it is unique
	// NOTE: review this, it may be a performance issue and may not be necessary
	sort.Slice(copiedStruct.Resources, func(i, j int) bool {
		return copiedStruct.Resources[i].Email < copiedStruct.Resources[j].Email
	})

	gm.HashCode = Hash(copiedStruct)
}

// GroupsMembersResult represents a group members result list entity.
type GroupsMembersResult struct {
	Items     int             `json:"items"`
	HashCode  string          `json:"hashCode,omitempty"`
	Resources []*GroupMembers `json:"resources"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for GroupsMembersResult entity.
func (gmr GroupsMembersResult) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(gmr.Items); err != nil {
		return nil, err
	}

	for _, group := range gmr.Resources {
		if err := enc.Encode(group); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for GroupsMembersResult entity.
func (gmr *GroupsMembersResult) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&gmr.Items); err != nil {
		return err
	}

	for i := 0; i < gmr.Items; i++ {
		var group GroupMembers
		if err := dec.Decode(&group); err != nil {
			return err
		}
		gmr.Resources = append(gmr.Resources, &group)
	}

	return nil
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
	// this copy is necessary to avoid changing the original data
	// with the sort.Slice function and always be consistent
	// when calculating the hash code
	c := deepcopy.SliceOfPointers(gmr.Resources)

	// only these fields are used in the hash calculation
	copiedStruct := GroupsMembersResult{
		Items:     gmr.Items,
		Resources: c,
	}

	// to order the members of the group we used the email of the members
	// because this never could be empty and it is unique
	// NOTE: review this, it may be a performance issue and may not be necessary
	sort.Slice(copiedStruct.Resources, func(i, j int) bool {
		return copiedStruct.Resources[i].HashCode < copiedStruct.Resources[j].HashCode
	})

	gmr.HashCode = Hash(copiedStruct)
}
