package model

type Choice struct {
	id    int
	title string
	desc  string
}

func (i Choice) ID() int             { return i.id }
func (i Choice) Title() string       { return i.title }
func (i Choice) Description() string { return i.desc }
func (i Choice) FilterValue() string { return i.title }
