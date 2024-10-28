package address

type GetAddressByIdReqStruct struct {
	Id int64 `json:"id" binding:"required"`
}

type AddAddressReqStruct struct {
	Id            int64  `json:"id"`
	MemberId      int64  `json:"memberId"`
	City          string `json:"city"`
	DefaultStatus int    `json:"defaultStatus"`
	DetailAddress string `json:"detailAddress"`
	Name          string `json:"name"`
	PhoneNumber   string `json:"phoneNumber"`
	PostCode      string `json:"postCode"`
	Province      string `json:"province"`
	Region        string `json:"region"`
}

type DeleteAddressReqStruct struct {
	Id int64 `json:"id"`
}

type UpdateAddressReqStruct struct {
	Id            int64  `json:"id"`
	MemberId      int64  `json:"memberId"`
	City          string `json:"city"`
	DefaultStatus int    `json:"defaultStatus"`
	DetailAddress string `json:"detailAddress"`
	Name          string `json:"name"`
	PhoneNumber   string `json:"phoneNumber"`
	PostCode      string `json:"postCode"`
	Province      string `json:"province"`
	Region        string `json:"region"`
}
