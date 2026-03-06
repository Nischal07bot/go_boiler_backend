package repository
import "github.com/Nischal07bot/go_boiler_backend/internal/server"

type Repositories struct {}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}