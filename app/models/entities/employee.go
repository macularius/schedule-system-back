package entities

// Employee this struct emulated employee's structure
type Employee struct {
	ID        string
	Name      string
	GroupID   string
	GroupName string
	Days      []Day
	Templates []Template
}
