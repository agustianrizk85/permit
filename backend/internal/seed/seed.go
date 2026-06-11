package seed

import (
	"log"

	"legalpermit/internal/config"
	"legalpermit/internal/model"
	"legalpermit/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// Accounts seeds the two milestone-1 accounts (DIROPS and KADEP) if the user
// table is empty. Passwords come from the environment so they can be rotated.
func Accounts(users *repository.UserRepository, cfg *config.Config) error {
	count, err := users.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	defaults := []struct {
		name     string
		email    string
		role     model.Role
		password string
	}{
		{"Direktur Operasional", "dirops@greenpark.id", model.RoleDirops, cfg.SeedDiropsPassword},
		{"Kepala Departemen", "kadep@greenpark.id", model.RoleKadep, cfg.SeedKadepPassword},
	}

	for _, d := range defaults {
		hash, err := bcrypt.GenerateFromPassword([]byte(d.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u := &model.User{Name: d.name, Email: d.email, Role: d.role, PasswordHash: string(hash)}
		if err := users.Create(u); err != nil {
			return err
		}
		log.Printf("seeded account: %s (%s)", d.email, d.role)
	}
	return nil
}
