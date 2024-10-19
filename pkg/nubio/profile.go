package nubio

type Profile struct {
	Name        string       `json:"name"`
	Contact     Contact      `json:"contact"`
	Links       []Link       `json:"links"`
	Experiences []Experience `json:"experiences"`
	Skills      []Skill      `json:"skills"`
	Languages   []Language   `json:"languages"`
	Education   []Education  `json:"education"`
	Interests   []string     `json:"interests"`
	Hobbies     []string     `json:"hobbies"`

	// Note:
	//	- For SSG: This field is required.
	//	- For server: This field is overwritten by corresponding app config field.
	Domain string `json:"domain"`
}

type Contact struct {
	EmailAddress string `json:"email_address"`

	// Public PGP key URL (without leading "https://").
	// This field is overwritten on startup if a PGP key is provided in the app config.
	PGP string `json:"pgp"`
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
