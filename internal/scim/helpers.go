package scim

import (
	"log/slog"
	"strings"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/pkg/aws"
)

// func buildUser creates a User model from a CreateUserRequest
func buildUser(user *aws.User) *model.User {
	if user == nil {
		return nil
	}

	if user.ID == "" {
		slog.Warn("scim: User ID is empty")
		return nil
	}

	if user.Name == nil {
		slog.Warn("scim: User name is nil")
		user.Name = &aws.Name{}
	}

	if user.Name.GivenName == "" {
		slog.Warn("scim: User given name is empty")
	}

	if user.Name.FamilyName == "" {
		slog.Warn("scim: User family name is empty")
	}

	if user.Emails == nil {
		slog.Warn("scim: User emails is nil, setting primary email as the only email")
		user.Emails = []aws.Email{}
	}

	var emails []model.Email
	if user.Emails != nil {
		for _, email := range user.Emails {
			if email.Primary {
				emails = append(emails,
					model.EmailBuilder().
						WithPrimary(email.Primary).
						WithType(strings.TrimSpace(email.Type)).
						WithValue(strings.TrimSpace(email.Value)).
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
					WithFormatted(strings.TrimSpace(user.Addresses[0].Formatted)).
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
		phoneNumbers = append(phoneNumbers,
			model.PhoneNumberBuilder().
				WithValue(strings.TrimSpace(user.PhoneNumbers[0].Value)).
				WithType(strings.TrimSpace(user.PhoneNumbers[0].Type)).
				Build(),
		)
	}

	var enterpriseData *model.EnterpriseData
	if user.SchemaEnterpriseUser != nil {

		var manager *model.Manager
		if user.SchemaEnterpriseUser.Manager != nil {
			manager = model.ManagerBuilder().
				WithValue(strings.TrimSpace(user.SchemaEnterpriseUser.Manager.Value)).
				WithRef(strings.TrimSpace(user.SchemaEnterpriseUser.Manager.Ref)).
				Build()
		}

		enterpriseData = model.EnterpriseDataBuilder().
			WithEmployeeNumber(strings.TrimSpace(user.SchemaEnterpriseUser.EmployeeNumber)).
			WithCostCenter(strings.TrimSpace(user.SchemaEnterpriseUser.CostCenter)).
			WithOrganization(strings.TrimSpace(user.SchemaEnterpriseUser.Organization)).
			WithDivision(strings.TrimSpace(user.SchemaEnterpriseUser.Division)).
			WithDepartment(strings.TrimSpace(user.SchemaEnterpriseUser.Department)).
			WithManager(manager).
			Build()

	}

	var name *model.Name
	if user.Name != nil {
		name = model.NameBuilder().
			WithGivenName(strings.TrimSpace(user.Name.GivenName)).
			WithFamilyName(strings.TrimSpace(user.Name.FamilyName)).
			WithFormatted(strings.TrimSpace(user.Name.Formatted)).
			WithMiddleName(user.Name.MiddleName).
			WithHonorificPrefix(user.Name.HonorificPrefix).
			WithHonorificSuffix(user.Name.HonorificSuffix).
			Build()
	}

	userModel := model.UserBuilder().
		WithIPID(strings.TrimSpace(user.ExternalID)).
		WithSCIMID(strings.TrimSpace(user.ID)).
		WithUserName(strings.TrimSpace(user.UserName)).
		WithDisplayName(strings.TrimSpace(user.DisplayName)).
		// WithNickName("Not Provided").
		// WithProfileURL("Not Provided").
		WithTitle(strings.TrimSpace(user.Title)).
		WithUserType(strings.TrimSpace(user.UserType)).
		WithPreferredLanguage(strings.TrimSpace(user.PreferredLanguage)).
		// WithLocale("Not Provided").
		// WithTimezone("Not Provided").
		WithActive(user.Active).
		// Arrays
		WithEmails(emails).
		WithAddresses(addresses).
		WithPhoneNumbers(phoneNumbers).
		// Pointers
		WithName(name).
		WithEnterpriseData(enterpriseData).
		Build()

	slog.Debug("scim: buildUser() converted user", "from", user, "to", userModel)

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
				Ref:   user.EnterpriseData.Manager.Ref,
			}
		}
	}

	slog.Debug("scim: buildCreateUserRequest()", "user", userRequest)

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
				Ref:   user.EnterpriseData.Manager.Ref,
			}
		}
	}

	slog.Debug("scim: buildPutUserRequest()", "user", userRequest)

	return userRequest
}
