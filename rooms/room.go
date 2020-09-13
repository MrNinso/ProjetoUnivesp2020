package rooms

type Room struct {
	id    string
	title string
}

func NewRoom(id, title string) Room {
	return Room{
		id:    id,
		title: title,
	}
}

func (r Room) getID() string {
	return r.id
}

func (r Room) getTitle() string {
	return r.title
}
