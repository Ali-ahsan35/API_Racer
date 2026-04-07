package request

type APIClient interface {
	Fetch(url string) (*APIResponse, error)
}