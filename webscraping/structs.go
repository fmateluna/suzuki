package webscraping

import (
	"fmt"
	"math/rand"

	"github.com/golang-jwt/jwt/v4"
)

type loginResponse struct {
	Status                   int    `json:"status"`
	SessionJwtToken          string `json:"sessionJwtToken"`
	ProductCodesLicenseCodes []struct {
		ProductCode string `json:"productCode"`
		LicenseCode string `json:"licenseCode"`
	} `json:"productCodesLicenseCodes"`
	LocalSessionPingInterval int `json:"localSessionPingInterval"`
}

type SBSEPC5S struct {
	jwt.RegisteredClaims
	Sid string `json:"SID"`
	Ts  string `json:"TS"`
	Pk  string `json:"PK"`
	Rd  string `json:"RD"`
}

type SBSEPC5CS struct {
	jwt.RegisteredClaims
	Sid string `json:"SID,omitempty"`
	Rd  string `json:"RD"`
	Ts  string `json:"TS"`
	Pk  string `json:"PK"`
}

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func (b *BotSuzuki) getDatasetId(datasetName string) string {
	for _, v := range b.UserBot.DatasetSettings {
		if v.DatasetName == datasetName {
			return v.DatasetID
		}
	}
	return ""
}

func (b *BotSuzuki) returnFirma() string {
	return "jobId=1|dataSetId=" + b.getDatasetId("Suzuki ATV") + "|locale=en-US|busReg=SUZ|LA|dataSetId=" + b.getDatasetId("Suzuki Automotive") + "|locale=en-US|busReg=SUZ|LA|dataSetId=" + b.getDatasetId("Suzuki Marine") + "|locale=en-US|busReg=SUZ|LA|dataSetId=" + b.getDatasetId("Suzuki Motorcycle") + "|locale=en-US|busReg=SUZ|LA|userId=" + b.AccountBot.UserDetails.UserID
}

func (rb *BotSuzuki) getCookieString() string {
	list := ""
	for key, cookie := range rb.Cookies {
		item := key + "=" + cookie.Value + ";"
		list = list + item
	}
	//list = list + rb.Csrf + ";" + rb.XfSession + ";xf_from_search=google; xf_csrf=5VOMNgavdKeCMhB8; _ga=GA1.2.864175125.1673533370; _gid=GA1.2.2061440800.1673533370;_ga=GA1.2.2005128399.1673553711; _gid=GA1.2.1758049683.1673553711; _gat_gtag_UA_62047822_1=1; xf_csrf=rQEBXUK_F0KBwmy7; __cf_bm=rBmeY2Kvbe5.uzpe9R12wH4Dwvz_PBbcut.9Gvb3T4c-1673553710-0-AVqClLrNn63m3XgpndbFY//mx/pb6owgX2jTPt/E0fxd7fUIB/OLMSImzDY3kdAu9wRmlRHX+u3Qv2KLitSsyl6+m67tqHnfUeBdYcveclgSN6nen+Ly6v9POi//bFkLAYP9eurO2dbOxC3np0cZNaM=;"
	fmt.Println("SETEANDO COOKIES " + list)
	return list
}

func (b *BotSuzuki) RandomString(length int) string {
	ll := len(chars)
	bit := make([]byte, length)
	rand.Read(bit) // generates len(b) random bytes
	for i := 0; i < length; i++ {
		bit[i] = chars[int(bit[i])%ll]
	}
	return string(bit)
}

