package stats

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlChangeEvent(t *testing.T) {
	var event ScriptChangeEvent

	var g sync.WaitGroup

	g.Add(1)
	go func() {
		defer g.Done()
		for i := 0; i < 199; i++ {
			event.Write(&ChangeEvent{
				Name: fmt.Sprintf("%d.p", i),
				NS:   fmt.Sprintf("%d", i),
				Op:   EventOpAdd,
			})
		}
	}()
	g.Add(1)
	go func() {
		defer g.Done()
		for i := 0; i < 299; i++ {
			event.Read()
		}
	}()
	g.Wait()

	event = ScriptChangeEvent{}

	tmp := []ChangeEvent{}
	for i := 0; i < 256; i++ {
		assert.Equal(t, tmp, event.Read())
		c := ChangeEvent{
			Name:     fmt.Sprint(i, ".p"),
			Category: fmt.Sprint(i),
			NS:       fmt.Sprint(i),
			Time:     time.Now(),
		}
		event.Write(&c)
		tmp = append(tmp, c)
		if len(tmp) > 100 {
			tmp = tmp[len(tmp)-100:]
		}
		assert.Equal(t, tmp, event.Read())
	}
}
