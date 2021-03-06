package gads

import (
	"encoding/xml"
	"fmt"
)

// CampaignService A campaignService holds the connection information for the
// campaign service.
type CampaignService struct {
	Auth
}

// NewCampaignService creates a new campaignService
func NewCampaignService(auth *Auth) *CampaignService {
	return &CampaignService{Auth: *auth}
}

// ConversionOptimizerEligibility
//
// RejectionReasons can be any of
//   "CAMPAIGN_IS_NOT_ACTIVE", "NOT_CPC_CAMPAIGN","CONVERSION_TRACKING_NOT_ENABLED",
//   "NOT_ENOUGH_CONVERSIONS", "UNKNOWN"
//
type conversionOptimizerEligibility struct {
	Eligible         bool     `xml:"eligible"`         // is eligible for optimization
	RejectionReasons []string `xml:"rejectionReasons"` // reason for why campaign is
	// not eligible for conversion optimization.
}

type FrequencyCap struct {
	Impressions int64  `xml:"impressions"`
	TimeUnit    string `xml:"timeUnit"`
	Level       string `xml:"level,omitempty"`
}

type CampaignSetting struct {
	XMLName xml.Name `xml:"settings"`
	Type    string   `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`

	// GeoTargetTypeSetting
	PositiveGeoTargetType *string `xml:"positiveGeoTargetType,omitempty"`
	NegativeGeoTargetType *string `xml:"negativeGeoTargetType,omitempty"`

	// RealTimeBiddingSetting
	OptIn *bool `xml:"optIn,omitempty"`

	// DynamicSearchAdsSetting
	DomainName   *string `xml:"domainName,omitempty"`
	LanguageCode *string `xml:"langaugeCode,omitempty"`

	// TrackingSetting
	TrackingUrl *string `xml:"trackingUrl,omitempty"`

	// ShoppingSetting
	MerchantId       *int64  `xml:"merchantId,omitempty"`
	SalesCountry     *string `xml:"salesCountry,omitempty"`
	CampaignPriority *int    `xml:"campaignPriority"`
	EnableLocal      *bool   `xml:"enableLocal,omitempty"`
}

func NewDynamicSearchAdsSetting(domainName, languageCode string) CampaignSetting {
	return CampaignSetting{
		Type:         "DynamicSearchAdsSetting",
		DomainName:   &domainName,
		LanguageCode: &languageCode,
	}
}

func NewGeoTargetTypeSetting(positiveGeoTargetType, negativeGeoTargetType string) CampaignSetting {
	return CampaignSetting{
		Type: "GeoTargetTypeSetting",
		PositiveGeoTargetType: &positiveGeoTargetType,
		NegativeGeoTargetType: &negativeGeoTargetType,
	}
}

func NewRealTimeBiddingSetting(optIn bool) CampaignSetting {
	return CampaignSetting{
		Type:  "RealTimeBiddingSetting",
		OptIn: &optIn,
	}
}

func NewTrackingSetting(trackingUrl string) CampaignSetting {
	return CampaignSetting{
		Type:        "TrackingSetting",
		TrackingUrl: &trackingUrl,
	}
}

func NewShoppingSetting(merchantID int64, salesCountry string, campaignPriority int, enableLocal bool) CampaignSetting {
	return CampaignSetting{
		Type:             "ShoppingSetting",
		MerchantId:       &merchantID,
		SalesCountry:     &salesCountry,
		CampaignPriority: &campaignPriority,
		EnableLocal:      &enableLocal,
	}
}

type NetworkSetting struct {
	TargetGoogleSearch         bool `xml:"targetGoogleSearch"`
	TargetSearchNetwork        bool `xml:"targetSearchNetwork"`
	TargetContentNetwork       bool `xml:"targetContentNetwork"`
	TargetPartnerSearchNetwork bool `xml:"targetPartnerSearchNetwork"`
}

/*
BIDDING SCHEME
*/

// BiddingSchemeInterface interface  for bidding scheme
type BiddingSchemeInterface interface {
	GetType() string
}

// BiddingScheme struct for ManualCpcBiddingScheme
type BiddingScheme struct {
	Type               string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	EnhancedCpcEnabled bool   `xml:"enhancedCpcEnabled"`
}

