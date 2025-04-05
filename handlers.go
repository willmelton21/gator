package main

import (
	"context"
  "log"
  "database/sql"
  "strings"
	"fmt"
  "strconv"
	"time"
	"github.com/willmelton21/gator/internal/database"
	"github.com/google/uuid"
  rss "github.com/willmelton21/gator/utils"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}

func scrapeFeeds(s *state) error {
  currFeed, err := s.db.GetNextFeedToFetch(context.Background())
  if err != nil {
     return fmt.Errorf("failed to get next feed row: %v",err)
  }
  
  err = s.db.MarkFeedFetched(context.Background(),currFeed.ID)
  if err != nil {
     return fmt.Errorf("failed to mark feed as fetched %v",err)
  }

  feed, err := rss.FetchFeed(context.Background(), currFeed.Url) 
  if err != nil {
    return fmt.Errorf("error fetching feed for current url %v",err)
  }
  
  for _, item := range feed.Channel.Item {
    publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    currFeed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", currFeed.Name, len(feed.Channel.Item))
  return nil
}
 

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

  return func(s *state, cmd command) error {
  user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
  if err != nil {
    return err
  }
  return handler(s,cmd,user) 
  }
}

func HandlerUnfollow(s *state, cmd command, user database.User) error {

  if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url for unfollow>", cmd.Name)
	}

  err := s.db.Unfollow(context.Background(),database.UnfollowParams{
   
    UserID: user.ID,
    Url:    cmd.Args[0],
    })

  if err != nil {
    
      return fmt.Errorf("failed to unfollow feeed: %v", err)
  }
  return nil
  }

func HandlerFollowing(s *state, cmd command, user database.User) error {


  followingList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
  if err != nil {
    return fmt.Errorf("couldn't get list of user followers %s",err)
  }

  if len(followingList) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range followingList{
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
  }

func AddFeed(s *state, cmd command, user database.User) error {
 	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
  fmt.Printf("Username: %s\n FeedName: %s\n",feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

  	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

  feed,err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
    if err != nil {
      return fmt.Errorf("couldn't get feed for url given %w",err)
  }


	feedFollows, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	  UserID:      user.ID,
    FeedID:      feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follows: %w", err)
	}

  fmt.Println("Feed Follows created")
  fmt.Printf("Username: %s\n FeedName: %s\n",feedFollows.UserName,feedFollows.FeedName)
  return nil
}

func PrintFeeds(s *state, cmd command) error {

  feeds, err := s.db.GetFeeds(context.Background())
  if err != nil {
      return fmt.Errorf("couldn't print feeds: %w",err)
  }
  for i := 0; i < len(feeds); i++ {
    fmt.Println(" Name:",feeds[i].Name)
    fmt.Println(" URL:",feeds[i].Url)
    fmt.Println(" User:",feeds[i].UserName)
    }
  
    return err
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)

}

func RunAggregator(s *state, cmd command) error {
  
    if len(cmd.Args) != 1 {
      return fmt.Errorf("time argument not provided for %v",cmd.Name) 
   }
    t, err := time.ParseDuration(cmd.Args[0])
    if err != nil {
     return fmt.Errorf("Proper time value was not entered %v",err)
       }

    fmt.Println("Collecting feeds every ",t)

    ticker := time.NewTicker(t)
    for ;; <-ticker.C {
      scrapeFeeds(s) 
     }
  
   return nil
}


func handlerDelete(s *state, cmd command) error {


  err := s.db.DeleteUsers(context.Background())
  if err != nil {
    return fmt.Errorf("couldn't delete useres %w",err)
  }
  return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func GetUsers(s *state, cmd command) error {
    
  users, err := s.db.GetUsers(context.Background())
  if err != nil {
      return fmt.Errorf("couldn't get current users %w",err)
    }
  for i := 0; i < len(users); i++ {

    if users[i].Name == s.cfg.CurrentUserName {
        fmt.Printf("%s (current)\n",users[i].Name)
    } else {
    fmt.Println(users[i].Name)
    }
  }

  return nil
}
func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

