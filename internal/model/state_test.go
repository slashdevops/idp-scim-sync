package model

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStateResources_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest *StateResources
	}{
		{
			name:   "empty",
			toTest: &StateResources{},
		},
		{
			name: "filled with Group",
			toTest: &StateResources{
				Groups: &GroupsResult{
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			enc := gob.NewEncoder(buf)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("User.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(buf)
			var got StateResources
			if err := dec.Decode(&got); err != nil {
				t.Errorf("User.GobEncode() error = %v", err)
			}

			var expectedGroupsResources *GroupsResult
			if tt.toTest.Groups != nil {
				expectedGroupsResources = &GroupsResult{
					Items:     tt.toTest.Groups.Items,
					Resources: make([]*Group, 0),
				}

				for _, g := range tt.toTest.Groups.Resources {
					expectedGroupsResources.Resources = append(expectedGroupsResources.Resources, &Group{
						IPID:  g.IPID,
						Name:  g.Name,
						Email: g.Email,
					})
				}
			}

			var expectedUsersResources *UsersResult
			if tt.toTest.Users != nil {
				expectedUsersResources = &UsersResult{
					Items:     tt.toTest.Users.Items,
					Resources: make([]*User, 0),
				}

				for _, u := range tt.toTest.Users.Resources {
					expectedUsersResources.Resources = append(expectedUsersResources.Resources, &User{
						IPID:        u.IPID,
						Name:        u.Name,
						DisplayName: u.DisplayName,
						Emails:      u.Emails,
						Active:      u.Active,
					})
				}
			}

			var expectedGroupsMembersResources *GroupsMembersResult
			if tt.toTest.GroupsMembers != nil {
				expectedGroupsMembersResources = &GroupsMembersResult{
					Items:     tt.toTest.GroupsMembers.Items,
					Resources: make([]*GroupMembers, 0),
				}

				for _, gm := range tt.toTest.GroupsMembers.Resources {
					expectedGroupsMembersResources.Resources = append(expectedGroupsMembersResources.Resources, &GroupMembers{
						Items: gm.Items,
						Group: &Group{
							IPID:  gm.Group.IPID,
							Name:  gm.Group.Name,
							Email: gm.Group.Email,
						},
						Resources: make([]*Member, 0),
					})

					for _, m := range gm.Resources {
						expectedGroupsMembersResources.Resources[len(expectedGroupsMembersResources.Resources)-1].Resources = append(expectedGroupsMembersResources.Resources[len(expectedGroupsMembersResources.Resources)-1].Resources, &Member{
							IPID:   m.IPID,
							Email:  m.Email,
							Status: m.Status,
						})
					}
				}
			}

			expected := StateResources{
				Groups:        expectedGroupsResources,
				Users:         expectedUsersResources,
				GroupsMembers: expectedGroupsMembersResources,
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(expected, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestState_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest *State
	}{
		{
			name:   "empty",
			toTest: &State{},
		},
		{
			name: "filled with Group",
			toTest: &State{
				LastSync: "2020-01-01T00:00:00Z",
				Resources: &StateResources{
					Groups: &GroupsResult{
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
				},
			},
		},
		{
			name: "filled with Group and Users",
			toTest: &State{
				LastSync: "2020-01-01T00:00:00Z",
				Resources: &StateResources{
					Groups: &GroupsResult{
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
					Users: &UsersResult{
						Items:    2,
						HashCode: "hashCode",
						Resources: []*User{
							{
								IPID:              "1",
								SCIMID:            "1",
								UserName:          "user.1",
								DisplayName:       "user 1",
								NickName:          "user 1",
								ProfileURL:        "https://profile.url",
								Title:             "title",
								UserType:          "user",
								Active:            true,
								Name:              &Name{FamilyName: "user", GivenName: "1", Formatted: "user 1"},
								Emails:            []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
								Addresses:         []Address{{Formatted: "address 1", StreetAddress: "street 1", Locality: "locality 1", Region: "region 1", PostalCode: "postalCode 1", Country: "country 1", Primary: true}},
								PreferredLanguage: "en",
								PhoneNumbers:      []PhoneNumber{{Value: "123456789", Type: "work"}},
								EnterpriseData:    &EnterpriseData{EmployeeNumber: "123456789", CostCenter: "123456789", Organization: "123456789", Division: "123456789", Department: "123456789"},
							},
							{
								IPID:              "2",
								SCIMID:            "2",
								UserName:          "user.2",
								DisplayName:       "user 2",
								NickName:          "user 2",
								ProfileURL:        "https://profile.url",
								Title:             "title",
								UserType:          "user",
								Active:            true,
								Name:              &Name{FamilyName: "user", GivenName: "2", Formatted: "user 2"},
								Emails:            []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}},
								Addresses:         []Address{{Formatted: "address 2", StreetAddress: "street 2", Locality: "locality 2", Region: "region 2", PostalCode: "postalCode 2", Country: "country 2", Primary: true}},
								PreferredLanguage: "en",
								PhoneNumbers:      []PhoneNumber{{Value: "123456789", Type: "work"}},
								EnterpriseData:    &EnterpriseData{EmployeeNumber: "123456789", CostCenter: "123456789", Organization: "123456789", Division: "123456789", Department: "123456789"},
							},
						},
					},
				},
			},
		},
		{
			name: "filled with Group, Users and GroupsMembers",
			toTest: &State{
				LastSync: "2020-01-01T00:00:00Z",
				Resources: &StateResources{
					Groups: &GroupsResult{
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
					Users: &UsersResult{
						Items:    2,
						HashCode: "hashCode",
						Resources: []*User{
							{
								IPID:        "1",
								SCIMID:      "1",
								UserName:    "user.1",
								DisplayName: "user 1",

								Active: true,
								Name:   &Name{FamilyName: "user", GivenName: "1", Formatted: "user 1"},
								Emails: []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
							},
							{
								IPID:              "2",
								SCIMID:            "2",
								UserName:          "user.2",
								DisplayName:       "user 2",
								ProfileURL:        "https://profile.url",
								Title:             "title",
								Active:            true,
								Name:              &Name{FamilyName: "user", GivenName: "2", Formatted: "user 2"},
								Emails:            []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}},
								PreferredLanguage: "en",
							},
						},
					},
					GroupsMembers: &GroupsMembersResult{
						Items: 1,
						Resources: []*GroupMembers{
							{
								Items: 2,
								Group: &Group{
									IPID:  "ipid",
									Name:  "name",
									Email: "email",
								},
								Resources: []*Member{
									{
										IPID:   "1",
										Email:  "user.1@mail.com",
										Status: "active",
									},
									{
										IPID:   "2",
										Email:  "user.2@mail.com",
										Status: "active",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			enc := gob.NewEncoder(buf)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("User.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(buf)
			var got State
			if err := dec.Decode(&got); err != nil {
				t.Errorf("User.GobEncode() error = %v", err)
			}

			var expectedStateResources *StateResources
			if tt.toTest.Resources != nil {

				// fill GroupsResult
				var expectedGroupsResources *GroupsResult
				if tt.toTest.Resources.Groups != nil {

					expectedGroupsResources = &GroupsResult{
						Items:     tt.toTest.Resources.Groups.Items,
						Resources: make([]*Group, 0),
					}

					if tt.toTest.Resources.Groups.Resources != nil {
						for _, g := range tt.toTest.Resources.Groups.Resources {
							expectedGroupsResources.Resources = append(expectedGroupsResources.Resources, &Group{
								IPID:  g.IPID,
								Name:  g.Name,
								Email: g.Email,
							})
						}
					}
				}
				// end fill GroupsResult

				// fill UsersResult
				var expectedUsersResources *UsersResult
				if tt.toTest.Resources.Users != nil {
					expectedUsersResources = &UsersResult{
						Items:     tt.toTest.Resources.Users.Items,
						Resources: make([]*User, 0),
					}

					if tt.toTest.Resources.Users.Resources != nil {
						for _, u := range tt.toTest.Resources.Users.Resources {
							expectedUsersResources.Resources = append(expectedUsersResources.Resources, &User{
								IPID:              u.IPID,
								UserName:          u.UserName,
								DisplayName:       u.DisplayName,
								NickName:          u.NickName,
								ProfileURL:        u.ProfileURL,
								Title:             u.Title,
								UserType:          u.UserType,
								PreferredLanguage: u.PreferredLanguage,
								Locale:            u.Locale,
								Timezone:          u.Timezone,
								Active:            u.Active,
								Emails:            u.Emails,
								Addresses:         u.Addresses,
								PhoneNumbers:      u.PhoneNumbers,
								Name:              u.Name,
								EnterpriseData:    u.EnterpriseData,
							})
						}
					}
				}
				// end fill UsersResult

				// fill GroupsMembersResult
				var expectedGroupsMembersResources *GroupsMembersResult
				if tt.toTest.Resources.GroupsMembers != nil {
					expectedGroupsMembersResources = &GroupsMembersResult{
						Items:     tt.toTest.Resources.GroupsMembers.Items,
						Resources: make([]*GroupMembers, 0),
					}

					if tt.toTest.Resources.GroupsMembers.Resources != nil {
						for _, gm := range tt.toTest.Resources.GroupsMembers.Resources {
							expectedMembers := make([]*Member, 0)
							for _, m := range gm.Resources {
								expectedMembers = append(expectedMembers, &Member{
									IPID:   m.IPID,
									Email:  m.Email,
									Status: m.Status,
								})
							}

							expectedGroupsMembersResources.Resources = append(expectedGroupsMembersResources.Resources, &GroupMembers{
								Items:     gm.Items,
								Group:     gm.Group,
								Resources: expectedMembers,
							})
						}
					}
				}
				// end fill GroupsMembersResult

				expectedStateResources = &StateResources{
					Groups:        expectedGroupsResources,
					Users:         expectedUsersResources,
					GroupsMembers: expectedGroupsMembersResources,
				}
			}

			expected := State{
				SchemaVersion: tt.toTest.SchemaVersion,
				Resources:     expectedStateResources,
			}

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(expected, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestState_MarshalJSON(t *testing.T) {
	type fields struct {
		LastSync  string
		Resources *StateResources
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
      "resources": []
    },
    "users": {
      "items": 0,
      "hashCode": "",
      "resources": []
    },
    "groupsMembers": {
      "items": 0,
      "resources": []
    }
  }
}`),
			wantErr: false,
		},
		{
			name: "empty GroupsMembersResult",
			fields: fields{
				LastSync: "2020-01-01T00:00:00Z",
				Resources: &StateResources{
					Groups: &GroupsResult{
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
					Users: &UsersResult{
						Items:    1,
						HashCode: "hashCode",
						Resources: []*User{
							{
								IPID:     "ipid",
								SCIMID:   "scimid",
								Name:     &Name{FamilyName: "lastName", GivenName: "name"},
								Emails:   []Email{{Value: "email", Type: "work", Primary: true}},
								HashCode: "hashCode",
							},
						},
					},
					GroupsMembers: &GroupsMembersResult{},
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
          "hashCode": "hashCode",
          "emails": [
            {
              "value": "email",
              "type": "work",
              "primary": true
            }
          ],
          "name": {
            "familyName": "lastName",
            "givenName": "name"
          }
        }
      ]
    },
    "groupsMembers": {
      "items": 0,
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

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("mismatch (-tt.want +got):\n%s", diff)
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
		u1 := User{IPID: "1", SCIMID: "1", Name: &Name{FamilyName: "user", GivenName: "1"}, DisplayName: "user.1", Active: true, Emails: []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}}
		u1.SetHashCode()

		u2 := User{IPID: "2", SCIMID: "2", Name: &Name{FamilyName: "user", GivenName: "2"}, DisplayName: "user.2", Active: true, Emails: []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}}
		u2.SetHashCode()

		u3 := User{IPID: "3", SCIMID: "3", Name: &Name{FamilyName: "user", GivenName: "3"}, DisplayName: "user.3", Active: false, Emails: []Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}}
		u3.SetHashCode()

		usrs := &UsersResult{Items: 3, Resources: []*User{&u1, &u2, &u3}}
		usrs.SetHashCode()

		g1 := &Group{IPID: "1", SCIMID: "1", Name: "group 1", Email: "group.1@mail.com"}
		g1.SetHashCode()

		g2 := &Group{IPID: "2", SCIMID: "2", Name: "group 2", Email: "group.2@mail.com"}
		g2.SetHashCode()

		g3 := &Group{IPID: "3", SCIMID: "3", Name: "group 3", Email: "group.3@mail.com"}
		g3.SetHashCode()

		grs := &GroupsResult{Items: 3, Resources: []*Group{g1, g2, g3}}
		grs.SetHashCode()

		m1 := Member{IPID: "1", SCIMID: "1", Email: u1.GetPrimaryEmailAddress(), Status: "ACTIVE"}
		m1.SetHashCode()

		m2 := Member{IPID: "2", SCIMID: "2", Email: u2.GetPrimaryEmailAddress(), Status: "ACTIVE"}
		m2.SetHashCode()

		m3 := Member{IPID: "2", SCIMID: "2", Email: u3.GetPrimaryEmailAddress(), Status: "ACTIVE"}
		m3.SetHashCode()

		gm1 := GroupMembers{Items: 1, Group: g1, Resources: []*Member{&m1}}
		gm1.SetHashCode()

		gm2 := GroupMembers{Items: 1, Group: g2, Resources: []*Member{&m2}}
		gm2.SetHashCode()

		gm3 := GroupMembers{Items: 1, Group: g3, Resources: []*Member{&m3}}
		gm3.SetHashCode()

		gmrs := &GroupsMembersResult{Items: 3, Resources: []*GroupMembers{&gm1, &gm2, &gm3}}
		gmrs.SetHashCode()

		sr := &StateResources{
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
