package main

import (
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
	"finnbank/internal-services/statement/server"
	"finnbank/internal-services/statement/service"
	"github.com/unidoc/unipdf/v3/common/license"
)


func init() {
	key := `b698286ffa1b2759cc3e8a89e29321c3ddda27236602d75e552588d8dbe22cea`
	err := license.SetMeteredKey(key)
	if err != nil {
		panic(err)
	}
}

func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")

	statementService := service.StatementService{
		Logger:                              logger,
		UnimplementedStatementServiceServer: statement.UnimplementedStatementServiceServer{},
	}

	logger.Info("Starting the server...")
	logger.Info("Server running on localhost:8084")
	err = server.StartGrpcServer(statementService, logger)
	if err != nil {
		logger.Fatal("Failed to start gRPC server")
	}
}


// func main() {
// 	http.HandleFunc("/pdf", pdfHandler)
// 	log.Println("Server started at http://localhost:8080/pdf")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// // }

// func drawFrontPage(c *creator.Creator, font, fontBold *model.PdfFont) {
// 	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
// 		p := c.NewStyledParagraph()
// 		p.SetMargins(0, 0, 300, 0)
// 		p.SetTextAlignment(creator.TextAlignmentCenter)

// 		chunk := p.Append("UniPDF")
// 		chunk.Style.Font = font
// 		chunk.Style.FontSize = 56
// 		chunk.Style.Color = creator.ColorRGBFrom8bit(56, 68, 77)

// 		p.Append("\n")
// 		chunk = p.Append("Table features")
// 		chunk.Style.Font = fontBold
// 		chunk.Style.FontSize = 40
// 		chunk.Style.Color = creator.ColorRGBFrom8bit(45, 148, 215)

// 		c.Draw(p)
// 	})
// }

// func pdfHandler() {
// 	// Load fonts.
// 	font, _ := model.NewStandard14Font("Helvetica")
// 	fontBold, _ := model.NewStandard14Font("Helvetica-Bold")

// 	c := creator.New()
// 	c.SetPageMargins(50, 50, 50, 50)

// 	// Add content (simplified here – reuse your own logic like drawFrontPage, etc.)
// 	p := c.NewParagraph("Hello PDF from UniPDF!")
// 	_ = c.NewPage()
// 	_ = c.Draw(p)

// 	drawFrontPage(c, font, fontBold)

// 	// Write PDF to a buffer.
// 	var buf bytes.Buffer
// 	if err := c.Write(&buf); err != nil {
// 		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set response headers and write buffer content to response.
// 	w.Header().Set("Content-Type", "application/pdf")
// 	w.Header().Set("Content-Disposition", "inline; filename=output.pdf")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(buf.Bytes())
// }