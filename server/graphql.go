package server

import (
	// 	context "context"
	// 	"encoding/json"
	// 	"fmt"
	// 	"log"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	//
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/segmentio/ksuid"
	"github.com/tinrab/retry"
	"github.com/vektah/gqlgen/handler"
)

// type contextKey string
//
// const (
// 	userContextKey = contextKey("user")
// )
//

type graphQLServer struct {
	redisClient     *redis.Client
	messageChannels map[string]chan Message
	userChannels    map[string]chan string
	mutex           sync.Mutex
}

// NewGraphQLServer returns new graphQLServer with redisURL.
func NewGraphQLServer(redisURL string) (*graphQLServer, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		_, err := client.Ping().Result()
		return err
	})

	return &graphQLServer{
		redisClient:     client,
		messageChannels: map[string]chan Message{},
		userChannels:    map[string]chan string{},
		mutex:           sync.Mutex{},
	}, nil
}

func (s *graphQLServer) Serve(route string, port int) error {
	mux := http.NewServeMux()
	mux.Handle(
		route,
		handler.GraphQL(MakeExecutableSchema(s),
			handler.WebsocketUpgrader(
				websocket.Upgrader{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			),
		),
	)

	mux.Handle("/playground", handler.Playground("GraphQL", route))

	handler := cors.AllowAll().Handler(mux)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func (s *graphQLServer) createUser(user string) error {
	// Upsert user
	if err := s.redisClient.SAdd("users", user).Err(); err != nil {
		return err
	}

	// Notify new user joined
	s.mutex.Lock()
	for _, ch := range s.userChannels {
		ch <- user
	}
	s.mutex.Unlock()
	return nil
}

func (s *graphQLServer) Mutation_postMessage(ctx context.Context, user string, text string) (*Message, error) {
	if err := s.createUser(user); err != nil {
		return nil, err
	}

	// Create message
	m := Message{
		ID:        ksuid.New().String(),
		User:      user,
		CreatedAt: time.Now().UTC(),
		Text:      text,
	}
	mj, _ := json.Marshal(m)
	if err := s.redisClient.LPush("message", mj).Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	// Notify new massage
	s.mutex.Lock()
	for _, ch := range s.messageChannels {
		ch <- m
	}
	s.mutex.Unlock()
	return &m, nil
}

func (s *graphQLServer) Query_messages(ctx context.Context) ([]Message, error) {
	cmd := s.redisClient.LRange("message", 0, -1)
	if cmd.Err() != nil {
		log.Println(cmd.Err())
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	messages := []Message{}
	for _, mj := range res {
		var m Message
		err = json.Unmarshal([]byte(mj), &m)
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *graphQLServer) Query_users(ctx context.Context) ([]string, error) {
	cmd := s.redisClient.SMembers("users")
	if cmd.Err() != nil {
		log.Println(cmd.Err())
		return nil, cmd.Err()
	}

	res, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}

/////////////////////////////////////////////////////////////////////////////////////
// TODO: implement methods for Resolvers
// type Resolvers interface {
// 	Mutation_postMessage(ctx context.Context, user string, text string) (*Message, error)
// 	Query_messages(ctx context.Context) ([]Message, error)
// 	Query_users(ctx context.Context) ([]string, error)
//
// 	Subscription_messagePosted(ctx context.Context, user string) (<-chan Message, error)
// 	Subscription_userJoined(ctx context.Context, user string) (<-chan string, error)
// }
//
// func (s *graphQLServer) Subscription_messagePosted(ctx context.Context, user string) (<-chan Message, error) {
// 	err := s.createUser(user)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Create new channel for request
// 	messages := make(chan Message, 1)
// 	s.mutex.Lock()
// 	s.messageChannels[user] = messages
// 	s.mutex.Unlock()
//
// 	// Delete channel when done
// 	go func() {
// 		<-ctx.Done()
// 		s.mutex.Lock()
// 		delete(s.messageChannels, user)
// 		s.mutex.Unlock()
// 	}()
//
// 	return messages, nil
// }
//
// func (s *graphQLServer) Subscription_userJoined(ctx context.Context, user string) (<-chan string, error) {
// 	err := s.createUser(user)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Create new channel for request
// 	users := make(chan string, 1)
// 	s.mutex.Lock()
// 	s.userChannels[user] = users
// 	s.mutex.Unlock()
//
// 	// Delete channel when done
// 	go func() {
// 		<-ctx.Done()
// 		s.mutex.Lock()
// 		delete(s.userChannels, user)
// 		s.mutex.Unlock()
// 	}()
//
// 	return users, nil
// }
//
