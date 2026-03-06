package service
import (
	"github.com/Nischal07bot/go_boiler_backend/internal/server"
	"github.com/Nischal07bot/go_boiler_backend/internal/repository"
	"github.com/Nischal07bot/go_boiler_backend/internal/lib/job"
)

type Services struct {
	Auth *AuthService
	Job *job.Jobservice
}

func NewServices(s *server.Server, r *repository.Repositories) (*Services,error){
	authService := NewAuthService(s)
	return &Services{
		Auth: authService,
		Job: s.Job,
	}, nil
}
