package tasks

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/marius004/phoenix/internal"
	"github.com/marius004/phoenix/services/eval"
)

type ExecuteTask struct {
	EvalConfig *internal.EvalConfig

	Request  *internal.ExecuteRequest
	Response *internal.ExecuteResponse
}

func (t *ExecuteTask) Run(ctx context.Context, sandbox internal.Sandbox) error {
	logger := internal.GetGlobalLoggerInstance()

	logger.Printf("Executing using sandbox %d\n", sandbox.GetID())
	lang, ok := t.EvalConfig.Languages[t.Request.Lang]

	if !ok {
		logger.Printf("Invalid language %s\n", t.Request.Lang)
		return internal.ErrLangNotFound
	}

	if err := sandbox.WriteToFile("box/"+strconv.Itoa(int(t.Request.ProblemId))+".in", t.Request.Input, 0644); err != nil {
		logger.Println("Can't write the input file in the sandbox", err)
		t.Response.Message = "Sandbox error: Cannot copy input file to the sandbox"
		return err
	}

	if err := sandbox.CreateFile("box/"+strconv.Itoa(int(t.Request.ProblemId))+".out", 0644); err != nil {
		logger.Println("Can't write the output file in the sandbox", err)
		t.Response.Message = "Sandbox error: Cannot copy output file to the sandbox"
		return err
	}

	binaryPath := path.Join(t.EvalConfig.CompilePath, fmt.Sprintf("%d.bin", t.Request.ID))
	binaryFile, err := os.OpenFile(binaryPath, os.O_RDONLY, 0644)

	if err != nil {
		logger.Println("Could not open the binary file", err)
		t.Response.Message = "Sandbox error: Could not open the binary file"
		return err
	}

	bin, err := ioutil.ReadAll(binaryFile)

	if err != nil {
		logger.Println("Could not open the read the binary file", err)
		t.Response.Message = "Sandbox error: Could not read the binary file"
		return err
	}

	if err := eval.CopyInSandbox(sandbox, lang.Executable, bin); err != nil {
		logger.Println(fmt.Sprintf("Could not copy the binary file in the sandbox %d", sandbox.GetID()), err)
		t.Response.Message = fmt.Sprintf("Could not copy the binary file in the sandbox %d", sandbox.GetID())
		return err
	}

	limit := internal.Limit{
		Time: t.Request.Time,

		Memory: t.Request.Memory,
		Stack:  t.Request.Stack,
	}

	metaFile, err := eval.ExecuteFile(ctx, sandbox, lang, int(t.Request.ProblemId), limit)
	if err != nil {
		logger.Println("could not execute the program", err)
		t.Response.Message = fmt.Sprintf("Could not execute the program %s", err.Error())
		return err
	}

	t.Response.TimeUsed = metaFile.Time
	t.Response.MemoryUsed = metaFile.Memory

	switch metaFile.Status {
	case "RO":
		t.Response.Message = metaFile.Message
	case "RE":
		t.Response.Message = metaFile.Message
	case "RG":
		t.Response.Message = metaFile.Message
	case "RX":
		t.Response.Message = metaFile.Message
	}

	if t.Response.ExitCode == 0 {
		path := fmt.Sprintf("%s/s%dt%d.out", t.EvalConfig.OutputPath, t.Request.ID, t.Request.TestId)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

		if err != nil {
			logger.Println(err)
			return err
		}

		if err := eval.CopyFromSandbox(sandbox, "box/"+strconv.Itoa(int(t.Request.ProblemId))+".out", file); err != nil {
			logger.Println("Could not copy the output file " + err.Error())
			t.Response.Message = "Could not copy the output file " + err.Error()
			return err
		}
	}

	return nil
}
