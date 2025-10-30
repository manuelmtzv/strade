package ingest

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strade/internal/config"
	"strade/internal/models"
	"strade/internal/store"
	"strade/internal/utils"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type BrowserIngestor struct {
	config config.IngestorConfig
	logger *zap.SugaredLogger
	store  *store.Storage
}

func NewBrowserIngestor(config config.IngestorConfig, logger *zap.SugaredLogger, store *store.Storage) *BrowserIngestor {
	return &BrowserIngestor{
		config: config,
		logger: logger,
		store:  store,
	}
}

func (s *BrowserIngestor) Ingest(ctx context.Context) error {
	start := time.Now()

	t0 := time.Now()
	data, err := s.getData(ctx)
	if err != nil {
		return err
	}
	getDataDuration := time.Since(t0)

	t1 := time.Now()
	rows, err := s.parseData(data)
	if err != nil {
		return err
	}
	parseDataDuration := time.Since(t1)

	s.logger.Infof("Parsed %d raw records", len(rows))

	t2 := time.Now()
	if err := s.transformAndStore(ctx, rows); err != nil {
		return err
	}
	transformAndStoreDuration := time.Since(t2)

	t3 := time.Now()
	if err := s.cleanup(); err != nil {
		return err
	}
	cleanupDuration := time.Since(t3)

	s.logger.Infof("Data ingestion completed in %s", time.Since(start))
	s.logger.Infof("  getData: %s", getDataDuration)
	s.logger.Infof("  parseData: %s", parseDataDuration)
	s.logger.Infof("  transformAndStore: %s", transformAndStoreDuration)
	s.logger.Infof("  cleanup: %s", cleanupDuration)

	return nil
}

