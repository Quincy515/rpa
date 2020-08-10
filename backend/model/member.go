package model

type OauthMemberBind struct {
	UserSn     uint64 `json:"user_sn" db:"user_sn"`
	ClientType string `json:"client_type" db:"client_type"`
	Type       int    `json:"type" db:"type"`
	Openid     string `json:"openid" db:"openid"`
	Unionid    string `json:"unionid" db:"unionid"`
	Extra      string `json:"extra" db:"extra"`
}
