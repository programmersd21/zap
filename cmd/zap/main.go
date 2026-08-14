package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/programmersd21/zap/internal/errs"
	"github.com/programmersd21/zap/internal/ops"
	"github.com/programmersd21/zap/internal/ui"
	"github.com/programmersd21/zap/internal/walk"
)

const version = "0.1.0"

var (
	flagMove      bool
	flagDelete    bool
	flagCopy      bool
	flagForce     bool
	flagRecursive bool
	flagVerbose   bool
	flagVersion   bool
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" {
			printVersion()
			os.Exit(0)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd := &cobra.Command{
		Use:   "zap [flags] <src>... <dst>",
		Short: "blazing fast file ops with gorgeous progress",
		Long: `zap is a modern alternative to cp, mv, and rm.
it shows real-time progress for any operation — from single files to huge directories.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx, args)
		},
	}

	rootCmd.Flags().BoolVarP(&flagMove, "move", "m", false, "move instead of copy")
	rootCmd.Flags().BoolVarP(&flagDelete, "delete", "d", false, "delete instead of copy")
	rootCmd.Flags().BoolVarP(&flagCopy, "copy", "c", false, "copy (explicit, same as default)")
	rootCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "overwrite without prompting")
	rootCmd.Flags().BoolVarP(&flagRecursive, "recursive", "r", false, "recursive (required for deleting directories)")
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "print each file as it completes")
	rootCmd.Flags().BoolVar(&flagVersion, "version", false, "show version")

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Print(ui.RenderHelp(cmd, ui.ThemeMocha))
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if flagVersion {
		printVersion()
		return nil
	}

	if flagDelete {
		if len(args) == 0 {
			return newHintError("delete mode requires at least one path", "zap -d -r <path>...")
		}
		return runDelete(ctx, args)
	}

	if len(args) < 2 {
		return newHintError("copy/move requires source and destination", "zap <src>... <dst>")
	}

	sources := args[:len(args)-1]
	dest := args[len(args)-1]

	if flagMove {
		return runMove(ctx, sources, dest)
	}

	return runCopy(ctx, sources, dest)
}

func runCopy(ctx context.Context, sources []string, dest string) error {
	stats, err := walk.ComputeStats(sources)
	if err != nil {
		return err
	}
	return runCopyWithProgress(ctx, sources, dest, stats)
}

func runCopyDirect(sources []string, dest string) error {
	ec := errs.NewCollector()
	for _, src := range sources {
		dst := destPath(dest, src)
		opts := ops.CopyOptions{Force: flagForce, Errors: ec}
		if err := ops.Copy(src, dst, opts); err != nil {
			ec.Add(src, err)
		}
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "  → %s\n", dst)
		}
	}
	return summarizeErrors(ec)
}

func runCopyWithProgress(ctx context.Context, sources []string, dest string, stats walk.Stats) error {
	if !flagForce {
		reader := bufio.NewReader(os.Stdin)
		for _, src := range sources {
			dst := destPath(dest, src)
			if _, err := os.Lstat(dst); err == nil {
				if !promptYN(reader, fmt.Sprintf("overwrite %s?", dst)) {
					return fmt.Errorf("operation cancelled")
				}
			}
		}
	}

	ec := errs.NewCollector()
	model := ui.NewModel(ui.ThemeMocha, ui.OpCopy, flagVerbose, stats.TotalBytes, stats.TotalFiles)
	p := tea.NewProgram(model)

	go func() {
		var cumulBytes, cumulFiles int64
		for _, src := range sources {
			dst := destPath(dest, src)
			opts := ops.CopyOptions{
				Force:      true,
				Program:    p,
				Errors:     ec,
				CumulBytes: &cumulBytes,
				CumulFiles: &cumulFiles,
			}
			if err := ops.Copy(src, dst, opts); err != nil {
				ec.Add(src, err)
			}
		}
		p.Send(ui.CompletedMsg{})
	}()

	if _, err := p.Run(); err != nil {
		return runCopyDirect(sources, dest)
	}

	if ctx.Err() != nil {
		printInterrupted("copy", stats.TotalFiles)
		os.Exit(130)
	}

	return summarizeErrors(ec)
}

func runMove(ctx context.Context, sources []string, dest string) error {
	ec := errs.NewCollector()

	for _, src := range sources {
		dst := destPath(dest, src)
		opts := ops.MoveOptions{Force: flagForce, Errors: ec}
		if err := ops.Move(src, dst, opts); err != nil {
			ec.Add(src, err)
		}
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "  %s → %s\n", src, dst)
		}
	}

	if ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "\ninterrupted during move")
		os.Exit(130)
	}

	return summarizeErrors(ec)
}

func runDelete(ctx context.Context, paths []string) error {
	stats, err := walk.ComputeStats(paths)
	if err != nil {
		return err
	}

	if !term.IsTerminal(os.Stdout.Fd()) {
		return runDeleteDirect(paths)
	}
	if w, _, err := term.GetSize(os.Stdout.Fd()); err != nil || w == 0 {
		return runDeleteDirect(paths)
	}

	return runDeleteWithProgress(ctx, paths, stats)
}

func runDeleteDirect(paths []string) error {
	ec := errs.NewCollector()
	for _, path := range paths {
		opts := ops.DeleteOptions{Recursive: flagRecursive, Errors: ec}
		if err := ops.Delete(path, opts); err != nil && !flagForce {
			ec.Add(path, err)
		}
		if flagVerbose {
			fmt.Fprintf(os.Stderr, "  deleted %s\n", path)
		}
	}
	return summarizeErrors(ec)
}

func runDeleteWithProgress(ctx context.Context, paths []string, stats walk.Stats) error {
	ec := errs.NewCollector()
	model := ui.NewModel(ui.ThemeMocha, ui.OpDelete, flagVerbose, 0, stats.TotalFiles)
	p := tea.NewProgram(model)

	go func() {
		for _, path := range paths {
			opts := ops.DeleteOptions{
				Recursive: flagRecursive,
				Program:   p,
				Errors:    ec,
			}
			if err := ops.Delete(path, opts); err != nil && !flagForce {
				ec.Add(path, err)
			}
		}
		p.Send(ui.CompletedMsg{})
	}()

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("ui error: %w", err)
	}

	if ctx.Err() != nil {
		printInterrupted("delete", stats.TotalFiles)
		os.Exit(130)
	}

	return summarizeErrors(ec)
}

func destPath(dest, src string) string {
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		return filepath.Join(dest, filepath.Base(src))
	}
	return dest
}

func promptYN(r *bufio.Reader, msg string) bool {
	styles := ui.NewStyles(ui.ThemeMocha)
	fmt.Fprintf(os.Stderr, "%s %s ",
		styles.Warning.Bold(true).Render("?"),
		styles.Primary.Render(msg),
	)
	resp, _ := r.ReadString('\n')
	resp = strings.ToLower(strings.TrimSpace(resp))
	return resp == "y" || resp == "yes"
}

func summarizeErrors(ec *errs.Collector) error {
	if !ec.HasErrors() {
		return nil
	}
	styles := ui.NewStyles(ui.ThemeMocha)
	fmt.Fprint(os.Stderr, "\n"+ec.FormatSummary(styles.Error))
	return fmt.Errorf("completed with %d error(s)", ec.Count())
}

type hintError struct {
	msg  string
	hint string
}

func (e *hintError) Error() string { return e.msg }

func newHintError(msg, hint string) error {
	return &hintError{msg: msg, hint: hint}
}

func printError(err error) {
	styles := ui.NewStyles(ui.ThemeMocha)
	icon := styles.Error.Bold(true).Render("✗")
	msg := styles.Error.Render(err.Error())
	fmt.Fprintf(os.Stderr, "\n%s  %s\n", icon, msg)

	if he, ok := err.(*hintError); ok {
		hint := styles.Muted.Render("→ " + he.hint)
		fmt.Fprintf(os.Stderr, "   %s\n", hint)
	}
	fmt.Fprintln(os.Stderr)
}

func printVersion() {
	fmt.Print(ui.RenderBanner(ui.ThemeMocha))
	styles := ui.NewStyles(ui.ThemeMocha)
	ver := styles.Muted.Render("version")
	num := styles.Blue.Bold(true).Render(version)
	fmt.Printf("  %s %s\n\n", ver, num)
}

func printInterrupted(op string, total int64) {
	styles := ui.NewStyles(ui.ThemeMocha)
	icon := styles.Warning.Bold(true).Render("⚡")
	msg := styles.Muted.Render(fmt.Sprintf("interrupted during %s — %s files in scope", op, ui.FormatNumber(total)))
	fmt.Fprintf(os.Stderr, "\n%s  %s\n\n", icon, msg)
}
