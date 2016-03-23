package engine

const (
	HtmlHeader = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>{{.Title}}</title>
    <style>
body{
  background-color:{{.BackgroundColor}};
  color:{{.TextColor}};
}
div{
  font-size:{{.FontSize}};
  font-family:monospace;
  white-space:pre;
  min-height:{{.FontSize}};
}
    </style>
  <head>
<body>
`
	HtmlFooter = `
<!-- {{.RowCount}} rows -->
</body>
</html>
`
)
