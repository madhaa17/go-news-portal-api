package seeds

import (
	"news-app/internal/core/domain/model"

	"news-app/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := conv.HashPassword("password123")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to hash password")
	}

	admin := model.User{
		Name:     "Admin",
		Email:    "admin@mail.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@mail.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to seed admin role")
	} else {
		log.Info().Msg("Admin role seeded successfully")
	}
}
