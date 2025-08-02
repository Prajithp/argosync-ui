package main

import (
	"github.com/Prajithp/argosync/pkg/models"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/repository/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.ApplyBasic(models.Deployment{}, models.Application{}, models.Environment{}, models.Region{})

	g.Execute()
}
