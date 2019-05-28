package entities

// Employee this struct emulated employee's structure
type Employee struct {
	eID         int
	Lastname    string
	Firstname   string
	Middlename  string
	Position    string
	Phonenumber string
	Email       string
	Days        []Day
	Templates   []Template
}
