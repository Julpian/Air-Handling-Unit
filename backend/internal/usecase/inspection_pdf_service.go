package usecase

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"ahu-backend/internal/repository"

	"github.com/jung-kurt/gofpdf"
)

type InspectionPDFService struct {
	inspectionRepo repository.InspectionRepository
}

func NewInspectionPDFService(repo repository.InspectionRepository) *InspectionPDFService {
	return &InspectionPDFService{inspectionRepo: repo}
}

func (s *InspectionPDFService) GenerateInspectionPDF(inspectionID string) error {
	fmt.Println("🔥 GENERATING COMPACT ONE-PAGE PDF:", inspectionID)

	os.MkdirAll("files/inspection", 0755)
	os.MkdirAll("files/tmp", 0755)

	report, err := s.inspectionRepo.GetInspectionReport(inspectionID)
	if err != nil || report == nil {
		return fmt.Errorf("report not found")
	}

	// Inisialisasi PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(false, 0) // Paksa satu halaman
	pdf.AddPage()

	// --- DEFINISI WARNA ---
	themeColor := func() { pdf.SetTextColor(13, 148, 136) } // Dark Teal
	headerBg := func() { pdf.SetFillColor(241, 245, 249) }  // Slate Light Gray
	textColor := func() { pdf.SetTextColor(30, 41, 59) }    // Dark Slate
	mutedText := func() { pdf.SetTextColor(100, 116, 139) } // Muted Gray

	// --- 1. HEADER (LOGO & JUDUL) ---
	logoPath := "public/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.ImageOptions(logoPath, 165, 8, 35, 0, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	}

	pdf.SetFont("Arial", "B", 14)
	themeColor()
	pdf.SetXY(10, 10)
	pdf.CellFormat(100, 6, "INSPECTION REPORT", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7.5)
	mutedText()
	pdf.CellFormat(100, 4, "AIR HANDLING UNIT SYSTEM - PREVENTIVE MAINTENANCE", "", 1, "L", false, 0, "")

	pdf.Ln(2)

	// --- 2. INTEGRATED INFO BOX (PADAT & SELARAS) ---
	pdf.SetDrawColor(226, 232, 240)
	headerBg()
	pdf.Rect(10, 22, 190, 32, "F")

	// 🔥 LOGIKA PERIODE OTOMATIS (Ambil langsung dari data record)
	displayPeriod := "-"
	p := strings.TrimSpace(strings.ToLower(report.Period))
	switch p {
	case "bulanan":
		displayPeriod = "1 Month"
	case "enam_bulan":
		displayPeriod = "6 Months"
	case "tahunan":
		displayPeriod = "1 Year"
	default:
		// Jika data di DB masih 'monthly' dll, tetap tercover
		if strings.Contains(p, "month") || strings.Contains(p, "bulan") {
			displayPeriod = "1 Month"
		} else if strings.Contains(p, "6") {
			displayPeriod = "6 Months"
		} else {
			displayPeriod = report.Period // Fallback ke teks asli jika tidak cocok
		}
	}

	textColor()
	// Row 1
	pdf.SetXY(12, 24)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "ID Unit AHU", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": "+report.UnitCode, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Document No.", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": F-I-EM-04-061-01/02", "", 1, "L", false, 0, "")

	// Row 2
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Area Name", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": "+report.AreaName, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Effective Date", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": 30 AUG 2023", "", 1, "L", false, 0, "")

	// Row 3
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Room Location", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": "+report.RoomName, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Period", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 7)
	themeColor()
	pdf.CellFormat(65, 5, ": "+displayPeriod, "", 1, "L", false, 0, "") // 🔥 Menampilkan Periode Otomatis

	// Row 4
	textColor()
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Manufacture", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": "+report.Vendor, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Cleanliness Class", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": "+report.CleanlinessClass, "", 1, "L", false, 0, "")

	// Row 5
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Date Executed", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	tglStr := "-"
	if report.InspectedAt != nil {
		tglStr = report.InspectedAt.Format("02 January 2006")
	}
	pdf.CellFormat(65, 5, ": "+tglStr, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(30, 5, "Revision", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 5, ": 02", "", 1, "L", false, 0, "")

	pdf.Ln(4)

	// --- 3. COMPACT INSPECTION TABLE ---
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(51, 65, 85) // Dark Slate
	pdf.SetTextColor(255, 255, 255)

	colWidths := []float64{10, 100, 40, 40}
	pdf.CellFormat(colWidths[0], 7, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[1], 7, "ACTIVITY ITEM", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], 7, "ACTUAL VALUE", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[3], 7, "RESULT", "1", 1, "C", true, 0, "")

	no := 1
	textColor()
	for _, sec := range report.Sections {
		pdf.SetFont("Arial", "B", 7.5)
		pdf.SetFillColor(226, 232, 240)
		pdf.CellFormat(190, 6, "  Section "+sec.Code+": "+sec.Title, "1", 1, "L", true, 0, "")

		pdf.SetFont("Arial", "", 7)
		for _, it := range sec.Items {
			if no%2 == 0 {
				pdf.SetFillColor(255, 255, 255)
			} else {
				pdf.SetFillColor(248, 250, 252)
			}

			h := 6.0
			pdf.CellFormat(colWidths[0], h, fmt.Sprint(no), "1", 0, "C", true, 0, "")
			pdf.CellFormat(colWidths[1], h, " "+it.Label, "1", 0, "L", true, 0, "")
			pdf.CellFormat(colWidths[2], h, it.Value, "1", 0, "C", true, 0, "")

			resClean := strings.ToLower(it.Result)
			if resClean == "pass" || resClean == "ok" {
				pdf.SetTextColor(13, 148, 136)
			} else {
				pdf.SetTextColor(220, 38, 38)
			}
			pdf.SetFont("Arial", "B", 7)
			pdf.CellFormat(colWidths[3], h, strings.ToUpper(it.Result), "1", 1, "C", true, 0, "")

			textColor()
			pdf.SetFont("Arial", "", 7)
			no++
		}
	}

	// --- 4. SIGNATURE SECTION (FIXED POSITION) ---
	pdf.SetY(248)
	pdf.SetDrawColor(203, 213, 225)
	ySign := pdf.GetY()

	// Box Inspector
	pdf.SetXY(20, ySign)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(60, 5, "Executed By (Inspector)", "", 1, "C", false, 0, "")
	if report.Signature != "" {
		s.embedSignature(pdf, report.Signature, 30, ySign+4, 40, 15, inspectionID+"-ins")
	}
	pdf.SetXY(20, ySign+20)
	pdf.CellFormat(60, 5, report.Inspector, "T", 1, "C", false, 0, "")
	mutedText()
	pdf.SetX(20)
	pdf.CellFormat(60, 4, "Date: "+tglStr, "", 1, "C", false, 0, "")

	// Box Supervisor
	textColor()
	pdf.SetXY(130, ySign)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(60, 5, "Verified By (Supervisor)", "", 1, "C", false, 0, "")
	if report.SPVSignature != "" {
		s.embedSignature(pdf, report.SPVSignature, 140, ySign+4, 40, 15, inspectionID+"-spv")
	}
	pdf.SetXY(130, ySign+20)
	pdf.CellFormat(60, 5, "Asep Suherman", "T", 1, "C", false, 0, "")
	mutedText()
	pdf.SetX(130)
	pdf.CellFormat(60, 4, "Digitally Approved", "", 1, "C", false, 0, "")

	// --- 5. FOOTER ---
	pdf.SetY(-10)
	pdf.SetFont("Arial", "I", 6.5)
	footerTxt := fmt.Sprintf("AIRA System - Generated on %s | PT Kimia Farma Industrial Hub", time.Now().Format("2006-01-02 15:04"))
	pdf.CellFormat(0, 10, footerTxt, "0", 0, "C", false, 0, "")

	path := fmt.Sprintf("files/inspection/%s.pdf", inspectionID)
	return pdf.OutputFileAndClose(path)
}

func (s *InspectionPDFService) embedSignature(pdf *gofpdf.Fpdf, base64Str string, x, y, w, h float64, name string) {
	clean := strings.TrimPrefix(base64Str, "data:image/png;base64,")
	img, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return
	}
	tmpPath := fmt.Sprintf("files/tmp/%s.png", name)
	os.WriteFile(tmpPath, img, 0644)
	pdf.ImageOptions(tmpPath, x, y, w, h, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
}
