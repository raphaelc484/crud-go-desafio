package db

import (
	"github.com/google/uuid"
)

type user struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Biography string
}

type application struct {
	data map[uuid.UUID]user
}

func NewApplication() *application {
	return &application{
		data: make(map[uuid.UUID]user),
	}
}

func (app *application) Insert(firstName, lastName, biography string) user {
	id := uuid.New()

	newUser := user{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Biography: biography,
	}

	app.data[id] = newUser

	return newUser
}

func (app *application) FindAll() []user {
	users := make([]user, 0, len(app.data))
	for _, user := range app.data {
		users = append(users, user)
	}

	return users
}

func (app *application) FindById(id uuid.UUID) (user, bool) {
	user, found := app.data[id]

	return user, found
}

func (app *application) Update(id uuid.UUID, firstName, lastName, biography string) (user, bool) {
	userExist, found := app.FindById(id)
	if !found {
		return user{}, found
	}

	userExist.FirstName = firstName
	userExist.LastName = lastName
	userExist.Biography = biography

	app.data[id] = userExist

	return userExist, found
}

func (app *application) Delete(id uuid.UUID) bool {
	_, found := app.FindById(id)
	if !found {
		return found
	}

	delete(app.data, id)

	return true
}
