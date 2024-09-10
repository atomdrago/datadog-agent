// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo.exe -godefs -- -fsigned-char types.go

package driver

const Signature = 0xddfd00000017

const (
	GetStatsIOCTL             = 0x122004
	SetFlowFilterIOCTL        = 0x122010
	SetDataFilterIOCTL        = 0x12200c
	GetFlowsIOCTL             = 0x122014
	SetMaxOpenFlowsIOCTL      = 0x122024
	SetMaxClosedFlowsIOCTL    = 0x122028
	FlushPendingHttpTxnsIOCTL = 0x122020
	EnableHttpIOCTL           = 0x122030
	EnableClassifyIOCTL       = 0x122040
	SetClosedFlowsLimitIOCTL  = 0x12203c
	GetOpenFlowsIOCTL         = 0x122036
	GetClosedFlowsIOCTL       = 0x12203a
)

type FilterAddress struct {
	Af         uint64
	V4_address [4]uint8
	V4_padding [4]uint8
	V6_address [16]uint8
	Mask       uint64
}

type FilterDefinition struct {
	FilterVersion  uint64
	Size           uint64
	FilterLayer    uint64
	Af             uint64
	LocalAddress   FilterAddress
	RemoteAddress  FilterAddress
	LocalPort      uint64
	RemotePort     uint64
	Protocol       uint64
	Direction      uint64
	InterfaceIndex uint64
}

const FilterDefinitionSize = 0x98

type FilterPacketHeader struct {
	FilterVersion    uint64
	Sz               uint64
	SkippedSinceLast uint64
	FilterId         uint64
	Direction        uint64
	PktSize          uint64
	Af               uint64
	OwnerPid         uint64
	Timestamp        uint64
}

const FilterPacketHeaderSize = 0x48

type FlowStats struct {
	Num_flow_collisions                      int64
	Num_flow_alloc_skipped_max_open_exceeded int64
	Num_flow_closed_dropped_max_exceeded     int64
	Num_flow_structures                      int64
	Peak_num_flow_structures                 int64
	Num_flow_closed_structures               int64
	Peak_num_flow_closed_structures          int64
	Open_table_adds                          int64
	Open_table_removes                       int64
	Closed_table_adds                        int64
	Closed_table_removes                     int64
	Num_flows_no_handle                      int64
	Peak_num_flows_no_handle                 int64
	Num_flows_missed_max_no_handle_exceeded  int64
	Num_packets_after_flow_closed            int64
	Classify_with_no_direction               int64
	Classify_multiple_request                int64
	Classify_multiple_response               int64
	Classify_response_no_request             int64
	No_state_at_ale_auth_connect             int64
	No_state_at_ale_auth_recv                int64
	No_state_at_ale_flow_established         int64
	No_state_at_ale_endpoint_closure         int64
	No_state_at_inbound_transport            int64
	No_state_at_outbound_transport           int64
}
type TransportStats struct {
	Packets_skipped int64
	Calls_requested int64
	Calls_completed int64
	Calls_cancelled int64
}
type HttpStats struct {
	Txns_captured              int64
	Txns_skipped_max_exceeded  int64
	Ndis_buffer_non_contiguous int64
	Flows_ignored_as_etw       int64
	Txn_zero_latency           int64
	Txn_batched_on_read        int64
}
type Stats struct {
	Flow_stats      FlowStats
	Transport_stats TransportStats
	Http_stats      HttpStats
}

const StatsSize = 0x118

