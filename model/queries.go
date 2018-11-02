package model

const (
	ADD_USER    = "INSERT INTO USERS(uuid, login, created) VALUES($1, $2, NOW())"
	GET_USER    = "SELECT * FROM USERS WHERE "
	UPDATE_USER = "UPDATE USERS "
)
