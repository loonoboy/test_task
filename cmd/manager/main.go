package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "Manager of workers",
		Usage: "Start multiple worker processes for tubes",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "count",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "Number of workers to launch per tube",
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			tubes := []string{"create_contact", "delete_contact", "update_contact"}
			count := command.Int("count")

			var wg sync.WaitGroup

			for _, tube := range tubes {
				for i := 1; i <= count; i++ {
					args := fmt.Sprintf("--tube=%s", tube)
					worker := exec.CommandContext(ctx, "./worker", args)
					worker.Stdout = os.Stdout
					worker.Stderr = os.Stderr

					if err := worker.Start(); err != nil {
						log.Printf("Failed to start worker #%d for tube %s: %v", i, tube, err)
						continue
					}

					log.Printf("Created worker #%d for tube %s (pid: %d)", i, tube, worker.Process.Pid)

					wg.Add(1)
					go func(cmd *exec.Cmd, idx int, t string) {
						defer wg.Done()
						if err := cmd.Wait(); err != nil {
							log.Printf("Worker #%d for tube %s exited with error: %v", idx, t, err)
						} else {
							log.Printf("Worker #%d for tube %s finished gracefully", idx, t)
						}
					}(worker, i, tube)
				}
			}

			wg.Wait()
			log.Println("All workers finished.")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
