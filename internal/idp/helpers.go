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
	if usr.Emails != nil {
		e, err := toEmails(usr.Emails)
		if err != nil {
			slog.Error("idp: error converting emails", "error", err)
		} else {
			emails = e
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
	if usr.Languages != nil {
		l, err := toLanguages(usr.Languages)
		if err != nil {
			slog.Error("idp: error converting languages", "error", err)
		} else {
			preferredLanguage = l
		}
	}

	// get the Addresses
	var addresses []model.Address
	if usr.Addresses != nil {
		a, err := toAddresses(usr.Addresses)
		if err != nil {
			slog.Error("idp: error converting addresses", "error", err)
		} else {
			addresses = a
		}
	}

	// get the phones
	var phoneNumbers []model.PhoneNumber
	if usr.Phones != nil {
		p, err := toPhones(usr.Phones)
		if err != nil {
			slog.Error("idp: error converting phones", "error", err)
		} else {
			phoneNumbers = p
		}
	}

	// get the relations (manager)
	var manager *model.Manager
	if usr.Relations != nil {
		m, err := toRelations(usr.Relations)
		if err != nil {
			slog.Error("idp: error converting relations", "error", err)
		} else {
			manager = m
		}
	}

	// get the organizations
	var mainOrganization *model.EnterpriseData
	var title string
	if usr.Organizations != nil {
		o, t, err := toOrganizations(usr.Organizations, manager)
		if err != nil {
			slog.Error("idp: error converting organizations", "error", err)
		} else {
			mainOrganization = o
			title = t
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
		WithTitle(title).
		WithUserType(strings.TrimSpace(usr.Kind)).
		WithPreferredLanguage(preferredLanguage).
		WithActive(!usr.Suspended).
		WithEmails(emails).
		WithAddresses(addresses).
		WithPhoneNumbers(phoneNumbers).
		WithName(name).
		WithEnterpriseData(mainOrganization).
		Build()

	slog.Debug("idp: buildUser() converted user", "from", usr, "to", userModel)

	return userModel
}

func toEmails(e any) ([]model.Email, error) {
	emails, ok := e.([]any)
	if !ok {
		return nil, fmt.Errorf("error converting emails: %v", e)
	}

	modelEmails := make([]model.Email, 0, len(emails))
	for _, email := range emails {
		emailMap, ok := email.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("error converting email: %v", email)
		}

		var primary bool
		if p, ok := emailMap["primary"].(bool); ok {
			primary = p
		}

		if primary {
			var emailType, address string
			if et, ok := emailMap["type"].(string); ok {
				emailType = et
			}
			if addr, ok := emailMap["address"].(string); ok {
				address = addr
			}

			modelEmails = append(modelEmails,
				model.EmailBuilder().
					WithPrimary(primary).
					WithType(emailType).
					WithValue(address).
					Build(),
			)
			break
		}
	}
	return modelEmails, nil
}

func toLanguages(l any) (string, error) {
	languages, ok := l.([]any)
	if !ok {
		return "", fmt.Errorf("error converting languages: %v", l)
	}

	var preferredLanguage string
	for _, language := range languages {
		languageMap, ok := language.(map[string]interface{})
		if !ok {
			// try to convert to *admin.UserLanguage
			if lang, ok := language.(*admin.UserLanguage); ok {
				if lang.Preference == "preferred" {
					preferredLanguage = lang.LanguageCode
					break
				}
				continue
			}
			return "", fmt.Errorf("error converting language: %v", language)
		}

		var preference string
		if p, ok := languageMap["preference"].(string); ok {
			preference = p
		}

		if preference == "preferred" {
			if lc, ok := languageMap["languageCode"].(string); ok {
				preferredLanguage = lc
				break
			}
		}
	}
	return preferredLanguage, nil
}

func toAddresses(a any) ([]model.Address, error) {
	addresses, ok := a.([]any)
	if !ok {
		return nil, fmt.Errorf("error converting addresses: %v", a)
	}

	modelAddresses := make([]model.Address, 0, len(addresses))
	for _, address := range addresses {
		addressMap, ok := address.(map[string]interface{})
		if !ok {
			// try to convert to *admin.UserAddress
			if addr, ok := address.(*admin.UserAddress); ok {
				if addr.Type == "work" || addr.Type == "home" {
					modelAddresses = append(modelAddresses,
						model.AddressBuilder().
							WithFormatted(addr.Formatted).
							Build())
					break
				}
				continue
			}
			return nil, fmt.Errorf("error converting address: %v", address)
		}

		var addressType, formatted string
		if at, ok := addressMap["type"].(string); ok {
			addressType = at
		}
		if f, ok := addressMap["formatted"].(string); ok {
			formatted = f
		}

		if addressType == "work" || addressType == "home" {
			modelAddresses = append(modelAddresses,
				model.AddressBuilder().
					WithFormatted(formatted).
					Build())
			break
		}
	}
	return modelAddresses, nil
}

func toPhones(p any) ([]model.PhoneNumber, error) {
	phones, ok := p.([]any)
	if !ok {
		return nil, fmt.Errorf("error converting phones: %v", p)
	}

	modelPhoneNumbers := make([]model.PhoneNumber, 0, len(phones))
	for _, phone := range phones {
		phoneMap, ok := phone.(map[string]interface{})
		if !ok {
			// try to convert to *admin.UserPhone
			if ph, ok := phone.(*admin.UserPhone); ok {
				if ph.Type == "work" || ph.Type == "home" {
					modelPhoneNumbers = append(modelPhoneNumbers,
						model.PhoneNumberBuilder().
							WithValue(ph.Value).
							WithType(ph.Type).
							Build())
					break
				}
				continue
			}
			return nil, fmt.Errorf("error converting phone: %v", phone)
		}

		var phoneType, value string
		if pt, ok := phoneMap["type"].(string); ok {
			phoneType = pt
		}
		if v, ok := phoneMap["value"].(string); ok {
			value = v
		}

		if phoneType == "work" || phoneType == "home" {
			modelPhoneNumbers = append(modelPhoneNumbers,
				model.PhoneNumberBuilder().
					WithValue(value).
					WithType(phoneType).
					Build())
			break
		}
	}
	return modelPhoneNumbers, nil
}

func toRelations(r any) (*model.Manager, error) {
	relations, ok := r.([]any)
	if !ok {
		return nil, fmt.Errorf("error converting relations: %v", r)
	}

	var manager *model.Manager
	for _, relation := range relations {
		relationMap, ok := relation.(map[string]interface{})
		if !ok {
			// try to convert to *admin.UserRelation
			if rel, ok := relation.(*admin.UserRelation); ok {
				if rel.Type == "manager" {
					manager = model.ManagerBuilder().
						WithValue(rel.Value).
						WithRef(rel.CustomType).
						Build()
					break
				}
				continue
			}
			return nil, fmt.Errorf("error converting relation: %v", relation)
		}

		var relationType, value, customType string
		if rt, ok := relationMap["type"].(string); ok {
			relationType = rt
		}
		if v, ok := relationMap["value"].(string); ok {
			value = v
		}
		if ct, ok := relationMap["customType"].(string); ok {
			customType = ct
		}

		if relationType == "manager" {
			manager = model.ManagerBuilder().
				WithValue(value).
				WithRef(customType).
				Build()
			break
		}
	}
	return manager, nil
}

func toOrganizations(o any, manager *model.Manager) (*model.EnterpriseData, string, error) {
	organizations, ok := o.([]any)
	if !ok {
		return nil, "", fmt.Errorf("error converting organizations: %v", o)
	}

	var mainOrganization *model.EnterpriseData
	var title string
	for _, organization := range organizations {
		organizationMap, ok := organization.(map[string]interface{})
		if !ok {
			// try to convert to *admin.UserOrganization
			if org, ok := organization.(*admin.UserOrganization); ok {
				if org.Primary {
					mainOrganization = model.EnterpriseDataBuilder().
						WithCostCenter(org.CostCenter).
						WithOrganization(org.Name).
						WithDivision(org.Domain).
						WithDepartment(org.Department).
						WithManager(manager).
						Build()
					title = org.Title
					break
				}
				continue
			}
			return nil, "", fmt.Errorf("error converting organization: %v", organization)
		}

		var primary bool
		if p, ok := organizationMap["primary"].(bool); ok {
			primary = p
		}

		if primary {
			var costCenter, name, domain, department, orgTitle string
			if cc, ok := organizationMap["costCenter"].(string); ok {
				costCenter = cc
			}
			if n, ok := organizationMap["name"].(string); ok {
				name = n
			}
			if d, ok := organizationMap["domain"].(string); ok {
				domain = d
			}
			if dep, ok := organizationMap["department"].(string); ok {
				department = dep
			}
			if t, ok := organizationMap["title"].(string); ok {
				orgTitle = t
			}

			mainOrganization = model.EnterpriseDataBuilder().
				WithCostCenter(costCenter).
				WithOrganization(name).
				WithDivision(domain).
				WithDepartment(department).
				WithManager(manager).
				Build()
			title = orgTitle
			break
		}
	}
	return mainOrganization, title, nil
}
