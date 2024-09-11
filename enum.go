package fonbnk

type (
	Country     string
	OffRampType string
)

const (
	KENYA  Country = "KE"
	UGANDA Country = "UG"

	BANK         OffRampType = "bank"
	AIRTIME      OffRampType = "airtime"
	MOBILE_MONEY OffRampType = "mobile_money"
	PAYBILL      OffRampType = "paybill"
)
