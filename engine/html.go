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
.l{
  font-size:10pt;
  font-family:monospace;
  white-space:pre;
  min-height: 1em;
}
.panel{
  display:none;
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
