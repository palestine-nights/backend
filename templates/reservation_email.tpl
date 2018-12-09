Hello, {{.FullName}}!
You reserved a table in Palestine Nights restaurant, here is your reservation details:
Date and time: {{.FormattedTime}}
Duration: {{.FormattedDuration}} minutes
Guests: {{.Guests}}

Please, confirm your reservation by clicking this link:
https://api.palestinenights.com/reservations/confirm?code={{.ConfirmationCode}}


You can cancel your reservation at any time by clicking this link:
https://api.palestinenights.com/reservations/cancel?code={{.CancellationCode}}