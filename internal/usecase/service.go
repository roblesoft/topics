package usecase

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis/v7"
	entity "github.com/roblesoft/topics/internal/entity"
	repo "github.com/roblesoft/topics/internal/usecase/repo"
	token "github.com/roblesoft/topics/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

const (
	usersKey       = "users"
	userChannelFmt = "user:%s:channels"
	ChannelsKey    = "channels"
)

type Service struct {
	UserRepo         *repo.UserRepository
	channelsHandler  *redis.PubSub
	stopListenerChan chan struct{}
	listening        bool
	MessageChan      chan redis.Message
	username         string
}

func Connect(rdb *redis.Client, username string) (*Service, error) {
	if _, err := rdb.SAdd(usersKey, username).Result(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	service := &Service{
		username:         username,
		stopListenerChan: make(chan struct{}),
		MessageChan:      make(chan redis.Message),
	}

	if err := service.connect(rdb); err != nil {
		return nil, err
	}

	return service, nil
}

func NewService(UserRepo *repo.UserRepository) *Service {
	return &Service{
		UserRepo: UserRepo,
	}
}

func (s *Service) GetUser(username string) (*entity.User, error) {
	return s.UserRepo.Get(username)
}

func (s *Service) CreateUser(b *entity.User) error {
	return s.UserRepo.Create(b)
}

func (s *Service) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) LoginCheck(username string, password string) (string, error) {
	user, err := s.GetUser(username)

	if err != nil {
		return "", err
	}

	err = s.verifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("verify")
		fmt.Println(err)
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		fmt.Println("generate")
		fmt.Println(err)
		return "", err
	}

	return token, nil

}

func (service *Service) Subscribe(rdb *redis.Client, channel string) error {

	userChannelsKey := fmt.Sprintf(userChannelFmt, service.username)

	if rdb.SIsMember(userChannelsKey, channel).Val() {
		return nil
	}
	if err := rdb.SAdd(userChannelsKey, channel).Err(); err != nil {
		return err
	}

	return service.connect(rdb)
}

func (service *Service) Unsubscribe(rdb *redis.Client, channel string) error {

	userChannelsKey := fmt.Sprintf(userChannelFmt, service.username)

	if !rdb.SIsMember(userChannelsKey, channel).Val() {
		return nil
	}
	if err := rdb.SRem(userChannelsKey, channel).Err(); err != nil {
		return err
	}

	return service.connect(rdb)
}

func (service *Service) connect(rdb *redis.Client) error {

	var c []string

	c1, err := rdb.SMembers(ChannelsKey).Result()
	if err != nil {
		return err
	}
	c = append(c, c1...)

	// get all user channels (from DB) and start subscribe
	c2, err := rdb.SMembers(fmt.Sprintf(userChannelFmt, service.username)).Result()
	if err != nil {
		return err
	}
	c = append(c, c2...)

	if len(c) == 0 {
		fmt.Println("no channels to connect to for user: ", service.username)
		return nil
	}

	if service.channelsHandler != nil {
		if err := service.channelsHandler.Unsubscribe(); err != nil {
			return err
		}
		if err := service.channelsHandler.Close(); err != nil {
			return err
		}
	}
	if service.listening {
		service.stopListenerChan <- struct{}{}
	}

	return service.doConnect(rdb, c...)
}

func (service *Service) doConnect(rdb *redis.Client, channels ...string) error {
	// subscribe all channels in one request
	pubSub := rdb.Subscribe(channels...)
	// keep channel handler to be used in unsubscribe
	service.channelsHandler = pubSub

	// The Listener
	go func() {
		service.listening = true
		fmt.Println("starting the listener for user:", service.username, "on channels:", channels)
		for {
			select {
			case msg, ok := <-pubSub.Channel():
				if !ok {
					return
				}
				service.MessageChan <- *msg

			case <-service.stopListenerChan:
				fmt.Println("stopping the listener for user:", service.username)
				return
			}
		}
	}()
	return nil
}

func (service *Service) Disconnect() error {
	if service.channelsHandler != nil {
		if err := service.channelsHandler.Unsubscribe(); err != nil {
			return err
		}
		if err := service.channelsHandler.Close(); err != nil {
			return err
		}
	}
	if service.listening {
		service.stopListenerChan <- struct{}{}
	}

	close(service.MessageChan)

	return nil
}

func Chat(rdb *redis.Client, channel string, content string) error {
	return rdb.Publish(channel, content).Err()
}

func List(rdb *redis.Client) ([]string, error) {
	return rdb.SMembers(usersKey).Result()
}

func GetChannels(rdb *redis.Client, username string) ([]string, error) {

	if !rdb.SIsMember(usersKey, username).Val() {
		return nil, errors.New("user not exists")
	}

	var c []string

	c1, err := rdb.SMembers(ChannelsKey).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c1...)

	// get all user channels (from DB) and start subscribe
	c2, err := rdb.SMembers(fmt.Sprintf(userChannelFmt, username)).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c2...)

	return c, nil
}
