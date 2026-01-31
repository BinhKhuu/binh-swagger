package commands

import (
	helper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/generate"
)

type Helpers struct {
	File     helper.FileHelper
	Generate generate.Orchestrator
}
