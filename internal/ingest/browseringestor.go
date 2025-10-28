package ingest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strade/internal/models"
	"strade/internal/store"
	"strade/internal/utils"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type BrowserIngestor struct {
	logger    *zap.SugaredLogger
	store     store.Storage
	sourceURL string
}

func NewBrowserIngestor(logger *zap.SugaredLogger, store store.Storage, sourceURL string) *BrowserIngestor {
	return &BrowserIngestor{
		logger:    logger,
		store:     store,
		sourceURL: sourceURL,
	}
}

func (s *BrowserIngestor) Ingest(ctx context.Context) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	rows, err := s.parseData(data)
	if err != nil {
		return err
	}

	s.logger.Infof("Ingested %d rows", len(rows))

	if err := s.cleanup(); err != nil {
		return err
	}

	return nil
}

func (s *BrowserIngestor) getData() ([]byte, error) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(s.sourceURL)
	page.MustWaitLoad()

	page.MustElement("#rblTipo_1").MustClick()
	page.MustElement("#btnDescarga").MustClick()

	wait := page.Browser().MustWaitDownload()
	data := wait()

	s.logger.Infof("Downloaded %d bytes (ZIP)", len(data))

	tmpDir := filepath.Join(".", "tmp", "ingest")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tmp dir: %w", err)
	}

	zipPath := filepath.Join(tmpDir, fmt.Sprintf("sepomex_%d.zip", time.Now().Unix()))
	if err := os.WriteFile(zipPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to save zip file: %w", err)
	}

	s.logger.Infof("Saved zip to %s", zipPath)

	if err := utils.UnzipFile(zipPath, tmpDir); err != nil {
		return nil, fmt.Errorf("failed to unzip file: %w", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "CPdescarga.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func (s *BrowserIngestor) parseData(data []byte) ([]models.RawDataRecord, error) {
	rows := []models.RawDataRecord{}

	for i, line := range strings.Split(string(data), "\n") {
		if i < 2 {
			continue
		}

		fields := strings.Split(line, "|")
		if len(fields) != 15 {
			continue
		}

		rows = append(rows, models.RawDataRecord{
			PostalCode:         fields[0],
			Settlement:         fields[1],
			SettlementType:     fields[2],
			Municipality:       fields[3],
			State:              fields[4],
			City:               fields[5],
			AdminPostalCode:    fields[6],
			StateCode:          fields[7],
			OfficePostalCode:   fields[8],
			EmptyField:         fields[9],
			SettlementTypeCode: fields[10],
			MunicipalityCode:   fields[11],
			SettlementID:       fields[12],
			Zone:               fields[13],
			CityCode:           fields[14],
		})
	}

	return rows, nil
}

func (s *BrowserIngestor) cleanup() error {
	tmpDir := filepath.Join(".", "tmp", "ingest")
	return os.RemoveAll(tmpDir)
}
