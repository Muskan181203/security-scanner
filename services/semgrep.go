package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github-security-scanner/models"
)

func RunSemgrep() ([]models.Vulnerability, error) {

	cmd := exec.Command(
		"semgrep",
		"scan",
		"./repos/project",
		"--json",
		"--quiet",
	)

	cmd.Env = append(os.Environ(),
		"PYTHONWARNINGS=ignore",
	)

	output, err := cmd.CombinedOutput()

	fmt.Println("OUTPUT LENGTH:", len(output))

	if err != nil {
		fmt.Println("SEMGREP ERROR:", err)
		fmt.Println(string(output))
		fmt.Println(string(output[:2000]))
	}

	var semgrepResponse models.SemgrepResponse

	err = json.Unmarshal(output, &semgrepResponse)
	fmt.Printf("%+v\n", semgrepResponse.Results[0])
	if err != nil {
		return nil, err
	}

	vulnerabilities := []models.Vulnerability{}

	for _, result := range semgrepResponse.Results {
		vulnerabilities = append(vulnerabilities, models.Vulnerability{
			Type:        result.CheckID,
			File:        result.Path,
			Line:        result.Start.Line,
			Description: result.Extra.Message,
			Severity:    result.Extra.Severity,
		})
	}

	return vulnerabilities, nil
}

// package services

// import (
// 	"encoding/json"
// 	"fmt"
// 	"github-security-scanner/models"
// 	"os/exec"
// )

// func RunSemgrep() ([]models.Vulnerability, error) {

// 	// cmd := exec.Command(
// 	// 	"semgrep",
// 	// 	"scan",
// 	// 	"./repos/project",
// 	// 	"--json",
// 	// )
// 	cmd := exec.Command(
// 		"semgrep",
// 		"scan",
// 		"./repos/project",
// 		"--json",
// 		"--quiet",
// 	)

// 	cmd.Env = append(cmd.Env,
// 		"PYTHONWARNINGS=ignore",
// 	)
// 	output, err := cmd.CombinedOutput()

// 	if err != nil {
// 		fmt.Println("SEMGREP ERROR:", err)
// 	}
// 	var semgrepResponse models.SemgrepResponse
// 	fmt.Println(string(output))
// 	fmt.Println("FIRST 500 CHARS:")

// 	err = json.Unmarshal(output, &semgrepResponse)
// 	// err := json.Unmarshal(output, &semgrepResponse)

// 	if err != nil {
// 	fmt.Println("SEMGREP ERROR:", err)
// 	fmt.Println("OUTPUT:")
// 	fmt.Println(string(output))
// }
// 	vulnerabilities := []models.Vulnerability{}
// 	for _, result := range semgrepResponse.Results {
// 		vulnerabilities = append(vulnerabilities, models.Vulnerability{
// 			Type:        result.CheckId,
// 			File:        result.Path,
// 			Description: result.Extra.Message,
// 		})
// 	}

// 	return vulnerabilities, nil
// }
