package git

import (
	"io/ioutil"
	"os"

	git2go "github.com/alleeclark/git2go"
	log "github.com/sirupsen/logrus"
)

type gitEvents struct {
	Topic    string
	Contents []byte
	Metadata map[string]string
}

type eventOptions struct {
	Deleted  bool
	Modified bool
	Added    bool
	Copied   bool
	Renamed  bool
}

type WithEventOptions func(*eventOptions) error

func WithIgornedfiles(path ...string) WithEventOptions {
	return func(e *eventOptions) error {
		return nil
	}
}

type EventsCollection struct {
	Event map[string]*gitEvents
}

//EventsFilterFunc as a template for implementing filter based functions
type EventsFilterFunc func(*EventsCollection) bool

//Filter function implements the function to be filtered if true
func (c *EventsCollection) Filter(fn EventsFilterFunc) *EventsCollection {
	if fn(c) {
		return c
	}
	return c
}

//ByCommitID filters by a given commit ID
func ByIngoredFiles(filepaths ...string) EventsFilterFunc {
	return func(e *EventsCollection) bool {
		for _, filename := range filepaths {
			delete(e.Event, filename)
		}
		return true
	}
}

func ByStatus(status string) EventsFilterFunc {
	return func(e *EventsCollection) bool {
		for path, event := range e.Event {
			if event.Topic == status {
				delete(e.Event, path)
			}
		}
		return true
	}
}

var defaultEventOptions = eventOptions{
	Deleted:  true,
	Modified: true,
	Copied:   true,
	Renamed:  true,
}

//ListFileChanges returns a map of files that have changed based on filtered commmits found along with the contents
func (c *Collection) ListDiffEvents(pullDir string, eventoptions ...WithEventOptions) map[string]*gitEvents {
	if len(c.Commits) < 1 {
		log.Infoln("No commits found to sync to contents")
		return nil
	}
	diffOptions, err := git2go.DefaultDiffOptions()
	if err != nil {
		log.Warningf("Error getting diff options %v", err)
		return nil
	}
	oldTree, err := c.Commits[0].Tree()
	if err != nil {
		return nil
	}
	newTree, err := c.Commits[len(c.Commits)-1].Tree()
	if err != nil {
		return nil
	}

	diff, err := c.DiffTreeToTree(oldTree, newTree, &diffOptions)
	if err != nil {
		log.Warningf("Error diffing tree %v", err)
		return nil
	}

	numOfDeltas, err := diff.NumDeltas()
	if err != nil {
		log.Warningf("Error getting num of deltas %v", err)
		return nil
	}
	events := make(map[string]*gitEvents, numOfDeltas)
	for d := 0; d < numOfDeltas; d++ {
		diffDelta, err := diff.GetDelta(d)
		if err != nil {
			log.Warningf("Error getting diff at %d %v", d, err)
		}
		contents, err := ioutil.ReadFile(pullDir + "/" + diffDelta.NewFile.Path)
		if err != nil || os.IsNotExist(err) {
			log.Warningf("Did not map contents %s becuase it does not exist %v", diffDelta.NewFile.Path, err)
			// make this better
			events[diffDelta.NewFile.Path] = nil
			continue
		}
		switch diffDelta.Status {
		case git2go.DeltaDeleted:
			events[diffDelta.NewFile.Path] = &gitEvents{Topic: "Deleted", Metadata: map[string]string{
				"old-path": diffDelta.OldFile.Path,
			}}
		case git2go.DeltaModified:
			events[diffDelta.NewFile.Path] = &gitEvents{Topic: "Modified", Contents: contents}
		case git2go.DeltaRenamed:
			events[diffDelta.NewFile.Path] = &gitEvents{Topic: "Renamed", Contents: contents, Metadata: map[string]string{
				"old-path": diffDelta.OldFile.Path,
			}}
		case git2go.DeltaAdded:
			events[diffDelta.NewFile.Path] = &gitEvents{Topic: "Added", Contents: contents}
		case git2go.DeltaCopied:
			events[diffDelta.NewFile.Path] = &gitEvents{Topic: "Copied", Contents: contents}
		}
	}
	return events
}
