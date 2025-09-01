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
}

type EducationItem struct {
	Title       string `json:"Title"`
	Institution string `json:"Institution"`
	Period      string `json:"Period"`
}

type EducationData struct {
	EducationItems []EducationItem `json:"EducationItems"`
}

type ProjectItem struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
	GitHubLink  string `json:"GitHubLink"`
}

type ProjectsData struct {
	ProjectItems []ProjectItem `json:"ProjectItems"`
}

type IndexData struct {
	Skills []string `json:"Skills"`
}
