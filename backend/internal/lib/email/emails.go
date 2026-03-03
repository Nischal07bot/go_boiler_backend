package email

func (c *Client) SendWelcomemail(to, firstName string) error {
	data := map[string]string{
		"UserfirstName":firstName,
	}

	return c.SendEmail(
		to,
		"welcome to alfred",
		TemplateWelcome,
		data,
	)
}