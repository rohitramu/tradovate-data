package db

type ColumnCollection []*Column

type Column struct {
	name         string
	dataType     Type
	isPrimaryKey bool
}

func NewColumn(name string, dataType Type) *Column {
	return &Column{
		name:         name,
		dataType:     dataType,
		isPrimaryKey: false,
	}
}

func NewPrimaryKeyColumn(name string, dataType Type) *Column {
	return &Column{
		name:         name,
		dataType:     dataType,
		isPrimaryKey: true,
	}
}
