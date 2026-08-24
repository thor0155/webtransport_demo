package cmder

import (
	"api/internal/frameworks/config"
	"api/internal/frameworks/obj"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	maxPhaseCmdLen int    = 5
	exitCmd        string = "exit"
)

var (
	ErrExceedingPhaseCmdLen            = errors.Newf("exceed phase command length %d", maxPhaseCmdLen)
	ErrCmdDuplicated                   = errors.New("command duplicated")
	ErrEmptyCommand                    = errors.New("empty command")
	ErrClosed                          = errors.New("closed")
	_                       obj.Server = (*Server)(nil)
)

type Server struct {
	name                     string
	lifecycleShutdownEmitter *obj.ShutdownEmitter
	cfg                      *CmderConfig
	logger                   *zap.Logger
	shutdownCancelFunc       context.CancelFunc
	shutdownCancelCtx        context.Context
	cmders                   Cmders
	rootCmd                  *cobra.Command
	phaseCmds                []map[string]*cobra.Command
}

// Close implements [obj.Server].
func (s *Server) Close(_ context.Context) error {
	s.shutdownCancelFunc()
	return nil
}

// Name implements [obj.Server].
func (s *Server) Name() string {
	return s.name
}

// Start implements [obj.Server].
func (s *Server) Start() {
	if err := s.handleCmders(); err != nil {
		s.logger.Error("stopped", zap.Error(err))
		return
	}
	if err := s.execute(); err != nil {
		s.logger.Error("stopped", zap.Error(err))
	}
}

func NewCmderServer(cfg *CmderConfig, cmders Cmders, logger *zap.Logger, lifecycleShutdownEmitter *obj.ShutdownEmitter) *Server {
	name := "cmder"
	server := &Server{
		cfg:                      cfg,
		name:                     name,
		logger:                   logger.Named(name),
		lifecycleShutdownEmitter: lifecycleShutdownEmitter,
		cmders:                   cmders,
		phaseCmds:                make([]map[string]*cobra.Command, maxPhaseCmdLen),
		rootCmd:                  &cobra.Command{Use: config.GetAppName(), Short: config.GetAppName() + " CLI"},
	}
	for i := 0; i < len(server.phaseCmds); i++ {
		server.phaseCmds[i] = make(map[string]*cobra.Command)
	}
	server.shutdownCancelCtx, server.shutdownCancelFunc = context.WithCancel(context.Background())
	return server
}

func (s *Server) handleCmders() error {
	for _, cmder := range s.cmders {
		command := cmder.Command()
		if command == "" {
			return errors.Wrapf(ErrEmptyCommand, "with %T", cmder)
		}
		commands := strings.Split(command, " ")
		if len(commands) > maxPhaseCmdLen {
			return errors.Wrapf(ErrExceedingPhaseCmdLen, "with %T", cmder)
		}

		lastIdx := len(commands) - 1
		var prevCmd *cobra.Command
		for i, command := range commands {
			command = strings.TrimSpace(command)
			cmd, exist := s.phaseCmds[i][command]
			if !exist {
				cmd = &cobra.Command{
					Use:   command,
					Short: command + " features",
				}
				s.phaseCmds[i][command] = cmd
				if i == 0 {
					s.rootCmd.AddCommand(cmd)
				} else {
					prevCmd.AddCommand(cmd)
				}
			}
			prevCmd = cmd

			if i != lastIdx {
				continue
			} else if exist {
				return errors.Wrapf(ErrCmdDuplicated, "[%T] -> %s", cmder, command)
			}

			cmd.Short = cmder.Description()
			cmd.Args = cmder.ProvideArgs()
			cmder.FlagSet(cmd.Flags())
			cmd.Run = func(cmd *cobra.Command, args []string) {
				ctx := newCmderContext(s.shutdownCancelCtx, cmd, args)
				cmder.Run(ctx)
				errs := ctx.GetErrors()
				if len(errs) > 0 {
					for _, err := range errs {
						s.logger.Error("after running",
							zap.String("cmder-handler", fmt.Sprintf("%T", cmder)),
							zap.Error(err))
					}
				}
			}
		}
	}
	return nil
}

func (s *Server) execute() error {

	if s.cfg.EnableReadEvalPrintLoop {

		reader := bufio.NewReader(os.Stdin)
		fmt.Println(config.GetAppName() + " started. Type 'help' or 'exit'.")

		for {
			select {
			case <-s.shutdownCancelCtx.Done():
				return s.shutdownCancelCtx.Err()
			default:
				fmt.Print(config.GetAppName() + "> ")
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)
				if line == exitCmd {
					s.lifecycleShutdownEmitter.Shutdown("exit")
					return ErrClosed
				}
				if line == "" {
					continue
				}
				s.rootCmd.SetArgs(strings.Split(line, " "))
				if err := s.rootCmd.ExecuteContext(s.shutdownCancelCtx); err != nil {
					s.logger.Error("execute", zap.Error(err))
					continue
				}
			}
		}
	} else {
		defer s.lifecycleShutdownEmitter.Shutdown("exit")
		return s.rootCmd.Execute()
	}
}
