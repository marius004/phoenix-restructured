package eval

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"strconv"
	"strings"

	"github.com/marius004/phoenix-algo/entities"
	"github.com/marius004/phoenix-algo/internal"
	"github.com/marius004/phoenix-algo/models"
)

const randomCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(size int) string {
	sb := strings.Builder{}
	sb.Grow(size)

	for ; size > 0; size-- {
		randIndex := rand.Intn(len(randomCharacters))
		sb.WriteByte(randomCharacters[randIndex])
	}

	return sb.String()
}

func CompileFile(ctx context.Context, sandbox internal.Sandbox, sourceCode []byte, lang internal.EvalConfigLanguage) (string, error) {
	if err := sandbox.WriteToFile(lang.SourceFile, sourceCode, 0644); err != nil {
		return "", err
	}

	var runConf models.RunConfig

	out := &bytes.Buffer{}

	runConf.Stdout = out
	runConf.Stderr = out
	runConf.MaxProcesses = 5

	if _, err := sandbox.ExecuteCommand(ctx, lang.Compile, &runConf); err != nil {
		return out.String(), err
	}

	return out.String(), nil
}

func CompiledSourceCode(sandbox internal.Sandbox, fileName string) bool {
	return sandbox.FileExists(fileName)
}

func CopyFromSandbox(sandbox internal.Sandbox, path string, w io.Writer) error {
	content, err := sandbox.ReadFile(path)

	if err != nil {
		return err
	}

	if _, err := w.Write(content); err != nil {
		return err
	}

	return nil
}

func CopyInSandbox(sandbox internal.Sandbox, path string, data []byte) error {
	return sandbox.WriteToFile(path, data, 7777)
}

func ExecuteFile(ctx context.Context, sandbox internal.Sandbox, lang internal.EvalConfigLanguage, problemId int, limit models.Limit) (*models.RunStatus, error) {
	var runConf models.RunConfig

	runConf.MaxProcesses = 10
	runConf.MemoryLimit = limit.Memory
	runConf.TimeLimit = limit.Time
	runConf.StackLimit = limit.Stack
	runConf.WallTimeLimit = 5 * limit.Time

	runConf.InputPath = strconv.Itoa(problemId) + ".in"
	runConf.OutputPath = strconv.Itoa(problemId) + ".out"

	return sandbox.ExecuteCommand(ctx, lang.Execute, &runConf)
}

func GetOutputFileName(evalConfig *internal.EvalConfig, submission *entities.Submission, problemTest *entities.ProblemTest) string {
	return evalConfig.OutputPath + "/s" + strconv.Itoa(int(submission.ID)) + "t" + strconv.Itoa(int(problemTest.ID)) + ".out"
}
