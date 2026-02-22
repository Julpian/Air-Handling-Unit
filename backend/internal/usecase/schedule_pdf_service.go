package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"

	"github.com/jung-kurt/gofpdf"
)

type SchedulePDFService struct {
	scheduleRepo repository.ScheduleRepository
}

func NewSchedulePDFService(repo repository.ScheduleRepository) *SchedulePDFService {
	return &SchedulePDFService{scheduleRepo: repo}
}

// Warna Material Design
var (
	colorBlue   = []int{30, 136, 229}  // Siap Diperiksa [cite: 3]
	colorOrange = []int{255, 179, 0}   // Dalam Pemeriksaan [cite: 4]
	colorGreen  = []int{67, 160, 71}   // Selesai [cite: 5]
	colorGray   = []int{235, 235, 235} // Weekend
	colorText   = []int{33, 33, 33}    // Text Utama
)

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

	// Setup PDF A3 Landscape (420 x 297 mm)
	pdf := gofpdf.New("L", "mm", "A3", "")
	pdf.SetMargins(15, 10, 15)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	// 1. HEADER & TITLE
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(0, 12, "KALENDER PENJADWALAN INSPEKSI AHU", "", 1, "C", false, 0, "") // [cite: 1]
	pdf.SetFont("Arial", "", 16)
	pdf.CellFormat(0, 8, fmt.Sprintf("Tahun %d", year), "", 1, "C", false, 0, "") // [cite: 2]

	// 2. LEGEND
	s.drawLegend(pdf)
	pdf.Ln(5)

	// 3. GRID 12 BULAN (3 Kolom x 4 Baris)
	months := []string{"JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI", "JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"} // [cite: 6, 7, 8, 119, 120, 121, 122, 123, 124, 125, 126, 127]
	weekHeaders := []string{"Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"}                                                                      // [cite: 9, 10, 11, 12, 13, 14, 15]

	pageW, _ := pdf.GetPageSize()
	marginL, _, marginR, _ := pdf.GetMargins()
	usableW := pageW - marginL - marginR

	monthW := usableW / 3
	monthH := 52.0 // DIKECILKAN agar tidak kepotong (52 * 4 baris = 208mm)
	startX, startY := marginL, pdf.GetY()

	for m := 1; m <= 12; m++ {
		col := (m - 1) % 3
		row := (m - 1) / 3
		x0 := startX + (float64(col) * monthW)
		y0 := startY + (float64(row) * monthH)

		// Border luar bulan
		pdf.SetDrawColor(180, 180, 180)
		pdf.Rect(x0+2, y0, monthW-4, monthH-4, "D")

		// Nama Bulan
		pdf.SetXY(x0, y0+1)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(monthW, 6, months[m-1], "", 1, "C", false, 0, "")

		// Header Hari
		colW := (monthW - 10) / 7
		pdf.SetFont("Arial", "B", 7)
		pdf.SetFillColor(245, 245, 245)
		pdf.SetXY(x0+5, y0+7)
		for _, h := range weekHeaders {
			pdf.CellFormat(colW, 5, h, "1", 0, "C", true, 0, "")
		}

		// Kalender Logic
		first := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		startW := int(first.Weekday())
		if startW == 0 {
			startW = 7
		} // Sun=0 ke 7
		offset := startW - 1
		daysInMonth := time.Date(year, time.Month(m)+1, 0, 0, 0, 0, 0, time.Local).Day()

		cellH := (monthH - 18) / 6
		for r := 0; r < 6; r++ {
			pdf.SetXY(x0+5, y0+12+(float64(r)*cellH))
			for c := 0; c < 7; c++ {
				dayIdx := r*7 + c - offset + 1
				currX, currY := pdf.GetX(), pdf.GetY()

				if dayIdx > 0 && dayIdx <= daysInMonth {
					key := time.Date(year, time.Month(m), dayIdx, 0, 0, 0, 0, time.Local).Format("2006-01-02")

					// Warna latar belakang berdasarkan status
					if statuses, ok := calendar[key]; ok && len(statuses) > 0 {
						s.setStatusFillColor(pdf, statuses[0]) // Pakai status pertama
						pdf.Rect(currX, currY, colW, cellH, "F")
						pdf.SetTextColor(255, 255, 255) // Text putih jika berwarna
					} else if c >= 5 {
						pdf.SetFillColor(colorGray[0], colorGray[1], colorGray[2])
						pdf.Rect(currX, currY, colW, cellH, "F")
						pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
					} else {
						pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
					}

					pdf.Rect(currX, currY, colW, cellH, "D")
					pdf.SetFont("Arial", "B", 8)
					pdf.Text(currX+1.5, currY+4, strconv.Itoa(dayIdx))
				} else {
					pdf.SetDrawColor(220, 220, 220)
					pdf.Rect(currX, currY, colW, cellH, "D")
				}
				pdf.SetXY(currX+colW, currY)
			}
		}
	}

	// 4. SIGNATURE AREA
	s.drawSignatureArea(pdf, approval)

	return pdf.OutputFileAndClose(path)
}

