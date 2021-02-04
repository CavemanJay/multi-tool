package music

type Video struct {
	ID    string
	Title string
}

func (v Video) Link() string {
	return "https://youtube.com/watch?v=" + v.ID
}
