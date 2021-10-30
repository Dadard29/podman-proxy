package models

type DomainName struct {
	Name string `json:"name"`
	Live bool   `json:"live"`
}

func NewDomainName(scan func(dest ...interface{}) error) (DomainName, error) {
	var name string
	var liveInt int
	err := scan(&name, &liveInt)
	if err != nil {
		return DomainName{}, err
	}

	live := liveInt == 1

	return DomainName{
		Name: name,
		Live: live,
	}, nil
}