// NewBiddingScheme returns new instance of BiddingScheme
func NewBiddingScheme(enhancedCpcEnabled bool) *BiddingScheme {
	return &BiddingScheme{Type: `ManualCpcBiddingScheme`, EnhancedCpcEnabled: enhancedCpcEnabled}
}

// GetType return type of bidding scheme
func (s *BiddingScheme) GetType() string {
	return s.Type
}

// TargetRoasBiddingScheme struct for TargetRoasBiddingScheme
type TargetRoasBiddingScheme struct {
	Type       string  `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	TargetRoas float64 `xml:"targetRoas"`
	BidCeiling *int64  `xml:"bidCeiling>microAmount"`
	BidFloor   *int64  `xml:"bidFloor>microAmount"`
}

// NewTargetRoasBiddingScheme returns new instance of TargetRoasBiddingScheme
func NewTargetRoasBiddingScheme(targetRoas float64, bidCeiling, bidFloor *int64) *TargetRoasBiddingScheme {
	return &TargetRoasBiddingScheme{
		Type:       `TargetRoasBiddingScheme`,
		TargetRoas: targetRoas,
		BidCeiling: bidCeiling,
		BidFloor:   bidFloor,
	}
}

// GetType return type of bidding scheme
func (s *TargetRoasBiddingScheme) GetType() string {
	return s.Type
}

func biddingSchemeUnmarshalXML(dec *xml.Decoder, start xml.StartElement) (BiddingSchemeInterface, error) {
	biddingSchemeType, err := findAttr(start.Attr, xml.Name{Space: "http://www.w3.org/2001/XMLSchema-instance", Local: "type"})
	if err != nil {
		return nil, err
	}
	switch biddingSchemeType {
	case "ManualCpcBiddingScheme":
		c := &BiddingScheme{Type: biddingSchemeType}
		return c, dec.DecodeElement(c, &start)
	case "TargetRoasBiddingScheme":
		c := &TargetRoasBiddingScheme{Type: biddingSchemeType}
		return c, dec.DecodeElement(c, &start)
	default:
		if StrictMode {
			return nil, fmt.Errorf("unknown bidding scheme type %#v", biddingSchemeType)
		}
		return nil, nil
	}
}

type Bid struct {
	Type         string  `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	Amount       int64   `xml:"bid>microAmount"`
	CpcBidSource *string `xml:"cpcBidSource"`
	CpmBidSource *string `xml:"cpmBidSource"`
}

type BiddingStrategyConfiguration struct {
	StrategyId     int64                  `xml:"biddingStrategyId,omitempty"`
	StrategyName   string                 `xml:"biddingStrategyName,omitempty"`
	StrategyType   string                 `xml:"biddingStrategyType,omitempty"`
	StrategySource string                 `xml:"biddingStrategySource,omitempty"`
	Scheme         BiddingSchemeInterface `xml:"biddingScheme,omitempty"`
	Bids           []Bid                  `xml:"bids"`
}

// UnmarshalXML special unmarshal for the different bidding schemes
func (b *BiddingStrategyConfiguration) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "biddingStrategyId":
				if err := dec.DecodeElement(&b.StrategyId, &start); err != nil {
					return err
				}
			case "biddingStrategyName":
				if err := dec.DecodeElement(&b.StrategyName, &start); err != nil {
					return err
				}
			case "biddingStrategyType":
				if err := dec.DecodeElement(&b.StrategyType, &start); err != nil {
					return err
				}
			case "biddingStrategySource":
				if err := dec.DecodeElement(&b.StrategySource, &start); err != nil {
					return err
				}
			case "biddingScheme":
				bs, err := biddingSchemeUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				b.Scheme = bs
			case "bids":
				if err := dec.DecodeElement(&b.Bids, &start); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type CustomParameter struct {
	Key      string `xml:"key"`
	Value    string `xml:"value"`
	IsRemove bool   `xml:"isRemove"`
}

type CustomParameters struct {
	CustomParameters []CustomParameter `xml:"parameters"`
	DoReplace        bool              `xml:"doReplace"`
}

