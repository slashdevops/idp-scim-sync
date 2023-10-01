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
