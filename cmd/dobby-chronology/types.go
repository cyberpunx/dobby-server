package main

import "localdev/dobby-server/internal/pkg/hogwartsforum/tool"

type session struct {
	Tool *tool.Tool
	Conf config
}