type PerFlowData struct {
	FlowHandle               uint64
	FlowCookie               uint64
	ProcessId                uint64
	AddressFamily            uint16
	Protocol                 uint16
	Flags                    uint32
	LocalAddress             [16]byte
	RemoteAddress            [16]byte
	PacketsOut               uint64
	MonotonicSentBytes       uint64
	TransportBytesOut        uint64
	PacketsIn                uint64
	MonotonicRecvBytes       uint64
	TransportBytesIn         uint64
	Timestamp                uint64
	LocalPort                uint16
	RemotePort               uint16
	ClassificationStatus     uint16
	ClassifyRequest          uint16
	ClassifyResponse         uint16
	HttpUpgradeToH2Requested uint8
	HttpUpgradeToH2Accepted  uint8
	Tls_versions_offered     uint16
	Tls_version_chosen       uint16
	Tls_alpn_requested       uint64
	Tls_alpn_chosen          uint64
	Protocol_u               [36]byte
}
type TCPFlowData struct {
	IRTT             uint64
	SRTT             uint64
	RttVariance      uint64
	RetransmitCount  uint64
	ConnectionStatus uint32
}
type UDPFlowData struct {
	Reserved uint64
}

const PerFlowDataSize = 0xbc

const (
	FlowDirectionMask     = 0x300
	FlowDirectionBits     = 0x8
	FlowDirectionInbound  = 0x1
	FlowDirectionOutbound = 0x2

	FlowClosedMask         = 0x10
	TCPFlowEstablishedMask = 0x20
)

type ConnectionStatus uint32

const (
	ConnectionStatusUnknown     ConnectionStatus = 0x0
	ConnectionStatusAttempted   ConnectionStatus = 0x1
	ConnectionStatusEstablished ConnectionStatus = 0x2
	ConnectionStatusACKRST      ConnectionStatus = 0x3
	ConnectionStatusTimeout     ConnectionStatus = 0x4
	ConnectionStatusSentRst     ConnectionStatus = 0x5
	ConnectionStatusRecvRst     ConnectionStatus = 0x6
)

const (
	DirectionInbound  = 0x0
	DirectionOutbound = 0x1
)

const (
	LayerTransport = 0x1
)

type HttpTransactionType struct {
	RequestStarted     uint64
	ResponseLastSeen   uint64
	Tup                ConnTupleType
	RequestMethod      uint32
	ResponseStatusCode uint16
	MaxRequestFragment uint16
	SzRequestFragment  uint16
	Pad                [6]uint8
	RequestFragment    *uint8
}
type HttpConfigurationSettings struct {
	MaxTransactions        uint64
	NotificationThreshold  uint64
	MaxRequestFragment     uint16
	EnableAutoETWExclusion uint16
}
type ConnTupleType struct {
	LocalAddr  [16]byte
	RemoteAddr [16]byte
	LocalPort  uint16
	RemotePort uint16
	Family     uint16
	Pad        uint16
}
type HttpMethodType uint32
type ClassificationSettings struct {
	Enabled uint64
}

type TcpConnectionStatus uint32

const (
	TcpStatusEstablished = 0x2
)
const (
	HttpTransactionTypeSize        = 0x50
	HttpSettingsTypeSize           = 0x14
	ClassificationSettingsTypeSize = 0x8
)

const (
	ClassificationUnclassified           = 0x0
	ClassificationClassified             = 0x1
	ClassificationUnableInsufficientData = 0x2
	ClassificationUnknown                = 0x3

	ClassificationRequestUnclassified = 0x0
	ClassificationRequestHTTPUnknown  = 0x1
	ClassificationRequestHTTPPost     = 0x2
	ClassificationRequestHTTPPut      = 0x3
	ClassificationRequestHTTPPatch    = 0x4
	ClassificationRequestHTTPGet      = 0x5
	ClassificationRequestHTTPHead     = 0x6
	ClassificationRequestHTTPOptions  = 0x7
	ClassificationRequestHTTPDelete   = 0x8
	ClassificationRequestHTTPLast     = 0x8

	ClassificationRequestHTTP2 = 0x9

	ClassificationRequestTLS  = 0xa
	ClassificationResponseTLS = 0x2

	ALPNProtocolHTTP2  = 0x1
	ALPNProtocolHTTP11 = 0x2

	ClassificationResponseUnclassified = 0x0
	ClassificationResponseHTTP         = 0x1
)