type Campaign struct {
	Id                             int64                           `xml:"id,omitempty"`
	Name                           string                          `xml:"name,omitempty"`
	Status                         string                          `xml:"status,omitempty"`        // Status: "ENABLED", "PAUSED", "REMOVED"
	ServingStatus                  *string                         `xml:"servingStatus,omitempty"` // ServingStatus: "SERVING", "NONE", "ENDED", "PENDING", "SUSPENDED"
	StartDate                      string                          `xml:"startDate,omitempty"`
	EndDate                        *string                         `xml:"endDate,omitempty"`
	BudgetId                       int64                           `xml:"budget>budgetId,omitempty"`
	ConversionOptimizerEligibility *conversionOptimizerEligibility `xml:"conversionOptimizerEligibility,omitempty"`
	AdServingOptimizationStatus    string                          `xml:"adServingOptimizationStatus,omitempty"`
	FrequencyCap                   *FrequencyCap                   `xml:"frequencyCap,omitempty"`
	Settings                       []CampaignSetting               `xml:"settings"`
	AdvertisingChannelType         string                          `xml:"advertisingChannelType,omitempty"`    // "UNKNOWN", "SEARCH", "DISPLAY", "SHOPPING"
	AdvertisingChannelSubType      *string                         `xml:"advertisingChannelSubType,omitempty"` // "UNKNOWN", "SEARCH_MOBILE_APP", "DISPLAY_MOBILE_APP", "SEARCH_EXPRESS", "DISPLAY_EXPRESS"
	NetworkSetting                 *NetworkSetting                 `xml:"networkSetting"`
	Labels                         []Label                         `xml:"labels,omitempty"`
	BiddingStrategyConfiguration   *BiddingStrategyConfiguration   `xml:"biddingStrategyConfiguration,omitempty"`
	ForwardCompatibilityMap        *map[string]string              `xml:"forwardCompatibilityMap,omitempty"`
	TrackingUrlTemplate            *string                         `xml:"trackingUrlTemplate,omitempty"`
	UrlCustomParameters            *CustomParameters               `xml:"urlCustomParameters,omitempty"`
	BaseCampaignID                 *int64                          `xml:"baseCampaignId,omitempty"`
	CampaignTrialType              *string                         `xml:"campaignTrialType,omitempty"`
	Errors                         []error                         `xml:"-"`
}

type CampaignOperations map[string][]Campaign

type CampaignLabel struct {
	CampaignId int64 `xml:"campaignId"`
	LabelId    int64 `xml:"labelId"`
}

type CampaignLabelOperations map[string][]CampaignLabel

