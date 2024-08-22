package domain

type RegistrationReq struct {
	Labels      []Label
	Resources   map[string]float64
	BindAddress string
}

type RegistrationResp struct {
	NodeId string
}
