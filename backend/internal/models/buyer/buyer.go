package buyer

type BuyerRegister struct{
	EmailAddress string `json:"email_address"`
	Password string `json:"password"`
}