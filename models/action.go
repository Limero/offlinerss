package models

type SyncToActions struct {
	Read      []string
	Unread    []string
	Starred   []string
	Unstarred []string
}

func (actions SyncToActions) Total() int {
	return len(actions.Read) + len(actions.Unread) + len(actions.Starred) + len(actions.Unstarred)
}
