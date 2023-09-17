package idp

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	admin "google.golang.org/api/admin/directory/v1"
)

func buildUser(usr *admin.User) *model.User {
	if usr == nil {
		return nil
	}

	if usr.Name == nil {
		usr.Name = &admin.UserName{}
	}

	// get the first language from the list of languages
	var preferredLanguage string
	var languages []admin.UserLanguage
	if languagesBytes, ok := usr.Languages.([]byte); ok {
		err := json.Unmarshal(languagesBytes, &languages)
		if err != nil {
			log.Warn("idp: Error unmarshalling languages: ", err)
		}
	}
	if len(languages) > 0 {
		preferredLanguage = languages[0].LanguageCode
	}

	// get the Addresses
	var mainAddress model.Address
	var locale string
	var addresses []admin.UserAddress
	if addressesBytes, ok := usr.Addresses.([]byte); ok {
		err := json.Unmarshal(addressesBytes, &addresses)
		if err != nil {
			log.Warn("idp: Error unmarshalling addresses: ", err)
		}
	}
	if len(addresses) > 0 {
		mainAddress = model.AddressBuilder().
			WithStreetAddress(addresses[0].StreetAddress).
			WithLocality(addresses[0].Locality).
			WithRegion(addresses[0].Region).
			WithPostalCode(addresses[0].PostalCode).
			WithCountry(addresses[0].Country).
			WithFormatted(addresses[0].Formatted).
			Build()

		locale = addresses[0].CountryCode
	}

	// get the phones
	var mainPhone model.PhoneNumber
	var phones []admin.UserPhone
	if phonesBytes, ok := usr.Phones.([]byte); ok {
		err := json.Unmarshal(phonesBytes, &phones)
		if err != nil {
			log.Warn("idp: Error unmarshalling phones: ", err)
		}
	}
	if len(phones) > 0 {
		mainPhone = model.PhoneNumberBuilder().
			WithValue(phones[0].Value).
			WithType(phones[0].Type).
			Build()
	}

	// get the organizations
	var mainOrganization model.EnterpriseData
	var title string
	var organizations []admin.UserOrganization
	if organizationsBytes, ok := usr.Organizations.([]byte); ok {
		err := json.Unmarshal(organizationsBytes, &organizations)
		if err != nil {
			log.Warn("idp: Error unmarshalling organizations: ", err)
		}
	}
	if len(organizations) > 0 {
		mainOrganization = *model.EnterpriseDataBuilder().
			WithCostCenter(organizations[0].CostCenter).
			WithOrganization(organizations[0].Name).
			WithDepartment(organizations[0].Department).
			Build()

		title = organizations[0].Title
	}

	email := model.EmailBuilder().
		WithValue(usr.PrimaryEmail).
		WithType("work").
		WithPrimary(true).
		Build()

	var displayName string
	if usr.Name.DisplayName != "" {
		displayName = usr.Name.DisplayName
	} else {
		displayName = fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName)
	}

	modUser := model.UserBuilder().
		WithIPID(usr.Id).
		WithUserName(usr.PrimaryEmail).
		WithDisplayName(displayName).
		// WithNickName("Not Provided").
		// WithTitle("Not Provided").
		// WithTimezone("Not Provided").
		WithProfileURL(usr.OrgUnitPath).
		WithUserType(usr.Kind).
		WithEmail(email).
		WithGivenName(usr.Name.GivenName).
		WithFamilyName(usr.Name.FamilyName).
		WithActive(!usr.Suspended).
		Build()

	if mainAddress != (model.Address{}) {
		modUser.Addresses = append(modUser.Addresses, mainAddress)
	}
	if preferredLanguage != "" {
		modUser.PreferredLanguage = preferredLanguage
	}
	if locale != "" {
		modUser.Locale = locale
	}
	if mainPhone != (model.PhoneNumber{}) {
		modUser.PhoneNumbers = append(modUser.PhoneNumbers, mainPhone)
	}
	if mainOrganization != (model.EnterpriseData{}) {
		modUser.EnterpriseData = &mainOrganization
	}
	if title != "" {
		modUser.Title = title
	}

	return modUser
}
