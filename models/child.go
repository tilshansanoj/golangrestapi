package models


type Child struct{
	User
	ParentID int `json:"id"`
}