package flagparsers

import (
	"github.com/timescale/outflux/internal/ingestion"
	"github.com/timescale/outflux/internal/schemamanagement"
)

// Flags used in outflux and their default values
const (
	InputHostFlag               = "input-host"
	InputUserFlag               = "input-user"
	InputPassFlag               = "input-pass"
	OutputConnFlag              = "output-conn"
	UseEnvVarsFlag              = "use-env-vars"
	SchemaStrategyFlag          = "schema-strategy"
	CommitStrategyFlag          = "commit-strategy"
	OutputSchemaFlag            = "output-schema"
	FromFlag                    = "from"
	ToFlag                      = "to"
	LimitFlag                   = "limit"
	ChunkSizeFlag               = "chunk-size"
	QuietFlag                   = "quiet"
	DataBufferFlag              = "data-buffer"
	MaxParallelFlag             = "max-parallel"
	RollbackOnExternalErrorFlag = "rollback-on-external-error"
	BatchSizeFlag               = "batch-size"

	DefaultInputHost               = "http://localhost:8086"
	DefaultInputUser               = ""
	DefaultInputPass               = ""
	DefaultOutputConn              = "sslmode=disable"
	DefaultUseEnvVars              = true
	DefaultOutputSchema            = ""
	DefaultSchemaStrategy          = schemamanagement.CreateIfMissing
	DefaultCommitStrategy          = ingestion.CommitOnEachBatch
	DefaultDataBufferSize          = 15000
	DefaultChunkSize               = 15000
	DefaultLimit                   = 0
	DefaultMaxParallel             = 2
	DefaultRollbackOnExternalError = true
	DefaultBatchSize               = 8000
)