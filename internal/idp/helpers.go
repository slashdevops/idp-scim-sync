package idp

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/convert"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	admin "google.golang.org/api/admin/directory/v1"
)

// buildUser builds a User model from a User coming from the IDP API
func buildUser(usr *admin.User) *model.User {
	if usr == nil {
		return nil
	}

	// these fields are required because the Constrains defined here:
	// https://docs.aws.amazon.com/singlesignon/latest/developerguide/createuser.html
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
					emails = append(emails,
						model.EmailBuilder().
							WithPrimary(v.(map[string]interface{})["primary"].(bool)).
							Build(),
					)

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
					Build()
				break
			} else if v.(map[string]interface{})["type"].(string) == "home" {
				mainAddress = model.AddressBuilder().
					WithFormatted(v.(map[string]interface{})["formatted"].(string)).
					WithType(v.(map[string]interface{})["type"].(string)).
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

	userModel := model.UserBuilder().
		WithIPID(usr.Id).
		WithUserName(usr.PrimaryEmail).
		WithDisplayName(displayName).
		// WithNickName("Not Provided").
		// WithProfileURL("Not Provided").
		WithTitle(title).
		WithUserType(usr.Kind).
		WithPreferredLanguage(preferredLanguage).
		// WithLocale("Not Provided").
		// WithTimezone("Not Provided").
		WithActive(!usr.Suspended).
		Build()

	if emails != nil {
		userModel.Emails = emails
	}

	if mainAddress != (model.Address{}) {
		userModel.Addresses = append(userModel.Addresses, mainAddress)
	}

	if mainPhone != (model.PhoneNumber{}) {
		userModel.PhoneNumbers = append(userModel.PhoneNumbers, mainPhone)
	}

	if usr.Name != nil {
		userModel.Name = model.NameBuilder().
			WithGivenName(usr.Name.GivenName).
			WithFamilyName(usr.Name.FamilyName).
			WithFormatted(usr.Name.FullName).
			Build()
	}

	if mainOrganization != (model.EnterpriseData{}) {
		userModel.EnterpriseData = &mainOrganization
	}

	// recalculate the hashcode because we have modified the user after building it
	userModel.SetHashCode()

	log.Tracef("idp: buildUser() from: %+v, --> to: %+v", convert.ToJSONString(usr), convert.ToJSONString(userModel))

	return userModel
}
