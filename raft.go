package main

import (
	"sync"
	"time"
)

type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)

type Raft struct {
	ID       string
	State    RaftState
	Term     int
	LeaderID string
	mu       sync.Mutex
}

func NewRaft(id string) *Raft {
	return &Raft{
		ID:    id,
		State: Follower,
		Term:  0,
	}
}

func (r *Raft) Run() {
	for {
		r.mu.Lock()
		state := r.State
		r.mu.Unlock()

		switch state {
		case Follower:
			time.Sleep(1 * time.Second) // Simplified: Wait for heartbeat
		case Candidate:
			r.startElection()
		case Leader:
			r.sendHeartbeats()
		}
	}
}

func (r *Raft) IsLeader() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.State == Leader
}

func (r *Raft) startElection() {
	// Simplified: Assume election succeeds
	r.mu.Lock()
	r.State = Leader
	r.LeaderID = r.ID
	r.mu.Unlock()
}

func (r *Raft) sendHeartbeats() {
	// Simplified: Send heartbeats to peers
	time.Sleep(100 * time.Millisecond)
}
