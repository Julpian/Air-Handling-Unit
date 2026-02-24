package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type SchedulePDFService struct {
	scheduleRepo repository.ScheduleRepository
}

func NewSchedulePDFService(repo repository.ScheduleRepository) *SchedulePDFService {
	return &SchedulePDFService{scheduleRepo: repo}
}

// Warna Professional & Clinical
var (
	clrPrimary   = []int{15, 23, 42}    // Slate 900
	clrSecondary = []int{100, 116, 139} // Slate 500

	clrBlue  = []int{219, 234, 254} // Background Siap Diperiksa
	clrBlueT = []int{30, 64, 175}   // Text Siap Diperiksa

	clrAmber  = []int{254, 243, 199} // Background Dalam Pemeriksaan
	clrAmberT = []int{146, 64, 14}   // Text Dalam Pemeriksaan

	clrGreen  = []int{220, 252, 231} // Background Selesai
	clrGreenT = []int{22, 101, 52}   // Text Selesai

	clrWeekend = []int{248, 250, 252}
	clrBorder  = []int{226, 232, 240}
)

func (s *SchedulePDFService) Generate(year int, approval *domain.ScheduleApproval, path string) error {
	schedules, err := s.scheduleRepo.ListWithDetailByYear(year)
	if err != nil {
		return err
	}

	calendar := map[string][]string{}
	for _, sch := range schedules {
		for d := sch.StartDate; !d.After(sch.EndDate); d = d.AddDate(0, 0, 1) {
			key := d.Format("2006-01-02")
			calendar[key] = append(calendar[key], sch.Status)
		}
	}

	pdf := gofpdf.New("L", "mm", "A3", "") // 420 x 297 mm
	pdf.SetMargins(20, 15, 20)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	s.drawHeader(pdf, year)
	s.drawLegend(pdf)
	pdf.Ln(8)

	pageW, pageH := pdf.GetPageSize()
	marginL, _, marginR, _ := pdf.GetMargins()
	usableW := pageW - marginL - marginR

	// Penentuan tata letak agar tidak tertimpa
	footerReservedH := 65.0
	startYGrid := pdf.GetY()
	availableGridH := pageH - startYGrid - footerReservedH

	months := []string{"JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI", "JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	weekHeaders := []string{"SEN", "SEL", "RAB", "KAM", "JUM", "SAB", "MIN"}

	// Ukuran kotak bulan
	monthW := (usableW / 3) - 6
	monthH := (availableGridH / 4) - 4

	for m := 1; m <= 12; m++ {
		col := (m - 1) % 3
		row := (m - 1) / 3
		x0 := marginL + (float64(col) * (monthW + 9))
		y0 := startYGrid + (float64(row) * (monthH + 5))

		// Judul Bulan
		pdf.SetXY(x0, y0)
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
		pdf.CellFormat(monthW, 7, months[m-1], "", 1, "L", false, 0, "")

		// Header Hari
		colW := monthW / 7
		pdf.SetFont("Arial", "B", 7)
		pdf.SetTextColor(clrSecondary[0], clrSecondary[1], clrSecondary[2])
		pdf.SetXY(x0, y0+7)
		for _, h := range weekHeaders {
			pdf.CellFormat(colW, 5, h, "B", 0, "C", false, 0, "")
		}

		// Kalender Logic
		first := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		startW := int(first.Weekday())
		if startW == 0 {
			startW = 7
		}
		offset := startW - 1
		daysInMonth := time.Date(year, time.Month(m)+1, 0, 0, 0, 0, 0, time.Local).Day()

		cellH := (monthH - 12) / 6
		for r := 0; r < 6; r++ {
			pdf.SetXY(x0, y0+12+(float64(r)*cellH))
			for c := 0; c < 7; c++ {
				dayIdx := r*7 + c - offset + 1
				currX, currY := pdf.GetX(), pdf.GetY()

				if dayIdx > 0 && dayIdx <= daysInMonth {
					key := time.Date(year, time.Month(m), dayIdx, 0, 0, 0, 0, time.Local).Format("2006-01-02")

					if statuses, ok := calendar[key]; ok && len(statuses) > 0 {
						s.setStatusColors(pdf, statuses[0])
						pdf.Rect(currX+0.3, currY+0.3, colW-0.6, cellH-0.6, "F")
					} else if c >= 5 {
						pdf.SetFillColor(clrWeekend[0], clrWeekend[1], clrWeekend[2])
						pdf.Rect(currX+0.3, currY+0.3, colW-0.6, cellH-0.6, "F")
						pdf.SetTextColor(clrSecondary[0], clrSecondary[1], clrSecondary[2])
					} else {
						pdf.SetTextColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
					}

					pdf.SetFont("Arial", "B", 8)
					pdf.Text(currX+1, currY+3.5, strconv.Itoa(dayIdx))
				}
				pdf.SetXY(currX+colW, currY)
			}
		}
	}

	// Signature Area diletakkan di koordinat absolut bawah
	s.drawSignatureArea(pdf, approval, year, pageH-45)
	s.drawFooter(pdf)

	if err := pdf.OutputFileAndClose(path); err != nil {
		return err
	}

	hash, _ := hashFile(path)
	_ = s.scheduleRepo.SetPDFHash(year, hash)
	return nil
}

func (s *SchedulePDFService) drawHeader(pdf *gofpdf.Fpdf, year int) {
	pdf.SetTextColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 15, "JADWAL INSPEKSI TAHUNAN AHU SYSTEM", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(clrSecondary[0], clrSecondary[1], clrSecondary[2])
	pdf.CellFormat(0, 5, fmt.Sprintf("Periode Tahun Penjadwalan %d Dokumen Resmi AIRA System", year), "", 1, "L", false, 0, "")
	pdf.Ln(8)
}

func (s *SchedulePDFService) setStatusColors(pdf *gofpdf.Fpdf, status string) {
	switch status {
	case "siap_diperiksa":
		pdf.SetFillColor(clrBlue[0], clrBlue[1], clrBlue[2])
		pdf.SetTextColor(clrBlueT[0], clrBlueT[1], clrBlueT[2])
	case "dalam_pemeriksaan":
		pdf.SetFillColor(clrAmber[0], clrAmber[1], clrAmber[2])
		pdf.SetTextColor(clrAmberT[0], clrAmberT[1], clrAmberT[2])
	case "selesai":
		pdf.SetFillColor(clrGreen[0], clrGreen[1], clrGreen[2])
		pdf.SetTextColor(clrGreenT[0], clrGreenT[1], clrGreenT[2])
	}
}

func (s *SchedulePDFService) drawLegend(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(clrSecondary[0], clrSecondary[1], clrSecondary[2])
	pdf.Text(20, pdf.GetY()+4, "KETERANGAN STATUS:")

	legends := []struct {
		label string
		bg    []int
		txt   []int
	}{
		{"SIAP DIPERIKSA", clrBlue, clrBlueT},
		{"DALAM PEMERIKSAAN", clrAmber, clrAmberT},
		{"SELESAI", clrGreen, clrGreenT},
	}

	x := 55.0
	y := pdf.GetY()
	for _, l := range legends {
		pdf.SetFillColor(l.bg[0], l.bg[1], l.bg[2])
		pdf.Rect(x, y+1, 30, 5, "F")
		pdf.SetTextColor(l.txt[0], l.txt[1], l.txt[2])
		pdf.SetXY(x, y+1)
		pdf.CellFormat(30, 5, l.label, "", 0, "C", false, 0, "")
		x += 35
	}
}

func (s *SchedulePDFService) drawSignatureArea(pdf *gofpdf.Fpdf, approval *domain.ScheduleApproval, year int, yBase float64) {
	pdf.SetDrawColor(clrBorder[0], clrBorder[1], clrBorder[2])
	pdf.Line(20, yBase-5, 400, yBase-5)

	if approval != nil && approval.VerifyToken != nil {
		verifyURL := fmt.Sprintf("http://10.9.118.16:8080/api/public/verify/%s", *approval.VerifyToken)
		qrPath := fmt.Sprintf("files/verify_%d.png", year)
		_ = generateQR(verifyURL, qrPath)

		pdf.Image(qrPath, 20, yBase, 30, 30, false, "", 0, "")
		pdf.SetXY(52, yBase+8)
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
		pdf.Cell(0, 4, "DIGITAL VERIFICATION SYSTEM")
		pdf.SetXY(52, yBase+12)
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(clrSecondary[0], clrSecondary[1], clrSecondary[2])
		pdf.MultiCell(70, 3, "Dokumen ini sah dan terverifikasi secara digital. Pindai kode QR untuk melihat riwayat approval asli pada server AIRA.", "", "L", false)
	}

	s.drawSignBox(pdf, 290, yBase, "Senior Vice President", "Asep", approval.SVPSignature)
	s.drawSignBox(pdf, 350, yBase, "Assistant Manager", "Hermawan", approval.AsmenSignature)
}

func (s *SchedulePDFService) drawSignBox(pdf *gofpdf.Fpdf, x, y float64, role, name string, sig *string) {
	pdf.SetTextColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
	pdf.SetFont("Arial", "B", 8)
	pdf.Text(x, y+4, strings.ToUpper(role))

	pdf.SetDrawColor(245, 245, 245)
	pdf.Rect(x, y+6, 50, 22, "D")

	if sig != nil && *sig != "" {
		path, _ := saveBase64Image(*sig, fmt.Sprintf("%s.png", name))
		pdf.ImageOptions(path, x+5, y+8, 0, 18, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	}

	pdf.SetFont("Arial", "B", 9)
	pdf.Text(x, y+34, name)
	pdf.SetDrawColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
	pdf.Line(x, y+35, x+50, y+35)
}

func (s *SchedulePDFService) drawFooter(pdf *gofpdf.Fpdf) {
	pdf.SetXY(20, 285)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(clrSecondary[0], clrSecondary[1], clrSecondary[2])
	now := time.Now().Format("02/01/2006 15:04:05")
	pdf.CellFormat(0, 10, fmt.Sprintf("Dicetak otomatis oleh AIRA System pada %s Halaman 1 dari 1", now), "", 0, "L", false, 0, "")
}

func saveBase64Image(data string, filename string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(
		strings.TrimPrefix(data, "data:image/png;base64,"),
	)
	if err != nil {
		return "", err
	}
	path := "files/" + filename
	return path, os.WriteFile(path, b, 0644)
}

func generateQR(url string, filename string) error {
	return qrcode.WriteFile(url, qrcode.Medium, 180, filename)
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
