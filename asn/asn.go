package asn

type Response struct {
	Success   bool      `json:"success"`
	ErrorText string    `json:"error_text,omitempty"`
	Values    []Counter `json:"values,omitempty"`
}

type Counter struct {
	ASN     		int64  `json:"asn_number"`
	IncomingPackets int64  `json:"incoming_packets"`
	IncomingBytes   int64  `json:"incoming_bytes"`
}
