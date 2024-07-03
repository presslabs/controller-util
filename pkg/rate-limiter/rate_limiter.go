package ratelimiter

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RateLimiter is a rate limiter for controllers.
type RateLimiter struct {
	c                     client.Client
	log                   logr.Logger
	mu                    sync.Mutex
	maxItems              int
	notReadyItems         int
	items                 map[types.NamespacedName]time.Time
	itemIsReady           func(context.Context, client.Client, types.NamespacedName, logr.Logger) bool
	durationToBecomeReady time.Duration
	logFrequency          time.Duration
	itemPoolingInterval   time.Duration
	itemTimeout           time.Duration
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(
	c client.Client,
	log logr.Logger,
	maxItems int,
	itemIsReady func(context.Context, client.Client, types.NamespacedName, logr.Logger) bool,
	durationToBecomeReady time.Duration,
	logFrequency time.Duration,
	itemPoolingInterval time.Duration,
	itemTimeout time.Duration,
) *RateLimiter {
	return &RateLimiter{
		c:                     c,
		log:                   log,
		mu:                    sync.Mutex{},
		maxItems:              maxItems,
		items:                 map[types.NamespacedName]time.Time{},
		itemIsReady:           itemIsReady,
		durationToBecomeReady: durationToBecomeReady,
		logFrequency:          logFrequency,
		itemPoolingInterval:   itemPoolingInterval,
		itemTimeout:           itemTimeout,
	}
}

// SetupWithManager instantiate the RateLimiter, managed by a given manager.
func (r *RateLimiter) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(r)
}

func (r *RateLimiter) checkAndUpdateItems(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.notReadyItems = 0

	for nsName, inTime := range r.items {
		// give the operators/controllers 10 seconds to update item status
		if time.Since(inTime) < r.durationToBecomeReady {
			continue
		}

		if time.Since(inTime) > r.itemTimeout {
			delete(r.items, nsName)

			r.log.V(0).Info("timeout exceeded", "item", nsName)

			continue
		}

		// check item. If it is ready, remove it from buffer
		if r.itemIsReady(ctx, r.c, nsName, r.log) {
			delete(r.items, nsName)

			continue
		}

		r.notReadyItems++
	}
}

func (r *RateLimiter) writeLog() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.items) > 0 {
		r.log.Info("stats",
			"noItems", fmt.Sprintf("%d/%d", len(r.items), r.maxItems),
			"notReadyItems", strconv.Itoa(r.notReadyItems),
		)
	}
}

// Start will start the RateLimiter.
func (r *RateLimiter) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	defer wg.Wait()

	go func() {
		defer wg.Done()

		for {
			select {
			case <-time.After(r.logFrequency):
				r.writeLog()
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-time.After(r.itemPoolingInterval):
			r.checkAndUpdateItems(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

// ShouldReconcile check the given item.
// if the item is ready for reconciliation, ShouldReconcile removes the item from buffer and returns true.
func (r *RateLimiter) ShouldReconcile(nsName types.NamespacedName) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.items) < r.maxItems {
		r.items[nsName] = time.Now()

		return true
	}

	return false
}
