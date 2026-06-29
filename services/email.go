package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, pdfFile string) error {

	// -----------------------------
	// Your Gmail account
	// -----------------------------

	from := "mneema@bestpeers.in"

	// Gmail App Password (NOT your Gmail password)
	password := "zukyiimpjopamjtk"

	// -----------------------------
	// Read PDF
	// -----------------------------

	pdfData, err := os.ReadFile(pdfFile)
	if err != nil {
		return err
	}

	// Encode PDF
	encoded := base64.StdEncoding.EncodeToString(pdfData)

	// -----------------------------
	// Email Headers
	// -----------------------------

	boundary := "my-boundary"

	header := make(map[string]string)

	header["From"] = from
	header["To"] = to
	header["Subject"] = "GitHub Security Report"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "multipart/mixed; boundary=" + boundary

	var message bytes.Buffer

	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	message.WriteString("\r\n")

	// -----------------------------
	// Email Body
	// -----------------------------

	message.WriteString("--" + boundary + "\r\n")
	message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")

	message.WriteString("Hello,\n\n")
	message.WriteString("Please find your GitHub Security Report attached.\n\n")
	message.WriteString("Regards,\n")
	message.WriteString("GitHub Security Scanner\n")

	// -----------------------------
	// PDF Attachment
	// -----------------------------

	message.WriteString("\r\n--" + boundary + "\r\n")

	message.WriteString("Content-Type: application/pdf\r\n")
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("Content-Disposition: attachment; filename=\"security-report.pdf\"\r\n\r\n")

	for i := 0; i < len(encoded); i += 76 {

		end := i + 76

		if end > len(encoded) {
			end = len(encoded)
		}

		message.WriteString(encoded[i:end] + "\r\n")
	}

	message.WriteString("--" + boundary + "--")

	// -----------------------------
	// Send Email
	// -----------------------------

	auth := smtp.PlainAuth(
		"",
		from,
		password,
		"smtp.gmail.com",
	)

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{to},
		message.Bytes(),
	)

	return err
}
