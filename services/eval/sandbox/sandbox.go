package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/marius004/phoenix-/internal"
	"github.com/marius004/phoenix-/models"
	"github.com/marius004/phoenix-/services/eval"
)

type Sandbox struct {
	path string
	id   int

	metaFile string
	mutex    sync.Mutex

	evalConfig *internal.EvalConfig
}

func (s *Sandbox) GetPath(path string) string {
	if path == "" {
		return s.path
	}

	return s.path + "/" + path
}

func (s *Sandbox) GetID() int {
	return s.id
}

func (s *Sandbox) CreateDirectory(path string, perm fs.FileMode) error {
	fullPath := s.GetPath(path)
	return os.Mkdir(fullPath, perm)
}

func (s *Sandbox) DeleteDirectory(path string) error {
	fullPath := s.GetPath(path)
	return os.RemoveAll(fullPath)
}

func (s *Sandbox) WriteToFile(path string, data []byte, perm fs.FileMode) error {
	fullPath := s.GetPath(path)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (s *Sandbox) FileExists(path string) bool {
	fullPath := s.GetPath(path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *Sandbox) CreateFile(path string, perm fs.FileMode) error {
	fullPath := s.GetPath(path)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY|os.O_TRUNC, perm)

	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}

func (s *Sandbox) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(s.GetPath(path))
}

func (s *Sandbox) DeleteFile(path string) error {
	return os.Remove(s.GetPath(path))
}

func (s *Sandbox) ExecuteCommand(ctx context.Context, command []string, config *models.RunConfig) (*models.RunStatus, error) {
	logger := internal.GetGlobalLoggerInstance()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	metaFile := path.Join(os.TempDir(), "pn-"+eval.RandomString(24))
	s.metaFile = metaFile

	defer func() { s.metaFile = "" }()

	params := append(s.buildRunFlags(config), command...)

	logger.Println("Command to be executed:", "isolate", params)
	cmd := exec.CommandContext(ctx, s.evalConfig.IsolatePath, params...)

	cmd.Stdin = config.Stdin
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr

	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		metaData, err := parseMetaFile(s.metaFile)
		if err != nil {
			logger.Println(err)
			return nil, err
		} else { // the program was stopped because of the time or memory constraints.
			return metaData, nil
		}
	}

	return parseMetaFile(s.metaFile)
}

func (s *Sandbox) Cleanup() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var params []string

	params = append(params, "--cg")
	params = append(params, "--box-id="+fmt.Sprintf("%d", s.id))
	params = append(params, "--cleanup")

	return exec.Command(s.evalConfig.IsolatePath, params...).Run()
}

func (s *Sandbox) buildRunFlags(config *models.RunConfig) (res []string) {
	res = append(res, "--box-id="+strconv.Itoa(s.id))
	res = append(res, "--cg", "--cg-timing")

	res = append(res, "--full-env")

	if config.TimeLimit != 0 {
		res = append(res, "--time="+strconv.FormatFloat(config.TimeLimit, 'f', -1, 64))
	}

	if config.WallTimeLimit != 0 {
		res = append(res, "--wall-time="+strconv.FormatFloat(config.WallTimeLimit, 'f', -1, 64))
	}

	if config.MemoryLimit != 0 {
		res = append(res, "--mem="+strconv.Itoa(config.MemoryLimit))
	}

	if config.StackLimit != 0 {
		res = append(res, "--stack="+strconv.Itoa(config.StackLimit))
	}

	if config.MaxProcesses == 0 {
		res = append(res, "--processes=1")
	} else {
		res = append(res, "--processes="+strconv.Itoa(config.MaxProcesses))
	}

	if config.InputPath != "" {
		res = append(res, "--stdin="+config.InputPath)
	}

	if config.OutputPath != "" {
		res = append(res, "--stdout="+config.OutputPath)
	}

	if s.metaFile != "" {
		res = append(res, "--meta="+s.metaFile)
	}

	res = append(res, "--silent", "--run", "--")
	return
}

// parseMetaFile parses the meta file that contains information about the execution of a particular task executed within the sandbox.
func parseMetaFile(path string) (*models.RunStatus, error) {
	r, err := os.OpenFile(path, os.O_RDONLY, 0777)

	if err != nil {
		return nil, err
	}

	var ret = new(models.RunStatus)
	s := bufio.NewScanner(r)

	for s.Scan() {
		if !strings.Contains(s.Text(), ":") {
			continue
		}

		l := strings.SplitN(s.Text(), ":", 2)
		switch l[0] {
		case "cg-mem":
			ret.Memory, _ = strconv.Atoi(l[1])
		case "exitcode":
			ret.ExitCode, _ = strconv.Atoi(l[1])
		case "exitsig":
			ret.ExitSignal, _ = strconv.Atoi(l[1])
		case "killed":
			ret.Killed = true
		case "message":
			ret.Message = l[1]
		case "status":
			ret.Status = l[1]
		case "time":
			ret.Time, _ = strconv.ParseFloat(l[1], 32)
		case "time-wall":
			ret.WallTime, _ = strconv.ParseFloat(l[1], 32)
		default:
			continue
		}
	}

	return ret, nil
}

func newSandbox(id int, evalConfig *internal.EvalConfig) (*Sandbox, error) {
	logger := internal.GetGlobalLoggerInstance()
	ret, err := exec.Command(evalConfig.IsolatePath, fmt.Sprintf("--box-id=%d", id), "--cg", "--init").CombinedOutput()

	if strings.HasPrefix(string(ret), "Box already exists") {
		exec.Command(evalConfig.IsolatePath, fmt.Sprintf("--box-id=%d", id), "--cg", "--cleanup").Run()
		return newSandbox(id, evalConfig)
	}

	if err != nil {
		logger.Println("Could not create sandbox", err)
		return nil, err
	}

	path := string(ret)
	return &Sandbox{path: strings.TrimSpace(path), id: id, evalConfig: evalConfig}, nil
}
