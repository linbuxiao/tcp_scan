package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
)

// 这种模式下，scanner 只负责两部分
// 1. 初始化 channel ports 段
// 2. 发送 port 给 channel
// 3. 接受线程池部分交给用户端调用 接受 一个 context，可以细粒度的取消
// 4. 在使用完后 defer cancel 退出，避免 协程泄露

func main() {
	app := &cli.App{
		Name: "tcp scanner",
		Authors: []*cli.Author{
			{
				Name:  "linbuxiao",
				Email: "linbuxiao@gmail.com",
			},
		},
		Usage: "Scanner some ports",
		Commands: []*cli.Command{
			{
				Name: "range",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "start",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "end",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					first := c.Int("start")
					end := c.Int("end")
					workers := c.Int("workers")
					cfg := NewConfig(WithWorkerNum(workers))
					scanner := NewRangeScanner(first, end)
					runScanner(scanner, cfg)
					return nil
				},
			},
			{
				Name: "ports",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					workers := c.Int("workers")
					cfg := NewConfig(WithWorkerNum(workers))
					var numSlice []int
					for _, v := range args {
						k, err := strconv.Atoi(v)
						if err != nil {
							log.Fatal("port must be number:", v)
						}
						numSlice = append(numSlice, k)
					}
					scanner := NewScanner(numSlice)
					runScanner(scanner, cfg)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "workers",
				Value: 10,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runScanner(scanner *Scanner, cfg *config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < cfg.workerNum; i++ {
		go scanner.StartWorker(ctx, cfg.timeout)
	}
	scanner.Run()
}
