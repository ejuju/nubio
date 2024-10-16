package nuage

type Profile struct {
	Name        string       `json:"name"`
	Experiences []Experience `json:"experiences"`
	Skills      []Skill      `json:"skills"`
	Languages   []Language   `json:"languages"`
	Education   []Education  `json:"education"`
	Interests   []string     `json:"interests"`
	Hobbies     []string     `json:"hobbies"`
	Links       []Link       `json:"links"`
	Contact     Contact      `json:"contact"`
}

type Experience struct {
	From         string   `json:"from"`
	To           string   `json:"to"`
	Title        string   `json:"title"`
	Organization string   `json:"organization"`
	Location     string   `json:"location"`
	Description  string   `json:"description"`
	Skills       []string `json:"skills"`
}

type Skill struct {
	Title string   `json:"title"`
	Tools []string `json:"tools"`
}

type Education struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Title        string `json:"title"`
	Organization string `json:"organization"`
}

type Language struct {
	Label       string `json:"label"`
	Proficiency string `json:"proficiency"`
}

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Contact struct {
	EmailAddress string `json:"email_address"`
	URL          string `json:"url"`
}
