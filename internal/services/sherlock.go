package services

type SherlockService interface {
}
type sherlockService struct {
}

func NewSherlockService() SherlockService {
	return &sherlockService{}
}
