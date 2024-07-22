package beat

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/event"
)

// Beat put a generic event on a channel after a certain time (e.g.: every day).
type Beat struct {
	C              chan event.GenericEvent
	TickerDuration time.Duration
}

// NewBeat instantiate a beat.
func NewBeat(tickerDuration time.Duration) *Beat {
	return &Beat{
		C:              make(chan event.GenericEvent),
		TickerDuration: tickerDuration,
	}
}

// Start watches for events on a channel after a certain time (e.g.: every day). It's designed to be run by a manager.
func (b *Beat) Start(ctx context.Context) error {
	ticker := time.NewTicker(b.TickerDuration)

	go func() {
		b.C <- event.GenericEvent{}

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				b.C <- event.GenericEvent{}
			}
		}
	}()

	return nil
}
