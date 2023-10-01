package scim

import (
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/convert"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// buildCreateUserRequest builds a CreateUserRequest from a User model
func buildCreateUserRequest(user *model.User) *aws.CreateUserRequest {
	if user == nil {
		return nil
	}

	userRequest := &aws.CreateUserRequest{
		ExternalID:        user.IPID,
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		UserType:          user.UserType,
		Title:             user.Title,
		PreferredLanguage: user.PreferredLanguage,
		Locale:            user.Locale,
		Timezone:          user.Timezone,
		Active:            user.Active,
	}

	if user.Name != nil {
		userRequest.Name = &aws.Name{
			FamilyName:      user.Name.FamilyName,
			GivenName:       user.Name.GivenName,
			Formatted:       user.Name.Formatted,
			MiddleName:      user.Name.MiddleName,
			HonorificPrefix: user.Name.HonorificPrefix,
			HonorificSuffix: user.Name.HonorificSuffix,
		}
	}

	if user.Emails != nil {
		for _, email := range user.Emails {
			if email.Primary {
				userRequest.Emails = append(userRequest.Emails, aws.Email{
					Value:   email.Value,
					Type:    email.Type,
					Primary: email.Primary,
				})
			}
		}
	}

	if user.Addresses != nil {
		for _, address := range user.Addresses {
			if address.Primary {
				userRequest.Addresses = []aws.Address{
					{
						Formatted:     user.Addresses[0].Formatted,
						Type:          user.Addresses[0].Type,
						StreetAddress: user.Addresses[0].StreetAddress,
						Locality:      user.Addresses[0].Locality,
						Region:        user.Addresses[0].Region,
						PostalCode:    user.Addresses[0].PostalCode,
						Country:       user.Addresses[0].Country,
						Primary:       user.Addresses[0].Primary,
					},
				}
			}
		}
	}

	if user.PhoneNumbers != nil {
		userRequest.PhoneNumbers = []aws.PhoneNumber{
			{
				Value: user.PhoneNumbers[0].Value,
				Type:  user.PhoneNumbers[0].Type,
			},
		}
	}

	if user.EnterpriseData != nil {
		userRequest.SchemaEnterpriseUser = &aws.SchemaEnterpriseUser{
			EmployeeNumber: user.EnterpriseData.EmployeeNumber,
			CostCenter:     user.EnterpriseData.CostCenter,
			Organization:   user.EnterpriseData.Organization,
			Division:       user.EnterpriseData.Division,
			Department:     user.EnterpriseData.Department,
		}

		if user.EnterpriseData.Manager != nil {
			userRequest.SchemaEnterpriseUser.Manager = &aws.Manager{
				Value: user.EnterpriseData.Manager.Value,
			}
		}
	}

	log.Tracef("scim buildCreateUserRequest(): %+v", convert.ToJSONString(userRequest))

	return userRequest
}

// buildPutUserRequest builds a PutUserRequest from a User model
func buildPutUserRequest(user *model.User) *aws.PutUserRequest {
	if user == nil {
		return nil
	}

	userRequest := &aws.PutUserRequest{
		ID:                user.SCIMID,
		ExternalID:        user.IPID,
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		UserType:          user.UserType,
		Title:             user.Title,
		PreferredLanguage: user.PreferredLanguage,
		Locale:            user.Locale,
		Timezone:          user.Timezone,
		Active:            user.Active,
	}

	if user.Name != nil {
		userRequest.Name = &aws.Name{
			FamilyName:      user.Name.FamilyName,
			GivenName:       user.Name.GivenName,
			Formatted:       user.Name.Formatted,
			MiddleName:      user.Name.MiddleName,
			HonorificPrefix: user.Name.HonorificPrefix,
			HonorificSuffix: user.Name.HonorificSuffix,
		}
	}

	if user.Emails != nil {
		for _, email := range user.Emails {
			if email.Primary {
				userRequest.Emails = append(userRequest.Emails, aws.Email{
					Value:   email.Value,
					Type:    email.Type,
					Primary: email.Primary,
				})
			}
		}
	}

	if user.Addresses != nil {
		for _, address := range user.Addresses {
			if address.Primary {
				userRequest.Addresses = []aws.Address{
					{
						Formatted:     user.Addresses[0].Formatted,
						Type:          user.Addresses[0].Type,
						StreetAddress: user.Addresses[0].StreetAddress,
						Locality:      user.Addresses[0].Locality,
						Region:        user.Addresses[0].Region,
						PostalCode:    user.Addresses[0].PostalCode,
						Country:       user.Addresses[0].Country,
						Primary:       user.Addresses[0].Primary,
					},
				}
			}
		}
	}

	if user.PhoneNumbers != nil {
		userRequest.PhoneNumbers = []aws.PhoneNumber{
			{
				Value: user.PhoneNumbers[0].Value,
				Type:  user.PhoneNumbers[0].Type,
			},
		}
	}

	if user.EnterpriseData != nil {
		userRequest.SchemaEnterpriseUser = &aws.SchemaEnterpriseUser{
			EmployeeNumber: user.EnterpriseData.EmployeeNumber,
			CostCenter:     user.EnterpriseData.CostCenter,
			Organization:   user.EnterpriseData.Organization,
			Division:       user.EnterpriseData.Division,
			Department:     user.EnterpriseData.Department,
		}

		if user.EnterpriseData.Manager != nil {
			userRequest.SchemaEnterpriseUser.Manager = &aws.Manager{
				Value: user.EnterpriseData.Manager.Value,
			}
		}
	}

	log.Tracef("scim buildPutUserRequest(): %+v", convert.ToJSONString(userRequest))

	return userRequest
}
