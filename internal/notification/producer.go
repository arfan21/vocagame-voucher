package notification

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/model"
)

func (n Notification) Produce(ctx context.Context, event model.Event) error {
	_, err := n.stream.Add(ctx, event)
	if err != nil {
		err = fmt.Errorf("failed to produce event: %w", err)
		return err
	}

	return nil
}