func (s *SchedulePDFService) setStatusFillColor(pdf *gofpdf.Fpdf, status string) {
	switch status {
	case "siap_diperiksa":
		pdf.SetFillColor(colorBlue[0], colorBlue[1], colorBlue[2])
	case "dalam_pemeriksaan":
		pdf.SetFillColor(colorOrange[0], colorOrange[1], colorOrange[2])
	case "selesai":
		pdf.SetFillColor(colorGreen[0], colorGreen[1], colorGreen[2])
	default:
		pdf.SetFillColor(200, 200, 200)
	}
}

func (s *SchedulePDFService) drawLegend(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "", 9)
	legends := []struct {
		label string
		rgb   []int
	}{
		{"Siap Diperiksa", colorBlue},
		{"Dalam Pemeriksaan", colorOrange},
		{"Selesai", colorGreen},
	}
	x := 150.0
	for _, l := range legends {
		pdf.SetFillColor(l.rgb[0], l.rgb[1], l.rgb[2])
		pdf.Rect(x, pdf.GetY()+1, 4, 4, "F")
		pdf.SetXY(x+5, pdf.GetY())
		pdf.Cell(40, 6, l.label)
		x += 45
	}
}

func (s *SchedulePDFService) drawSignatureArea(
	pdf *gofpdf.Fpdf,
	approval *domain.ScheduleApproval,
) {

	if approval == nil || approval.VerifyToken == nil || *approval.VerifyToken == "" {
		return
	}
	// Geser y sedikit ke atas (252) agar nama di bawah kotak tidak terpotong
	// di batas bawah kertas A3 (297mm)
	y := 252.0

	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.SetFont("Arial", "B", 11)
	pdf.Text(300, y, "Disetujui,") // [cite: 18]

	// --- BAGIAN SVP ---
	pdf.SetFont("Arial", "B", 10)
	pdf.Text(300, y+5, "SVP") // [cite: 20]

	// Kotak Tanda Tangan
	pdf.SetDrawColor(180, 180, 180)
	pdf.Rect(300, y+7, 45, 20, "D")

	// Nama di bawah kotak
	pdf.Text(300, y+31, "Asep")

	if approval.SVPSignature != nil {
		path, _ := saveBase64Image(*approval.SVPSignature, "svp.png")
		// PERBAIKAN: Gunakan tinggi (height) 16mm sebagai batasan agar pas di kotak 20mm
		// Lebar (width) diatur 0 agar skala gambar tetap proporsional
		pdf.ImageOptions(path, 302, y+9, 0, 16, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	}

	// --- BAGIAN ASMEN ---
	pdf.Text(355, y+5, "ASMEN") // [cite: 21]

	// Kotak Tanda Tangan
	pdf.Rect(355, y+7, 45, 20, "D")

	// Nama di bawah kotak
	pdf.Text(355, y+31, "Hermawan")

	if approval.AsmenSignature != nil {
		path, _ := saveBase64Image(*approval.AsmenSignature, "asmen.png")
		// PERBAIKAN: Batasi tinggi gambar 16mm agar tidak meluber keluar kotak
		pdf.ImageOptions(path, 357, y+9, 0, 16, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
	}

	verifyURL := fmt.Sprintf(
		"http://192.168.0.127:8080/api/public/verify/%s",
		*approval.VerifyToken,
	)

	qrPath := "files/verify.png"
	_ = generateQR(verifyURL, qrPath)

	pdf.Image(qrPath, 250, y+5, 30, 0, false, "", 0, "")
}
