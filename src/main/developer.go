package main

type Developer struct {
	Name      string `json:"name"`
	Followers int    `json:"followers"`
	Stars     int    `json:"stars"`
	Commits   int    `json:"commits"`
	N_repos   int    `json:"n_repos"`
}

type Developers []Developer
