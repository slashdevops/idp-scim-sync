package idp

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/convert"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	admin "google.golang.org/api/admin/directory/v1"
)

func buildUser(usr *admin.User) *model.User {
	if usr == nil {
		return nil
	}

	if usr.Name == nil {
		log.Warn("idp: User name is nil")
		return nil
	}

	if usr.Name.GivenName == "" {
		log.Warn("idp: User given name is empty")
		return nil
	}

	if usr.Name.FamilyName == "" {
		log.Warn("idp: User family name is empty")
		return nil
	}

	if usr.PrimaryEmail == "" {
		log.Warn("idp: User primary email is empty")
		return nil
	}

	if usr.Emails == nil {
		log.Warn("idp: User emails is nil, setting primary email as the only email")
		usr.Emails = []interface{}{
			map[string]interface{}{
				"address": usr.PrimaryEmail,
				"type":    "work",
				"primary": true,
			},
		}
	}

	var emails []model.Email
	if m, ok := usr.Emails.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["primary"] != nil {
				if v.(map[string]interface{})["primary"].(bool) {
					emails = append(emails, model.EmailBuilder().
						// WithValue(v.(map[string]interface{})["address"].(string)).
						// WithType(v.(map[string]interface{})["type"].(string)).
						WithPrimary(v.(map[string]interface{})["primary"].(bool)).
						Build())

					if v.(map[string]interface{})["address"] != nil {
						emails[0].Value = v.(map[string]interface{})["address"].(string)
					}

					if v.(map[string]interface{})["type"] != nil {
						emails[0].Type = v.(map[string]interface{})["type"].(string)
					}

					break
				}
			}
		}
	}

	// get the first language from the list of languages
	var preferredLanguage string
	if m, ok := usr.Languages.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["preference"].(string) == "preferred" {
				preferredLanguage = v.(map[string]interface{})["languageCode"].(string)
				break
			}
		}
	}

	// get the Addresses
	var mainAddress model.Address
	if m, ok := usr.Addresses.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["type"].(string) == "work" {
				mainAddress = model.AddressBuilder().
					WithFormatted(v.(map[string]interface{})["formatted"].(string)).
					WithType(v.(map[string]interface{})["type"].(string)).
					WithPrimary(true).
					Build()
				break
			} else if v.(map[string]interface{})["type"].(string) == "home" {
				mainAddress = model.AddressBuilder().
					WithFormatted(v.(map[string]interface{})["formatted"].(string)).
					WithType(v.(map[string]interface{})["type"].(string)).
					WithPrimary(true).
					Build()
				break
			}
		}
	}

	// get the phones
	var mainPhone model.PhoneNumber
	if m, ok := usr.Phones.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["type"].(string) == "work" {
				mainPhone = model.PhoneNumberBuilder().
					WithValue(v.(map[string]interface{})["value"].(string)).
					WithType(v.(map[string]interface{})["type"].(string)).
					Build()
				break
			} else if v.(map[string]interface{})["type"].(string) == "home" {
				mainPhone = model.PhoneNumberBuilder().
					WithValue(v.(map[string]interface{})["value"].(string)).
					WithType(v.(map[string]interface{})["type"].(string)).
					Build()
				break
			}
		}
	}

	// get the organizations
	var mainOrganization model.EnterpriseData
	var title string
	if m, ok := usr.Organizations.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["primary"] != nil && v.(map[string]interface{})["primary"].(bool) {

				mainOrganization = *model.EnterpriseDataBuilder().Build()

				if v.(map[string]interface{})["costCenter"] != nil {
					mainOrganization.CostCenter = v.(map[string]interface{})["costCenter"].(string)
				}

				if v.(map[string]interface{})["department"] != nil {
					mainOrganization.Department = v.(map[string]interface{})["department"].(string)
				}

				if v.(map[string]interface{})["division"] != nil {
					mainOrganization.Division = v.(map[string]interface{})["division"].(string)
				}

				if v.(map[string]interface{})["employeeNumber"] != nil {
					mainOrganization.EmployeeNumber = v.(map[string]interface{})["employeeNumber"].(string)
				}

				if v.(map[string]interface{})["name"] != nil {
					mainOrganization.Organization = v.(map[string]interface{})["name"].(string)
				}

				if v.(map[string]interface{})["primary"] != nil {
					mainOrganization.Primary = v.(map[string]interface{})["primary"].(bool)
				}

				if v.(map[string]interface{})["title"] != nil {
					title = v.(map[string]interface{})["title"].(string)
					mainOrganization.Title = title
				}

				var manager *model.Manager
				if v.(map[string]interface{})["manager"] != nil {
					manager = model.ManagerBuilder().
						WithValue("").
						Build()

					mainOrganization.Manager = manager
				}

				break
			}
		}
	}

	var displayName string
	if usr.Name.FullName != "" {
		displayName = usr.Name.FullName
	} else {
		displayName = fmt.Sprintf("%s %s", usr.Name.GivenName, usr.Name.FamilyName)
	}

	createdUser := model.UserBuilder().
		WithIPID(usr.Id).
		WithUserName(usr.PrimaryEmail).
		WithDisplayName(displayName).
		// WithNickName("Not Provided").
		// WithTitle("Not Provided").
		// WithTimezone("Not Provided").
		// WithProfileURL("Not Provided").
		WithUserType(usr.Kind).
		WithActive(!usr.Suspended).
		WithPreferredLanguage(preferredLanguage).
		WithTitle(title).
		Build()

	if emails != nil {
		createdUser.Emails = emails
	}

	if usr.Name != nil {
		createdUser.Name = model.NameBuilder().
			WithGivenName(usr.Name.GivenName).
			WithFamilyName(usr.Name.FamilyName).
			WithFormatted(usr.Name.FullName).
			Build()
	}

	if mainOrganization != (model.EnterpriseData{}) {
		createdUser.EnterpriseData = &mainOrganization
	}

	if mainAddress != (model.Address{}) {
		createdUser.Addresses = append(createdUser.Addresses, mainAddress)
	}

	if mainPhone != (model.PhoneNumber{}) {
		createdUser.PhoneNumbers = append(createdUser.PhoneNumbers, mainPhone)
	}

	// recalculate the hashcode because we have modified the user after building it
	createdUser.SetHashCode()

	log.WithFields(log.Fields{
		"object": convert.ToJSONString(usr),
	}).Trace("idp: building user from")

	log.WithFields(log.Fields{
		"object": convert.ToJSONString(createdUser),
	}).Trace("idp: building user to")

	return createdUser
}