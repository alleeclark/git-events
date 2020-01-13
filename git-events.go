package main

import (
	"git-events/git"
	"time"

	"os"

	"fmt"
	"github.com/jasonlvhit/gocron"
)

func main() {
	app := New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "git-events:  %s\n", err)
		os.Exit(1)
	}
}

type server struct {
	interval uint64
	lasttime time.Time
	name     string
	branch   string
}

func (s *server) Event(in *git.EventRequest, srv git.EventsService_EventServer) error {
	gocron.Every(s.interval).Seconds().Do(s.getEvent, in, srv)
	<-gocron.Start()
	return nil
}

func (s *server) getEvent(in *git.EventRequest, srv git.EventsService_EventServer) {
	fmt.Fprintf(os.Stdout, "The last time ran was %s", s.lasttime.String())
	if s.lasttime.Second() > time.Now().Second() {
		s.lasttime = time.Now()
	} else {
		repo := git.Open(s.name)
		repo = repo.Fetch(git.CloneOptions("", nil), s.branch).Filter(git.ByDate(s.lasttime))
		s.lasttime = time.Now()
		gitEvents := repo.ListDiffEvents(s.name, nil)
		for filename, event := range gitEvents {
			for _, topic := range in.Topics {
				if topic == event.Topic {
					resp := &git.EventResponse{
						Filename: filename,
						Topic:    event.Topic,
						Contents: event.Contents,
						Metadata: event.Metadata,
					}
					if err := srv.Send(resp); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to send response: %v", err)
					}
				}
			}
		}
	}
}
