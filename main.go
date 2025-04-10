package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/willmelton21/gator/internal/config"
	"github.com/willmelton21/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
  cmds.register("browse",middlewareLoggedIn(handlerBrowse))
  cmds.register("unfollow",middlewareLoggedIn(HandlerUnfollow))
  cmds.register("follow",middlewareLoggedIn(handlerFollow))
  cmds.register("following",middlewareLoggedIn(HandlerFollowing))
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
  cmds.register("reset",handlerDelete)
  cmds.register("users",GetUsers)
  cmds.register("agg", RunAggregator)
  cmds.register("addfeed",middlewareLoggedIn(AddFeed))
  cmds.register("feeds",PrintFeeds)
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

