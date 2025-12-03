package api

// GetBalance retrieves the current account balance
func (c *Client) GetBalance() (float64, error) {
	var resp BalanceResponse
	err := c.do("GET", "/balance", nil, &resp)
	return resp.Balance, err
}
