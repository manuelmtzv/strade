package ingest

import "context"

type Ingestor interface {
	Ingest(ctx context.Context) error
}
