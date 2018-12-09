Hello, {{ .FullName }}!

You reserved the {{ .TableID }} table in Palestine Nights restaurant.

Details:

Date and time: {{ .FormattedTime }}
Duration: {{ .FormattedDuration }} minutes.
Guests: {{ .Guests }}

Please, confirm your reservation by clicking this link:
http://localhost:4000/confirm/{{ .ConfirmationCode }}

You can cancel your reservation at any time by clicking this link:
http://localhost:4000/cancel/{{ .CancellationCode }}
