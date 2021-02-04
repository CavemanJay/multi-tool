package music

type PlayList struct {
	Name string
	ID   string
}

func (p PlayList) Link() string {
	return "https://youtube.com/playlist?list=" + p.ID
}
