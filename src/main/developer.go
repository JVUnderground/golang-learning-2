package main

type Developer struct {
	Name      string `json:"name"`
	Followers int64  `json:"followers"`
	Stars     int64  `json:"stars"`
	Commits   int64  `json:"commits"`
	N_repos   int64  `json:"n_repos"`
	Avatar    string `json:"avatar"`
	Price     string `json:"price"`
}

type Developers []Developer
