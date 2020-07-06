package sample

import (
	"github.com/louisevanderlith/husk"
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
	return husk.ValidateStruct(&j)
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

/*
[{"last_updated": "2020-02-03T17:25:49Z", "bibjson": {"country": "BR", "submission_charges_url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#custom-4", "keywords": ["sociology", "anthropology", "history", "political science", "international relations"], "subject": [{"code": "D", "scheme": "LCC", "term": "History (General) and history of Europe"}, {"code": "D1-2009", "scheme": "LCC", "term": "History (General)"}], "link": [{"type": "homepage", "url": "https://periodicos.furg.br/rbhcs"}, {"type": "editorial_board", "url": "https://periodicos.furg.br/rbhcs/about/editorialTeam"}, {"type": "aims_scope", "url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#focusAndScope"}, {"type": "author_instructions", "url": "https://periodicos.furg.br/rbhcs/about/submissions#authorGuidelines"}, {"type": "oa_statement", "url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#openAccessPolicy"}], "language": ["PT", "ES"], "title": "Revista Brasileira de Hist\u00f3ria & Ci\u00eancias Sociais", "archiving_policy": {"url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#sectionPolicies"}, "plagiarism_detection": {"detection": true, "url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#custom-6"}, "institution": "Universidade Federal do Rio Grande", "editorial_review": {"process": "Blind peer review", "url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#peerReviewProcess"}, "allows_fulltext_indexing": true, "article_statistics": {"url": "", "statistics": false}, "author_publishing_rights": {"publishing_rights": "False", "url": ""}, "alternative_title": "RBHCS", "deposit_policy": ["Diadorim"], "identifier": [{"id": "2175-3423", "type": "eissn"}], "format": ["PDF", "HTML"], "active": true, "license": [{"embedded_example_url": "", "NC": false, "ND": false, "BY": true, "open_access": true, "title": "CC BY", "type": "CC BY", "embedded": false, "url": "https://periodicos.furg.br/rbhcs", "SA": false}], "oa_start": {"year": 2009}, "author_copyright": {"copyright": "False", "url": ""}, "publisher": "Universidade Federal do Rio Grande", "apc_url": "https://periodicos.furg.br/rbhcs/about/editorialPolicies#custom-4", "publication_time": 16}, "admin": {"ticked": true, "seal": false}, "id": "e750089e1efe4d39891a8259d4120a39", "created_date": "2010-03-03T15:34:02Z"},
{"last_updated": "2020-02-03T14:12:42Z", "bibjson": {"country": "US", "submission_charges_url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/ForAuthors.html", "keywords": ["hiv", "aids", "public health", "health policy", "health economics"], "subject": [{"code": "RC581-607", "scheme": "LCC", "term": "Immunologic diseases. Allergy"}], "link": [{"type": "homepage", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652"}, {"type": "waiver_policy", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/article_publication_charges.htm"}, {"type": "editorial_board", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/EditorialBoard.html"}, {"type": "aims_scope", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/ProductInformation.html"}, {"type": "author_instructions", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/ForAuthors.html"}, {"type": "oa_statement", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/publishing_with_jias.htm"}], "language": ["EN"], "title": "Journal of the International AIDS Society ", "plagiarism_detection": {"detection": true, "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/ForAuthors.html"}, "archiving_policy": {"known": ["LOCKSS", "Portico", "PMC/Europe PMC/PMC Canada"], "nat_lib": "Koninklijke Bibliotheek", "url": "http://olabout.wiley.com/WileyCDA/Section/id-406156.html"}, "institution": "International AIDS Society", "editorial_review": {"process": "Blind peer review", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/ForAuthors.html"}, "allows_fulltext_indexing": true, "author_publishing_rights": {"publishing_rights": "True", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/open_access_license_and_copyright.htm"}, "article_statistics": {"url": "", "statistics": false}, "alternative_title": "JIAS", "apc": {"average_price": 2400, "currency": "USD"}, "provider": "Wiley Online Library", "deposit_policy": ["Sherpa/Romeo"], "identifier": [{"id": "1758-2652", "type": "eissn"}], "format": ["PDF", "HTML"], "persistent_identifier_scheme": ["DOI"], "license": [{"embedded_example_url": "http://onlinelibrary.wiley.com/doi/10.7448/IAS.20.1.21188/abstract", "NC": false, "ND": false, "BY": true, "open_access": true, "title": "CC BY", "type": "CC BY", "embedded": true, "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/open_access_license_and_copyright.htm", "SA": false}], "oa_start": {"year": 2004}, "author_copyright": {"copyright": "True", "url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/open_access_license_and_copyright.htm"}, "publisher": "Wiley", "apc_url": "http://onlinelibrary.wiley.com/journal/10.1002/(ISSN)1758-2652/homepage/article_publication_charges.htm", "publication_time": 25}, "admin": {"ticked": true, "seal": false}, "created_date": "2010-12-01T07:05:59Z", "id": "39837fcc443146b29211a0a399f64133"},
*/
