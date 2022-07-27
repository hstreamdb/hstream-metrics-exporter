package plugin_test

import (
	"context"
	"fmt"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/gen/hstreamdb/hstream/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

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
	req := server.CommandQuery{StmtText: "select * from vv where x = 0;"}

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

	xs := plugin.FlattenResponse(resultSet)
	for k, v := range xs {
		fmt.Println(k, v)
	}

}