func (s *BrowserIngestor) getData(ctx context.Context) ([]byte, error) {
	browser := rod.New().Context(ctx).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(s.config.SourceURL)
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

	file, err := os.Open(filepath.Join(tmpDir, "CPdescarga.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	utf8Data, err := io.ReadAll(decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to read and decode file: %w", err)
	}

	return utf8Data, nil
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

func (s *BrowserIngestor) transformAndStore(ctx context.Context, rows []models.RawDataRecord) error {
	s.logger.Info("Starting data transformation and upsert...")

	states := s.extractStates(rows)
	settlementTypes := s.extractSettlementTypes(rows)
	municipalities := s.extractMunicipalities(rows)
	cities := s.extractCities(rows)

	s.logger.Infof("Extracted: %d states, %d settlement types, %d municipalities, %d cities",
		len(states), len(settlementTypes), len(municipalities), len(cities))

	if err := (utils.WithTx(ctx, s.store.DB, func(tx *sql.Tx) error {
		if err := s.store.StateStore.BulkUpsertTx(ctx, tx, states); err != nil {
			return err
		}

		if err := s.store.SettlementTypeStore.BulkUpsertTx(ctx, tx, settlementTypes); err != nil {
			return err
		}

		if err := s.store.MunicipalityStore.BulkUpsertTx(ctx, tx, municipalities); err != nil {
			return err
		}

		if err := s.store.CityStore.BulkUpsertTx(ctx, tx, cities); err != nil {
			return err
		}

		return nil
	})); err != nil {
		return err
	}

	settlements := s.transformToSettlements(rows)
	s.logger.Infof("Transformed %d settlements", len(settlements))

	if err := s.bulkUpsertSettlementsWithWorkers(ctx, settlements); err != nil {
		return err
	}

	s.logger.Info("Data transformation and upsert completed")
	return nil
}

func (s *BrowserIngestor) extractStates(rows []models.RawDataRecord) []*models.State {
	stateMap := make(map[string]*models.State)

	for _, row := range rows {
		if row.StateCode == "" {
			continue
		}
		if _, exists := stateMap[row.StateCode]; !exists {
			stateMap[row.StateCode] = &models.State{
				ID:   row.StateCode,
				Name: row.State,
				Slug: slug.Make(row.State),
			}
		}
	}

	states := make([]*models.State, 0, len(stateMap))
	for _, state := range stateMap {
		states = append(states, state)
	}

	return states
}

func (s *BrowserIngestor) extractSettlementTypes(rows []models.RawDataRecord) []*models.SettlementType {
	typeMap := make(map[string]*models.SettlementType)

	for _, row := range rows {
		if row.SettlementTypeCode == "" {
			continue
		}
		if _, exists := typeMap[row.SettlementTypeCode]; !exists {
			typeMap[row.SettlementTypeCode] = &models.SettlementType{
				ID:   row.SettlementTypeCode,
				Name: row.SettlementType,
				Slug: slug.Make(row.SettlementType),
			}
		}
	}

	settlementTypes := make([]*models.SettlementType, 0, len(typeMap))
	for _, st := range typeMap {
		settlementTypes = append(settlementTypes, st)
	}

	return settlementTypes
}

func (s *BrowserIngestor) extractMunicipalities(rows []models.RawDataRecord) []*models.Municipality {
	municipalityMap := make(map[string]*models.Municipality)

	for _, row := range rows {
		if row.MunicipalityCode == "" || row.StateCode == "" {
			continue
		}
		compositeID := row.StateCode + row.MunicipalityCode
		if _, exists := municipalityMap[compositeID]; !exists {
			municipalityMap[compositeID] = &models.Municipality{
				ID:      compositeID,
				Name:    row.Municipality,
				Slug:    slug.Make(row.Municipality),
				StateID: row.StateCode,
			}
		}
	}

	municipalities := make([]*models.Municipality, 0, len(municipalityMap))
	for _, m := range municipalityMap {
		municipalities = append(municipalities, m)
	}

	return municipalities
}

func (s *BrowserIngestor) extractCities(rows []models.RawDataRecord) []*models.City {
	cityMap := make(map[string]*models.City)

	for _, row := range rows {
		cityCode := strings.TrimSpace(row.CityCode)
		if cityCode == "" || row.StateCode == "" {
			continue
		}

		compositeID := row.StateCode + cityCode
		if _, exists := cityMap[compositeID]; !exists {
			cityName := row.City
			if cityName == "" {
				cityName = "Sin nombre"
			}
			cityMap[compositeID] = &models.City{
				ID:      compositeID,
				Name:    cityName,
				Slug:    slug.Make(cityName),
				StateID: row.StateCode,
			}
		}
	}

	cities := make([]*models.City, 0, len(cityMap))
	for _, c := range cityMap {
		cities = append(cities, c)
	}

	return cities
}

func (s *BrowserIngestor) transformToSettlements(rows []models.RawDataRecord) []*models.Settlement {
	settlementMap := make(map[string]*models.Settlement)
	skippedCount := 0

	for _, row := range rows {
		if row.PostalCode == "" || row.Settlement == "" || row.MunicipalityCode == "" || row.SettlementTypeCode == "" || row.SettlementID == "" {
			skippedCount++
			continue
		}

		cityCode := strings.TrimSpace(row.CityCode)
		cityID := "000"
		if cityCode != "" {
			cityID = row.StateCode + cityCode
		}

		key := row.StateCode + row.MunicipalityCode + row.SettlementID

		if _, exists := settlementMap[key]; exists {
			continue
		}

		municipalityID := row.StateCode + row.MunicipalityCode

		settlementMap[key] = &models.Settlement{
			ID:               key,
			PostalCode:       row.PostalCode,
			Name:             row.Settlement,
			Slug:             slug.Make(row.Settlement),
			SettlementTypeID: row.SettlementTypeCode,
			MunicipalityID:   municipalityID,
			CityID:           cityID,
			StateID:          row.StateCode,
			OfficePostalCode: row.OfficePostalCode,
			Zone:             row.Zone,
		}
	}

	settlements := make([]*models.Settlement, 0, len(settlementMap))
	for _, settlement := range settlementMap {
		settlements = append(settlements, settlement)
	}

	return settlements
}

func chunkSettlements(settlements []*models.Settlement, batchSize int) [][]*models.Settlement {
	var chunks [][]*models.Settlement
	for i := 0; i < len(settlements); i += batchSize {
		end := min(i+batchSize, len(settlements))
		chunks = append(chunks, settlements[i:end])
	}
	return chunks
}

func (s *BrowserIngestor) bulkUpsertSettlementsWithWorkers(ctx context.Context, settlements []*models.Settlement) error {
	batchSize := s.config.SettlementsBatchSize
	numWorkers := s.config.SettlementsWorkers

	batches := chunkSettlements(settlements, batchSize)

	var wg sync.WaitGroup
	errChan := make(chan error, len(batches))

	semaphore := make(chan struct{}, numWorkers)

	for _, batch := range batches {
		wg.Add(1)
		go func(b []*models.Settlement) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := utils.WithTx(ctx, s.store.DB, func(tx *sql.Tx) error {
				return s.store.SettlementStore.BulkUpsertTx(ctx, tx, b)
			})
			if err != nil {
				errChan <- err
			}
		}(batch)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BrowserIngestor) cleanup() error {
	tmpDir := filepath.Join(".", "tmp", "ingest")
	return os.RemoveAll(tmpDir)
}
