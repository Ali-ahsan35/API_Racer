package request

type RealClient struct{}

func (c *RealClient) Fetch(url string) (*APIResponse, error) {
	return FetchAPI(url)
}