package aws

type SCIMService interface {
}

type scim struct {
}

func NewSCIMService() (SCIMService, error) {

	return &scim{}, nil
}
