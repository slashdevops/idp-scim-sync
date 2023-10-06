package scim

import (
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/convert"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// func buildUser creates a User model from a CreateUserRequest
func buildUser(user *aws.User) *model.User {
	if user == nil {
		return nil
	}

	var emails []model.Email
	if user.Emails != nil {
		for _, email := range user.Emails {
			if email.Primary {
				emails = append(emails,
					model.EmailBuilder().
						WithPrimary(email.Primary).
						WithType(email.Type).
						WithValue(email.Value).
						Build(),
				)
			}
		}
	}

	var addresses []model.Address
	if user.Addresses != nil {
		if len(user.Addresses) > 0 {
			addresses = append(addresses,
				model.AddressBuilder().
					WithType(user.Addresses[0].Type).
					WithFormatted(user.Addresses[0].Formatted).
					WithStreetAddress(user.Addresses[0].StreetAddress).
					WithLocality(user.Addresses[0].Locality).
					WithRegion(user.Addresses[0].Region).
					WithPostalCode(user.Addresses[0].PostalCode).
					WithCountry(user.Addresses[0].Country).
					Build(),
			)
		}
	}

	var phoneNumbers []model.PhoneNumber
	if user.PhoneNumbers != nil {
		for _, phoneNumber := range user.PhoneNumbers {
			phoneNumbers = append(phoneNumbers,
				model.PhoneNumberBuilder().
					WithType(phoneNumber.Type).
					WithValue(phoneNumber.Value).
					Build(),
			)
		}
	}

	var enterpriseData *model.EnterpriseData
	if user.SchemaEnterpriseUser != nil {

		var manager *model.Manager
		if user.SchemaEnterpriseUser.Manager != nil {
			manager = model.ManagerBuilder().
				WithValue(user.SchemaEnterpriseUser.Manager.Value).
				WithRef(user.SchemaEnterpriseUser.Manager.Ref).
				Build()
		}

		enterpriseData = model.EnterpriseDataBuilder().
			WithEmployeeNumber(user.SchemaEnterpriseUser.EmployeeNumber).
			WithCostCenter(user.SchemaEnterpriseUser.CostCenter).
			WithOrganization(user.SchemaEnterpriseUser.Organization).
			WithDivision(user.SchemaEnterpriseUser.Division).
			WithDepartment(user.SchemaEnterpriseUser.Department).
			WithManager(manager).
			Build()

	}

	var name *model.Name
	if user.Name != nil {
		name = model.NameBuilder().
			WithFamilyName(user.Name.FamilyName).
			WithGivenName(user.Name.GivenName).
			WithFormatted(user.Name.Formatted).
			WithMiddleName(user.Name.MiddleName).
			WithHonorificPrefix(user.Name.HonorificPrefix).
			WithHonorificSuffix(user.Name.HonorificSuffix).
			Build()
	}

	userModel := model.UserBuilder().
		WithIPID(user.ExternalID).
		WithSCIMID(user.ID).
		WithUserName(user.UserName).
		WithDisplayName(user.DisplayName).
		// WithNickName("Not Provided").
		// WithProfileURL("Not Provided").
		WithUserType(user.UserType).
		WithTitle(user.Title).
		WithPreferredLanguage(user.PreferredLanguage).
		// WithLocale("Not Provided").
		// WithTimezone("Not Provided").
		WithActive(user.Active).
		// arrays
		WithEmails(emails).
		WithAddresses(addresses).
		WithPhoneNumbers(phoneNumbers).
		// Pointers
		WithName(name).
		WithEnterpriseData(enterpriseData).
		Build()

	log.Tracef("scim: buildUser() from: %+v, --> to: %+v", convert.ToJSONString(user), convert.ToJSONString(userModel))

	return userModel
}

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
		userRequest.Addresses = []aws.Address{
			{
				Formatted:     user.Addresses[0].Formatted,
				Type:          user.Addresses[0].Type,
				StreetAddress: user.Addresses[0].StreetAddress,
				Locality:      user.Addresses[0].Locality,
				Region:        user.Addresses[0].Region,
				PostalCode:    user.Addresses[0].PostalCode,
				Country:       user.Addresses[0].Country,
			},
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
		userRequest.Addresses = []aws.Address{
			{
				Formatted:     user.Addresses[0].Formatted,
				Type:          user.Addresses[0].Type,
				StreetAddress: user.Addresses[0].StreetAddress,
				Locality:      user.Addresses[0].Locality,
				Region:        user.Addresses[0].Region,
				PostalCode:    user.Addresses[0].PostalCode,
				Country:       user.Addresses[0].Country,
			},
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
