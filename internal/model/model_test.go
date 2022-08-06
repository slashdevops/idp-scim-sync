package model

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func TestUsersResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		Resources []*User
		HashCode  string
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
  "hashCode": "",
  "resources": []
}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []*User{
					{
						IPID:   "1",
						SCIMID: "1",
						Name: Name{
							GivenName:  "user",
							FamilyName: "1",
						},
						DisplayName: "user 1",
						Active:      true,
						Email:       "user.1@mail.com",
						HashCode:    "1111",
					},
				},
			},
			want: []byte(`{
  "items": 1,
  "hashCode": "test",
  "resources": [
    {
      "ipid": "1",
      "scimid": "1",
      "name": {
        "familyName": "1",
        "givenName": "user"
      },
      "displayName": "user 1",
      "active": true,
      "email": "user.1@mail.com",
      "hashCode": "1111"
    }
  ]
}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &UsersResult{
				Items:     tt.fields.Items,
				Resources: tt.fields.Resources,
				HashCode:  tt.fields.HashCode,
			}
			got, err := ur.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersResult.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UsersResult.MarshalJSON() = %s, want %s", string(got), string(tt.want))
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
  "hashCode": "",
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
        "email": "group.1@mai.com",
        "hashCode": ""
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UsersResult.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestGroupsResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		HashCode  string
		Resources []*Group
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
  "hashCode": "",
  "resources": []
}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []*Group{
					{
						IPID:     "1",
						SCIMID:   "1",
						Name:     "group",
						HashCode: "1111",
					},
				},
			},
			want: []byte(`{
  "items": 1,
  "hashCode": "test",
  "resources": [
    {
      "ipid": "1",
      "scimid": "1",
      "name": "group",
      "email": "",
      "hashCode": "1111"
    }
  ]
}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GroupsResult{
				Items:     tt.fields.Items,
				HashCode:  tt.fields.HashCode,
				Resources: tt.fields.Resources,
			}
			got, err := gr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsResult.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupsResult.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestUser_SetHashCode(t *testing.T) {
	tests := []struct {
		name string
		user User
		want User
	}{
		{
			name: "success with SCIM field and hashcode",
			user: User{
				IPID:   "1",
				SCIMID: "1",
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
				HashCode:    "test",
			},
			want: User{
				IPID: "1",
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
			},
		},
		{
			name: "success with default field values",
			user: User{
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
				HashCode:    "test",
			},
			want: User{
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
			},
		},
		{
			name: "success empty",
			user: User{},
			want: User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.SetHashCode()
			tt.want.SetHashCode()
			got := tt.user.HashCode
			if got != tt.want.HashCode {
				t.Errorf("User.SetHashCode() = %s, want %s", got, tt.want.HashCode)
			}
		})
	}
}

func TestUser_SetHashCode_pointer(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want *User
	}{
		{
			name: "success",
			user: &User{
				IPID:   "1",
				SCIMID: "1",
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
				HashCode:    "test",
			},
			want: &User{
				IPID: "1",
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.SetHashCode()
			tt.want.SetHashCode()
			got := tt.user.HashCode
			if got != tt.want.HashCode {
				t.Errorf("User.SetHashCode() = %s, want %s", got, tt.want.HashCode)
			}
		})
	}
}

