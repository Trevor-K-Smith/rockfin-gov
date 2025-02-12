package db

import (
	"gorm.io/gorm"
)

type SamGovData struct {
	gorm.Model
	TotalRecords  int64               `json:"totalRecords"`
	Limit         int                 `json:"limit"`
	Offset        int                 `json:"offset"`
	Opportunities []SamGovOpportunity `json:"opportunitiesData"`
	Links         []SamGovLink        `json:"links"`
}

type SamGovOpportunity struct {
	gorm.Model
	NoticeId                  string               `json:"noticeId"`
	Title                     string               `json:"title"`
	SolicitationNumber        string               `json:"solicitationNumber"`
	FullParentPathName        string               `json:"fullParentPathName"`
	FullParentPathCode        string               `json:"fullParentPathCode"`
	PostedDate                string               `json:"postedDate"`
	Type                      string               `json:"type"`
	BaseType                  string               `json:"baseType"`
	ArchiveType               string               `json:"archiveType"`
	ArchiveDate               string               `json:"archiveDate"`
	TypeOfSetAsideDescription *string              `json:"typeOfSetAsideDescription"`
	TypeOfSetAside            *string              `json:"typeOfSetAside"`
	ResponseDeadLine          string               `json:"responseDeadLine"`
	NaicsCode                 string               `json:"naicsCode"`
	NaicsCodes                []SamGovNaicsCode    `json:"naicsCodes"`
	ClassificationCode        string               `json:"classificationCode"`
	Active                    string               `json:"active"`
	Award                     *string              `json:"award"`
	PointOfContact            []SamGovContact      `json:"pointOfContact"`
	Description               string               `json:"description"`
	OrganizationType          string               `json:"organizationType"`
	OfficeAddress             SamGovAddress        `json:"officeAddress"`
	PlaceOfPerformance        *string              `json:"placeOfPerformance"`
	AdditionalInfoLink        *string              `json:"additionalInfoLink"`
	UILink                    string               `json:"uiLink"`
	Links                     []SamGovLink         `json:"links"`
	ResourceLinks             []SamGovResourceLink `json:"resourceLinks"`
	SamGovDataID              uint
}

type SamGovLink struct {
	gorm.Model
	Rel                 string `json:"rel"`
	Href                string `json:"href"`
	SamGovDataID        uint
	SamGovOpportunityID uint
}

type SamGovNaicsCode struct {
	gorm.Model
	Code                string `json:"code"`
	SamGovOpportunityID uint
}

type SamGovContact struct {
	gorm.Model
	Fax                 string  `json:"fax"`
	Type                string  `json:"type"`
	Email               string  `json:"email"`
	Phone               string  `json:"phone"`
	Title               *string `json:"title"`
	FullName            string  `json:"fullName"`
	SamGovOpportunityID uint
}

type SamGovAddress struct {
	gorm.Model
	Zipcode             string `json:"zipcode"`
	City                string `json:"city"`
	CountryCode         string `json:"countryCode"`
	State               string `json:"state"`
	SamGovOpportunityID uint
}

type SamGovResourceLink struct {
	gorm.Model
	URL                 string `json:"url"`
	SamGovOpportunityID uint
}

type Contact struct {
	Fax      string  `json:"fax"`
	Type     string  `json:"type"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Title    *string `json:"title"`
	FullName string  `json:"fullName"`
}

type Address struct {
	Zipcode     string `json:"zipcode"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	State       string `json:"state"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func PersistSamGovData(db *gorm.DB, data SamGovData) error {
	result := db.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