// https://suzuki.snaponepc.com/epc-services/equipment/search
type SearchVIN struct {
	VinSearchResults []struct {
		DatasetName string `json:"datasetName"`
		Vins        []VIN  `json:"vins"`
		Columns     []struct {
			Field  string `json:"field"`
			Header string `json:"header"`
		} `json:"columns"`
	} `json:"vinSearchResults"`
}
type VIN struct {
	DatasetID              string        `json:"datasetId"`
	SerializedPath         string        `json:"serializedPath"`
	DatasetName            string        `json:"datasetName"`
	ModelName              string        `json:"modelName"`
	ModelQualifierName     string        `json:"modelQualifierName"`
	EquipmentName          string        `json:"equipmentName"`
	ID                     string        `json:"id"`
	Vin                    string        `json:"vin"`
	FormattedVin           string        `json:"formattedVin"`
	BusinessRegion         int           `json:"businessRegion"`
	BusinessRegionName     string        `json:"businessRegionName"`
	ExternalBusinessRegion bool          `json:"externalBusinessRegion"`
	VinNote                []interface{} `json:"vinNote"`
	EquipmentRefID         string        `json:"equipmentRefId"`
	EquipmentKey           string        `json:"equipmentKey"`
	RangeLookup            bool          `json:"rangeLookup"`
	VinResolution          string        `json:"vinResolution"`
	EinID                  string        `json:"einId"`
	ValidationString       string        `json:"validationString"`
	HasUserVinNote         bool          `json:"hasUserVinNote"`
	VinRecalled            bool          `json:"vinRecalled"`
	Categories             []Category    `json:"categories"`
}