// Get returns an array of Campaign's and the total number of campaign's matching
// the selector.
//
// Example
//
//   campaigns, totalCount, err := campaignService.Get(
//     gads.Selector{
//       Fields: []string{
//         "AdGroupId",
//         "Status",
//         "AdGroupCreativeApprovalStatus",
//         "AdGroupAdDisapprovalReasons",
//         "AdGroupAdTrademarkDisapproved",
//       },
//       Predicates: []gads.Predicate{
//         {"AdGroupId", "EQUALS", []string{adGroupId}},
//       },
//     },
//   )
//
// Selectable fields are
//   "Id", "Name", "Status", "ServingStatus", "StartDate", "EndDate", "AdServingOptimizationStatus",
//   "Settings", "AdvertisingChannelType", "AdvertisingChannelSubType", "Labels", "TrackingUrlTemplate",
//   "UrlCustomParameters"
//
// filterable fields are
//   "Id", "Name", "Status", "ServingStatus", "StartDate", "EndDate", "AdvertisingChannelType",
//   "AdvertisingChannelSubType", "Labels", "TrackingUrlTemplate"
//
// Relevant documentation
//
//     https://developers.google.com/adwords/api/docs/reference/v201806/CampaignService#get
//
func (s *CampaignService) Get(selector Selector) (campaigns []Campaign, totalCount int64, err error) {
	selector.XMLName = xml.Name{"", "serviceSelector"}
	respBody, err := s.Auth.request(
		campaignServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return campaigns, totalCount, err
	}
	getResp := struct {
		Size      int64      `xml:"rval>totalNumEntries"`
		Campaigns []Campaign `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaigns, totalCount, err
	}
	return getResp.Campaigns, getResp.Size, err
}

// Mutate allows you to add and modify campaigns, returning the
// campaigns.  Note that the "REMOVE" operator is not supported.
// To remove a campaign set its Status to "REMOVED".
//
// Example
//
//  campaignNeedingRemoval.Status = "REMOVED"
//  ads, err := campaignService.Mutate(
//    gads.CampaignOperations{
//      "ADD": {
//        gads.Campaign{
//          Name: "my campaign name",
//          Status: "PAUSED",
//          StartDate: time.Now().Format("20060102"),
//          BudgetId: 321543214,
//          AdServingOptimizationStatus: "ROTATE_INDEFINITELY",
//          Settings: []gads.CampaignSetting{
//            gads.NewRealTimeBiddingSetting(true),
//          },
//          AdvertisingChannelType: "SEARCH",
//          BiddingStrategyConfiguration: &gads.BiddingStrategyConfiguration{
//            StrategyType: "MANUAL_CPC",
//          },
//        },
//        campaignNeedingRemoval,
//      },
//      "SET": {
//        modifiedCampaign,
//      },
//    }
//
// Relevant documentation
//
//     https://developers.google.com/adwords/api/docs/reference/v201806/CampaignService#mutate
//
func (s *CampaignService) Mutate(campaignOperations CampaignOperations) (campaigns []Campaign, err error) {
	type campaignOperation struct {
		Action   string   `xml:"operator"`
		Campaign Campaign `xml:"operand"`
	}
	operations := []campaignOperation{}
	for action, campaigns := range campaignOperations {
		for _, campaign := range campaigns {
			// you can't mutate those fields
			// if you want to perform campaign mutate from campaign get
			// you can't.
			campaign.CampaignTrialType = nil
			campaign.AdServingOptimizationStatus = ""
			// you can't mutate this field too
			//if campaign.BiddingStrategyConfiguration != nil {
			//	campaign.BiddingStrategyConfiguration.StrategyType = ""
			//}
			operations = append(operations,
				campaignOperation{
					Action:   action,
					Campaign: campaign,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []campaignOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations}
	respBody, err := s.Auth.request(campaignServiceUrl, "mutate", mutation)
	if err != nil {
		return campaigns, err
	}
	mutateResp := struct {
		BaseResponse
		Campaigns []Campaign `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return campaigns, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.Campaigns, err
}

// Mutate allows you to add and removes labels from campaigns.
//
// Example
//
//  cls, err := campaignService.MutateLabel(
//    gads.CampaignOperations{
//      "ADD": {
//        gads.CampaignLabel{CampaignId: 3200, LabelId: 5353},
//        gads.CampaignLabel{CampaignId: 4320, LabelId: 5643},
//      },
//      "REMOVE": {
//        gads.CampaignLabel{CampaignId: 3653, LabelId: 5653},
//      },
//    }
//
// Relevant documentation
//
//     https://developers.google.com/adwords/api/docs/reference/v201806/CampaignService#mutateLabel
//
func (s *CampaignService) MutateLabel(campaignLabelOperations CampaignLabelOperations) (campaignLabels []CampaignLabel, err error) {
	type campaignLabelOperation struct {
		Action        string        `xml:"operator"`
		CampaignLabel CampaignLabel `xml:"operand"`
	}
	operations := []campaignLabelOperation{}
	for action, campaignLabels := range campaignLabelOperations {
		for _, campaignLabel := range campaignLabels {
			operations = append(operations,
				campaignLabelOperation{
					Action:        action,
					CampaignLabel: campaignLabel,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []campaignLabelOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutateLabel",
		},
		Ops: operations}
	respBody, err := s.Auth.request(campaignServiceUrl, "mutateLabel", mutation)
	if err != nil {
		return campaignLabels, err
	}
	mutateResp := struct {
		BaseResponse
		CampaignLabels []CampaignLabel `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return campaignLabels, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.CampaignLabels, err
}

// Query is not yet implemented
//
// Relevant documentation
//
//     https://developers.google.com/adwords/api/docs/reference/v201806/CampaignService#query
//
func (s *CampaignService) Query(query string) (campaigns []Campaign, totalCount int64, err error) {
	return campaigns, totalCount, ERROR_NOT_YET_IMPLEMENTED
}
