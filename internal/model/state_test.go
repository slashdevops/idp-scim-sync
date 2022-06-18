package model

import (
	"reflect"
	"testing"
)

func TestState_MarshalJSON(t *testing.T) {
	type fields struct {
		LastSync  string
		Resources StateResources
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
  "schemaVersion": "",
  "codeVersion": "",
  "lastSync": "",
  "hashCode": "",
  "resources": {
    "groups": {
      "items": 0,
      "hashCode": "",
      "resources": []
    },
    "users": {
      "items": 0,
      "hashCode": "",
      "resources": []
    },
    "groupsMembers": {
      "items": 0,
      "hashCode": "",
      "resources": []
    }
  }
}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				LastSync: "2020-01-01T00:00:00Z",
				Resources: StateResources{
					Groups: GroupsResult{
						Items:    1,
						HashCode: "hashCode",
						Resources: []*Group{
							{
								IPID:     "ipid",
								SCIMID:   "scimid",
								Name:     "name",
								Email:    "email",
								HashCode: "hashCode",
							},
						},
					},
					Users: UsersResult{
						Items:    1,
						HashCode: "hashCode",
						Resources: []*User{
							{
								IPID:     "ipid",
								SCIMID:   "scimid",
								Name:     Name{FamilyName: "lastName", GivenName: "name"},
								Email:    "email",
								HashCode: "hashCode",
							},
						},
					},
				},
			},
			want: []byte(`{
  "schemaVersion": "",
  "codeVersion": "",
  "lastSync": "2020-01-01T00:00:00Z",
  "hashCode": "",
  "resources": {
    "groups": {
      "items": 1,
      "hashCode": "hashCode",
      "resources": [
        {
          "ipid": "ipid",
          "scimid": "scimid",
          "name": "name",
          "email": "email",
          "hashCode": "hashCode"
        }
      ]
    },
    "users": {
      "items": 1,
      "hashCode": "hashCode",
      "resources": [
        {
          "ipid": "ipid",
          "scimid": "scimid",
          "name": {
            "familyName": "lastName",
            "givenName": "name"
          },
          "displayName": "",
          "active": false,
          "email": "email",
          "hashCode": "hashCode"
        }
      ]
    },
    "groupsMembers": {
      "items": 0,
      "hashCode": "",
      "resources": []
    }
  }
}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				LastSync:  tt.fields.LastSync,
				Resources: tt.fields.Resources,
			}
			got, err := s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("State.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestState_SetHashCode(t *testing.T) {
	t.Run("empty state", func(t *testing.T) {
		st := &State{}
		st.SetHashCode()

		if st.HashCode == "" {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.HashCode, "")
		}
		if st.Resources.Users.HashCode != "" {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Users.HashCode, "")
		}
		if st.Resources.Groups.HashCode != "" {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Groups.HashCode, "")
		}
		if st.Resources.GroupsMembers.HashCode != "" {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.GroupsMembers.HashCode, "")
		}
	})

	t.Run("full state", func(t *testing.T) {
		u1 := User{IPID: "1", SCIMID: "1", Name: Name{FamilyName: "user", GivenName: "1"}, DisplayName: "user.1", Active: true, Email: "user.1@mail.com"}
		u1.SetHashCode()

		u2 := User{IPID: "2", SCIMID: "2", Name: Name{FamilyName: "user", GivenName: "2"}, DisplayName: "user.2", Active: true, Email: "user.2@mail.com"}
		u2.SetHashCode()

		u3 := User{IPID: "3", SCIMID: "3", Name: Name{FamilyName: "user", GivenName: "3"}, DisplayName: "user.3", Active: false, Email: "user.3@mail.com"}
		u3.SetHashCode()

		usrs := UsersResult{Items: 3, Resources: []*User{&u1, &u2, &u3}}
		usrs.SetHashCode()

		g1 := Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"}
		g1.SetHashCode()

		g2 := Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"}
		g2.SetHashCode()

		g3 := Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"}
		g3.SetHashCode()

		grs := GroupsResult{Items: 3, Resources: []*Group{&g1, &g2, &g3}}
		grs.SetHashCode()

		m1 := Member{IPID: "1", SCIMID: "1", Email: u1.Email, Status: "ACTIVE"}
		m1.SetHashCode()

		m2 := Member{IPID: "2", SCIMID: "2", Email: u2.Email, Status: "ACTIVE"}
		m2.SetHashCode()

		m3 := Member{IPID: "2", SCIMID: "2", Email: u3.Email, Status: "ACTIVE"}
		m3.SetHashCode()

		gm1 := GroupMembers{Items: 1, Group: g1, Resources: []*Member{&m1}}
		gm1.SetHashCode()

		gm2 := GroupMembers{Items: 1, Group: g2, Resources: []*Member{&m2}}
		gm2.SetHashCode()

		gm3 := GroupMembers{Items: 1, Group: g3, Resources: []*Member{&m3}}
		gm3.SetHashCode()

		gmrs := GroupsMembersResult{Items: 3, Resources: []*GroupMembers{&gm1, &gm2, &gm3}}
		gmrs.SetHashCode()

		sr := StateResources{
			Users:         usrs,
			Groups:        grs,
			GroupsMembers: gmrs,
		}

		st := State{
			SchemaVersion: "1.0.0",
			CodeVersion:   "1.0.0",
			LastSync:      "",
			Resources:     sr,
		}
		st.SetHashCode()

		if st.HashCode == "" {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.HashCode, "")
		}
		if st.Resources.Users.HashCode != usrs.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Users.HashCode, usrs.HashCode)
		}
		if st.Resources.Groups.HashCode != grs.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Groups.HashCode, grs.HashCode)
		}
		if st.Resources.GroupsMembers.HashCode != gmrs.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.GroupsMembers.HashCode, gmrs.HashCode)
		}

		if st.Resources.Users.Resources[0].HashCode != u1.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Users.Resources[0].HashCode, u1.HashCode)
		}

		if st.Resources.Users.Resources[1].HashCode != u2.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Users.Resources[1].HashCode, u2.HashCode)
		}

		if st.Resources.Users.Resources[2].HashCode != u3.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Users.Resources[2].HashCode, u3.HashCode)
		}

		if st.Resources.Groups.Resources[0].HashCode != g1.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Groups.Resources[0].HashCode, g1.HashCode)
		}

		if st.Resources.Groups.Resources[1].HashCode != g2.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Groups.Resources[1].HashCode, g2.HashCode)
		}

		if st.Resources.Groups.Resources[2].HashCode != g3.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.Groups.Resources[2].HashCode, g3.HashCode)
		}

		if st.Resources.GroupsMembers.Resources[0].HashCode != gm1.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.GroupsMembers.Resources[0].HashCode, gm1.HashCode)
		}

		if st.Resources.GroupsMembers.Resources[1].HashCode != gm2.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.GroupsMembers.Resources[1].HashCode, gm2.HashCode)
		}

		if st.Resources.GroupsMembers.Resources[2].HashCode != gm3.HashCode {
			t.Errorf("State.SetHashCode() error = %v, wantErr %v", st.Resources.GroupsMembers.Resources[2].HashCode, gm3.HashCode)
		}
	})
}
