package nubio

// Holds the resume information that is actually public.
// This type definition is mainly needed for JSON exports,
// to "select" which fields are exported.
type ResumeExport struct {
	Slug           string           `json:"slug"`
	Name           string           `json:"name"`
	Domain         string           `json:"domain"`
	EmailAddress   string           `json:"email_address"`
	PGPKeyURL      string           `json:"pgp_key_url"`
	Links          []Link           `json:"links"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Skills         []Skill          `json:"skills"`
	Languages      []Language       `json:"languages"`
	Education      []Education      `json:"education"`
	Interests      []string         `json:"interests"`
	Hobbies        []string         `json:"hobbies"`
}

func (conf *ResumeConfig) ToResumeExport() *ResumeExport {
	return &ResumeExport{
		Slug:           conf.Slug,
		Name:           conf.Name,
		Domain:         conf.Domain,
		EmailAddress:   conf.EmailAddress,
		PGPKeyURL:      conf.PGPKeyURL,
		Links:          conf.Links,
		WorkExperience: conf.WorkExperience,
		Skills:         conf.Skills,
		Languages:      conf.Languages,
		Education:      conf.Education,
		Interests:      conf.Interests,
		Hobbies:        conf.Hobbies,
	}
}
