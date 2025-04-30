package service

import (
	"context"
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
	"bytes"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

type StatementService struct {
	Logger *utils.Logger
	Grpc   statement.StatementServiceServer
	statement.UnimplementedStatementServiceServer
}

// mustEmbedUnimplementedStatementServiceServer implements statement.StatementServiceServer.


func drawFrontPage(c *creator.Creator, font, fontBold *model.PdfFont) {
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {
		p := c.NewStyledParagraph()
		p.SetMargins(0, 0, 300, 0)
		p.SetTextAlignment(creator.TextAlignmentCenter)

		chunk := p.Append("UniPDF")
		chunk.Style.Font = font
		chunk.Style.FontSize = 56
		chunk.Style.Color = creator.ColorRGBFrom8bit(56, 68, 77)

		p.Append("\n")
		chunk = p.Append("Table features")
		chunk.Style.Font = fontBold
		chunk.Style.FontSize = 40
		chunk.Style.Color = creator.ColorRGBFrom8bit(45, 148, 215)

		c.Draw(p)
	})
}

func (s *StatementService) GenerateStatement(context.Context, *statement.ClientRequest) (*statement.ClientResponse, error) {
	// Load fonts.
	font, _ := model.NewStandard14Font("Helvetica")
	fontBold, _ := model.NewStandard14Font("Helvetica-Bold")

	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Add content (simplified here – reuse your own logic like drawFrontPage, etc.)
	p := c.NewParagraph("Hello PDF from UniPDF!")
	_ = c.NewPage()
	_ = c.Draw(p)

	drawFrontPage(c, font, fontBold)

	// Write PDF to a buffer.
	var buf bytes.Buffer

	if err := c.Write(&buf); err != nil {
		s.Logger.Error("Failed to write PDF to buffer: %v", err)
		return nil, err
	}

	return &statement.ClientResponse{
		PdfBuffer: buf.Bytes(),
	}, nil
}

