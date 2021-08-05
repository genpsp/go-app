package admin_response

import (
	entities "github.com/genpsp/go-app/domain/entities"
)

type ItemResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

type ItemsResponse struct {
	Items []*ItemResponse `json:"items"`
}

func ConvertItemResponse(entity entities.Item) *ItemResponse {
	return &ItemResponse{
		ID:   entity.ID,
		Name: entity.Name,
	}
}

func ConvertItemsResponse(entities *[]entities.Item) []*ItemResponse {
	list := make([]*ItemResponse, len(*entities), len(*entities))
	for i, entity := range *entities {
		list[i] = ConvertItemResponse(entity)
	}
	return list
}
