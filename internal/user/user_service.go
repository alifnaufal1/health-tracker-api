package user

import (
	"context"
	"errors"
	"fmt"
	"health-tracker-api/pkg/apperror"
	"health-tracker-api/pkg/logging"
	"health-tracker-api/pkg/transaction"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, request UserCreateRequest) (*UserResponse, error)
	VerifyCredentials(ctx context.Context, username, password string) (*UserResponse, error)
	Update(ctx context.Context, request UserUpdateRequest, userID uuid.UUID) (*UserResponse, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}

type UserServiceImpl struct {
	repo   UserRepository
	tx     transaction.Manager
	hasher PasswordHasher
	log    logging.Logger
}

func NewUserService(repo UserRepository, tx transaction.Manager, hasher PasswordHasher, log logging.Logger) UserService {
	return &UserServiceImpl{
		repo:   repo,
		tx:     tx,
		hasher: hasher,
		log:    log.With("layer", "service").With("module", "user"),
	}
}

func errInvalidCredentials() error {
	return &apperror.UnauthorizedError{Message: "invalid username or password"}
}

func parseUserID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, &apperror.ValidationError{Message: "invalid user_id"}
	}
	return id, nil
}

func (s *UserServiceImpl) Create(ctx context.Context, request UserCreateRequest) (*UserResponse, error) {
	log := s.log.WithContext(ctx).With("username", request.Username)
	log.Debug("received create user request")

	_, err := s.repo.FindByUsername(ctx, request.Username)
	if err == nil {
		log.Warn("request denied: this user already registered")
		return nil, &apperror.ConflictError{Message: "this user already registered"}
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("check registered user: %w", err)
	}

	hash, err := s.hasher.Hash(request.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &User{ID: uuid.New(), Username: request.Username, Password: hash, Name: request.Name}

	if err := s.repo.Create(ctx, u); err != nil {
		if errors.Is(err, ErrDuplicateUsername) {
			log.Warn("request denied: this user already registered (concurrent request)")
			return nil, &apperror.ConflictError{Message: "this user already registered"}
		}
		return nil, fmt.Errorf("save user: %w", err)
	}

	log.With("user_id", u.ID).Debug("user created successfully")
	return ToUserResponse(u), nil
}

func (s *UserServiceImpl) VerifyCredentials(ctx context.Context, username, password string) (*UserResponse, error) {
	log := s.log.WithContext(ctx).With("username", username)
	log.Debug("received verify credentials request")

	u, err := s.repo.FindByUsername(ctx, username)
	if errors.Is(err, ErrNotFound) {
		log.Warn("credentials rejected: user not found")
		return nil, errInvalidCredentials()
	}
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}

	if !s.hasher.Compare(password, u.Password) {
		log.Warn("credentials rejected: invalid password")
		return nil, errInvalidCredentials()
	}

	return ToUserResponse(u), nil
}

func (s *UserServiceImpl) Update(ctx context.Context, request UserUpdateRequest, userID uuid.UUID) (*UserResponse, error) {
	log := s.log.WithContext(ctx).With("user_id", userID)
	log.Debug("received update user request")

	var newHash string
	var err error
	if request.Password != "" {
		if newHash, err = s.hasher.Hash(request.Password); err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
	}

	var updated *User
	err = s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		u, err := s.repo.FindByID(ctx, userID)
		if errors.Is(err, ErrNotFound) {
			log.Warn("not found: this user is not found")
			return &apperror.NotFoundError{Message: "user not found"}
		}
		if err != nil {
			return fmt.Errorf("find user by id: %w", err)
		}

		if newHash != "" {
			u.Password = newHash
		}
		if request.Name != "" {
			u.Name = request.Name
		}

		if err := s.repo.Update(ctx, u); err != nil {
			return fmt.Errorf("update user %s: %w", userID, err)
		}
		updated = u
		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Debug("user updated successfully")
	return ToUserResponse(updated), nil
}

func (s *UserServiceImpl) FindByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	log := s.log.WithContext(ctx).With("user_id", userID)
	log.Debug("received find user by id request")

	u, err := s.repo.FindByID(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		log.Warn("not found: this user is not found")
		return nil, &apperror.NotFoundError{Message: "user not found"}
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id %s: %w", userID, err)
	}

	log.Debug("user found successfully")
	return ToUserResponse(u), nil
}
