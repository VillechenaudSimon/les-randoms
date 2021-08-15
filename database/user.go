package database

import "database/sql"

type User struct {
	Id       int
	Name     string
	Password string
}

func User_SelectAll() ([]User, error) {
	rows, err := SelectDatabase("id, name, password FROM User")
	if err != nil {
		return nil, err
	}
	users := make([]User, 0)
	for rows.Next() {
		var id int
		var name string
		var password string
		err = rows.Scan(&id, &name, &password)
		if err != nil {
			return nil, err
		}
		users = append(users, User{Id: id, Name: name, Password: password})
	}
	return users, nil
}

func User_CreateNew(name string, password string) (User, sql.Result, error) {
	result, err := InsertDatabase("User VALUES(" + name + ", " + password + ")")
	if err != nil {
		return User{}, result, err
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return User{}, result, err
	}
	return User{Id: int(newId), Name: name, Password: password}, result, err
}
