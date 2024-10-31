package idp

import (
	"fmt"
	"log/slog"
	"strings"

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
		slog.Warn("idp: User name is nil")
		return nil
	}

	if usr.Name.GivenName == "" {
		slog.Warn("idp: User given name is empty")
		return nil
	}

	if usr.Name.FamilyName == "" {
		slog.Warn("idp: User family name is empty")
		return nil
	}

	if usr.PrimaryEmail == "" {
		slog.Warn("idp: User primary email is empty")
		return nil
	}

	var emails []model.Email
	if m, ok := usr.Emails.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["primary"] != nil {
				if v.(map[string]interface{})["primary"].(bool) {

					var emailValue string
					if v.(map[string]interface{})["address"] != nil {
						emailValue = strings.TrimSpace(v.(map[string]interface{})["address"].(string))
					}

					var emailType string
					if v.(map[string]interface{})["type"] != nil {
						emailType = strings.TrimSpace(v.(map[string]interface{})["type"].(string))
					}

					emails = append(emails,
						model.EmailBuilder().
							WithPrimary(v.(map[string]interface{})["primary"].(bool)).
							WithType(emailType).
							WithValue(emailValue).
							Build(),
					)

					break
				}
			}
		}
	}

	if len(emails) == 0 {
		slog.Warn("idp: User emails is empty, setting primary email as the only email")
		emails = append(emails,
			model.EmailBuilder().
				WithPrimary(true).
				WithType("work").
				WithValue(strings.TrimSpace(usr.PrimaryEmail)).
				Build(),
		)
	}

	// get the first language from the list of languages
	var preferredLanguage string
	if m, ok := usr.Languages.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["preference"].(string) == "preferred" {
				preferredLanguage = strings.TrimSpace(v.(map[string]interface{})["languageCode"].(string))
				break
			}
		}
	}

	// get the Addresses
	var addresses []model.Address
	if m, ok := usr.Addresses.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["type"].(string) == "work" {
				addresses = append(addresses,
					model.AddressBuilder().
						WithFormatted(strings.TrimSpace(v.(map[string]interface{})["formatted"].(string))).
						Build())
				break
			} else if v.(map[string]interface{})["type"].(string) == "home" {
				addresses = append(addresses,
					model.AddressBuilder().
						WithFormatted(strings.TrimSpace(v.(map[string]interface{})["formatted"].(string))).
						Build())
				break
			}
		}
	}

	// get the phones
	var phoneNumbers []model.PhoneNumber
	if m, ok := usr.Phones.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["type"].(string) == "work" {
				phoneNumbers = append(phoneNumbers,
					model.PhoneNumberBuilder().
						WithValue(strings.TrimSpace(v.(map[string]interface{})["value"].(string))).
						WithType(strings.TrimSpace(v.(map[string]interface{})["type"].(string))).
						Build())
				break
			} else if v.(map[string]interface{})["type"].(string) == "home" {
				phoneNumbers = append(phoneNumbers,
					model.PhoneNumberBuilder().
						WithValue(strings.TrimSpace(v.(map[string]interface{})["value"].(string))).
						WithType(strings.TrimSpace(v.(map[string]interface{})["type"].(string))).
						Build())
				break
			}
		}
	}

	// get the relations (manager)
	var manager *model.Manager
	if m, ok := usr.Relations.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["type"].(string) == "manager" {
				manager = model.ManagerBuilder().
					WithValue(strings.TrimSpace(v.(map[string]interface{})["value"].(string))).
					WithRef(strings.TrimSpace(v.(map[string]interface{})["customType"].(string))).
					Build()
				break
			}
		}
	}

	// get the organizations
	var mainOrganization *model.EnterpriseData
	var title string
	if m, ok := usr.Organizations.([]interface{}); ok {
		for _, v := range m {
			if v.(map[string]interface{})["primary"] != nil && v.(map[string]interface{})["primary"].(bool) {

				var employeeNumber string
				if v.(map[string]interface{})["employeeNumber"] != nil {
					employeeNumber = strings.TrimSpace(v.(map[string]interface{})["employeeNumber"].(string))
				}

				var costCenter string
				if v.(map[string]interface{})["costCenter"] != nil {
					costCenter = strings.TrimSpace(v.(map[string]interface{})["costCenter"].(string))
				}

				var organization string
				if v.(map[string]interface{})["name"] != nil {
					organization = strings.TrimSpace(v.(map[string]interface{})["name"].(string))
				}

				var division string
				if v.(map[string]interface{})["division"] != nil {
					division = strings.TrimSpace(v.(map[string]interface{})["division"].(string))
				}

				var department string
				if v.(map[string]interface{})["department"] != nil {
					department = strings.TrimSpace(v.(map[string]interface{})["department"].(string))
				}

				if v.(map[string]interface{})["title"] != nil {
					title = strings.TrimSpace(v.(map[string]interface{})["title"].(string))
				}

				mainOrganization = model.EnterpriseDataBuilder().
					WithEmployeeNumber(employeeNumber).
					WithCostCenter(costCenter).
					WithOrganization(organization).
					WithDivision(division).
					WithDepartment(department).
					WithManager(manager).
					Build()
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

	var name *model.Name
	if usr.Name != nil {
		name = model.NameBuilder().
			WithGivenName(strings.TrimSpace(usr.Name.GivenName)).
			WithFamilyName(strings.TrimSpace(usr.Name.FamilyName)).
			WithFormatted(strings.TrimSpace(usr.Name.FullName)).
			Build()
	}

	userModel := model.UserBuilder().
		WithIPID(strings.TrimSpace(usr.Id)).
		WithUserName(strings.TrimSpace(usr.PrimaryEmail)).
		WithDisplayName(strings.TrimSpace(displayName)).
		// WithNickName("Not Provided").
		// WithProfileURL("Not Provided").
		WithTitle(title).
		WithUserType(strings.TrimSpace(usr.Kind)).
		WithPreferredLanguage(preferredLanguage).
		// WithLocale("Not Provided").
		// WithTimezone("Not Provided").
		WithActive(!usr.Suspended).
		// Arrays
		WithEmails(emails).
		WithAddresses(addresses).
		WithPhoneNumbers(phoneNumbers).
		// Pointers
		WithName(name).
		WithEnterpriseData(mainOrganization).
		Build()

	slog.Debug("idp: buildUser() converted user", "from", usr, "to", userModel)

	return userModel
}
