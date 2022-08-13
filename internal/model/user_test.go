package model

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

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
