/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registrytest

import (
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
)

type MinionRegistry struct {
	Err     error
	Minion  string
	Minions api.NodeList
	sync.Mutex
}

func MakeMinionList(minions []string, nodeResources api.NodeResources) *api.NodeList {
	list := api.NodeList{
		Items: make([]api.Node, len(minions)),
	}
	for i := range minions {
		list.Items[i].Name = minions[i]
		list.Items[i].Spec.Capacity = nodeResources.Capacity
	}
	return &list
}

func NewMinionRegistry(minions []string, nodeResources api.NodeResources) *MinionRegistry {
	return &MinionRegistry{
		Minions: *MakeMinionList(minions, nodeResources),
	}
}

func (r *MinionRegistry) ListMinions(ctx api.Context) (*api.NodeList, error) {
	r.Lock()
	defer r.Unlock()
	return &r.Minions, r.Err
}

func (r *MinionRegistry) CreateMinion(ctx api.Context, minion *api.Node) error {
	r.Lock()
	defer r.Unlock()
	r.Minion = minion.Name
	r.Minions.Items = append(r.Minions.Items, *minion)
	return r.Err
}

func (r *MinionRegistry) UpdateMinion(ctx api.Context, minion *api.Node) error {
	r.Lock()
	defer r.Unlock()
	for i, node := range r.Minions.Items {
		if node.Name == minion.Name {
			r.Minions.Items[i] = *minion
			return r.Err
		}
	}
	return r.Err
}

func (r *MinionRegistry) GetMinion(ctx api.Context, minionID string) (*api.Node, error) {
	r.Lock()
	defer r.Unlock()
	for _, node := range r.Minions.Items {
		if node.Name == minionID {
			return &node, r.Err
		}
	}
	return nil, r.Err
}

func (r *MinionRegistry) DeleteMinion(ctx api.Context, minionID string) error {
	r.Lock()
	defer r.Unlock()
	var newList []api.Node
	for _, node := range r.Minions.Items {
		if node.Name != minionID {
			newList = append(newList, api.Node{ObjectMeta: api.ObjectMeta{Name: node.Name}})
		}
	}
	r.Minions.Items = newList
	return r.Err
}
