package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	hstream "github.com/grafana/grafana-starter-datasource-backend/pkg/gen/hstreamdb/hstream/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"math/rand"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*SampleDatasource)(nil)
	_ backend.CheckHealthHandler    = (*SampleDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*SampleDatasource)(nil)
)

// NewSampleDatasource creates a new datasource instance.
func NewSampleDatasource(_ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &SampleDatasource{}, nil
}

// SampleDatasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type SampleDatasource struct {
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *SampleDatasource) Dispose() {
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *SampleDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	WithStreaming bool `json:"withStreaming"`
}

func (d *SampleDatasource) query(ctx context.Context, _ backend.PluginContext, query backend.DataQuery) backend.DataResponse {

	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	var objMap map[string]json.RawMessage
	response.Error = json.Unmarshal(query.JSON, &objMap)
	if response.Error != nil {
		response.Error = fmt.Errorf("error when unmarshalling %s: %s", string(query.JSON), response.Error)
		return response
	}
	cmd := string(objMap["queryText"])
	if cmd == "\"\"" || cmd == "" {
		return response
	}

	// trim quotes
	cmd = cmd[1:]
	cmd = cmd[:len(cmd)-1]

	serverUrl := "127.0.0.1:6570"
	conn, err := grpc.DialContext(ctx, serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		response.Error = errors.New(fmt.Sprintf("Error when connect to %s: %s", serverUrl, err))
		return response
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	c := hstream.NewHStreamApiClient(conn)
	req := hstream.CommandQuery{StmtText: cmd}

	grpcQueryResult, err := c.ExecuteQuery(ctx, &req)
	if err != nil {
		response.Error = fmt.Errorf("error when executing %s: %s", cmd, err)
		return response
	}

	framesMap := FlattenResponse(
		grpcQueryResult.GetResultSet())

	lenResults := 0
	for _, v := range framesMap {
		lenResults = len(v)
		break
	}

	timeIs := make([]time.Time, lenResults)
	values := make([]string, lenResults)

	frame := data.NewFrame("response")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, timeIs),
		data.NewField("values", nil, values),
	)

	for k, v := range framesMap {
		frame.Fields = append(frame.Fields,
			data.NewField(k, nil, v),
		)
	}

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)
	return response
}

func FlattenResponse(xs []*structpb.Struct) map[string][]string {
	keySet := map[string]struct{}{}
	for _, x := range xs {
		fields := x.GetFields()["SELECTVIEW"].GetStructValue().GetFields()
		for k := range fields {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, len(keySet))
	i := 0
	for k := range keySet {
		keys[i] = k
		i++
	}

	frameMap := map[string][]string{}
	for _, x := range xs {
		fields := x.GetFields()["SELECTVIEW"].GetStructValue().GetFields()
		for _, key := range keys {
			if val, ok := fields[key]; ok {
				frameMap[key] = append(frameMap[key], val.String())
			} else {
				frameMap[key] = append(frameMap[key], "")
			}
		}
	}

	return frameMap
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *SampleDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	if rand.Int()%2 == 0 {
		status = backend.HealthStatusError
		message = "randomized error"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

// SubscribeStream is called when a client wants to connect to a stream. This callback
// allows sending the first message.
func (d *SampleDatasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	log.DefaultLogger.Info("SubscribeStream called", "request", req)

	status := backend.SubscribeStreamStatusPermissionDenied
	if req.Path == "stream" {
		// Allow subscribing only on expected path.
		status = backend.SubscribeStreamStatusOK
	}
	return &backend.SubscribeStreamResponse{
		Status: status,
	}, nil
}

// RunStream is called once for any open channel.  Results are shared with everyone
// subscribed to the same channel.
func (d *SampleDatasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	log.DefaultLogger.Info("RunStream called", "request", req)

	// Create the same data frame as for query data.
	frame := data.NewFrame("response")

	// Add fields (matching the same schema used in QueryData).
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, make([]time.Time, 1)),
		data.NewField("values", nil, make([]int64, 1)),
	)

	counter := 0

	// Stream data frames periodically till stream closed by Grafana.
	for {
		select {
		case <-ctx.Done():
			log.DefaultLogger.Info("Context done, finish streaming", "path", req.Path)
			return nil
		case <-time.After(time.Second):
			// Send new data periodically.
			frame.Fields[0].Set(0, time.Now())
			frame.Fields[1].Set(0, int64(10*(counter%2+1)))

			counter++

			err := sender.SendFrame(frame, data.IncludeAll)
			if err != nil {
				log.DefaultLogger.Error("Error sending frame", "error", err)
				continue
			}
		}
	}
}

// PublishStream is called when a client sends a message to the stream.
func (d *SampleDatasource) PublishStream(_ context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	log.DefaultLogger.Info("PublishStream called", "request", req)

	// Do not allow publishing at all.
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}