func TestGroup_SetHashCode(t *testing.T) {
	tests := []struct {
		name  string
		group *Group
		want  *Group
	}{
		{
			name: "success",
			group: &Group{
				IPID:     "1",
				SCIMID:   "1",
				Name:     "group 1",
				Email:    "user.1@mail.com",
				HashCode: "test",
			},
			want: &Group{
				IPID:  "1",
				Name:  "group 1",
				Email: "user.1@mail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.group.SetHashCode()
			tt.want.SetHashCode()
			got := tt.group.HashCode
			if got != tt.want.HashCode {
				t.Errorf("Group.SetHashCode() = %s, want %s", got, tt.want.HashCode)
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
				IPID:  "1",
				Email: "user.1@mail.com",
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

func TestUsersResult_SetHashCode(t *testing.T) {
	u2 := &User{IPID: "2", SCIMID: "2", Name: Name{GivenName: "User", FamilyName: "2"}, Email: "user.2@mail.com"}
	u1 := &User{IPID: "1", SCIMID: "1", Name: Name{GivenName: "User", FamilyName: "1"}, Email: "user.1@mail.com"}
	u3 := &User{IPID: "3", SCIMID: "3", Name: Name{GivenName: "User", FamilyName: "3"}, Email: "user.3@mail.com"}

	u1.SetHashCode()
	u2.SetHashCode()
	u3.SetHashCode()

	ur1 := UsersResult{
		Items:     3,
		Resources: []*User{u1, u2, u3},
	}
	ur1.SetHashCode()

	ur2 := UsersResult{
		Items:     3,
		Resources: []*User{u2, u3, u1},
	}
	ur2.SetHashCode()

	ur3 := UsersResult{
		Items:     3,
		Resources: []*User{u3, u2, u1},
	}
	ur3.SetHashCode()

	ur4 := MergeUsersResult(&ur2, &ur1, &ur3)
	ur4.SetHashCode()
	ur5 := MergeUsersResult(&ur3, &ur2, &ur1)
	ur5.SetHashCode()

	t.Logf("ur4: %s\n", utils.ToJSON(ur4))
	t.Logf("ur5: %s\n", utils.ToJSON(ur5))

	t.Logf("ur4.HashCode: %s\n", ur4.HashCode)
	t.Logf("ur5.HashCode: %s\n", ur5.HashCode)

	if ur1.HashCode != ur2.HashCode {
		t.Errorf("UsersResult.HashCode should be equal")
	}
	if ur1.HashCode != ur3.HashCode {
		t.Errorf("UsersResult.HashCode should be equal")
	}
	if ur2.HashCode != ur3.HashCode {
		t.Errorf("UsersResult.HashCode should be equal")
	}

	if ur5.HashCode != ur4.HashCode {
		t.Errorf("UsersResult.HashCode should be equal: ur5-> %s, ur4-> %s", ur5.HashCode, ur4.HashCode)
	}
}

func TestGroupsResult_SetHashCode(t *testing.T) {
	g1 := &Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"}
	g2 := &Group{IPID: "2", SCIMID: "2", Name: "group", Email: "group.2@mail.com"}
	g3 := &Group{IPID: "3", SCIMID: "3", Name: "group", Email: "group.3@mail.com"}

	g1.SetHashCode()
	g2.SetHashCode()
	g3.SetHashCode()

	gr1 := GroupsResult{
		Items:     3,
		Resources: []*Group{g1, g2, g3},
	}
	gr1.SetHashCode()

	gr2 := GroupsResult{
		Items:     3,
		Resources: []*Group{g2, g3, g1},
	}
	gr2.SetHashCode()

	gr3 := GroupsResult{
		Items:     3,
		Resources: []*Group{g3, g2, g1},
	}
	gr3.SetHashCode()

	gr4 := MergeGroupsResult(&gr2, &gr1, &gr3)
	gr4.SetHashCode()
	gr5 := MergeGroupsResult(&gr3, &gr2, &gr1)
	gr5.SetHashCode()

	t.Logf("gr4: %s\n", utils.ToJSON(gr4))
	t.Logf("gr5: %s\n", utils.ToJSON(gr5))

	t.Logf("gr4.HashCode: %s\n", gr4.HashCode)
	t.Logf("gr5.HashCode: %s\n", gr5.HashCode)

	if gr1.HashCode != gr2.HashCode {
		t.Errorf("GroupsResult.HashCode should be equal")
	}
	if gr1.HashCode != gr3.HashCode {
		t.Errorf("GroupsResult.HashCode should be equal")
	}
	if gr2.HashCode != gr3.HashCode {
		t.Errorf("GroupsResult.HashCode should be equal")
	}

	if gr5.HashCode != gr4.HashCode {
		t.Errorf("GroupsResult.HashCode should be equal: gr5-> %s, gr4-> %s", gr5.HashCode, gr4.HashCode)
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

	t.Logf("mr1: %s\n", utils.ToJSON(mr1))
	t.Logf("mr2: %s\n", utils.ToJSON(mr2))
	t.Logf("mr3: %s\n", utils.ToJSON(mr3))

	t.Logf("mr1.HashCode: %s\n", mr1.HashCode)
	t.Logf("mr2.HashCode: %s\n", mr2.HashCode)
	t.Logf("mr3.HashCode: %s\n", mr3.HashCode)

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

	t.Logf("gmr4: %s\n", utils.ToJSON(gmr4))
	t.Logf("gmr5: %s\n", utils.ToJSON(gmr5))

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

func TestGroup_GobEncode(t *testing.T) {
	tests := []struct {
		name    string
		g       *Group
		wantErr bool
	}{
		{
			name: "Test Group GobEncode",
			g: &Group{
				IPID:     "1",
				SCIMID:   "1",
				Name:     "group",
				Email:    "user.1@mail.com",
				HashCode: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GobEncode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.GobEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)
			if err := enc.Encode(tt.g.IPID); err != nil {
				t.Fatal(err)
			}
			if err := enc.Encode(tt.g.Name); err != nil {
				t.Fatal(err)
			}
			if err := enc.Encode(tt.g.Email); err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, b.Bytes()) {
				t.Errorf("Group.GobEncode() = %v\n, want %v\n", got, b.Bytes())
			}
		})
	}
}

func TestMember_GobEncode(t *testing.T) {
	tests := []struct {
		name    string
		m       Member
		wantErr bool
	}{
		{
			name: "Test Member GobEncode",
			m: Member{
				IPID:     "1",
				SCIMID:   "1",
				Email:    "member.1@mail.com",
				HashCode: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GobEncode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Member.GobEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)
			if err := enc.Encode(tt.m.IPID); err != nil {
				t.Fatal(err)
			}

			if err := enc.Encode(tt.m.Email); err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, b.Bytes()) {
				t.Errorf("Group.GobEncode() = %v\n, want %v\n", got, b.Bytes())
			}
		})
	}
}

func TestUser_GobEncode(t *testing.T) {
	tests := []struct {
		name    string
		u       User
		wantErr bool
	}{
		{
			name: "Test User GobEncode",
			u: User{
				IPID:     "1",
				SCIMID:   "1",
				Name:     Name{FamilyName: "user", GivenName: "1"},
				Email:    "user.1@mail.com",
				HashCode: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.GobEncode()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.GobEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)
			if err := enc.Encode(tt.u.IPID); err != nil {
				panic(err)
			}
			if err := enc.Encode(tt.u.Name); err != nil {
				panic(err)
			}
			if err := enc.Encode(tt.u.DisplayName); err != nil {
				panic(err)
			}
			if err := enc.Encode(tt.u.Active); err != nil {
				panic(err)
			}
			if err := enc.Encode(tt.u.Email); err != nil {
				panic(err)
			}
			if !bytes.Equal(got, b.Bytes()) {
				t.Errorf("Group.GobEncode() = %v\n, want %v\n", got, b.Bytes())
			}
		})
	}
}
