package models

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	PassHash  string `json:"-"`
}

func (u *User) Table() string {
	return "users"
}

func (u *User) Insert() (columns []string, values []interface{}) {
	columns = []string{"id", "username", "firstname", "lastname", "passhash"}
	values = []interface{}{u.ID, u.Username, u.FirstName, u.LastName, u.PassHash}
	return
}

func (u *User) Update() (columns []string, values []interface{}) {
	columns = []string{"firstname", "lastname"}
	values = []interface{}{u.FirstName, u.LastName}
	return
}

func (u *User) Select() (columns []string, into []interface{}) {
	columns = []string{"id", "username", "firstname", "lastname", "passhash"}
	into = []interface{}{&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.PassHash}
	return
}

func (u *User) KeyColumns() []string     { return []string{"id"} }
func (u *User) KeyValues() []interface{} { return []interface{}{u.ID} }
