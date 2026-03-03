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

var emailClient *email.Client

func (j *Jobservice) InitHandlers(config *config.Config, logger *zerolog.Logger){
	emailClient = email.NewClient(config , logger)
}

func (j *Jobservice) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload//it is a struct 
	if err := json.Unmarshal(t.Payload(), &p); err != nil {//unmarshal the payload of the task asynq provides .payload() function which gives json unmarshall to get go bytes[]
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	j.Logger.Info().
	    Str("type","welcome").//log entry at info level adding field type=welcome 
		Str("to",p.To).//add field to ="" email
		Msg("Processing welcome email task")//the messgae
	
	err := emailClient.SendWelcomemail(p.To,p.FirstName)//sending the welcome mail which further triggers the rest of the things
	//not directly sendemail

	if err != nil {
		j.Logger.Error().
		    Str("type","welcome").
			Str("to",p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.Logger.Info().Str("type","welcome").
	Str("to",p.To).
	Msg("Successfully sent welcome email")
	return nil
}//hence handler helping in sending email 