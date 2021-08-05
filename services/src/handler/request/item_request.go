package request

type CreateItemRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type GetItemRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
