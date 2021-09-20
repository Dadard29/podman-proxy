package models

type DomainName struct {
	Name string `json:"name"`
}

func NewDomainName(scan func(dest ...interface{}) error) (DomainName, error) {
	var name string
	err := scan(&name)
	if err != nil {
		return DomainName{}, err
	}

	return DomainName{
		Name: name,
	}, nil
}
