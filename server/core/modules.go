// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"sync"

	"github.com/maxlandon/wiregost/data_service/remote"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	// Modules - Map of all modules available in Wiregost (map["path/to/module"] = Module)
	Modules = &modules{
		Loaded: &map[string]Module{},
		mutex:  &sync.RWMutex{},
	}
)

type modules struct {
	Loaded *map[string]Module
	mutex  *sync.RWMutex
}

// Module - Represents a module, providing access to its methods
// All Wiregost modules must implement this interface
type Module interface {
	Init(int32) error
	Run(string) (string, error)
	SetOption(string, string)
	ToProtobuf() *pb.Module
}

// GetModule - Get module by path, (load it if needed)
func GetModule(path string) (Module, error) {

	Modules.mutex.Lock()
	defer Modules.mutex.Unlock()

	mod, ok := (*Modules.Loaded)[path]
	if !ok {
		return nil, errors.New("No module for given path")
	}
	return mod, nil
}

// AddModule - Add a module to Wiregost list of available modules
func AddModule(path string, mod Module) error {

	Modules.mutex.Lock()
	defer Modules.mutex.Unlock()

	(*Modules.Loaded)[path] = mod
	return nil
}

var (
	// Stacks - All module stacks (one per workspace, per user),
	Stacks = &map[uint]map[string]*stack{}

	storageLog = log.ServerLogger("core", "stack")
)

type stack struct {
	Loaded *map[string]Module
	mutex  *sync.RWMutex
}

// LoadModule - Load a module onto the stack, by fetching it into Modules
func (s *stack) LoadModule(userID int32, path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	mod, err := GetModule(path)
	if err != nil {
		return err
	}

	// Init and load onto stack
	mod.Init(userID)
	(*s.Loaded)[path] = mod

	return nil
}

// PopModule - Unload a module from stack
func (s *stack) PopModule(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(*s.Loaded, path)

	return nil
}

// Module - Get a module by path, (load it onto the stack if needed)
func (s *stack) Module(userID int32, path string) (Module, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	mod, ok := (*s.Loaded)[path]
	if !ok {
		s.LoadModule(userID, path)
		return (*s.Loaded)[path], nil
	}
	return mod, nil

}

// LoadStacks - Inits an empty stack for each user, in each workspace (N users x N workspaces = P stacks)
func LoadStacks() error {
	// Users
	clientCerts := certs.UserClientListCertificates()
	users := []string{}
	for _, c := range clientCerts {
		users = append(users, c.Subject.CommonName)
	}
	users = unique(users)

	// Workspaces
	workspaces, _ := remote.Workspaces(nil)

	// Init stack maps, to modify them later if needed
	for _, w := range workspaces {
		workspaceStacks := &map[string]*stack{}
		(*Stacks)[w.ID] = (*workspaceStacks)
		for _, user := range users {
			(*workspaceStacks)[user] = &stack{}
			(*workspaceStacks)[user].Loaded = &map[string]Module{}
			(*workspaceStacks)[user].mutex = &sync.RWMutex{}
		}
	}

	return nil
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
