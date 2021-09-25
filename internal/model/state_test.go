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
			name:    "empty",
			fields:  fields{},
			want:    []byte(`{"lastSync":"","resources":{"groups":{"items":0,"hashCode":"","resources":[]},"users":{"items":0,"hashCode":"","resources":[]},"groupsUsers":{"items":0,"hashCode":"","resources":[]}}}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				LastSync: "2020-01-01T00:00:00Z",
				Resources: StateResources{
					Groups: &GroupsResult{
						Items:    1,
						HashCode: "hashCode",
						Resources: []Group{
							{
								ID:       "id",
								Name:     "name",
								Email:    "email",
								HashCode: "hashCode",
							},
						},
					},
					Users: &UsersResult{
						Items:    1,
						HashCode: "hashCode",
						Resources: []User{
							{
								ID:       "id",
								Name:     Name{FamilyName: "lastName", GivenName: "name"},
								Email:    "email",
								HashCode: "hashCode",
							},
						},
					},
					GroupsUsers: &GroupsUsersResult{
						Items:    1,
						HashCode: "hashCode",
						Resources: []GroupUsers{
							{
								Items:    1,
								HashCode: "hashCode",
								Group: Group{
									ID:       "id",
									Name:     "name",
									Email:    "email",
									HashCode: "hashCode",
								},
								Resources: []User{
									{
										ID:       "id",
										Name:     Name{FamilyName: "lastName", GivenName: "name"},
										Email:    "email",
										HashCode: "hashCode",
									},
								},
							},
						},
					},
				},
			},
			want:    []byte(`{"lastSync":"2020-01-01T00:00:00Z","resources":{"groups":{"items":1,"hashCode":"hashCode","resources":[{"id":"id","name":"name","email":"email","hashCode":"hashCode"}]},"users":{"items":1,"hashCode":"hashCode","resources":[{"id":"id","name":{"familyName":"lastName","givenName":"name"},"displayName":"","active":false,"email":"email","hashCode":"hashCode"}]},"groupsUsers":{"items":1,"hashCode":"hashCode","resources":[{"items":1,"group":{"id":"id","name":"name","email":"email","hashCode":"hashCode"},"hashCode":"hashCode","resources":[{"id":"id","name":{"familyName":"lastName","givenName":"name"},"displayName":"","active":false,"email":"email","hashCode":"hashCode"}]}]}}}`),
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
