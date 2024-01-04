package polyfea

import "context"

type PolyfeaApiService struct {
}

func NewPolyfeaAPIService() *PolyfeaApiService {
	return &PolyfeaApiService{}
}

func (s *PolyfeaApiService) GetContextArea(context.Context, string, string, float32) (ImplResponse, error) {
	return ImplResponse{Code: 501}, nil
}

func (s *PolyfeaApiService) GetStaticConfig(context.Context) (ImplResponse, error) {
	return ImplResponse{Code: 501}, nil
}