// https://suzuki.snaponepc.com/epc-services/auth/account
type Account struct {
	UserDetails struct {
		UserName     string `json:"userName"`
		UserID       string `json:"userId"`
		LastAccess   int64  `json:"lastAccess"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		EmailAddress string `json:"emailAddress"`
	} `json:"userDetails"`
	DealerDetails struct {
		Name         string `json:"name"`
		Address1     string `json:"address1"`
		Address2     string `json:"address2"`
		Address3     string `json:"address3"`
		City         string `json:"city"`
		PostalCode   string `json:"postalCode"`
		Country      string `json:"country"`
		EmailAddress string `json:"emailAddress"`
		Phone1       string `json:"phone1"`
		Fax          string `json:"fax"`
	} `json:"dealerDetails"`
}

// https://suzuki.snaponepc.com/epc-services/settings/user/
type User struct {
	ApplicationSettings struct {
		Locale                          string  `json:"locale"`
		DateFormat                      string  `json:"dateFormat"`
		TimeFormat                      string  `json:"timeFormat"`
		NavigationStyle                 string  `json:"navigationStyle"`
		PartsPanelPosition              string  `json:"partsPanelPosition"`
		EulaVersionAccepted             float64 `json:"eulaVersionAccepted"`
		ConfirmBeforeExitingApplication bool    `json:"confirmBeforeExitingApplication"`
		ConfirmBeforeClearingPicklist   bool    `json:"confirmBeforeClearingPicklist"`
		ConfirmBeforeClosingJob         bool    `json:"confirmBeforeClosingJob"`
		WarnMeOnUpdatesAvailable        bool    `json:"warnMeOnUpdatesAvailable"`
		SelectedPriceSource             struct {
		} `json:"selectedPriceSource"`
		AddPartsToBottomOfPicklist  bool   `json:"addPartsToBottomOfPicklist"`
		ShowPicklistByDefault       bool   `json:"showPicklistByDefault"`
		SelectedPicklistPriceSource string `json:"selectedPicklistPriceSource"`
		AutoClearPicklist           bool   `json:"autoClearPicklist"`
		HideQtyPrompt               bool   `json:"hideQtyPrompt"`
		DebugEnabled                bool   `json:"debugEnabled"`
		ShowAllIndicators           bool   `json:"showAllIndicators"`
		StatisticsEnabled           bool   `json:"statisticsEnabled"`
	} `json:"applicationSettings"`
	DatasetSettings []struct {
		DatasetID         string `json:"datasetId"`
		DatasetName       string `json:"datasetName"`
		Locale            string `json:"locale"`
		BusinessRegionKey string `json:"businessRegionKey"`
		BusinessRegion    int    `json:"businessRegion"`
	} `json:"datasetSettings"`
	EstimateSettings struct {
		Contact         string `json:"contact"`
		PriceMultiplier string `json:"priceMultiplier"`
		Currency        string `json:"currency"`
		LaborRate       string `json:"laborRate"`
		TaxRate         string `json:"taxRate"`
		HidePartNumbers bool   `json:"hidePartNumbers"`
		TaxLabor        bool   `json:"taxLabor"`
	} `json:"estimateSettings"`
}

type ResponseFilters struct {
	DatasetID       string `json:"datasetId"`
	LevelColumnSort string `json:"levelColumnSort"`
	Children        struct {
		ChildLevelTitle              string        `json:"childLevelTitle"`
		ChildLevelType               string        `json:"childLevelType"`
		ChildLevelSection            string        `json:"childLevelSection"`
		ChildLevelIllustrated        bool          `json:"childLevelIllustrated"`
		ChildLevelIllustrationWidth  int           `json:"childLevelIllustrationWidth"`
		ChildLevelIllustrationHeight int           `json:"childLevelIllustrationHeight"`
		ChildNodes                   []SubCategory `json:"childNodes"`
	} `json:"children"`
	Error bool `json:"error"`
}

type SubCategory struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	HasNotes       bool           `json:"hasNotes"`
	LeafNode       bool           `json:"leafNode"`
	ImageID        string         `json:"imageId"`
	SerializedPath string         `json:"serializedPath"`
	Filtered       bool           `json:"filtered"`
	Parts          ResponsePiezas `json:"parts"`
}

type Category struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	HasNotes       bool          `json:"hasNotes"`
	LeafNode       bool          `json:"leafNode"`
	ImageID        string        `json:"imageId"`
	SerializedPath string        `json:"serializedPath"`
	Filtered       bool          `json:"filtered"`
	SubCategory    []SubCategory `json:"subCategory"`
}

type ResponsePiezas struct {
	PageID            string           `json:"pageId"`
	Illustrated       bool             `json:"illustrated"`
	ImageID           string           `json:"imageId"`
	PageCode          string           `json:"pageCode"`
	PartItems         []PartItems      `json:"partItems"`
	PageImages        []PageImages     `json:"pageImages"`
	ColumnConfigs     []ColumnConfigs  `json:"columnConfigs"`
	PageLimitExceeded bool             `json:"pageLimitExceeded"`
	HasPageNotes      bool             `json:"hasPageNotes"`
	SubParts          []ResponsePiezas `json:"subparts,omitempty"`
}
type PartItems struct {
	PartID              string        `json:"partId"`
	ParentPartID        string        `json:"parentPartId"`
	SecondaryPartID     string        `json:"secondaryPartId"`
	Manufacturer        string        `json:"manufacturer"`
	PartNumber          string        `json:"partNumber"`
	FormattedPartNumber string        `json:"formattedPartNumber"`
	PartItemID          string        `json:"partItemId"`
	ParentPartItemID    string        `json:"parentPartItemId"`
	CalloutLabel        string        `json:"calloutLabel"`
	PaddedCallout       string        `json:"paddedCallout"`
	CrossCatKey         string        `json:"crossCatKey"`
	Description         string        `json:"description"`
	Quantity            string        `json:"quantity"`
	PartType            string        `json:"partType"`
	Indicators          []interface{} `json:"indicators"`
	AlphaSort           string        `json:"alphaSort"`
	AlphaSortSequence   string        `json:"alphaSortSequence"`
	Filtered            bool          `json:"filtered"`
	AddedManually       bool          `json:"addedManually"`
	Remarks             string        `json:"remarks,omitempty"`
}
type PageImages struct {
	ImageID    string `json:"imageId"`
	PageID     string `json:"pageId"`
	ImageTitle string `json:"imageTitle"`
}
type ColumnConfigs struct {
	Key         string  `json:"key"`
	Order       int     `json:"order"`
	Width       int     `json:"width"`
	PdfWidth    int     `json:"pdfWidth"`
	MaxPdfWidth float64 `json:"maxPdfWidth"`
	MinWidth    int     `json:"minWidth"`
	MaxWidth    int     `json:"maxWidth"`
	Resizable   bool    `json:"resizable"`
	Title       string  `json:"title"`
	Visible     bool    `json:"visible"`
	Override    string  `json:"override"`
}
