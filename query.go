package smooth

type Query struct {
	Where      *[]Where
	Raw        *Raw
	Limit      *int
	Offset     *int
	Debug      bool
	Unscoped   bool
	With       *[]With
	InnerJoins *[]InnerJoins
	//Future

	// Group  *GroupBy
	// Select *Select
	// Model  *Model
	// Order  *[]OrderBy
	// Join       *[]Join
}

type OrderBy struct {
	Field interface{}
}

type Where struct {
	Column    string
	Condition string
	Value     any
}

type With struct {
	Field string
}

type InnerJoins struct {
	Field string
	Where *[]Where
}
type GroupBy struct {
	Field string
}

type Model struct {
	Interface interface{}
}
type Select struct {
	Arg string
}

type Join struct {
	Query string
}

type Raw struct {
	Query      string
	Interfaces []interface{}
}
