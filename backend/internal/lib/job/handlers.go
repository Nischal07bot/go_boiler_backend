package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/Nischal07bot/go_boiler_backend/internal/config"
	"github.com/Nischal07bot/go_boiler_backend/internal/lib/email"
)