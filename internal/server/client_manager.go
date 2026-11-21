package server

import "sync"

type ClientManager struct {
    mu      sync.RWMutex
    clients map[string]*Client
}

func NewClientManager() *ClientManager {
    return &ClientManager{
        clients: make(map[string]*Client),
    }
}

func (m *ClientManager) Add(client *Client) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.clients[client.ID] = client
}

func (m *ClientManager) Remove(id string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    delete(m.clients, id)
}