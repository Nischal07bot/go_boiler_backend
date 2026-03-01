package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Database struct {
	Pool *pgxpool.Pool //this will be exported to be used in other packages do convention capital letters
	log  *zerolog.Logger
}

// allows chaining multiple tracers together to create a single tracer that can be used for logging and monitoring purposes. This is useful for applications that need to support multiple tracing systems or want to combine different types of tracers for more comprehensive observability.
type multitracer struct {
	tracers []any
}

func (mt *multitracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryStart(context.Context, *pgx.Conn, pgx.TraceQueryStartData) context.Context
		}); ok {
			ctx=t.TraceQueryStart(ctx, conn, data)		
		}
	}
}
func (mt *multitracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData) context.Context
		}); ok {
			ctx=t.TraceQueryEnd(ctx, conn, data)		

		}
	}
}
const DatabasePingTimeout = 10
