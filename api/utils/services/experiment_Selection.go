package services

import (
	"archive-api/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Response struct {
	VariableName       string                 `json:"variable_name"`
	Path_ts            []string               `json:"paths_ts"`
	Path_mean          []string               `json:"paths_mean"`
	Levels             int                    `json:"levels"`
	Timesteps          int                    `json:"timesteps"`
	Xsize              int                    `json:"xsize"`
	Xfirst             float32                `json:"xfirst"`
	Yinc               float32                `json:"xinc"`
	Ysize              int                    `json:"ysize"`
	Yfirst             float32                `json:"yfirst"`
	Xinc               float32                `json:"yinc"`
	Metadata           map[string]interface{} `json:"metadata"`
	Created_at         time.Time              `json:"created_at"`
	Config_name        string                 `json:"config_name"`
	Extension          string                 `json:"extension"`
	Lossless           bool                   `json:"lossless"`
	Nan_value_encoding int                    `json:"nan_value_encoding"`
	Chunks             int                    `json:"chunks"`
	Rx                 float64                `json:"rx"`
	Ry                 float64                `json:"ry"`
	Exp_id             string                 `json:"exp_id"`
	Threshold          float32                `json:"threshold"`
}

func GetExperimentByIDWithConfig(id string, c *fiber.Ctx, pool *pgxpool.Pool, psql *squirrel.StatementBuilderType) error {

}
func GetExperimentByID(id string, c *fiber.Ctx, pool *pgxpool.Pool, psql *squirrel.StatementBuilderType) error {
	pl := new(utils.Placholder)
	pl.Build(0, 1)

	sql := fmt.Sprintf(`with nimbus_run as 
	(
		select *
		from table_nimbus_execution 
		where exp_id = %s
		order by created_at desc
		limit 1
	)
	select 
		name as variable_name,
		paths_ts,
		paths_mean,levels,
		timesteps,
		xsize,
		xfirst,
		xinc,
		ysize,
		yfirst,
		yinc,
		metadata,
		created_at,
		config_name,
		extension,
		lossless,
		nan_value_encoding,
		chunks,
		rx,
		ry,
		exp_id,
		threshold
	from table_variable
	inner join 
		( 
			select * 
			from join_nimbus_execution_variables
			inner join nimbus_run 
			on join_nimbus_execution_variables.id_nimbus_execution = nimbus_run.id
		) as joined
	on table_variable.id = joined.variable_id`, pl.Get(id))
	rows, err := pool.Query(context.Background(), sql, pl.Args...)
	if err != nil {
		fmt.Println(pl.Args...)
		log.Default().Println("Unable to query:", sql, "error :", err)
		return err
	}
	defer rows.Close()
	responses, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Response, error) {
		var res Response
		err := row.Scan(
			&res.VariableName,
			&res.Path_ts,
			&res.Path_mean,
			&res.Levels,
			&res.Timesteps,
			&res.Xsize,
			&res.Xfirst,
			&res.Yinc,
			&res.Ysize,
			&res.Yfirst,
			&res.Xinc,
			&res.Metadata,
			&res.Created_at,
			&res.Config_name,
			&res.Extension,
			&res.Lossless,
			&res.Nan_value_encoding,
			&res.Chunks,
			&res.Rx,
			&res.Ry,
			&res.Exp_id,
			&res.Threshold,
		)
		if err != nil {
			log.Default().Println(err)
		}
		return res, err
	})
	return c.JSON(responses)
}