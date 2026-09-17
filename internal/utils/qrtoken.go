package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

// QRPayload adalah struktur data yang di-encode ke dalam QR code kasir
// (order_qr_payload). Berisi info minimal yang dibutuhkan kasir untuk
// verifikasi & konfirmasi pembayaran, plus expiry singkat agar tidak bisa
// direplay setelah pesanan selesai.
type QRPayload struct {
	OrderID     string `json:"order_id"`
	TableNumber string `json:"table_number"`
	TotalAmount int64  `json:"total_amount"` // dalam rupiah, integer, hindari float di payload
	IssuedAt    int64  `json:"issued_at"`
	ExpiresAt   int64  `json:"expires_at"`
}

// GenerateSignedToken membuat token dalam format "base64(payload).base64(signature)"
// mirip struktur JWT sederhana, ditandatangani dengan HMAC-SHA256.
func GenerateSignedToken(payload interface{}, secret string) (string, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payloadB64))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return payloadB64 + "." + signature, nil
}

// VerifySignedToken memvalidasi signature lalu unmarshal payload ke `out`.
// Selalu gunakan fungsi ini di sisi kasir/server sebelum mempercayai isi QR —
// jangan pernah trust payload dari client tanpa verifikasi signature.
func VerifySignedToken(token, secret string, out interface{}) error {
	parts := splitToken(token)
	if len(parts) != 2 {
		return errors.New("format token tidak valid")
	}
	payloadB64, signature := parts[0], parts[1]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payloadB64))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return errors.New("signature tidak valid, kemungkinan token dipalsukan")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return err
	}
	return json.Unmarshal(payloadBytes, out)
}

// NewCashierQRPayload adalah helper khusus untuk generate payload QR pembayaran kasir,
// dengan expiry default 2 jam dari sekarang.
func NewCashierQRPayload(orderID, tableNumber string, totalAmount int64) QRPayload {
	now := time.Now()
	return QRPayload{
		OrderID:     orderID,
		TableNumber: tableNumber,
		TotalAmount: totalAmount,
		IssuedAt:    now.Unix(),
		ExpiresAt:   now.Add(2 * time.Hour).Unix(),
	}
}

func splitToken(token string) []string {
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '.' {
			return []string{token[:i], token[i+1:]}
		}
	}
	return []string{token}
}
