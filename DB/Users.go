package DB

type userDB map[string]string

var _username, _password string = "admin", "admin"
var Users userDB = userDB{_username: _password}
