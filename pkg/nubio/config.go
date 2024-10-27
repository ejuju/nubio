package nubio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"
	"unicode/utf8"

	"github.com/ejuju/nubio/pkg/httpmux"
)

type Config struct {
	Address         string  `json:"address"`        // Local HTTP server address.
	TrueIPHeader    string  `json:"true_ip_header"` // Optional: (ex: "X-Forwarded-For", use when reverse proxying).
	TLSDirpath      string  `json:"tls_dirpath"`    // Path to TLS certificate directory.
	TLSEmailAddress string  `json:"tls_email_addr"` // Email address in TLS certificate.
	Resume          *Resume `json:"resume"`         // Resume content.
}

// Publicly exportable user resume info.
type Resume struct {
	Slug string `json:"slug"` // Optional: Name as a URI-compatible slug (ex: "alex-doe").
	Name string `json:"name"` // Full name (ex: "Alex Doe").

	// Public domain name (ex: "alexdoe.example")
	// Note:
	//	- For SSG: This field is required.
	//	- For server: This field is overwritten by corresponding app config field.
	Domain string `json:"domain"`

	PGPKeyPath    string `json:"pgp_key_path"`    // Path to PGP public key. Not exported.
	PGPKey        string `json:"pgp_key"`         // Literal value or populated by the corresponding file's content on load.
	CustomCSSPath string `json:"custom_css_path"` // Path to custom CSS stylesheet. Not exported.
	CustomCSS     string `json:"custom_css"`      // Literal value or populated by the corresponding file's content on load.

	Contact        Contact          `json:"contact"`
	Links          []Link           `json:"links"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Skills         []Skill          `json:"skills"`
	Languages      []Language       `json:"languages"`
	Education      []Education      `json:"education"`
	Interests      []string         `json:"interests"`
	Hobbies        []string         `json:"hobbies"`
}

func LoadConfig(path string) (conf *Config, err error) {
	// Read whole file into memory.
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read whole file: %w", err)
	}

	// Decode from JSON.
	conf = &Config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	// Create name slug if none is provided.
	if conf.Resume.Slug == "" {
		conf.Resume.Slug = httpmux.Slugify(conf.Resume.Name)
	}

	// Load PGP key if provided.
	if conf.Resume.PGPKeyPath != "" {
		b, err = os.ReadFile(conf.Resume.PGPKeyPath)
		if err != nil {
			return nil, fmt.Errorf("load PGP key: %w", err)
		}
		conf.Resume.PGPKey = string(b)
	}

	// Load custom CSS if provided.
	if conf.Resume.CustomCSSPath != "" {
		b, err = os.ReadFile(conf.Resume.CustomCSSPath)
		if err != nil {
			return nil, fmt.Errorf("load custom CSS: %w", err)
		}
		conf.Resume.CustomCSS = string(b)
	}

	return conf, nil
}

func (conf *Config) Check() (errs []error) {
	// Check HTTPS-related fields.
	if conf.TLSDirpath != "" && conf.TLSEmailAddress == "" {
		errs = append(errs, errors.New("missing TLS email address"))
	}

	if conf.Resume == nil {
		errs = append(errs, errors.New("missing resume"))
	} else {
		errs = append(errs, conf.Resume.Check()...)
	}

	return errs
}

func (p *Resume) Check() (errs []error) {
	// Check name and domain.
	if p.Name == "" {
		errs = append(errs, errors.New("missing name"))
	} else if nameSize := utf8.RuneCountInString(p.Name); nameSize > 100 {
		errs = append(errs, fmt.Errorf("name is too big: %d characters", nameSize))
	}
	if p.Domain == "" {
		errs = append(errs, errors.New("missing domain"))
	}

	// Check contact info.
	for _, err := range p.Contact.Check() {
		errs = append(errs, fmt.Errorf("contact: %w", err))
	}

	// Check links.
	if len(p.Links) == 0 {
		errs = append(errs, errors.New("missing links"))
	}
	for i, v := range p.Links {
		for _, err := range v.Check() {
			errs = append(errs, fmt.Errorf("link %d: %w", i, err))
		}
	}

	// Check experiences.
	if len(p.WorkExperience) == 0 {
		errs = append(errs, errors.New("missing experiences"))
	}
	for i, v := range p.WorkExperience {
		for _, err := range v.Check() {
			errs = append(errs, fmt.Errorf("experience %d: %w", i, err))
		}
	}

	// Check skills.
	if len(p.Skills) == 0 {
		errs = append(errs, errors.New("missing skills"))
	}
	for i, v := range p.Skills {
		for _, err := range v.Check() {
			errs = append(errs, fmt.Errorf("skill %d: %w", i, err))
		}
	}

	// Check languages.
	if len(p.Languages) == 0 {
		errs = append(errs, errors.New("missing languages"))
	}
	for i, v := range p.Languages {
		for _, err := range v.Check() {
			errs = append(errs, fmt.Errorf("language %d: %w", i, err))
		}
	}

	// Check education.
	if len(p.Education) == 0 {
		errs = append(errs, errors.New("missing education"))
	}
	for i, v := range p.Education {
		for _, err := range v.Check() {
			errs = append(errs, fmt.Errorf("education %d: %w", i, err))
		}
	}

	// Check interests.
	if len(p.Interests) == 0 {
		errs = append(errs, errors.New("missing interests"))
	}
	for i, v := range p.Interests {
		for v == "" {
			errs = append(errs, fmt.Errorf("interest %d: empty text", i))
		}
	}

	// Check hobbies.
	if len(p.Hobbies) == 0 {
		errs = append(errs, errors.New("missing hobbies"))
	}
	for i, v := range p.Hobbies {
		for v == "" {
			errs = append(errs, fmt.Errorf("hobby %d: empty text", i))
		}
	}

	return errs
}

type Contact struct {
	EmailAddress string `json:"email_address"`

	// Public PGP key URL (without leading "https://").
	// This field is overwritten on startup if a PGP key is provided in the app config.
	PGP string `json:"pgp"`
}

func (v *Contact) Check() (errs []error) {
	if v.EmailAddress == "" {
		errs = append(errs, errors.New("missing email address"))
	}
	if v.PGP != "" {
		u, err := url.Parse(v.PGP)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid PGP key URL"))
		} else if u.Scheme == "" {
			errs = append(errs, errors.New("missing scheme in PGP URL"))
		}
	}
	return errs
}

type WorkExperience struct {
	From         string   `json:"from"`
	To           string   `json:"to"`
	Title        string   `json:"title"`
	Organization string   `json:"organization"`
	Location     string   `json:"location"`
	Description  string   `json:"description"`
	Skills       []string `json:"skills"`
}

const DateLayout = "January 2006"

var (
	minExpDate = time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC)
	maxExpDate = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
)

func (v *WorkExperience) Check() (errs []error) {
	if v.From == "" {
		errs = append(errs, errors.New("missing start date"))
	}
	from, err := parseDateMinMax(DateLayout, v.From, minExpDate, maxExpDate)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid start date: %w", err))
	}
	if v.To == "" {
		errs = append(errs, errors.New("missing end date"))
	}
	_, err = parseDateMinMax(DateLayout, v.To, from, maxExpDate)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid end date: %w", err))
	}

	if v.Title == "" {
		errs = append(errs, errors.New("missing title"))
	}
	if v.Organization == "" {
		errs = append(errs, errors.New("missing organization"))
	}
	if v.Location == "" {
		errs = append(errs, errors.New("missing location"))
	}
	if v.Description == "" {
		errs = append(errs, errors.New("missing description"))
	}
	if len(v.Skills) == 0 {
		errs = append(errs, errors.New("missing skills"))
	}

	return errs
}

type Skill struct {
	Title string   `json:"title"`
	Tools []string `json:"tools"`
}

func (v *Skill) Check() (errs []error) {
	if v.Title == "" {
		errs = append(errs, errors.New("missing title"))
	}
	if len(v.Tools) == 0 {
		errs = append(errs, errors.New("missing tools"))
	}
	for i, v := range v.Tools {
		if v == "" {
			errs = append(errs, fmt.Errorf("empty text at index %d", i))
		}
	}
	return errs
}

type Education struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Title        string `json:"title"`
	Organization string `json:"organization"`
}

func (v *Education) Check() (errs []error) {
	if v.From == "" {
		errs = append(errs, errors.New("missing start date"))
	}
	from, err := parseDateMinMax(DateLayout, v.From, minExpDate, maxExpDate)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid start date: %w", err))
	}
	if v.To == "" {
		errs = append(errs, errors.New("missing end date"))
	}
	_, err = parseDateMinMax(DateLayout, v.To, from, maxExpDate)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid end date: %w", err))
	}

	if v.Title == "" {
		errs = append(errs, errors.New("missing title"))
	}
	if v.Organization == "" {
		errs = append(errs, errors.New("missing organization"))
	}
	return errs
}

type Language struct {
	Label       string `json:"label"`
	Proficiency string `json:"proficiency"`
}

func (v *Language) Check() (errs []error) {
	if v.Label == "" {
		errs = append(errs, errors.New("missing label"))
	}
	if v.Proficiency == "" {
		errs = append(errs, errors.New("missing proficiency"))
	}
	return errs
}

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

func (v *Link) Check() (errs []error) {
	if v.Label == "" {
		errs = append(errs, errors.New("missing label"))
	}
	if v.URL == "" {
		errs = append(errs, errors.New("missing URL"))
	} else if _, err := url.Parse(v.URL); err != nil {
		errs = append(errs, fmt.Errorf("invalid URL: %w", err))
	}
	return errs
}

var (
	errMissingDate  = errors.New("missing date")
	errDateTooEarly = errors.New("date is too early")
	errDateTooLate  = errors.New("date is too late")
)

// Special case: raw == "now" is accepted.
func parseDateMinMax(layout, raw string, min, max time.Time) (t time.Time, err error) {
	if raw == "" {
		return t, errMissingDate
	}

	// Parse date.
	if raw == "now" {
		t = time.Now()
	} else {
		t, err = time.Parse(layout, raw)
		if err != nil {
			return t, err
		}
	}

	// Check min/max.
	if t.Before(min) {
		return t, errDateTooEarly
	} else if t.After(max) {
		return t, errDateTooLate
	}

	return t, nil
}
