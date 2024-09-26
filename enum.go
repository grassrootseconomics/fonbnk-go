package fonbnk

type (
	Country            string
	Network            string
	OffRampType        string
	OffRampCurrency    string
	OffRampAsset       string
	OffRampPaymentType string
	KYCStatus          string
	KYCIDType          string
)

const (
	KENYA  Country = "KE"
	UGANDA Country = "UG"

	CELO Network = "CELO"

	BANK         OffRampType = "bank"
	AIRTIME      OffRampType = "airtime"
	MOBILE_MONEY OffRampType = "mobile_money"
	PAYBILL      OffRampType = "paybill"

	LOCAL OffRampCurrency = "local"
	USD   OffRampCurrency = "usd"

	CUSD OffRampAsset = "CUSD"
	USDT OffRampAsset = "USDT"
	USDC OffRampAsset = "USDC"

	CRYPTO OffRampPaymentType = "CRYPTO_WALLET"

	INITIATED KYCStatus = "initiated"
	APPROVED  KYCStatus = "approved"
	REJECTED  KYCStatus = "rejected"
	INVALID   KYCStatus = "invalid"

	KE_NATIONAL_ID KYCIDType = "NATIONAL_ID"
	KE_PASSPORT    KYCIDType = "PASSPORT"
	KE_ALIEN_CARD  KYCIDType = "ALIEN_CARD"
)
