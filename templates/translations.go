package templates

// Translation maps for static text
var Translations = map[string]map[string]string{
	"name_label": {
		"en": "Name",
		"fr": "Nom",
	},
	"email_label": {
		"en": "Email",
		"fr": "Email",
	},
	"message_label": {
		"en": "Message",
		"fr": "Message",
	},
	"send_message": {
		"en": "Send Message",
		"fr": "Envoyer le Message",
	},
	"loading_experience": {
		"en": "Loading experience...",
		"fr": "Chargement de l'expérience...",
	},
	"loading_education": {
		"en": "Loading education...",
		"fr": "Chargement de l'éducation...",
	},
	"loading_projects": {
		"en": "Loading projects...",
		"fr": "Chargement des projets...",
	},
	"loading_contact": {
		"en": "Loading contact form...",
		"fr": "Chargement du formulaire de contact...",
	},
	"filter_skills": {
		"en": "Filter skills...",
		"fr": "Filtrer les compétences...",
	},
	"professional_experience": {
		"en": "Professional Experience",
		"fr": "Expérience Professionnelle",
	},
	"education": {
		"en": "Education",
		"fr": "Éducation",
	},
	"personal_projects": {
		"en": "Personal Projects",
		"fr": "Projets Personnels",
	},
	"skills": {
		"en": "Skills",
		"fr": "Compétences",
	},
	"contact_me": {
		"en": "Contact Me",
		"fr": "Contactez-Moi",
	},
	"experience": {
		"en": "Experience",
		"fr": "Expérience",
	},
	"projects": {
		"en": "Projects",
		"fr": "Projets",
	},
	"contact": {
		"en": "Contact",
		"fr": "Contact",
	},
	"message_sent": {
		"en": "Message sent successfully!",
		"fr": "Message envoyé avec succès !",
	},
	"failed_send": {
		"en": "Failed to send email",
		"fr": "Échec de l'envoi de l'email",
	},
	"all_fields_required": {
		"en": "All fields are required",
		"fr": "Tous les champs sont requis",
	},
	"loading_stats": {
		"en": "Loading stats...",
		"fr": "Chargement des statistiques...",
	},
	"view_on_github": {
		"en": "View on GitHub",
		"fr": "Voir sur GitHub",
	},
}

// Helper to get translation
func GetTranslation(key, lang string) string {
	if val, ok := Translations[key][lang]; ok {
		return val
	}
	return Translations[key]["en"] // Fallback to English
}
