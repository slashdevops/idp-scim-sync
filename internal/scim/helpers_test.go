package scim

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

func Test_buildCreateUserRequest(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name string
		args args
		want *aws.CreateUserRequest
	}{
		{
			name: "nil user",
			args: args{
				user: nil,
			},
			want: nil,
		},
		{
			name: "empty user",
			args: args{
				user: &model.User{},
			},
			want: &aws.CreateUserRequest{},
		},
		{
			name: "user with name",
			args: args{
				user: &model.User{
					IPID: "ipid",
					Name: &model.Name{
						FamilyName:      "familyName",
						GivenName:       "givenName",
						Formatted:       "formatted",
						MiddleName:      "middleName",
						HonorificPrefix: "honorificPrefix",
						HonorificSuffix: "honorificSuffix",
					},
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
				Name: &aws.Name{
					FamilyName:      "familyName",
					GivenName:       "givenName",
					Formatted:       "formatted",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
			},
		},
		{
			name: "user with name and email",
			args: args{
				user: &model.User{
					IPID: "ipid",
					Name: &model.Name{
						FamilyName:      "familyName",
						GivenName:       "givenName",
						Formatted:       "formatted",
						MiddleName:      "middleName",
						HonorificPrefix: "honorificPrefix",
						HonorificSuffix: "honorificSuffix",
					},
					Emails: []model.Email{
						{
							Value:   "email",
							Type:    "work",
							Primary: true,
						},
					},
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
				Name: &aws.Name{
					FamilyName:      "familyName",
					GivenName:       "givenName",
					Formatted:       "formatted",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
				Emails: []aws.Email{
					{
						Value:   "email",
						Type:    "work",
						Primary: true,
					},
				},
			},
		},
		{
			name: "user with name and email and phone",
			args: args{
				user: &model.User{
					IPID: "ipid",
					Name: &model.Name{
						FamilyName:      "familyName",
						GivenName:       "givenName",
						Formatted:       "formatted",
						MiddleName:      "middleName",
						HonorificPrefix: "honorificPrefix",
						HonorificSuffix: "honorificSuffix",
					},
					Emails: []model.Email{
						{
							Value:   "email",
							Type:    "work",
							Primary: true,
						},
					},
					PhoneNumbers: []model.PhoneNumber{
						{
							Value: "phone",
							Type:  "work",
						},
					},
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
				Name: &aws.Name{
					FamilyName:      "familyName",
					GivenName:       "givenName",
					Formatted:       "formatted",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
				Emails: []aws.Email{
					{
						Value:   "email",
						Type:    "work",
						Primary: true,
					},
				},
				PhoneNumbers: []aws.PhoneNumber{
					{
						Value: "phone",
						Type:  "work",
					},
				},
			},
		},
		{
			name: "user with nil phone numbers",
			args: args{
				user: &model.User{
					IPID:         "ipid",
					PhoneNumbers: nil,
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
			},
		},
		{
			name: "user with empty (non-nil) phone numbers",
			args: args{
				user: &model.User{
					IPID:         "ipid",
					PhoneNumbers: []model.PhoneNumber{},
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
			},
		},
		{
			name: "user with nil addresses",
			args: args{
				user: &model.User{
					IPID:      "ipid",
					Addresses: nil,
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
			},
		},
		{
			name: "user with empty (non-nil) addresses",
			args: args{
				user: &model.User{
					IPID:      "ipid",
					Addresses: []model.Address{},
				},
			},
			want: &aws.CreateUserRequest{
				ExternalID: "ipid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCreateUserRequest(tt.args.user)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("buildCreateUserRequest() (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_buildPutUserRequest(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name string
		args args
		want *aws.PutUserRequest
	}{
		{
			name: "nil user",
			args: args{
				user: nil,
			},
			want: nil,
		},
		{
			name: "empty user",
			args: args{
				user: &model.User{},
			},
			want: &aws.PutUserRequest{},
		},
		{
			name: "user with name",
			args: args{
				user: &model.User{
					SCIMID: "scimid",
					Name: &model.Name{
						FamilyName:      "familyName",
						GivenName:       "givenName",
						Formatted:       "formatted",
						MiddleName:      "middleName",
						HonorificPrefix: "honorificPrefix",
						HonorificSuffix: "honorificSuffix",
					},
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
				Name: &aws.Name{
					FamilyName:      "familyName",
					GivenName:       "givenName",
					Formatted:       "formatted",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
			},
		},
		{
			name: "user with name and email",
			args: args{
				user: &model.User{
					SCIMID: "scimid",
					Name: &model.Name{
						FamilyName:      "familyName",
						GivenName:       "givenName",
						Formatted:       "formatted",
						MiddleName:      "middleName",
						HonorificPrefix: "honorificPrefix",
						HonorificSuffix: "honorificSuffix",
					},
					Emails: []model.Email{
						{
							Value:   "email",
							Type:    "work",
							Primary: true,
						},
					},
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
				Name: &aws.Name{
					FamilyName:      "familyName",
					GivenName:       "givenName",
					Formatted:       "formatted",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
				Emails: []aws.Email{
					{
						Value:   "email",
						Type:    "work",
						Primary: true,
					},
				},
			},
		},
		{
			name: "user with name and email and phone",
			args: args{
				user: &model.User{
					SCIMID: "scimid",
					Name: &model.Name{
						FamilyName:      "familyName",
						GivenName:       "givenName",
						Formatted:       "formatted",
						MiddleName:      "middleName",
						HonorificPrefix: "honorificPrefix",
						HonorificSuffix: "honorificSuffix",
					},
					Emails: []model.Email{
						{
							Value:   "email",
							Type:    "work",
							Primary: true,
						},
					},
					PhoneNumbers: []model.PhoneNumber{
						{
							Value: "phone",
							Type:  "work",
						},
					},
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
				Name: &aws.Name{
					FamilyName:      "familyName",
					GivenName:       "givenName",
					Formatted:       "formatted",
					MiddleName:      "middleName",
					HonorificPrefix: "honorificPrefix",
					HonorificSuffix: "honorificSuffix",
				},
				Emails: []aws.Email{
					{
						Value:   "email",
						Type:    "work",
						Primary: true,
					},
				},
				PhoneNumbers: []aws.PhoneNumber{
					{
						Value: "phone",
						Type:  "work",
					},
				},
			},
		},
		{
			name: "user with nil phone numbers",
			args: args{
				user: &model.User{
					SCIMID:       "scimid",
					PhoneNumbers: nil,
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
			},
		},
		{
			name: "user with empty (non-nil) phone numbers",
			args: args{
				user: &model.User{
					SCIMID:       "scimid",
					PhoneNumbers: []model.PhoneNumber{},
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
			},
		},
		{
			name: "user with nil addresses",
			args: args{
				user: &model.User{
					SCIMID:    "scimid",
					Addresses: nil,
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
			},
		},
		{
			name: "user with empty (non-nil) addresses",
			args: args{
				user: &model.User{
					SCIMID:    "scimid",
					Addresses: []model.Address{},
				},
			},
			want: &aws.PutUserRequest{
				ID: "scimid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPutUserRequest(tt.args.user)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("buildPutUserRequest() (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_buildUser(t *testing.T) {
	type args struct {
		user *aws.User
	}
	tests := []struct {
		name string
		args args
		want *model.User
	}{
		{
			name: "nil user",
			args: args{user: nil},
			want: nil,
		},
		{
			name: "user with nil phone numbers",
			args: args{
				user: &aws.User{
					ID:           "scimid",
					Name:         &aws.Name{GivenName: "John", FamilyName: "Doe"},
					PhoneNumbers: nil,
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					Build()).
				Build(),
		},
		{
			name: "user with empty (non-nil) phone numbers",
			args: args{
				user: &aws.User{
					ID:           "scimid",
					Name:         &aws.Name{GivenName: "John", FamilyName: "Doe"},
					PhoneNumbers: []aws.PhoneNumber{},
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					Build()).
				Build(),
		},
		{
			name: "user with work phone number",
			args: args{
				user: &aws.User{
					ID:   "scimid",
					Name: &aws.Name{GivenName: "John", FamilyName: "Doe"},
					PhoneNumbers: []aws.PhoneNumber{
						{Value: "123", Type: "work"},
					},
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					Build()).
				WithPhoneNumbers([]model.PhoneNumber{
					model.PhoneNumberBuilder().WithValue("123").WithType("work").Build(),
				}).
				Build(),
		},
		{
			name: "user with empty ID",
			args: args{
				user: &aws.User{
					ID: "",
				},
			},
			want: nil,
		},
		{
			name: "user with nil name",
			args: args{
				user: &aws.User{
					ID:   "scimid",
					Name: nil,
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().Build()).
				Build(),
		},
		{
			name: "user with addresses",
			args: args{
				user: &aws.User{
					ID:   "scimid",
					Name: &aws.Name{GivenName: "John", FamilyName: "Doe"},
					Addresses: []aws.Address{
						{
							Formatted:     "123 Main St",
							StreetAddress: "123 Main St",
							Locality:      "Springfield",
							Region:        "IL",
							PostalCode:    "62701",
							Country:       "US",
						},
					},
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					Build()).
				WithAddresses([]model.Address{
					model.AddressBuilder().
						WithFormatted("123 Main St").
						WithStreetAddress("123 Main St").
						WithLocality("Springfield").
						WithRegion("IL").
						WithPostalCode("62701").
						WithCountry("US").
						Build(),
				}).
				Build(),
		},
		{
			name: "user with enterprise data and manager",
			args: args{
				user: &aws.User{
					ID:   "scimid",
					Name: &aws.Name{GivenName: "John", FamilyName: "Doe"},
					SchemaEnterpriseUser: &aws.SchemaEnterpriseUser{
						EmployeeNumber: "E001",
						CostCenter:     "CC100",
						Organization:   "Acme",
						Division:       "Tech",
						Department:     "Engineering",
						Manager: &aws.Manager{
							Value: "boss-id",
							Ref:   "Users/boss-id",
						},
					},
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					Build()).
				WithEnterpriseData(
					model.EnterpriseDataBuilder().
						WithEmployeeNumber("E001").
						WithCostCenter("CC100").
						WithOrganization("Acme").
						WithDivision("Tech").
						WithDepartment("Engineering").
						WithManager(
							model.ManagerBuilder().
								WithValue("boss-id").
								WithRef("Users/boss-id").
								Build(),
						).
						Build(),
				).
				Build(),
		},
		{
			name: "user with enterprise data without manager",
			args: args{
				user: &aws.User{
					ID:   "scimid",
					Name: &aws.Name{GivenName: "Jane", FamilyName: "Doe"},
					SchemaEnterpriseUser: &aws.SchemaEnterpriseUser{
						Department: "Sales",
					},
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("Jane").
					WithFamilyName("Doe").
					Build()).
				WithEnterpriseData(
					model.EnterpriseDataBuilder().
						WithDepartment("Sales").
						Build(),
				).
				Build(),
		},
		{
			name: "user with nil emails",
			args: args{
				user: &aws.User{
					ID:     "scimid",
					Name:   &aws.Name{GivenName: "John", FamilyName: "Doe"},
					Emails: nil,
				},
			},
			want: model.UserBuilder().
				WithSCIMID("scimid").
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					Build()).
				Build(),
		},
		{
			name: "full user with all fields populated",
			args: args{
				user: &aws.User{
					ID:                "scimid",
					ExternalID:        "ext-id",
					UserName:          "jdoe",
					DisplayName:       "John Doe",
					Title:             "Engineer",
					UserType:          "Employee",
					PreferredLanguage: "en",
					Active:            true,
					Name: &aws.Name{
						GivenName:       "John",
						FamilyName:      "Doe",
						Formatted:       "John Doe",
						MiddleName:      "M",
						HonorificPrefix: "Mr",
						HonorificSuffix: "Jr",
					},
					Emails: []aws.Email{
						{Value: "jdoe@example.com", Type: "work", Primary: true},
					},
					Addresses: []aws.Address{
						{Formatted: "123 Main St", Country: "US"},
					},
					PhoneNumbers: []aws.PhoneNumber{
						{Value: "+1234", Type: "work"},
					},
					SchemaEnterpriseUser: &aws.SchemaEnterpriseUser{
						Department: "Eng",
					},
				},
			},
			want: model.UserBuilder().
				WithIPID("ext-id").
				WithSCIMID("scimid").
				WithUserName("jdoe").
				WithDisplayName("John Doe").
				WithTitle("Engineer").
				WithUserType("Employee").
				WithPreferredLanguage("en").
				WithActive(true).
				WithName(model.NameBuilder().
					WithGivenName("John").
					WithFamilyName("Doe").
					WithFormatted("John Doe").
					WithMiddleName("M").
					WithHonorificPrefix("Mr").
					WithHonorificSuffix("Jr").
					Build()).
				WithEmails([]model.Email{
					model.EmailBuilder().WithValue("jdoe@example.com").WithType("work").WithPrimary(true).Build(),
				}).
				WithAddresses([]model.Address{
					model.AddressBuilder().WithFormatted("123 Main St").WithCountry("US").Build(),
				}).
				WithPhoneNumbers([]model.PhoneNumber{
					model.PhoneNumberBuilder().WithValue("+1234").WithType("work").Build(),
				}).
				WithEnterpriseData(
					model.EnterpriseDataBuilder().WithDepartment("Eng").Build(),
				).
				Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildUser(tt.args.user)

			sort := func(x, y string) bool { return x > y }
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(sort)); diff != "" {
				t.Errorf("buildUser() (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_buildCreateUserRequest_withAllFields(t *testing.T) {
	t.Run("user with addresses and enterprise data", func(t *testing.T) {
		user := &model.User{
			IPID:        "ipid",
			UserName:    "jdoe",
			DisplayName: "John Doe",
			Name: &model.Name{
				FamilyName: "Doe",
				GivenName:  "John",
			},
			Emails: []model.Email{
				{Value: "jdoe@example.com", Type: "work", Primary: true},
			},
			Addresses: []model.Address{
				{
					Formatted:     "123 Main St",
					StreetAddress: "123 Main St",
					Locality:      "Springfield",
					Region:        "IL",
					PostalCode:    "62701",
					Country:       "US",
				},
			},
			PhoneNumbers: []model.PhoneNumber{
				{Value: "+1234", Type: "work"},
			},
			EnterpriseData: &model.EnterpriseData{
				EmployeeNumber: "E001",
				CostCenter:     "CC100",
				Organization:   "Acme",
				Division:       "Tech",
				Department:     "Engineering",
				Manager: &model.Manager{
					Value: "boss-id",
					Ref:   "Users/boss-id",
				},
			},
		}

		got := buildCreateUserRequest(user)

		if got == nil {
			t.Fatal("expected non-nil result")
		}
		if len(got.Addresses) != 1 {
			t.Fatalf("expected 1 address, got %d", len(got.Addresses))
		}
		if got.Addresses[0].Country != "US" {
			t.Errorf("expected country 'US', got %q", got.Addresses[0].Country)
		}
		if got.SchemaEnterpriseUser == nil {
			t.Fatal("expected enterprise user data")
		}
		if got.SchemaEnterpriseUser.EmployeeNumber != "E001" {
			t.Errorf("expected employee number 'E001', got %q", got.SchemaEnterpriseUser.EmployeeNumber)
		}
		if got.SchemaEnterpriseUser.Manager == nil {
			t.Fatal("expected manager")
		}
		if got.SchemaEnterpriseUser.Manager.Value != "boss-id" {
			t.Errorf("expected manager value 'boss-id', got %q", got.SchemaEnterpriseUser.Manager.Value)
		}
	})
}

func Test_buildPutUserRequest_withAllFields(t *testing.T) {
	t.Run("user with addresses and enterprise data", func(t *testing.T) {
		user := &model.User{
			SCIMID:      "scimid",
			IPID:        "ipid",
			UserName:    "jdoe",
			DisplayName: "John Doe",
			Name: &model.Name{
				FamilyName: "Doe",
				GivenName:  "John",
			},
			Emails: []model.Email{
				{Value: "jdoe@example.com", Type: "work", Primary: true},
			},
			Addresses: []model.Address{
				{
					Formatted:     "456 Oak Ave",
					StreetAddress: "456 Oak Ave",
					Locality:      "Chicago",
					Region:        "IL",
					PostalCode:    "60601",
					Country:       "US",
				},
			},
			PhoneNumbers: []model.PhoneNumber{
				{Value: "+5678", Type: "home"},
			},
			EnterpriseData: &model.EnterpriseData{
				Department: "Sales",
				Manager: &model.Manager{
					Value: "mgr-id",
					Ref:   "Users/mgr-id",
				},
			},
		}

		got := buildPutUserRequest(user)

		if got == nil {
			t.Fatal("expected non-nil result")
		}
		if got.ID != "scimid" {
			t.Errorf("expected ID 'scimid', got %q", got.ID)
		}
		if len(got.Addresses) != 1 {
			t.Fatalf("expected 1 address, got %d", len(got.Addresses))
		}
		if got.Addresses[0].Locality != "Chicago" {
			t.Errorf("expected locality 'Chicago', got %q", got.Addresses[0].Locality)
		}
		if len(got.PhoneNumbers) != 1 {
			t.Fatalf("expected 1 phone, got %d", len(got.PhoneNumbers))
		}
		if got.SchemaEnterpriseUser == nil {
			t.Fatal("expected enterprise user data")
		}
		if got.SchemaEnterpriseUser.Department != "Sales" {
			t.Errorf("expected department 'Sales', got %q", got.SchemaEnterpriseUser.Department)
		}
		if got.SchemaEnterpriseUser.Manager == nil {
			t.Fatal("expected manager")
		}
		if got.SchemaEnterpriseUser.Manager.Ref != "Users/mgr-id" {
			t.Errorf("expected manager ref 'Users/mgr-id', got %q", got.SchemaEnterpriseUser.Manager.Ref)
		}
	})
}
