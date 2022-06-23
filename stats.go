package main

type Stats struct {
	methods   []string
	intervals []int
}

const DefaultIntervalStr = "10s"
const SubscriptionId = "subscription_id"
const StreamName = "stream_name"
const Subscription = "subscription"
const Stream = "stream"
const StreamCounter = "stream_counter"

// https://github.com/hstreamdb/hstream/blob/main/common/stats/include/per_subscription_time_series.inc

var subscriptionStats = []Stats{
	subscriptionSendOutBytes, subscriptionSendOutRecords, subscriptionRequestMessages, subscriptionResponseMessages, subscriptionAcks}

func GetSubscriptionStats() []Stats {
	return subscriptionStats
}

var subscriptionSendOutBytes = Stats{
	methods:   []string{"send_out_bytes", "sends"},
	intervals: []int{4, 60, 300, 600},
}

var subscriptionAcks = Stats{
	methods: []string{"acks", "acknowledgements"},
}

func GetSubscriptionSendOutBytes() Stats {
	return subscriptionSendOutBytes
}

var subscriptionSendOutRecords = Stats{
	methods:   []string{"send_out_records"},
	intervals: []int{4, 60, 300, 600},
}

func GetSubscriptionSendOutRecords() Stats {
	return subscriptionSendOutRecords
}

var subscriptionRequestMessages = Stats{
	methods:   []string{"request_messages"},
	intervals: []int{4, 60, 300, 600},
}

func GetSubscriptionRequestMessages() Stats {
	return subscriptionRequestMessages
}

var subscriptionResponseMessages = Stats{
	methods:   []string{"response_messages"},
	intervals: []int{4, 60, 300, 600},
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
	methods:   []string{"append_in_bytes", "appends"},
	intervals: []int{4, 60, 300, 600},
}

func GetStreamAppendInBytes() Stats {
	return streamAppendInBytes
}

var streamAppendInRecords = Stats{
	methods:   []string{"append_in_records"},
	intervals: []int{4, 60, 300, 600},
}

func GetStreamAppendInRecords() Stats {
	return streamAppendInRecords
}

var streamAppendInRequests = Stats{
	methods:   []string{"append_in_requests"},
	intervals: []int{4, 60, 300, 600},
}

func GetStreamAppendInRequests() Stats {
	return streamAppendInRequests
}

var streamAppendFailedRequests = Stats{
	methods:   []string{"append_failed_requests"},
	intervals: []int{4, 60, 300, 600},
}

func GetStreamAppendFailedRequests() Stats {
	return streamAppendFailedRequests
}

var streamRecordBytes = Stats{
	methods:   []string{"record_bytes", "reads"},
	intervals: []int{900, 1800, 3600},
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
