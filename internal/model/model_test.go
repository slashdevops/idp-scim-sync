package model

import (
	"reflect"
	"testing"
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
			name:    "empty",
			fields:  fields{},
			want:    []byte(`{"items":0,"hashCode":"","resources":[]}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []User{
					{
						ID: "1",
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
			want:    []byte(`{"items":1,"hashCode":"test","resources":[{"id":"1","name":{"familyName":"1","givenName":"user"},"displayName":"user 1","active":true,"email":"user.1@mail.com","hashCode":"1111"}]}`),
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
			name:    "empty",
			fields:  fields{},
			want:    []byte(`{"items":0,"hashCode":"","resources":[]}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []Group{
					{
						ID:       "1",
						Name:     "group",
						HashCode: "1111",
					},
				},
			},
			want:    []byte(`{"items":1,"hashCode":"test","resources":[{"id":"1","name":"group","email":"","hashCode":"1111"}]}`),
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

func TestGroupsUsersResult_MarshalJSON(t *testing.T) {
	type fields struct {
		Items     int
		HashCode  string
		Resources []GroupUsers
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
			want:    []byte(`{"items":0,"hashCode":"","resources":[]}`),
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Items:    1,
				HashCode: "test",
				Resources: []GroupUsers{
					{
						Items: 1,
						Group: Group{
							ID:       "1",
							Name:     "group 1",
							Email:    "group.1@mail.com",
							HashCode: "test",
						},
						HashCode: "test",
						Resources: []User{
							{
								ID: "1",
								Name: Name{
									GivenName:  "user",
									FamilyName: "1",
								},
								DisplayName: "user 1",
								Active:      true,
								Email:       "user.1@mail.com",
								HashCode:    "test",
							},
						},
					},
				},
			},
			want: []byte(`{"items":1,"hashCode":"test","resources":[{"items":1,"group":{"id":"1","name":"group 1","email":"group.1@mail.com","hashCode":"test"},"hashCode":"test","resources":[{"id":"1","name":{"familyName":"1","givenName":"user"},"displayName":"user 1","active":true,"email":"user.1@mail.com","hashCode":"test"}]}]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gur := &GroupsUsersResult{
				Items:     tt.fields.Items,
				HashCode:  tt.fields.HashCode,
				Resources: tt.fields.Resources,
			}
			got, err := gur.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsUsersResult.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupsUsersResult.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}
