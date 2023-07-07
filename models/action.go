package models

const (
	ActionStoryRead      = 1
	ActionStoryUnread    = 2
	ActionStoryStarred   = 3
	ActionStoryUnstarred = 4
)

type SyncToAction struct {
	Id     string
	Action int
}

type SyncToActions []SyncToAction

func (actions SyncToActions) SumActionTypes() (read int, unread int, starred int, unstarred int) {
	for _, action := range actions {
		switch action.Action {
		case ActionStoryRead:
			read++
		case ActionStoryUnread:
			unread++
		case ActionStoryStarred:
			starred++
		case ActionStoryUnstarred:
			unstarred++
		}
	}
	return read, unread, starred, unstarred
}
