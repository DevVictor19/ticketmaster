package cdc

type DebeziumEnvelope[T any] struct {
	// Schema  interface{}        `json:"schema"`
	Payload DebeziumPayload[T] `json:"payload"`
}

type DebeziumPayload[T any] struct {
	Before      *T                   `json:"before"`
	After       *T                   `json:"after"`
	Source      DebeziumSource       `json:"source"`
	Transaction *DebeziumTransaction `json:"transaction"`
	Op          DebeziumOperation    `json:"op"` // c=create, u=update, d=delete, r=read
	TsMs        int64                `json:"ts_ms"`
	TsUs        int64                `json:"ts_us,omitempty"`
	TsNs        int64                `json:"ts_ns,omitempty"`
}

type DebeziumOperation string

var (
	CreateOp DebeziumOperation = "c"
	UpdateOp DebeziumOperation = "u"
	DeleteOp DebeziumOperation = "d"
)

type DebeziumSource struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	Snapshot  string `json:"snapshot"`
	Db        string `json:"db"`
	Sequence  string `json:"sequence,omitempty"`
	TsUs      int64  `json:"ts_us,omitempty"`
	TsNs      int64  `json:"ts_ns,omitempty"`
	Schema    string `json:"schema"`
	Table     string `json:"table"`
	TxId      int64  `json:"txId,omitempty"`
	Lsn       int64  `json:"lsn,omitempty"`
	Xmin      *int64 `json:"xmin,omitempty"`
	Origin    string `json:"origin,omitempty"`
	OriginLsn int64  `json:"origin_lsn,omitempty"`
}

type DebeziumTransaction struct {
	ID                  string `json:"id"`
	TotalOrder          int64  `json:"total_order"`
	DataCollectionOrder int64  `json:"data_collection_order"`
}
