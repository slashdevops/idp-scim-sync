package model

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

func TestName_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest Name
	}{
		{
			name:   "empty Name",
			toTest: Name{},
		},
		{
			name: "filled Name",
			toTest: Name{
				Formatted:       "formatted",
				FamilyName:      "familyName",
				GivenName:       "givenName",
				MiddleName:      "middleName",
				HonorificPrefix: "honorificPrefix",
				HonorificSuffix: "honorificSuffix",
			},
		},
		{
			name: "filled Name with empty values",
			toTest: Name{
				FamilyName:      "familyName",
				GivenName:       "givenName",
				HonorificSuffix: "honorificSuffix",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("Name.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got Name
			if err := dec.Decode(&got); err != nil {
				t.Errorf("Name.GobEncode() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.toTest) {
				t.Errorf("Name.GobEncode() = %v, want %v", got, tt.toTest)
			}
		})
	}
}

func TestEmail_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest Email
	}{
		{
			name:   "empty Email",
			toTest: Email{},
		},
		{
			name: "filled Email",
			toTest: Email{
				Value:   "value",
				Type:    "type",
				Primary: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("Email.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got Email
			if err := dec.Decode(&got); err != nil {
				t.Errorf("Email.GobEncode() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.toTest) {
				t.Errorf("Email.GobEncode() = %v, want %v", got, tt.toTest)
			}
		})
	}
}

func TestAddress_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest Address
	}{
		{
			name:   "empty Address",
			toTest: Address{},
		},
		{
			name: "filled Address",
			toTest: Address{
				Formatted:     "formatted",
				StreetAddress: "streetAddress",
				Locality:      "locality",
				Region:        "region",
				PostalCode:    "postalCode",
				Country:       "country",
				Type:          "type",
				Primary:       true,
			},
		},
		{
			name: "filled Address with empty values",
			toTest: Address{
				StreetAddress: "streetAddress",
				Locality:      "locality",
				Region:        "region",
				PostalCode:    "postalCode",
				Country:       "country",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("Address.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got Address
			if err := dec.Decode(&got); err != nil {
				t.Errorf("Address.GobEncode() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.toTest) {
				t.Errorf("Address.GobEncode() = %v, want %v", got, tt.toTest)
			}
		})
	}
}

func TestPhoneNumber_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest PhoneNumber
	}{
		{
			name:   "empty PhoneNumber",
			toTest: PhoneNumber{},
		},
		{
			name: "filled PhoneNumber",
			toTest: PhoneNumber{
				Value: "value",
				Type:  "type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("PhoneNumber.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got PhoneNumber
			if err := dec.Decode(&got); err != nil {
				t.Errorf("PhoneNumber.GobEncode() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.toTest) {
				t.Errorf("PhoneNumber.GobEncode() = %v, want %v", got, tt.toTest)
			}
		})
	}
}

