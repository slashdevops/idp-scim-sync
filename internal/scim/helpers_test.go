package scim

import (
	"reflect"
	"testing"

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
			if got := buildCreateUserRequest(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildCreateUserRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
