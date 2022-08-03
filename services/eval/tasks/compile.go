package tasks

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/marius004/phoenix-/internal"
	"github.com/marius004/phoenix-/models"
	"github.com/marius004/phoenix-/services/eval"
)

type CompileTask struct {
	EvalConfig *internal.EvalConfig

	Request  *models.CompileRequest
	Response *models.CompileResponse
}

func (task *CompileTask) Run(ctx context.Context, sandbox internal.Sandbox) error {
	logger := internal.GetGlobalLoggerInstance()

	logger.Printf("Compiling using sandbox %d\n", sandbox.GetID())
	lang, ok := task.EvalConfig.Languages[task.Request.Lang]

	if !ok {
		logger.Printf("Invalid language %s\n", task.Request.Lang)
		return internal.ErrLangNotFound
	}

	binaryPath := path.Join(task.EvalConfig.CompilePath, fmt.Sprintf("%d.bin", task.Request.ID))
	task.Response.Success = true

	if lang.IsCompiled {
		message, err := eval.CompileFile(ctx, sandbox, task.Request.SourceCode, lang)
		task.Response.Message = message

		if err != nil {
			task.Response.Success = false
			task.Response.Message = err.Error()
			logger.Printf("Could not compile %s\n", err.Error())
			return err
		}

		if !eval.CompiledSourceCode(sandbox, lang.Executable) {
			task.Response.Success = false
			return nil
		}

		file, err := os.OpenFile(binaryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)

		if err != nil {
			task.Response.Message = err.Error()
			task.Response.Success = false
			logger.Printf("Could not create the binary file %s\n", err.Error())
			return err
		}

		if err := eval.CopyFromSandbox(sandbox, lang.Executable, file); err != nil {
			task.Response.Message = err.Error()
			task.Response.Success = false
			logger.Printf("Could not copy the binary file from sandbox %d %s\n", sandbox.GetID(), err.Error())
			return err
		}

		if err := file.Close(); err != nil {
			task.Response.Message = err.Error()
			task.Response.Success = false
			logger.Printf("Could not close the binary file %s\n", err.Error())
			return err
		}

		return nil
	}

	// TODO dynamic typed languages

	return nil
}
