package logger

import (
	"context"
	"fmt"
)

// TODO: Use Library
func Error(ctx context.Context, txt string, err error) {
	fmt.Printf("%s: %s\n", txt, err.Error())
}

func Warn(ctx context.Context, txt string) {
	fmt.Println(txt)
}

func Info(ctx context.Context, txt string) {
	fmt.Println(txt)
}

func LogService(ctx context.Context, serviceName string, err error) {
	if err != nil {
		Error(ctx, serviceName, err)
		return
	}

	Info(ctx, fmt.Sprintf("%s success", serviceName))
}
