package boot

import (
	"fmt"
	"os"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("cannot load environment variables")
	}

	models.Setup(os.Getenv("DSN"))

	return err
}
