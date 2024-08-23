package model

import "errors"

type User struct {
	FirstName string
	LastName  string
	Biography string
}

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

func (v *User) Validate() error {
	if len(v.FirstName) < 2 || len(v.FirstName) > 20 {
		return errors.New("first_name must be between 2 and 20 characters")
	}
	if len(v.LastName) < 2 || len(v.LastName) > 20 {
		return errors.New("last_name must be between 2 and 20 characters")
	}
	if len(v.Biography) < 20 || len(v.Biography) > 450 {
		return errors.New("biography must be between 20 and 450 characters")
	}
	return nil
}

func (req *UserRequest) ToUser() User {
	return User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Biography: req.Biography,
	}
}
