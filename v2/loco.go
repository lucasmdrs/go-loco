package loco

import (
	"context"
	"errors"
	"log"

	"github.com/lucasmdrs/ctxpoller"
)

const baseURL = "https://localise.biz/api/"
const authEndpoint = "auth/verify"
const filenameTemplate = "%s.json"
const endpointTemplate = "export/" + filenameTemplate
const authParameter = "?key=%s"

// goloco is a structure that holds the required information and implements the GoLoco service interface
type goloco struct {
	notifier chan interface{}
	poller   ctxpoller.Poller
	projects map[uint]*project
}

// GoLoco defines the service interface
type GoLoco interface {
	AddProject(key, assetsPath string) error
	StartPoller() (chan interface{}, error)
	StopPoller()
	FetchTranslations(context.Context)
}

// Init initializes the service
func Init() GoLoco {
	return &goloco{
		notifier: make(chan interface{}),
		projects: make(map[uint]*project, 0),
	}
}

// AddProject includes a new Loco project from a API Key
// and sets the assets destination.
func (g *goloco) AddProject(key string, assetsDestinationPath string) error {
	p, err := getProjectInformation(key)
	if err != nil {
		return err
	}
	if _, exists := g.projects[p.ID]; exists {
		return errors.New("Duplicated project")
	}
	p.assetsPath = assetsDestinationPath
	g.projects[p.ID] = &p
	g.notifier = make(chan interface{}, len(g.projects))
	return nil
}

// StartPoller keep a poller for any changes in the projects translations
// notifying it in the returned channel whenever the assets changed
func (g *goloco) StartPoller() (chan interface{}, error) {
	g.poller = ctxpoller.DefaultPoller(g.FetchTranslations)
	return g.notifier, g.poller.Start()
}

// StopPoller stops looking for changes in the translations
func (g *goloco) StopPoller() {
	g.poller.Stop()
}

// FetchTranslations fetches the translations from Loco and save it assets
func (g *goloco) FetchTranslations(context.Context) {
	for _, p := range g.projects {
		if err := p.fetchProjectTranslations(g.notifier); err != nil {
			log.Fatalf("Failed to retrieve translations for %s: %s", p.Name, err.Error())
		}
	}
	if g.poller == nil || !g.poller.IsActive() {
		emptyChannel(g.notifier)
	}
}

func emptyChannel(ch chan interface{}) {
	for len(ch) > 0 {
		<-ch
	}
}
