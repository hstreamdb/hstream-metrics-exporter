package plugin_test

import (
	"context"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/gen/hstreamdb/hstream/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/plugin"
)

// This is where the tests for the datasource backend live.
func TestQueryData(t *testing.T) {
	ds := plugin.SampleDatasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}

func TestMainFunc(t *testing.T) {
	serverUrl := "127.0.0.1:6570"
	conn, err := grpc.DialContext(context.Background(), serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	c := server.NewHStreamApiClient(conn)
	req := server.CommandQuery{StmtText: "select * from vvv where i = 1;"}

	resp, err := c.ExecuteQuery(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}

	resultSet := resp.GetResultSet()
	for _, xs := range resultSet {
		fields := xs.GetFields()
		_, ok := fields["SELECTVIEW"]
		fmt.Println(ok)
		for k, x := range fields["SELECTVIEW"].GetStructValue().GetFields() {
			fmt.Println(k, x)
		}
	}

	framesMap := plugin.FlattenResponse(resultSet)
	for k, v := range framesMap {
		fmt.Println(k, v)
	}

	{
		response := backend.DataResponse{}

		lenResults := 1
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

		for _, frame := range response.Frames {
			tab, _ := frame.StringTable(-1, -1)
			fmt.Println(tab)
		}

	}

}
