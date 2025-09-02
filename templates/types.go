package templates

// Define structs for data
type ExperienceItem struct {
	Title       string   `json:"Title"`
	Company     string   `json:"Company"`
	Period      string   `json:"Period"`
	Description []string `json:"Description"`
}

type ExperienceData struct {
	ExperienceItems []ExperienceItem `json:"ExperienceItems"`
	Language        string           `json:"-"`
	Translations    map[string]map[string]string `json:"-"`
}

type EducationItem struct {
	Title       string `json:"Title"`
	Institution string `json:"Institution"`
	Period      string `json:"Period"`
}

type EducationData struct {
	EducationItems []EducationItem `json:"EducationItems"`
	Language       string          `json:"-"`
	Translations   map[string]map[string]string `json:"-"`
}

type ProjectItem struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
	GitHubLink  string `json:"GitHubLink"`
}

type ProjectsData struct {
	ProjectItems []ProjectItem `json:"ProjectItems"`
	Language     string        `json:"-"`
	Translations map[string]map[string]string `json:"-"`
}

type ProfileData struct {
	Title    string `json:"Title"`
	Text     string `json:"Text"`
	Language string `json:"-"`
	Translations map[string]map[string]string `json:"-"`
}

type IndexData struct {
	Skills       []string                      `json:"Skills"`
	Profile      ProfileData                   `json:"Profile"`
	Language     string                        `json:"-"`
	Translations map[string]map[string]string `json:"-"`
}
