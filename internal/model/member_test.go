package model

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMember_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest Member
	}{
		{
			name:   "empty",
			toTest: Member{},
		},
		{
			name: "filled",
			toTest: Member{
				IPID:     "1",
				SCIMID:   "1",
				Email:    "user.1@mail.com",
				Status:   "ACTIVE",
				HashCode: "",
			},
		},
		{
			name: "partial filled",
			toTest: Member{
				IPID:     "1",
				Email:    "user.1@mail.com",
				HashCode: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("Member.MarshalBinary() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got Member
			if err := dec.Decode(&got); err != nil {
				t.Errorf("Member.UnmarshalBinary() error = %v", err)
			}

			// SCIMID is not exported, so it will not be encoded
			// HashCode is not exported, so it will not be encoded
			expected := Member{
				IPID:   tt.toTest.IPID,
				Email:  tt.toTest.Email,
				Status: tt.toTest.Status,
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(expected, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("UsersResult.GobEncode() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestMember_SetHashCode(t *testing.T) {
	tests := []struct {
		name   string
		member Member
		want   Member
	}{
		{
			name: "success",
			member: Member{
				IPID:     "1",
				SCIMID:   "1",
				Email:    "user.1@mail.com",
				Status:   "ACTIVE",
				HashCode: "test",
			},
			want: Member{
				IPID:   "1",
				Email:  "user.1@mail.com",
				Status: "ACTIVE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.member.SetHashCode()
			tt.want.SetHashCode()

			got := tt.member.HashCode

			if got != tt.want.HashCode {
				t.Errorf("Member.SetHashCode() = %s, want %s", got, tt.want.HashCode)
			}
		})
	}
}

func TestMembersResult_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest MembersResult
	}{
		{
			name:   "empty",
			toTest: MembersResult{},
		},
		{
			name: "filled",
			toTest: MembersResult{
				Items:    2,
				HashCode: "test",
				Resources: []*Member{
					{IPID: "1", SCIMID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
					{IPID: "2", SCIMID: "2", Email: "user.2@mail.com", Status: "ACTIVE"},
				},
			},
		},
		{
			name: "partial filled",
			toTest: MembersResult{
				Items:     1,
				HashCode:  "test",
				Resources: []*Member{{IPID: "1", Email: "user.1@mail.com"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("MembersResult.MarshalBinary() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got MembersResult
			if err := dec.Decode(&got); err != nil {
				t.Errorf("MembersResult.UnmarshalBinary() error = %v", err)
			}

			var expectedMembers []*Member
			for _, m := range tt.toTest.Resources {
				expectedMembers = append(expectedMembers, &Member{
					IPID:   m.IPID,
					Email:  m.Email,
					Status: m.Status,
				})
			}

			// HashCode is not exported, so it will not be encoded
			expected := MembersResult{
				Items:     tt.toTest.Items,
				Resources: expectedMembers,
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(expected, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("MembersResult.GobEncode() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestMembersResult_SetHashCode(t *testing.T) {
	m1 := &Member{IPID: "1", SCIMID: "1", Email: "group.1@mail.com", Status: "ACTIVE"}
	m2 := &Member{IPID: "2", SCIMID: "2", Email: "group.2@mail.com", Status: "ACTIVE"}
	m3 := &Member{IPID: "3", SCIMID: "3", Email: "group.3@mail.com", Status: "ACTIVE"}

	m1.SetHashCode()
	m2.SetHashCode()
	m3.SetHashCode()

	mr1 := MembersResult{
		Items:     3,
		Resources: []*Member{m1, m2, m3},
	}
	mr1.SetHashCode()

	mr2 := MembersResult{
		Items:     3,
		Resources: []*Member{m2, m3, m1},
	}
	mr2.SetHashCode()

	mr3 := MembersResult{
		Items:     3,
		Resources: []*Member{m3, m2, m1},
	}
	mr3.SetHashCode()

	if mr1.HashCode != mr2.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if mr1.HashCode != mr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if mr2.HashCode != mr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
}

func TestGroupsMembersResult_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest GroupsMembersResult
	}{
		{
			name:   "empty",
			toTest: GroupsMembersResult{},
		},
		{
			name: "filled",
			toTest: GroupsMembersResult{
				Items:    2,
				HashCode: "test",
				Resources: []*GroupMembers{
					{
						Items: 2,
						Group: &Group{
							IPID:   "1",
							SCIMID: "1",
							Name:   "group 1",
							Email:  "user.1@mail.com",
						},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
							{IPID: "2", SCIMID: "2", Email: "user.2@mail.com", Status: "ACTIVE"},
						},
					},
					{
						Items: 3,
						Group: &Group{
							IPID:   "2",
							SCIMID: "2",
							Name:   "group 2",
							Email:  "user.2@mail.com",
						},
						Resources: []*Member{
							{IPID: "3", SCIMID: "3", Email: "user.3@mail.com", Status: "ACTIVE"},
							{IPID: "4", SCIMID: "4", Email: "user.4@mail.com", Status: "ACTIVE"},
							{IPID: "5", SCIMID: "5", Email: "user.5@mail.com", Status: "ACTIVE"},
						},
					},
				},
			},
		},
		{
			name: "partial filled",
			toTest: GroupsMembersResult{
				Items:    1,
				HashCode: "test",
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: &Group{
							IPID:  "1",
							Email: "user.1@mail.com",
						},
						Resources: []*Member{
							{IPID: "1", SCIMID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("GroupsMembersResult.MarshalBinary() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got GroupsMembersResult
			if err := dec.Decode(&got); err != nil {
				t.Errorf("GroupsMembersResult.UnmarshalBinary() error = %v", err)
			}

			var expectedGroups []*GroupMembers
			for _, m := range tt.toTest.Resources {
				var expectedMembers []*Member
				for _, member := range m.Resources {
					expectedMembers = append(expectedMembers, &Member{
						IPID:   member.IPID,
						Email:  member.Email,
						Status: member.Status,
					})
				}

				expectedGroups = append(expectedGroups, &GroupMembers{
					Items: m.Items,
					Group: &Group{
						IPID:  m.Group.IPID,
						Name:  m.Group.Name,
						Email: m.Group.Email,
					},
					Resources: expectedMembers,
				})
			}

			// HashCode is not exported, so it will not be encoded
			expected := GroupsMembersResult{
				Items:     tt.toTest.Items,
				Resources: expectedGroups,
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(expected, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupsMembersResult.GobEncode() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestGroupsMembersResult_SetHashCode(t *testing.T) {
	m1 := &Member{IPID: "1", SCIMID: "1", Email: "group.1@mail.com"}
	m2 := &Member{IPID: "2", SCIMID: "2", Email: "group.2@mail.com"}
	m3 := &Member{IPID: "3", SCIMID: "3", Email: "group.3@mail.com"}

	m1.SetHashCode()
	m2.SetHashCode()
	m3.SetHashCode()

	gm1 := &GroupMembers{Group: &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"}, Resources: []*Member{m1, m2, m3}}
	gm2 := &GroupMembers{Group: &Group{IPID: "2", SCIMID: "2", Name: "group", Email: "group.2@mail.com"}, Resources: []*Member{m2, m1, m3}}
	gm3 := &GroupMembers{Group: &Group{IPID: "3", SCIMID: "3", Name: "group", Email: "group.3@mail.com"}, Resources: []*Member{m1, m3, m2}}

	gm1.SetHashCode()
	gm2.SetHashCode()
	gm3.SetHashCode()

	gmr1 := GroupsMembersResult{
		Items:     3,
		Resources: []*GroupMembers{gm1, gm2, gm3},
	}
	gmr1.SetHashCode()

	gmr2 := GroupsMembersResult{
		Items:     3,
		Resources: []*GroupMembers{gm2, gm3, gm1},
	}
	gmr2.SetHashCode()

	gmr3 := GroupsMembersResult{
		Items:     3,
		Resources: []*GroupMembers{gm3, gm2, gm1},
	}
	gmr3.SetHashCode()

	gmr4 := MergeGroupsMembersResult(&gmr2, &gmr1, &gmr3)
	gmr4.SetHashCode()
	gmr5 := MergeGroupsMembersResult(&gmr3, &gmr2, &gmr1)
	gmr5.SetHashCode()

	t.Logf("gmr4.HashCode: %s\n", gmr4.HashCode)
	t.Logf("gmr5.HashCode: %s\n", gmr5.HashCode)

	if gmr1.HashCode != gmr2.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if gmr1.HashCode != gmr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}
	if gmr2.HashCode != gmr3.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal")
	}

	if gmr5.HashCode != gmr4.HashCode {
		t.Errorf("GroupsMembersResult.HashCode should be equal: gmr5-> %s, gmr4-> %s", gmr5.HashCode, gmr4.HashCode)
	}
}

func TestGroupMembers_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest GroupMembers
	}{
		{
			name:   "empty",
			toTest: GroupMembers{},
		},
		{
			name: "filled",
			toTest: GroupMembers{
				Items:    2,
				HashCode: "test",
				Group: &Group{
					IPID:   "1",
					SCIMID: "1",
					Name:   "group 1",
					Email:  "user.1@mail.com",
				},
				Resources: []*Member{
					{IPID: "1", SCIMID: "1", Email: "user.1@mail.com", Status: "ACTIVE"},
					{IPID: "2", SCIMID: "2", Email: "user.2@mail.com", Status: "ACTIVE"},
				},
			},
		},
		{
			name: "partial filled",
			toTest: GroupMembers{
				Items:    1,
				HashCode: "test",
				Group: &Group{
					IPID:  "1",
					Name:  "group 1",
					Email: "user.1@mail.com",
				},
				Resources: []*Member{
					{IPID: "1", Email: "user.1@mail.com"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("GroupMembers.MarshalBinary() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got GroupMembers
			if err := dec.Decode(&got); err != nil {
				t.Errorf("GroupMembers.UnmarshalBinary() error = %v", err)
			}

			var expectedGroup *Group
			if tt.toTest.Group != nil {
				expectedGroup = &Group{
					IPID:  tt.toTest.Group.IPID,
					Name:  tt.toTest.Group.Name,
					Email: tt.toTest.Group.Email,
				}
			}

			var expectedMembers []*Member
			for _, m := range tt.toTest.Resources {
				expectedMembers = append(expectedMembers, &Member{
					IPID:   m.IPID,
					Email:  m.Email,
					Status: m.Status,
				})
			}

			// HashCode is not exported, so it will not be encoded
			expected := GroupMembers{
				Items:     tt.toTest.Items,
				Group:     expectedGroup,
				Resources: expectedMembers,
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(expected, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("GroupMembers.GobEncode() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestGroupMembers_SetHashCode(t *testing.T) {
	tests := []struct {
		name         string
		groupMembers GroupMembers
		want         GroupMembers
	}{
		{
			name: "success",
			groupMembers: GroupMembers{
				Items:    3,
				HashCode: "test",
				Group:    &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"},
				Resources: []*Member{
					{IPID: "1", SCIMID: "1", Email: "group.1@mail.com", Status: "ACTIVE"},
					{IPID: "2", SCIMID: "2", Email: "group.2@mail.com", Status: "ACTIVE"},
					{IPID: "3", SCIMID: "3", Email: "group.3@mail.com", Status: "ACTIVE"},
				},
			},
			want: GroupMembers{
				Items: 3,
				Group: &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"},
				Resources: []*Member{
					{IPID: "3", SCIMID: "3", Email: "group.3@mail.com", Status: "ACTIVE"},
					{IPID: "1", SCIMID: "1", Email: "group.1@mail.com", Status: "ACTIVE"},
					{IPID: "2", SCIMID: "2", Email: "group.2@mail.com", Status: "ACTIVE"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.groupMembers.SetHashCode()
			tt.want.SetHashCode()

			got := tt.groupMembers.HashCode

			if got != tt.want.HashCode {
				t.Errorf("GroupMembers.SetHashCode() = %s, want %s", got, tt.want.HashCode)
				t.Errorf("GroupMembers.SetHashCode() = %+v, want %+v", tt.groupMembers, tt.want)
			}
		})
	}
}

func TestGroupsMembersResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		HashCode  string
		Resources []*GroupMembers
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:   "empty",
			fields: fields{},
			want: []byte(`{
  "items": 0,
  "resources": []
}`),
			wantErr: false,
		},
		{
			name: "not empty",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []*GroupMembers{
					{
						Items: 1,
						Group: &Group{
							IPID:   "1",
							SCIMID: "1",
							Name:   "group 1",
							Email:  "group.1@mai.com",
						},
					},
				},
			},
			want: []byte(`{
  "items": 1,
  "hashCode": "test",
  "resources": [
    {
      "items": 1,
      "group": {
        "ipid": "1",
        "scimid": "1",
        "name": "group 1",
        "email": "group.1@mai.com"
      },
      "resources": null
    }
  ]
}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &GroupsMembersResult{
				Items:     tt.fields.Items,
				Resources: tt.fields.Resources,
				HashCode:  tt.fields.HashCode,
			}
			got, err := ur.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersResult.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("MarshalJSON() (-want +got):\n%s", diff)
			}
		})
	}
}
