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
span.blink{
  animation:         blink-animation 1s infinite;
  -webkit-animation: blink-animation 1s infinite;
}
@keyframes blink-animation {
  0% { visibility: hidden; }
  50% { visibility: hidden; }
}
@-webkit-keyframes blink-animation {
  0% { visibility: hidden; }
  50% { visibility: hidden; }
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
