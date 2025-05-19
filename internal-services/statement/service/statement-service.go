package service

import (
	"bytes"
	"context"
	"encoding/json"
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

type StatementService struct {
	Logger *utils.Logger
	Grpc   statement.StatementServiceServer
	statement.UnimplementedStatementServiceServer
}

type FindBothAccountNumber []struct {
	OpenedAccountID int    `json:"openedaccount_id"`
	AccountNumber   string `json:"account_number"`
}

// Transaction struct must match GraphQL response fields exactly
type Transaction struct {
	Amount            float64   `json:"amount"`
	RefNo             string    `json:"ref_no"`
	Notes             string    `json:"notes"`
	DateTransaction   time.Time `json:"date_transaction"`
	ReceiverID        int       `json:"receiver_id"`
	SenderID          int       `json:"sender_id"`
	TransactionType   string    `json:"transaction_type"`
	TransactionStatus string    `json:"transaction_status"`
	TransactionID     int       `json:"transaction_id"`
	TransactionFee    float64   `json:"transaction_fee"`
}

func fetchTransactions(ctx context.Context, creditId, savingsId, debitId int32, startTime, endTime string) ([]Transaction, error) {
	query := `
		query($creditId: Int!, $debitId: Int!, $savingsId: Int!, $startTime: DateTime!, $endTime: DateTime!) {
			getTransactionsByTimeStampByUserId(
				creditId: $creditId,
				debitId: $debitId,
				savingsId: $savingsId,
				startTime: $startTime,
				endTime: $endTime
			) {
				amount
				ref_no
				notes
				date_transaction
				receiver_id
				sender_id
				transaction_type
				transaction_status
				transaction_id
				transaction_fee
			}
		}
	`

	fmt.Println("Parsed Start Time:", startTime)
	fmt.Println("Parsed End Time:", endTime)

	vars := map[string]interface{}{
		"creditId":  creditId,
		"debitId":   debitId,
		"savingsId": savingsId,
		"startTime": startTime,
		"endTime":   endTime,
	}

	payload := map[string]interface{}{
		"query":     query,
		"variables": vars,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", "http://localhost:8083/graphql/transaction", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var result struct {
		Data struct {
			GetTransactionsByTimeStampByUserId []Transaction `json:"getTransactionsByTimeStampByUserId"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %+v", result.Errors)
	}

	// Filter out failed transactions
	successfulTransactions := make([]Transaction, 0)
	for _, tx := range result.Data.GetTransactionsByTimeStampByUserId {
		if tx.TransactionStatus != "FAILED" {
			successfulTransactions = append(successfulTransactions, tx)
		}
	}

	return successfulTransactions, nil
}

func fetchAccountNumbers(ctx context.Context, sender_id int, receiver_id int) (FindBothAccountNumber, error) {

	query := fmt.Sprintf(`{
		find_both_account_num(sender_id: %d, receiver_id: %d) {
			openedaccount_id
			account_number
		}
	}`, sender_id, receiver_id)

	vars := map[string]interface{}{
		"sender_id":   sender_id,
		"receiver_id": receiver_id,
	}

	payload := map[string]interface{}{
		"query":     query,
		"variables": vars,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", "http://localhost:8083/graphql/opened-account", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var result struct {
		Data struct {
			FindBothAccountNumber FindBothAccountNumber `json:"find_both_account_num"`
		} `json:"data"`
		Errors []interface{} `json:"errors"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %+v", result.Errors)
	}

	return result.Data.FindBothAccountNumber, nil
}

func drawFrontPage(c *creator.Creator, img *creator.Image, font, fontBold *model.PdfFont) {
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		if err := c.Draw(img); err != nil {
			fmt.Printf("Error drawing image: %v\n", err)
		}

		p := c.NewStyledParagraph()
		p.SetMargins(0, 0, 300, 0)
		p.SetTextAlignment(creator.TextAlignmentCenter)

		chunk := p.Append("   FinnBank")
		chunk.Style.Font = fontBold
		chunk.Style.FontSize = 56
		chunk.Style.Color = creator.ColorRGBFrom8bit(45, 148, 215)

		c.Draw(p)
	})
}

func drawTitle(c *creator.Creator, fontBold *model.PdfFont) {
	title := c.NewParagraph("Transaction Statement")
	title.SetFont(fontBold)
	title.SetFontSize(22)
	title.SetMargins(0, 0, 20, 20)
	title.SetTextAlignment(creator.TextAlignmentLeft)
	_ = c.Draw(title)
}

