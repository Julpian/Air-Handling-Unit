package dto

type ScanNFCRequest struct {
	NFCUID string `json:"nfc_uid"`
}

type ScanNFCResponse struct {
	InspectionID string `json:"inspection_id"`
	ScanToken    string `json:"scan_token"`
}
