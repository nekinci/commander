package listener

import (
	"commander/src/job"
	"commander/src/specification"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Event struct {
	Type string
	Data any
}

type Func func(chan Event)

type fileChangeType string

func (l *Listener) initWatch() {
	fmt.Println("Watching files...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	l.fileWatcher = watcher
}

const (
	created fileChangeType = "create"
	deleted fileChangeType = "delete"
	saved   fileChangeType = "save"
	change  fileChangeType = "change"
)

type Listener struct {
	jobManager  *job.Manager
	listener    chan Event
	lDone       chan bool
	listenerFns []Func
	fileWatcher *fsnotify.Watcher
	includes    map[fileChangeType][]string
}

func New(spec *specification.Specification, jobManager *job.Manager) *Listener {
	l := &Listener{
		jobManager: jobManager,
		listener:   make(chan Event),
		lDone:      make(chan bool),
		includes: map[fileChangeType][]string{
			created: {},
			deleted: {},
			saved:   {},
			change:  {},
		},
	}
	l.initWatch()
	l.registerListeners(spec)
	return l
}

func (l *Listener) addListener(fn Func) *Listener {
	l.listenerFns = append(l.listenerFns, fn)
	return l
}

func getPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	wd, _ := os.Getwd()
	join := filepath.Join(wd, path)
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		join += string(os.PathSeparator)
	}
	return join
}

func (l *Listener) registerDirectoryRecursiveOrFile(path string) {
	path = getPath(path)
	startIsDir := false
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p == path {
				startIsDir = true
			}

			err := l.fileWatcher.Add(p)
			if err != nil {
				return err
			}
		}

		if !startIsDir {
			err := l.fileWatcher.Add(p)
			return err
		}

		return nil

	})
	if err != nil {
		log.Fatalln(err)
	}

}

func (l *Listener) registerListeners(specification *specification.Specification) {
	set := make(map[string]bool)

	for _, i := range specification.On.Create.Directories {
		p := getPath(i)
		l.includes[created] = append(l.includes[created], p)
		set[p] = true
	}

	for _, i := range specification.On.Delete.Files {
		p := getPath(i)
		l.includes[deleted] = append(l.includes[deleted], p)
		set[p] = true
	}

	for _, i := range specification.On.Save.Files {
		p := getPath(i)
		l.includes[saved] = append(l.includes[saved], p)
		set[p] = true
	}

	for _, i := range specification.On.Change.Files {
		p := getPath(i)
		l.includes[change] = append(l.includes[change], p)
		set[p] = true
	}

	for s := range set {
		err := l.fileWatcher.Add(s)
		if err != nil {
			log.Fatalln(err)
		}
	}

	l.addListener(func(c chan Event) {
		for {
			select {
			case event := <-l.fileWatcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					c <- Event{Type: string(created), Data: event}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					c <- Event{Type: string(deleted), Data: event}
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					c <- Event{Type: string(saved), Data: event}
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					c <- Event{Type: string(change), Data: event}
				}
			case err := <-l.fileWatcher.Errors:
				log.Println("error:", err)
			}
		}
	})

}

func (l *Listener) Listen() {

	defer func(fileWatcher *fsnotify.Watcher) {
		err := fileWatcher.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(l.fileWatcher)

	for _, fn := range l.listenerFns {
		go fn(l.listener)
	}

	for {
		select {
		case event := <-l.listener:
			fmt.Println(event)
			if event.Type == string(created) {
				for _, p := range l.includes[created] {
					if strings.HasPrefix(event.Data.(fsnotify.Event).Name, p) {
						l.jobManager.Run()
					}
				}
			}
			if event.Type == string(deleted) {
				for _, p := range l.includes[deleted] {
					if strings.HasPrefix(event.Data.(fsnotify.Event).Name, p) {
						l.jobManager.Run()
					}
				}
			}
			if event.Type == string(saved) {
				for _, p := range l.includes[saved] {
					if strings.HasPrefix(event.Data.(fsnotify.Event).Name, p) {
						l.jobManager.Run()
					}
				}
			}
			if event.Type == string(change) {
				for _, p := range l.includes[change] {
					if strings.HasPrefix(event.Data.(fsnotify.Event).Name, p) {
						l.jobManager.Run()
					}
				}
			}

		case <-l.lDone:
			fmt.Println("done received")
			return
		}
	}
}
