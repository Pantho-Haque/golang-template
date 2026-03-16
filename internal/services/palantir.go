package services

type PalantirService interface {
}
type palantirService struct {
}

func NewPalantirService() PalantirService {
	return &palantirService{}
}
