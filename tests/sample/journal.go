package sample

import (
	"github.com/louisevanderlith/husk/validation"
	"time"
)

type Journal struct {
	LastUpdated time.Time `json:"last_updated"`
	Created     time.Time `json:"created_date"`
	Entry       Entry     `json:"bibjson"`
	Admin       Admin     `json:"admin"`
	ID          string    `json:"id"`
}

func (j Journal) Valid() error {
	return validation.Struct(j)
}

type Entry struct {
	Country                string              `json:"country"`
	ChargesUrl             string              `json:"submission_charges_url"`
	Keywords               []string            `json:"keywords"`
	Subjects               []Subject           `json:"subject"`
	Links                  []Link              `json:"link"`
	Languages              []string            `json:"language"`
	Title                  string              `json:"title"`
	ArchivingPolicy        Policy              `json:"archiving_policy"`
	PlagiarismDetection    PlagiarismDetection `json:"plagiarism_detection"`
	Institution            string              `json:"institution"`
	EditorialReview        Review              `json:"editorial_review"`
	AllowsIndex            bool                `json:"allows_fulltext_indexing"`
	Statistics             Statistics          `json:"statistics"`
	AuthorPublishingRights PublishingRights    `json:"author_publishing_rights"`
	AlternativeTitle       string              `json:"alternative_title"`
	DepositPolicy          []string            `json:"deposit_policy"`
	Identifier             []Identifier        `json:"identifier"`
	Formats                []string            `json:"format"`
	Active                 bool                `json:"active"`
	Licenses               []License           `json:"license"`
	OAStart                OA                  `json:"oa_start"`
	Publisher              string              `json:"publisher"`
	AuthorCopyright        Copyright           `json:"author_copyright"`
	APCUrl                 string              `json:"apc_url"`
	PublicationTime        int                 `json:"publication_time"`
}

type Subject struct {
	Code   string `json:"code"`
	Scheme string `json:"scheme"`
	Term   string `json:"term"`
}

type Link struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Policy struct {
	Url string `json:"url"`
}

type PlagiarismDetection struct {
	Detect bool   `json:"detection"`
	Url    string `json:"url"`
}

type Review struct {
	Process string `json:"process"`
	Url     string `json:"url"`
}

type Statistics struct {
	Statistics bool   `json:"statistics"`
	Url        string `json:"url"`
}

type PublishingRights struct {
	PublishingRights string `json:"publishing_rights"`
	Url              string `json:"url"`
}

type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type License struct {
	EmbeddedExample string `json:"embedded_example_url"`
	NC              bool   `json:"NC"`
	ND              bool   `json:"ND"`
	BY              bool   `json:"BY"`
	SA              bool   `json:"SA"`
	OpenAccess      bool   `json:"open_access"`
	Title           string `json:"title"`
	Type            string `json:"type"`
	Embedded        bool   `json:"embedded"`
	Url             string `json:"url"`
}

type OA struct {
	Year int `json:"year"`
}

type Copyright struct {
	Copyright string `json:"copyright"`
	Url       string `json:"url"`
}

type Admin struct {
	Ticked bool `json:"ticked"`
	Seal   bool `json:"seal"`
}
