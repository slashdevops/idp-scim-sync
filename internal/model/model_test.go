package model

import (
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func TestUsersResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		Resources []User
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
				Resources: []User{
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

func TestGroupsResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		HashCode  string
		Resources []Group
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
				Resources: []Group{
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
		want string
	}{
		{
			name: "success",
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
			want: hash.Get(User{
				IPID: "1",
				Name: Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Email:       "user.1@mail.com",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.SetHashCode()
			got := tt.user.HashCode
			if got != tt.want {
				t.Errorf("User.SetHashCode() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestGroupsResult_SetHashCode(t *testing.T) {
	g1 := Group{IPID: "1", SCIMID: "1", Name: "group", Email: "group.1@mail.com"}
	g2 := Group{IPID: "2", SCIMID: "2", Name: "group", Email: "group.2@mail.com"}
	g3 := Group{IPID: "3", SCIMID: "3", Name: "group", Email: "group.3@mail.com"}

	g1.SetHashCode()
	g2.SetHashCode()
	g3.SetHashCode()

	gr1 := GroupsResult{
		Items:     3,
		Resources: []Group{g1, g2, g3},
	}
	gr1.SetHashCode()

	gr2 := GroupsResult{
		Items:     3,
		Resources: []Group{g2, g3, g1},
	}
	gr2.SetHashCode()

	gr3 := GroupsResult{
		Items:     3,
		Resources: []Group{g3, g2, g1},
	}
	gr3.SetHashCode()

	gr4 := mergeGroupsResult(&gr2, &gr1, &gr3)
	gr4.SetHashCode()
	gr5 := mergeGroupsResult(&gr3, &gr2, &gr1)
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

func mergeGroupsResult(grs ...*GroupsResult) (merged GroupsResult) {
	groups := make([]Group, 0)

	for _, gr := range grs {
		groups = append(groups, gr.Resources...)
	}

	merged = GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	if merged.Items > 0 {
		merged.SetHashCode()
	}

	return
}
