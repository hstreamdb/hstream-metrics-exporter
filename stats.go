package main

type Stats struct {
	methods []string
}

const SubscriptionId = "subscription_id"
const StreamName = "stream_name"
const Subscription = "subscription"
const Stream = "stream"
const StreamCounter = "stream_counter"
const ServerHistogram = "server_histogram"

// https://github.com/hstreamdb/hstream/blob/main/common/stats/include/per_subscription_time_series.inc

var subscriptionStats = []Stats{
	subscriptionSendOutBytes, subscriptionSendOutRecords, subscriptionRequestMessages, subscriptionResponseMessages, subscriptionAcks}

func GetSubscriptionStats() []Stats {
	return subscriptionStats
}

var subscriptionSendOutBytes = Stats{
	methods: []string{"send_out_bytes", "sends"},
}

var subscriptionAcks = Stats{
	methods: []string{"acks", "acknowledgements"},
}

func GetSubscriptionSendOutBytes() Stats {
	return subscriptionSendOutBytes
}

var subscriptionSendOutRecords = Stats{
	methods: []string{"send_out_records"},
}

func GetSubscriptionSendOutRecords() Stats {
	return subscriptionSendOutRecords
}

var subscriptionRequestMessages = Stats{
	methods: []string{"request_messages"},
}

func GetSubscriptionRequestMessages() Stats {
	return subscriptionRequestMessages
}

var subscriptionResponseMessages = Stats{
	methods: []string{"response_messages"},
}

func GetSubscriptionResponseMessages() Stats {
	return subscriptionResponseMessages
}

// https://github.com/hstreamdb/hstream/blob/main/common/stats/include/per_stream_time_series.inc

var streamStats = []Stats{
	streamAppendInBytes, streamAppendInRecords, streamAppendInRequests, streamAppendFailedRequests, streamRecordBytes}

func GetStreamStats() []Stats {
	return streamStats
}

var streamAppendInBytes = Stats{
	methods: []string{"append_in_bytes", "appends"},
}

func GetStreamAppendInBytes() Stats {
	return streamAppendInBytes
}

// append_in_records can only have 0s as interval
var streamAppendInRecords = Stats{
	methods: []string{"append_in_record"},
}

func GetStreamAppendInRecords() Stats {
	return streamAppendInRecords
}

var streamAppendInRequests = Stats{
	methods: []string{"append_in_requests"},
}

func GetStreamAppendInRequests() Stats {
	return streamAppendInRequests
}

var streamAppendFailedRequests = Stats{
	methods: []string{"append_failed_requests"},
}

func GetStreamAppendFailedRequests() Stats {
	return streamAppendFailedRequests
}

var streamRecordBytes = Stats{
	methods: []string{"record_bytes", "reads"},
}

func GetStreamRecordBytes() Stats {
	return streamRecordBytes
}

var serverHistogramStats = []Stats{
	appendRequestLatencyStats, appendLatencyStats,
}

var appendRequestLatencyStats = Stats{
	methods: []string{"append_request_latency"},
}

var appendLatencyStats = Stats{
	methods: []string{"append_latency"},
}

func GetServerHistogramStats() []Stats {
	return serverHistogramStats
}

func GetStreamCounterStats() []Stats {
	return streamCounterStats
}

var streamCounterStats = []Stats{
	counterAppendTotal, counterAppendFailed,
}

var counterAppendTotal = Stats{
	methods: []string{"append_total"},
}

var counterAppendFailed = Stats{
	methods: []string{"append_failed"},
}
