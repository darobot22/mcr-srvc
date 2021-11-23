package models

type Service struct {
	Id   int    `json:"Id,omitempty"`
	Name string `json:"Name,omitempty"`
	Code string `json:"Code,omitempty"`
}

type RequestedService struct {
	UserId    int
	ServiceId int
	Params    string
}

type User struct {
	Id       int    `json:"Id,omitempty"`
	Name     string `json:"Name,omitempty"`
	Password string `json:"Password,omitempty"`
	Email    string `json:"Email,omitempty"`
	Phone    string `json:"Phone,omitempty"`
}

type ServiceHistory struct {
	Id            int
	ServiceCode   string
	ServiceName   string
	UserId        int
	UserName      string
	CreateDate    string
	ResultData    string
	ExecutionDate string
}
