package usecase

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"ahu-backend/internal/repository"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type InspectionPDFService struct {
	inspectionRepo repository.InspectionRepository
}

func NewInspectionPDFService(repo repository.InspectionRepository) *InspectionPDFService {
	return &InspectionPDFService{inspectionRepo: repo}
}

func (s *InspectionPDFService) GenerateInspectionPDF(inspectionID string) error {
	fmt.Println("🔥 GENERATING COMPACT ONE-PAGE PDF WITH QR:", inspectionID)

	os.MkdirAll("files/inspection", 0755)
	os.MkdirAll("files/tmp", 0755)

	report, err := s.inspectionRepo.GetInspectionReport(inspectionID)
	if err != nil || report == nil {
		return fmt.Errorf("report not found")
	}

	// --- 1. GENERATE QR CODE ---
	localIP := "10.9.118.16"

	// QR diarahkan ke aplikasi Next.js (Port 3000)
	verifyURL := fmt.Sprintf("http://%s:3000/verify/inspection/%s", localIP, inspectionID)

	qrPath := fmt.Sprintf("files/tmp/qr-%s.png", inspectionID)
	err = qrcode.WriteFile(verifyURL, qrcode.Medium, 256, qrPath)
	if err != nil {
		return fmt.Errorf("failed to generate QR: %v", err)
	}

	// Inisialisasi PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	// --- DEFINISI WARNA ---
	themeColor := func() { pdf.SetTextColor(13, 148, 136) }
	headerBg := func() { pdf.SetFillColor(241, 245, 249) }
	textColor := func() { pdf.SetTextColor(30, 41, 59) }
	mutedText := func() { pdf.SetTextColor(100, 116, 139) }

	// --- 2. HEADER (LOGO & JUDUL) ---
	logoPath := "public/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.ImageOptions(logoPath, 10, 10, 35, 0, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	}

	pdf.SetXY(48, 11)
	pdf.SetFont("Arial", "B", 14)
	themeColor()
	pdf.CellFormat(100, 6, "INSPECTION REPORT", "", 1, "L", false, 0, "")

	pdf.SetX(48)
	pdf.SetFont("Arial", "B", 7.5)
	mutedText()
	pdf.CellFormat(100, 4, "AIR HANDLING UNIT SYSTEM - PREVENTIVE MAINTENANCE", "", 1, "L", false, 0, "")

	// 🔥 SISIPKAN QR CODE (Di pojok kanan atas agar tidak tertutup kotak info)
	pdf.ImageOptions(qrPath, 170, 39, 28, 28, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	pdf.SetFont("Arial", "B", 5)
	pdf.SetXY(170, 67)
	pdf.CellFormat(28, 3, "SCAN TO VERIFY", "", 0, "C", false, 0, "")

	pdf.Ln(10)

	// --- 3. DYNAMIC INFO BOX ---
	pdf.SetDrawColor(226, 232, 240)
	headerBg()
	pdf.Rect(10, 39, 155, 28, "F") // Lebar dikurangi (155mm) agar tidak menabrak QR

	// Mapping Periode Otomatis
	displayPeriod := "-"
	p := strings.TrimSpace(strings.ToLower(report.Period))
	if strings.Contains(p, "bulan") || strings.Contains(p, "month") {
		displayPeriod = "Monthly (1 Month)"
	} else if strings.Contains(p, "6") {
		displayPeriod = "6 Months"
	} else if strings.Contains(p, "tahunan") || strings.Contains(p, "year") {
		displayPeriod = "Yearly (1 Year)"
	} else {
		displayPeriod = report.Period
	}

	textColor()
	pdf.SetFont("Arial", "B", 7)

	// Data Row 1
	pdf.SetXY(12, 41)
	pdf.CellFormat(25, 5, "ID Unit AHU", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(55, 5, ": "+report.UnitCode, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Document No.", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(45, 5, ": F-I-EM-04-061-01/02", "", 1, "L", false, 0, "")

	// Data Row 2
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Area Name", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(55, 5, ": "+report.AreaName, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Effective Date", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(45, 5, ": 30 AUG 2023", "", 1, "L", false, 0, "")

	// Data Row 3
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Room Location", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(55, 5, ": "+report.RoomName, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Period", "", 0, "L", false, 0, "")
	themeColor()
	pdf.CellFormat(45, 5, ": "+displayPeriod, "", 1, "L", false, 0, "")

	// Data Row 4
	textColor()
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Manufacture", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(55, 5, ": "+report.Vendor, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Cleanliness Class", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(45, 5, ": "+report.CleanlinessClass, "", 1, "L", false, 0, "")

	// Data Row 5
	pdf.SetX(12)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Date Executed", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	tglStr := "-"
	if report.InspectedAt != nil {
		tglStr = report.InspectedAt.Format("02 January 2006")
	}
	pdf.CellFormat(55, 5, ": "+tglStr, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(25, 5, "Revision", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(45, 5, ": 02", "", 1, "L", false, 0, "")

	pdf.SetXY(10, 72) // Mulai tabel di bawah kotak info

	// --- 4. INSPECTION TABLE ---
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(51, 65, 85)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(10, 7, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(100, 7, "ACTIVITY ITEM", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, "ACTUAL VALUE", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, "RESULT", "1", 1, "C", true, 0, "")

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
			h := 5.8
			pdf.CellFormat(10, h, fmt.Sprint(no), "1", 0, "C", true, 0, "")
			pdf.CellFormat(100, h, " "+it.Label, "1", 0, "L", true, 0, "")
			pdf.CellFormat(40, h, it.Value, "1", 0, "C", true, 0, "")

			resClean := strings.ToLower(it.Result)
			if resClean == "pass" || resClean == "ok" {
				pdf.SetTextColor(13, 148, 136)
			} else {
				pdf.SetTextColor(220, 38, 38)
			}
			pdf.SetFont("Arial", "B", 7)
			pdf.CellFormat(40, h, strings.ToUpper(it.Result), "1", 1, "C", true, 0, "")

			textColor()
			pdf.SetFont("Arial", "", 7)
			no++
		}
	}

	// --- 5. SIGNATURE SECTION ---
	pdf.SetY(250) // Posisikan di dasar halaman
	ySign := pdf.GetY()
	pdf.SetDrawColor(203, 213, 225)

	// Box Inspector
	pdf.SetXY(20, ySign)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(60, 5, "Executed By (Inspector)", "", 1, "C", false, 0, "")
	if report.Signature != "" {
		// Pastikan koordinat Y adalah ySign + 4 agar pas di tengah
		s.embedSignature(pdf, report.Signature, 30, ySign+4, 40, 15, inspectionID+"-ins")
	}
	pdf.SetXY(20, ySign+20)
	pdf.CellFormat(60, 5, report.Inspector, "T", 1, "C", false, 0, "")
	mutedText()
	pdf.SetX(20)
	pdf.CellFormat(60, 4, "Date: "+tglStr, "", 1, "C", false, 0, "")

	// Box Supervisor (Kanan)
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

	// --- 6. FOOTER ---
	pdf.SetY(-8)
	pdf.SetFont("Arial", "I", 6)
	footerTxt := fmt.Sprintf("AIRA SECURE DOC ID: %s | PT Kimia Farma Hub", inspectionID)
	pdf.CellFormat(0, 5, footerTxt, "0", 0, "C", false, 0, "")

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
