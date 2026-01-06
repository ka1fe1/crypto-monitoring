package polymarket

// Market is the raw API response from Gamma API
type Market struct {
	ID                           string  `json:"id"`
	Question                     string  `json:"question"`
	ConditionID                  string  `json:"conditionId"`
	Slug                         string  `json:"slug"`
	TwitterCardImage             string  `json:"twitterCardImage,omitempty"`
	ResolutionSource             string  `json:"resolutionSource,omitempty"`
	EndDate                      string  `json:"endDate"`
	Category                     string  `json:"category"`
	AmmType                      string  `json:"ammType"`
	Liquidity                    string  `json:"liquidity"`
	SponsorName                  string  `json:"sponsorName,omitempty"`
	SponsorImage                 string  `json:"sponsorImage,omitempty"`
	StartDate                    string  `json:"startDate"`
	XAxisValue                   string  `json:"xAxisValue,omitempty"`
	YAxisValue                   string  `json:"yAxisValue,omitempty"`
	DenominationToken            string  `json:"denominationToken"`
	Fee                          string  `json:"fee"`
	Image                        string  `json:"image,omitempty"`
	Icon                         string  `json:"icon,omitempty"`
	LowerBound                   string  `json:"lowerBound,omitempty"`
	UpperBound                   string  `json:"upperBound,omitempty"`
	Description                  string  `json:"description,omitempty"`
	Outcomes                     string  `json:"outcomes"`
	OutcomePrices                string  `json:"outcomePrices"`
	Volume                       string  `json:"volume"`
	Active                       bool    `json:"active"`
	MarketType                   string  `json:"marketType"`
	FormatType                   string  `json:"formatType"`
	LowerBoundDate               string  `json:"lowerBoundDate,omitempty"`
	UpperBoundDate               string  `json:"upperBoundDate,omitempty"`
	Closed                       bool    `json:"closed"`
	MarketMakerAddress           string  `json:"marketMakerAddress"`
	CreatedBy                    int     `json:"createdBy"`
	UpdatedBy                    int     `json:"updatedBy"`
	CreatedAt                    string  `json:"createdAt"`
	UpdatedAt                    string  `json:"updatedAt"`
	ClosedTime                   string  `json:"closedTime,omitempty"`
	WideFormat                   bool    `json:"wideFormat"`
	New                          bool    `json:"new"`
	MailchimpTag                 string  `json:"mailchimpTag,omitempty"`
	Featured                     bool    `json:"featured"`
	Archived                     bool    `json:"archived"`
	ResolvedBy                   string  `json:"resolvedBy,omitempty"`
	Restricted                   bool    `json:"restricted"`
	MarketGroup                  int     `json:"marketGroup,omitempty"`
	GroupItemTitle               string  `json:"groupItemTitle,omitempty"`
	AutomaticallyResolved        bool    `json:"automaticallyResolved"`
	OneDayPriceChange            float64 `json:"oneDayPriceChange,omitempty"`
	OneHourPriceChange           float64 `json:"oneHourPriceChange,omitempty"`
	OneWeekPriceChange           float64 `json:"oneWeekPriceChange,omitempty"`
	OneMonthPriceChange          float64 `json:"oneMonthPriceChange,omitempty"`
	OneYearPriceChange           float64 `json:"oneYearPriceChange,omitempty"`
	LastTradePrice               float64 `json:"lastTradePrice,omitempty"`
	BestBid                      float64 `json:"bestBid,omitempty"`
	BestAsk                      float64 `json:"bestAsk,omitempty"`
	AutomaticallyActive          bool    `json:"automaticallyActive"`
	ClearBookOnStart             bool    `json:"clearBookOnStart"`
	ChartColor                   string  `json:"chartColor,omitempty"`
	ShowGmpOutcome               bool    `json:"showGmpOutcome"`
	ManualActivation             bool    `json:"manualActivation"`
	NegRiskOther                 bool    `json:"negRiskOther"`
	GameID                       string  `json:"gameId,omitempty"`
	GroupItemRange               string  `json:"groupItemRange,omitempty"`
	SportsMarketType             string  `json:"sportsMarketType,omitempty"`
	Line                         float64 `json:"line,omitempty"`
	UmaResolutionStatuses        string  `json:"umaResolutionStatuses,omitempty"`
	PendingDeployment            bool    `json:"pendingDeployment"`
	Deploying                    bool    `json:"deploying"`
	DeployingTimestamp           string  `json:"deployingTimestamp,omitempty"`
	ScheduledDeploymentTimestamp string  `json:"scheduledDeploymentTimestamp,omitempty"`
}

// MarketDetail is the refined response structure
type MarketDetail struct {
	Question           string             `json:"question"`
	Slug               string             `json:"slug"`
	Volume             float64            `json:"volume"`
	OutcomePrices      map[string]float64 `json:"outcome_prices"`
	Closed             bool               `json:"closed"`
	OneHourPriceChange float64            `json:"one_hour_price_change"`
	OneWeekPriceChange float64            `json:"one_week_price_change"`
}
