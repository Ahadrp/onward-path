package usr

type LoginParam struct {
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

type AddClientParam struct {
	Server  int `json:"server"`
	ExpiryTime int `json:"expiry_time"`
	Flow string `json:"flow"`
	Total int `json:"total"`
	Email  string `json:"email"`
}
