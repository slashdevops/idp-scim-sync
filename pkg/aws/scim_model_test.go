package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	type fields struct {
		ID                   string
		ExternalID           string
		Meta                 *Meta
		Schemas              []string
		UserName             string
		Name                 *Name
		DisplayName          string
		NickName             string
		ProfileURL           string
		Title                string
		UserType             string
		PreferredLanguage    string
		Locale               string
		Timezone             string
		Active               bool
		Emails               []Email
		Addresses            []Address
		PhoneNumbers         []PhoneNumber
		SchemaEnterpriseUser *SchemaEnterpriseUser
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid user",
			wantErr: false,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "Empty UserName",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "Empty DisplayName",
			wantErr: true,
			fields: fields{
				ID:         "2819c223-7f76-453a-919d-413861904646",
				ExternalID: "701984",
				UserName:   "bjensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "empty GivenName",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "empty FamilyName",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "empty Emails",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "too many Emails",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
					{
						Value:   "bjensen@mail.com",
						Type:    "home",
						Primary: false,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "no primary Emails",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: false,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "too many Addresses",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
		{
			name:    "too many PhoneNumbers",
			wantErr: true,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				UserName:    "bjensen",
				ExternalID:  "701984",
				DisplayName: "Babs Jensen",
				Name: &Name{
					GivenName:       "Barbara",
					FamilyName:      "Jensen",
					Formatted:       "Ms. Barbara J Jensen, III",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				NickName:   "Babs",
				ProfileURL: "https://login.example.com/bjensen",
				Emails: []Email{
					{
						Value:   "bjensen@example.com",
						Type:    "work",
						Primary: true,
					},
				},
				Addresses: []Address{
					{
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza Hollywood, CA 91608 USA",
					},
				},
				PhoneNumbers: []PhoneNumber{
					{
						Value: "555-555-5555",
						Type:  "work",
					},
					{
						Value: "555-555-5555",
						Type:  "home",
					},
				},
				UserType:          "Employee",
				Title:             "Tour Guide",
				PreferredLanguage: "en-US",
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            true,
				SchemaEnterpriseUser: &SchemaEnterpriseUser{
					EmployeeNumber: "701984",
					CostCenter:     "4130",
					Organization:   "Universal Studios",
					Division:       "Theme Park",
					Department:     "Tour Operations",
					Manager: &Manager{
						Value: "9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
						Ref:   "../Users/9067729b3d-ee533c18-538a-4cd3-a572-63fb863ed734",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:                   tt.fields.ID,
				ExternalID:           tt.fields.ExternalID,
				Meta:                 tt.fields.Meta,
				Schemas:              tt.fields.Schemas,
				UserName:             tt.fields.UserName,
				Name:                 tt.fields.Name,
				DisplayName:          tt.fields.DisplayName,
				NickName:             tt.fields.NickName,
				ProfileURL:           tt.fields.ProfileURL,
				Title:                tt.fields.Title,
				UserType:             tt.fields.UserType,
				PreferredLanguage:    tt.fields.PreferredLanguage,
				Locale:               tt.fields.Locale,
				Timezone:             tt.fields.Timezone,
				Active:               tt.fields.Active,
				Emails:               tt.fields.Emails,
				Addresses:            tt.fields.Addresses,
				PhoneNumbers:         tt.fields.PhoneNumbers,
				SchemaEnterpriseUser: tt.fields.SchemaEnterpriseUser,
			}
			if err := u.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroup_Validate(t *testing.T) {
	type fields struct {
		ID          string
		Meta        Meta
		Schemas     []string
		DisplayName string
		ExternalID  string
		Members     []*Member
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid group",
			wantErr: false,
			fields: fields{
				ID:          "2819c223-7f76-453a-919d-413861904646",
				DisplayName: "Tour Guides",
				ExternalID:  "701984",
				Members: []*Member{
					{
						Value: "2819c223-7f76-453a-919d-413861904646",
						Ref:   "../Users/2819c223-7f76-453a-919d-413861904646",
					},
				},
			},
		},
		{
			name:    "empty DisplayName",
			wantErr: true,
			fields: fields{
				ID:         "2819c223-7f76-453a-919d-413861904646",
				ExternalID: "701984",
				Members: []*Member{
					{
						Value: "2819c223-7f76-453a-919d-413861904646",
						Ref:   "../Users/2819c223-7f76-453a-919d-413861904646",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				ID:          tt.fields.ID,
				Meta:        tt.fields.Meta,
				Schemas:     tt.fields.Schemas,
				DisplayName: tt.fields.DisplayName,
				ExternalID:  tt.fields.ExternalID,
				Members:     tt.fields.Members,
			}
			if err := g.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Group.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Enhanced test cases for improved coverage

func TestUser_ValidateEnhanced(t *testing.T) {
	t.Run("should fail when user name is nil", func(t *testing.T) {
		user := &User{
			DisplayName: "Test User",
			Name: &Name{
				GivenName:  "Test",
				FamilyName: "User",
			},
			Emails: []Email{
				{Value: "test@example.com", Primary: true},
			},
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUserNameEmpty)
	})

	t.Run("should fail when Name is nil", func(t *testing.T) {
		user := &User{
			UserName:    "testuser",
			DisplayName: "Test User",
			Name:        nil,
			Emails: []Email{
				{Value: "test@example.com", Primary: true},
			},
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNameEmpty)
	})

	t.Run("should fail when email value is empty", func(t *testing.T) {
		user := &User{
			UserName:    "testuser",
			DisplayName: "Test User",
			Name: &Name{
				GivenName:  "Test",
				FamilyName: "User",
			},
			Emails: []Email{
				{Value: "", Primary: true},
			},
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrEmailValueEmpty)
	})

	t.Run("should fail when no primary email", func(t *testing.T) {
		user := &User{
			UserName:    "testuser",
			DisplayName: "Test User",
			Name: &Name{
				GivenName:  "Test",
				FamilyName: "User",
			},
			Emails: []Email{
				{Value: "test@example.com", Primary: false},
			},
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPrimaryEmailEmpty)
	})

	t.Run("should fail when multiple primary emails", func(t *testing.T) {
		user := &User{
			UserName:    "testuser",
			DisplayName: "Test User",
			Name: &Name{
				GivenName:  "Test",
				FamilyName: "User",
			},
			Emails: []Email{
				{Value: "test1@example.com", Primary: true},
				{Value: "test2@example.com", Primary: true},
			},
		}
		err := user.Validate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrEmailsTooMany)
	})

	t.Run("should pass with valid user", func(t *testing.T) {
		user := &User{
			UserName:    "testuser",
			DisplayName: "Test User",
			Name: &Name{
				GivenName:  "Test",
				FamilyName: "User",
			},
			Emails: []Email{
				{Value: "test@example.com", Primary: true},
			},
		}
		err := user.Validate()
		assert.NoError(t, err)
	})
}

func TestUser_GetPrimaryEmail(t *testing.T) {
	t.Run("should return primary email", func(t *testing.T) {
		user := &User{
			Emails: []Email{
				{Value: "test1@example.com", Primary: false},
				{Value: "test2@example.com", Primary: true},
			},
		}
		primaryEmail := user.GetPrimaryEmail()
		assert.NotNil(t, primaryEmail)
		assert.Equal(t, "test2@example.com", primaryEmail.Value)
		assert.True(t, primaryEmail.Primary)
	})

	t.Run("should return nil when no primary email", func(t *testing.T) {
		user := &User{
			Emails: []Email{
				{Value: "test1@example.com", Primary: false},
				{Value: "test2@example.com", Primary: false},
			},
		}
		primaryEmail := user.GetPrimaryEmail()
		assert.Nil(t, primaryEmail)
	})

	t.Run("should return nil when no emails", func(t *testing.T) {
		user := &User{
			Emails: []Email{},
		}
		primaryEmail := user.GetPrimaryEmail()
		assert.Nil(t, primaryEmail)
	})
}

func TestUser_String(t *testing.T) {
	t.Run("should return JSON string representation", func(t *testing.T) {
		user := &User{
			UserName:    "testuser",
			DisplayName: "Test User",
		}
		str := user.String()
		assert.Contains(t, str, "testuser")
		assert.Contains(t, str, "Test User")
	})

	t.Run("should include all populated fields in JSON", func(t *testing.T) {
		user := &User{
			ID:          "user-id-123",
			UserName:    "jdoe",
			DisplayName: "Jane Doe",
			Active:      true,
			Name: &Name{
				GivenName:  "Jane",
				FamilyName: "Doe",
			},
			Emails: []Email{
				{Value: "jdoe@example.com", Type: "work", Primary: true},
			},
		}
		str := user.String()
		assert.Contains(t, str, "user-id-123")
		assert.Contains(t, str, "jdoe")
		assert.Contains(t, str, "Jane Doe")
		assert.Contains(t, str, "jdoe@example.com")
	})
}

func TestUser_GetPrimaryAddress(t *testing.T) {
	t.Run("should return first address when addresses exist", func(t *testing.T) {
		user := &User{
			Addresses: []Address{
				{Formatted: "123 Main St, NY", Country: "US"},
				{Formatted: "456 Oak Ave, LA", Country: "US"},
			},
		}
		addr := user.GetPrimaryAddress()
		assert.NotNil(t, addr)
		assert.Equal(t, "123 Main St, NY", addr.Formatted)
		assert.Equal(t, "US", addr.Country)
	})

	t.Run("should return nil when no addresses", func(t *testing.T) {
		user := &User{
			Addresses: []Address{},
		}
		addr := user.GetPrimaryAddress()
		assert.Nil(t, addr)
	})

	t.Run("should return nil when addresses is nil", func(t *testing.T) {
		user := &User{}
		addr := user.GetPrimaryAddress()
		assert.Nil(t, addr)
	})
}

func TestUser_GetPrimaryEmailAddress(t *testing.T) {
	t.Run("should return primary email address string", func(t *testing.T) {
		user := &User{
			Emails: []Email{
				{Value: "secondary@example.com", Primary: false},
				{Value: "primary@example.com", Primary: true},
			},
		}
		got := user.GetPrimaryEmailAddress()
		assert.Equal(t, "primary@example.com", got)
	})

	t.Run("should return empty when no primary", func(t *testing.T) {
		user := &User{
			Emails: []Email{
				{Value: "noprimary@example.com", Primary: false},
			},
		}
		got := user.GetPrimaryEmailAddress()
		assert.Equal(t, "", got)
	})

	t.Run("should return empty when no emails", func(t *testing.T) {
		user := &User{}
		got := user.GetPrimaryEmailAddress()
		assert.Equal(t, "", got)
	})
}

func TestGroup_String(t *testing.T) {
	t.Run("should return JSON string representation", func(t *testing.T) {
		group := &Group{
			ID:          "group-id-456",
			DisplayName: "Engineering",
		}
		str := group.String()
		assert.Contains(t, str, "group-id-456")
		assert.Contains(t, str, "Engineering")
	})

	t.Run("should include members in JSON", func(t *testing.T) {
		group := &Group{
			DisplayName: "Team",
			Members: []*Member{
				{Value: "user-1"},
			},
		}
		str := group.String()
		assert.Contains(t, str, "Team")
		assert.Contains(t, str, "user-1")
	})
}