func TestManager_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest Manager
	}{
		{
			name:   "empty Manager",
			toTest: Manager{},
		},
		{
			name: "filled Manager",
			toTest: Manager{
				Value: "value",
				Ref:   "ref",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("Manager.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got Manager
			if err := dec.Decode(&got); err != nil {
				t.Errorf("Manager.GobEncode() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.toTest) {
				t.Errorf("Manager.GobEncode() = %v, want %v", got, tt.toTest)
			}
		})
	}
}

func TestEnterpriseData_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest EnterpriseData
	}{
		{
			name:   "empty EnterpriseData",
			toTest: EnterpriseData{},
		},
		{
			name: "filled EnterpriseData",
			toTest: EnterpriseData{
				EmployeeNumber: "employeeNumber",
				CostCenter:     "costCenter",
				Organization:   "organization",
				Manager:        Manager{Value: "123456789", Ref: "https://idp.example.com/idp/user/123456789"},
				Department:     "department",
				Division:       "division",
			},
		},
		{
			name: "filled EnterpriseData with empty values",
			toTest: EnterpriseData{
				EmployeeNumber: "employeeNumber",
				CostCenter:     "costCenter",
				Department:     "department",
				Division:       "division",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("EnterpriseData.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got EnterpriseData
			if err := dec.Decode(&got); err != nil {
				t.Errorf("EnterpriseData.GobEncode() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.toTest) {
				t.Errorf("EnterpriseData.GobEncode() = %v, want %v", got, tt.toTest)
			}
		})
	}
}

func TestUser_GobEncode(t *testing.T) {
	tests := []struct {
		name   string
		toTest User
	}{
		{
			name:   "empty User",
			toTest: User{},
		},
		{
			name: "filled User",
			toTest: User{
				IPID:              "IPID",
				SCIMID:            "this should not be encoded",
				UserName:          "username",
				DisplayName:       "displayname",
				NickName:          "nickname",
				ProfileURL:        "profileURL",
				Title:             "title",
				UserType:          "userType",
				PreferredLanguage: "preferredLanguage",
				Locale:            "locale",
				Timezone:          "timezone",
				Emails:            []Email{{Value: "email value", Type: "email type", Primary: true}},
				Addresses:         []Address{{Formatted: "formatted", StreetAddress: "street address", Locality: "locality", Region: "region", PostalCode: "postal code", Country: "country", Type: "type", Primary: true}},
				PhoneNumbers:      []PhoneNumber{{Value: "phone value", Type: "phone type"}},
				Name: &Name{
					Formatted:       "formatted",
					FamilyName:      "familyName",
					GivenName:       "givenName",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
				EnterpriseData: &EnterpriseData{
					EmployeeNumber: "employeeNumber",
					CostCenter:     "costCenter",
					Organization:   "organization",
					Manager:        Manager{Value: "manager value", Ref: "manager ref"},
					Department:     "department",
					Division:       "division",
				},
				Active:   true,
				HashCode: "this should not be encoded",
			},
		},
		{
			name: "filled User with empty values",
			toTest: User{
				IPID:        "1",
				UserName:    "user1",
				DisplayName: "user 1",
				ProfileURL:  "https://idp.example.com/idp/user/1",
				Name: &Name{
					FamilyName: "1",
					GivenName:  "user",
				},
				Active:       true,
				Emails:       []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
				PhoneNumbers: []PhoneNumber{{Value: "value", Type: "type"}},
				HashCode:     "this should not be encoded",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)

			if err := enc.Encode(tt.toTest); err != nil {
				t.Errorf("User.GobEncode() error = %v", err)
			}

			dec := gob.NewDecoder(b)
			var got User
			if err := dec.Decode(&got); err != nil {
				t.Errorf("User.GobEncode() error = %v", err)
			}

			// SCIMID is not exported, so it will not be encoded
			// HashCode is not exported, so it will not be encoded
			expected := User{
				IPID:              tt.toTest.IPID,
				UserName:          tt.toTest.UserName,
				DisplayName:       tt.toTest.DisplayName,
				NickName:          tt.toTest.NickName,
				ProfileURL:        tt.toTest.ProfileURL,
				Title:             tt.toTest.Title,
				UserType:          tt.toTest.UserType,
				PreferredLanguage: tt.toTest.PreferredLanguage,
				Locale:            tt.toTest.Locale,
				Timezone:          tt.toTest.Timezone,
				Emails:            tt.toTest.Emails,
				Addresses:         tt.toTest.Addresses,
				PhoneNumbers:      tt.toTest.PhoneNumbers,
				Name:              tt.toTest.Name,
				EnterpriseData:    tt.toTest.EnterpriseData,
				Active:            tt.toTest.Active,
			}

			if !reflect.DeepEqual(got, expected) {
				t.Errorf("User.GobEncode() = %v, want %v", got, expected)
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
				Name: &Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
				HashCode:    "test",
			},
			want: User{
				IPID: "1",
				Name: &Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
			},
		},
		{
			name: "success with default field values",
			user: User{
				Name: &Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
				HashCode:    "test",
			},
			want: User{
				Name: &Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
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
				Name: &Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
				HashCode:    "test",
			},
			want: &User{
				IPID: "1",
				Name: &Name{
					GivenName:  "user",
					FamilyName: "1",
				},
				DisplayName: "user 1",
				Active:      true,
				Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
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
						Name: &Name{
							GivenName:  "user",
							FamilyName: "1",
						},
						DisplayName: "user 1",
						Active:      true,
						Emails:      []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
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
      "displayName": "user 1",
      "hashCode": "1111",
      "emails": [
        {
          "value": "user.1@mail.com",
          "type": "work",
          "primary": true
        }
      ],
      "name": {
        "familyName": "1",
        "givenName": "user"
      },
      "active": true
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
	u2 := &User{IPID: "2", SCIMID: "2", Name: &Name{GivenName: "User", FamilyName: "2"}, Emails: []Email{{Value: "user.2@mail.com", Type: "work", Primary: true}}}
	u1 := &User{IPID: "1", SCIMID: "1", Name: &Name{GivenName: "User", FamilyName: "1"}, Emails: []Email{{Value: "user.1@mail.com", Type: "work", Primary: true}}}
	u3 := &User{IPID: "3", SCIMID: "3", Name: &Name{GivenName: "User", FamilyName: "3"}, Emails: []Email{{Value: "user.3@mail.com", Type: "work", Primary: true}}}

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
