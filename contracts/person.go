package contracts

type PersonForCreate struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type PersonForUpdate struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

type PersonForFetch struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}
