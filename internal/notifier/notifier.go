package notifier

import (
	"sync"

	"github.com/google/uuid"

	"github.com/kaigoh/monate/v2/internal/data"
)

// InvoiceNotifier fans out invoice updates to subscribers that back the GraphQL subscriptions.
type InvoiceNotifier struct {
	mu          sync.RWMutex
	subscribers map[uuid.UUID]map[chan *data.Invoice]struct{}
}

func NewInvoiceNotifier() *InvoiceNotifier {
	return &InvoiceNotifier{
		subscribers: make(map[uuid.UUID]map[chan *data.Invoice]struct{}),
	}
}

// Subscribe registers a channel for invoice updates and returns it alongside a cleanup func.
func (n *InvoiceNotifier) Subscribe(id uuid.UUID) (<-chan *data.Invoice, func()) {
	ch := make(chan *data.Invoice, 1)

	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.subscribers[id]; !ok {
		n.subscribers[id] = make(map[chan *data.Invoice]struct{})
	}
	n.subscribers[id][ch] = struct{}{}

	return ch, func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		if subs, ok := n.subscribers[id]; ok {
			if _, exists := subs[ch]; exists {
				delete(subs, ch)
				close(ch)
			}
			if len(subs) == 0 {
				delete(n.subscribers, id)
			}
		}
	}
}

// Publish sends an update to every subscriber waiting on the invoice ID.
func (n *InvoiceNotifier) Publish(inv *data.Invoice) {
	if inv == nil {
		return
	}
	n.mu.RLock()
	defer n.mu.RUnlock()
	if subs, ok := n.subscribers[inv.ID]; ok {
		for ch := range subs {
			select {
			case ch <- inv:
			default:
			}
		}
	}
}
