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

	fromEmails := make([]*admin.UserEmail, 0, len(emails))
	for _, email := range emails {
		fromEmails = append(fromEmails, email.(*admin.UserEmail))
	}

	modelEmails := make([]model.Email, 0, len(fromEmails))
	for _, v := range fromEmails {
		if v.Primary {
			modelEmails = append(modelEmails,
				model.EmailBuilder().
					WithPrimary(v.Primary).
					WithType(v.Type).
					WithValue(v.Address).
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

	fromLanguages := make([]*admin.UserLanguage, 0, len(languages))
	for _, language := range languages {
		fromLanguages = append(fromLanguages, language.(*admin.UserLanguage))
	}

	var preferredLanguage string
	for _, v := range fromLanguages {
		if v.Preference == "preferred" {
			preferredLanguage = v.LanguageCode
			break
		}
	}
	return preferredLanguage, nil
}

func toAddresses(a any) ([]model.Address, error) {
	addresses, ok := a.([]any)
	if !ok {
		return nil, fmt.Errorf("error converting addresses: %v", a)
	}

	fromAddresses := make([]*admin.UserAddress, 0, len(addresses))
	for _, address := range addresses {
		fromAddresses = append(fromAddresses, address.(*admin.UserAddress))
	}

	modelAddresses := make([]model.Address, 0, len(fromAddresses))
	for _, v := range fromAddresses {
		if v.Type == "work" {
			modelAddresses = append(modelAddresses,
				model.AddressBuilder().
					WithFormatted(v.Formatted).
					Build())
			break
		} else if v.Type == "home" {
			modelAddresses = append(modelAddresses,
				model.AddressBuilder().
					WithFormatted(v.Formatted).
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

	fromPhones := make([]*admin.UserPhone, 0, len(phones))
	for _, phone := range phones {
		fromPhones = append(fromPhones, phone.(*admin.UserPhone))
	}

	modelPhoneNumbers := make([]model.PhoneNumber, 0, len(fromPhones))
	for _, v := range fromPhones {
		if v.Type == "work" {
			modelPhoneNumbers = append(modelPhoneNumbers,
				model.PhoneNumberBuilder().
					WithValue(v.Value).
					WithType(v.Type).
					Build())
			break
		} else if v.Type == "home" {
			modelPhoneNumbers = append(modelPhoneNumbers,
				model.PhoneNumberBuilder().
					WithValue(v.Value).
					WithType(v.Type).
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

	fromRelations := make([]*admin.UserRelation, 0, len(relations))
	for _, relation := range relations {
		fromRelations = append(fromRelations, relation.(*admin.UserRelation))
	}

	var manager *model.Manager
	for _, v := range fromRelations {
		if v.Type == "manager" {
			manager = model.ManagerBuilder().
				WithValue(v.Value).
				WithRef(v.CustomType).
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

	fromOrganizations := make([]*admin.UserOrganization, 0, len(organizations))
	for _, organization := range organizations {
		fromOrganizations = append(fromOrganizations, organization.(*admin.UserOrganization))
	}

	var mainOrganization *model.EnterpriseData
	var title string
	for _, v := range fromOrganizations {
		if v.Primary {
			mainOrganization = model.EnterpriseDataBuilder().
				WithCostCenter(v.CostCenter).
				WithOrganization(v.Name).
				WithDivision(v.Domain).
				WithDepartment(v.Department).
				WithManager(manager).
				Build()
			title = v.Title
			break
		}
	}
	return mainOrganization, title, nil
}
