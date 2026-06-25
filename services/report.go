package services
import (
	"fmt"
	"os"
	"github-security-scanner/models"
)
var rows string
var summary string

func GenerateReport(scan models.ScanResponse) (string,error){
	fileName := "report.html"
	for vulnType, count := range scan.SummaryByType {

	summary += fmt.Sprintf(
		"<li>%s : %d</li>",
		vulnType,
		count,
	)
}
for _, vuln := range scan.Vulnerabilities {

	rows += fmt.Sprintf(`
	<tr>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%d</td>
		<td>%s</td>
	</tr>
	`,
		vuln.Severity,
		vuln.Type,
		vuln.File,
		vuln.Line,
		vuln.Description,
	)
}
   html := fmt.Sprintf(`
<!DOCTYPE html>

<html>

<head>

<title>Security Report</title>

<style>

body{
	font-family: Arial;
	margin:40px;
}

table{
	width:100%%;
	border-collapse:collapse;
}

th,td{
	border:1px solid #ccc;
	padding:8px;
}

th{
	background:#f5f5f5;
}

</style>

</head>

<body>

<h1>GitHub Security Report</h1>

<hr>

<h2>Repository</h2>

<p>%s</p>

<h2>Risk Summary</h2>

<p><b>Risk Score:</b> %d</p>

<p><b>Risk Level:</b> %s</p>

<p><b>Total Vulnerabilities:</b> %d</p>

<p><b>Errors:</b> %d</p>

<p><b>Warnings:</b> %d</p>

<p><b>Info:</b> %d</p>

<p><b>Scan Duration:</b> %s</p>

<h2>Summary By Type</h2>

<ul>
%s
</ul>

<h2>Vulnerabilities</h2>

<table>

<tr>
<th>Severity</th>
<th>Type</th>
<th>File</th>
<th>Line</th>
<th>Description</th>
</tr>

%s

</table>

</body>

</html>
`,
	scan.RepoURL,
	scan.RiskScore,
	scan.RiskLevel,
	scan.TotalVulnerabilities,
	scan.ErrorCount,
	scan.WarningCount,
	scan.InfoCount,
	scan.ScanDuration,
	summary,
	rows,
)


		err := os.WriteFile(
			fileName,
			[]byte(html),
			0644,
		)
		if err!= nil{
			return "", err
		}

		

	
	return fileName , nil
}