package objets

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

func (r Room) GetID() string {
	return r.id
}

func (r Room) GetTitle() string {
	return r.title
}
