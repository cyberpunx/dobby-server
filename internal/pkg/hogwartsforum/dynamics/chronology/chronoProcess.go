package chronology

import "localdev/dobby-server/internal/pkg/hogwartsforum/parser"

type ChronoReport struct {
	ChronoThreads []*ChronoThread
}

type ChronoThread struct {
	Title       string
	ChronoPosts []*ChronoPost
}

type ChronoPost struct {
	Author  string
	Message string
}

func ChronoThreadProcessor(Thread parser.Thread) *ChronoThread {

	var chronoPosts []*ChronoPost
	for _, post := range Thread.Posts {
		chronoPost := ChronoPost{
			Author:  post.Author.Username,
			Message: post.Content,
		}
		chronoPosts = append(chronoPosts, &chronoPost)
	}
	chronoThread := ChronoThread{
		Title:       Thread.Title,
		ChronoPosts: chronoPosts,
	}
	return &chronoThread
}
