package db

import (
	"context"
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
)

const networkLogTableName = "network_logs"

func (db *Db) InsertNetworkLog(netLog models.NetworkLog) error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	_, err := db.conn.ExecContext(
		ctx,
		fmt.Sprintf(`insert into
		%s(timestamp,
			request_host, request_method, request_path, request_body, request_args,
			response_status_code, response_duration, response_body)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?)`, networkLogTableName),
		netLog.Timestamp.Unix(),
		netLog.RequestHost, netLog.RequestMethod, netLog.RequestPath, netLog.RequestBody, netLog.RequestArgs,
		netLog.ResponseStatusCode, netLog.ResponseDuration, netLog.ResponseBody,
	)

	return err
}
