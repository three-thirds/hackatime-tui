type Client struct {
	creds      *Creds
	httpClient *http.Client
}

func NewClient(creds *Creds) *Client {
	return &Client{
		creds:      creds,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetStatusToday() (*StatusToday, error) {
	url := c.creds.APIURL + "/users/current/statusbar/today"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+c.creds.APIKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bleh bleh bleh skill issue")
	}

	var status StatusToday
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}