func drawTableHeaders(table *creator.Table, c *creator.Creator, fontBold *model.PdfFont) {

	drawCell := func(text string, font *model.PdfFont, align creator.TextAlignment) {
		p := c.NewStyledParagraph()
		p.SetTextAlignment(align)
		p.SetMargins(2, 2, 0, 0)
		chunk := p.Append(text)
		chunk.Style.Font = font
		chunk.Style.FontSize = 8
		chunk.Style.Color = creator.ColorRGBFrom8bit(255, 255, 255)

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetContent(p)
		cell.SetIndent(0)
		cell.SetBackgroundColor(creator.ColorRGBFrom8bit(51, 153, 255))
	}

	drawCell("Date", fontBold, creator.TextAlignmentCenter)
	drawCell("Ref No", fontBold, creator.TextAlignmentCenter)
	drawCell("Transaction", fontBold, creator.TextAlignmentCenter)
	drawCell("Description", fontBold, creator.TextAlignmentCenter)
	drawCell("Amount(PhP)", fontBold, creator.TextAlignmentCenter)
	drawCell("Fee", fontBold, creator.TextAlignmentCenter)
	drawCell("Sender", fontBold, creator.TextAlignmentCenter)
	drawCell("Receiver", fontBold, creator.TextAlignmentCenter)

}

func drawTable(c *creator.Creator, font, fontBold *model.PdfFont, transactions []Transaction, accountNumbers FindBothAccountNumber) {
	table := c.NewTable(8)
	//
	table.SetColumnWidths(0.11, 0.12, 0.14, 0.15, 0.11, 0.09, 0.15, 0.15)
	table.SetMargins(10, 10, 10, 10)

	drawTableHeaders(table, c, fontBold)

	drawCell := func(text string, font *model.PdfFont, align creator.TextAlignment) {
		p := c.NewStyledParagraph()
		p.SetTextAlignment(align)
		p.SetMargins(2, 2, 0, 0)
		chunk := p.Append(text)
		chunk.Style.Font = font
		chunk.Style.FontSize = 7

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetContent(p)
		cell.SetIndent(0)

		if table.CurRow()%2 == 0 {
			cell.SetBackgroundColor(creator.ColorRGBFrom8bit(240, 240, 255))
		} else {
			cell.SetBackgroundColor(creator.ColorRGBFrom8bit(250, 250, 255))
		}
	}

	for _, tx := range transactions {
		var sender, receiver string

		if accountNumbers[0].OpenedAccountID == tx.SenderID {
			sender = accountNumbers[0].AccountNumber
			receiver = accountNumbers[1].AccountNumber
		} else {
			sender = accountNumbers[1].AccountNumber
			receiver = accountNumbers[0].AccountNumber
		}

		drawCell(tx.DateTransaction.Format("2006-01-02"), font, creator.TextAlignmentCenter)
		drawCell(tx.RefNo, font, creator.TextAlignmentCenter)
		drawCell(fmt.Sprintf("%s (%d)", tx.TransactionType, tx.TransactionID), font, creator.TextAlignmentLeft)
		drawCell(tx.Notes, font, creator.TextAlignmentLeft)
		drawCell(fmt.Sprintf("Php %.2f", tx.Amount), font, creator.TextAlignmentLeft)
		drawCell(fmt.Sprintf("Php %.2f", tx.TransactionFee), font, creator.TextAlignmentLeft)
		drawCell(sender, font, creator.TextAlignmentCenter)
		drawCell(receiver, font, creator.TextAlignmentCenter)
	}

	_ = c.Draw(table)
}

func (s *StatementService) GenerateStatement(ctx context.Context, req *statement.ClientRequest) (*statement.ClientResponse, error) {

	// font
	font, err := model.NewStandard14Font("Helvetica")
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}
	fontBold, err := model.NewStandard14Font("Helvetica-Bold")
	if err != nil {
		return nil, fmt.Errorf("failed to load bold font: %w", err)
	}

	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	img, err := c.NewImageFromFile("./assets/finnbank-logo.png") // adjust path if needed
	if err != nil {
		s.Logger.Error("Failed to load image: %v", err)
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	// Get page and image dimensions
	pageWidth := c.Context().PageWidth
	pageHeight := c.Context().PageHeight
	imgWidth := img.Width()
	imgHeight := img.Height()
	// Calculate the center position
	centerX := (pageWidth - imgWidth) / 2
	centerY := (pageHeight - imgHeight) / 2

	// Set the image position to the center of the page
	img.SetPos(centerX, centerY)
	img.ScaleToWidth(200) // adjust size as needed
	img.SetPos(50, 290)   // adjust position on page
	encoder := core.NewFlateEncoder()
	img.SetEncoder(encoder)

	// Draw the front page
	// This function creates a styled front page with the title "UniPDF" and a subtitle "Table features"
	drawFrontPage(c, img, font, fontBold)

	// Adjust as per correct userId
	transactions, err := fetchTransactions(ctx, req.CreditId, req.SavingsId, req.DebitId, req.StartDate, req.EndDate)
	accountNumber, err1 := fetchAccountNumbers(ctx, int(req.CreditId), int(req.DebitId))

	if err != nil || err1 != nil {
		s.Logger.Error("Failed to fetch transactions or account numbers: %v %v", err, err1)
		return nil, err
	}

	// Draw title and transaction table
	drawTitle(c, fontBold)
	drawTable(c, font, fontBold, transactions, accountNumber)

	var buf bytes.Buffer
	if err := c.Write(&buf); err != nil {
		s.Logger.Error("Failed to write PDF: %v", err)
		return nil, err
	}

	return &statement.ClientResponse{
		PdfBuffer: buf.Bytes(),
	}, nil
}
