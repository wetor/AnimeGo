package mikan

type MikanInfo struct {
	ID         int    `json:"id"`
	SubGroupID int    `json:"sub_group_id"`
	PubGroupID int    `json:"pub_group_id"`
	GroupName  string `json:"group_name"`
}
