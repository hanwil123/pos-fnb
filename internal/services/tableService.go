package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/yourorg/pos-fnb-backend/internal/dto"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"github.com/yourorg/pos-fnb-backend/internal/repository"
)

type TableService interface {
	ResolveTable(qrToken string) (*dto.TableSessionResponse, error)
	ListTables(restaurantID string) ([]dto.TableResponse, error)
	CreateTable(req *dto.CreateTableRequest) (*dto.TableResponse, error)
}

type tableService struct {
	repo repository.TableRepository
}

func NewTableService(repo repository.TableRepository) TableService {
	return &tableService{repo: repo}
}

func (s *tableService) ResolveTable(qrToken string) (*dto.TableSessionResponse, error) {
	table, err := s.repo.FindByQRToken(qrToken)
	if err != nil {
		return nil, err
	}

	// Cek session aktif yang masih berlaku
	session, err := s.repo.FindActiveSession(table.ID)
	if err != nil {
		// Belum ada session aktif -> buat baru
		session = &models.TableSession{
			TableID:      table.ID,
			SessionToken: uuid.New().String(),
			Status:       models.SessionActive,
		}
		if err := s.repo.CreateSession(session); err != nil {
			return nil, err
		}

		// Update status meja jadi occupied
		table.Status = models.TableOccupied
		if err := s.repo.Update(table); err != nil {
			return nil, err
		}
	}

	return &dto.TableSessionResponse{
		Table: dto.TableResponse{
			ID:           table.ID.String(),
			RestaurantID: table.RestaurantID.String(),
			TableNumber:  table.TableNumber,
			QRToken:      table.QRToken,
			Status:       string(table.Status),
		},
		SessionToken: session.SessionToken,
		SessionID:    session.ID.String(),
	}, nil
}

func (s *tableService) ListTables(restaurantID string) ([]dto.TableResponse, error) {
	var tables []models.Table
	var err error

	if restaurantID != "" {
		restaurantUUID, parseErr := uuid.Parse(restaurantID)
		if parseErr != nil {
			return nil, parseErr
		}
		tables, err = s.repo.FindByRestaurantID(restaurantUUID)
	} else {
		tables, err = s.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	result := make([]dto.TableResponse, 0, len(tables))
	for _, table := range tables {
		result = append(result, dto.TableResponse{
			ID:           table.ID.String(),
			RestaurantID: table.RestaurantID.String(),
			TableNumber:  table.TableNumber,
			QRToken:      table.QRToken,
			Status:       string(table.Status),
		})
	}

	return result, nil
}

func (s *tableService) CreateTable(req *dto.CreateTableRequest) (*dto.TableResponse, error) {
	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		return nil, err
	}

	// qr_token dibuat dari UUID + timestamp lalu di-hash pendek agar unik & sulit ditebak
	rawToken := uuid.New().String() + time.Now().String()
	qrToken := uuid.NewSHA1(uuid.NameSpaceOID, []byte(rawToken)).String()

	table := &models.Table{
		RestaurantID: restaurantID,
		TableNumber:  req.TableNumber,
		QRToken:      qrToken,
		Status:       models.TableAvailable,
	}

	if err := s.repo.Create(table); err != nil {
		return nil, err
	}
	// 1. Buat URL tujuan untuk di-scan customer
	baseUrlQrCode := fmt.Sprintf("https://pos-fnb-nu.vercel.app/%s", table.QRToken)

	// 2. Generate QR Code langsung ke bentuk []byte di memori (Ukuran 512px agar tajam saat di-print)
	qrPngBytes, err := qrcode.Encode(baseUrlQrCode, qrcode.Medium, 512)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code: %w", err)
	}

	// 3. Konversi byte gambar PNG menjadi string Base64 dengan prefix Data URI
	qrBase64Str := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrPngBytes))

	return &dto.TableResponse{
		ID:           table.ID.String(),
		RestaurantID: table.RestaurantID.String(),
		TableNumber:  table.TableNumber,
		QRToken:      table.QRToken,
		Status:       string(table.Status),
		QRCodeImage:  qrBase64Str,
	}, nil
}
