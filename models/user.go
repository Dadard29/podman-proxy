package models

type User struct {
	Name           string `json:"user"`
	HashedPassword string `json:"hashed_password"`
}

func NewUser(scan func(dest ...interface{}) error) (User, error) {
	var name string
	var hashedPassword string
	err := scan(&name, &hashedPassword)
	if err != nil {
		return User{}, err
	}

	return User{
		Name:           name,
		HashedPassword: hashedPassword,
	}, nil
}
