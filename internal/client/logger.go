package fractalCloud

func (c *Client) logDebug(s string) {
	if c.Logger == nil || c.Logger.Debug == nil {
		return
	}
	c.Logger.Debug(s)
}

func (c *Client) logInformation(s string) {
	if c.Logger == nil || c.Logger.Information == nil {
		return
	}
	c.Logger.Information(s)
}

func (c *Client) logWarning(s string) {
	if c.Logger == nil || c.Logger.Warning == nil {
		return
	}
	c.Logger.Warning(s)
}

func (c *Client) logError(s string) {
	if c.Logger == nil || c.Logger.Error == nil {
		return
	}
	c.Logger.Error(s)
}
