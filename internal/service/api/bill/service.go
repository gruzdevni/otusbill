package bill

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"otusbill/internal/models"
	query "otusbill/internal/repo"
)

var (
	ErrNotCorrectData   = errors.New("Not correct password or email")
	ErrEmailAlreadyUsed = errors.New("Email is already registered. Please login")
	ErrSuchUser         = errors.New("No such user")
)

type repo interface {
	GetUserBalance(ctx context.Context, userGuid uuid.UUID) (query.GetUserBalanceRow, error)
	ChangeBalance(ctx context.Context, arg query.ChangeBalanceParams) error
	InsertUser(ctx context.Context, guid uuid.UUID) error
}

type service struct {
	repo repo
}

type Service interface {
	IncreaseBalance(ctx context.Context, param models.BalanceIncreaseParams) error
	ReduceBalance(ctx context.Context, param models.BalanceReduceParams) error
	GetUserBalance(ctx context.Context, guid uuid.UUID) (models.UserBalance, error)
	InsertUser(ctx context.Context, guid uuid.UUID) error
}

func NewService(repo repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) InsertUser(ctx context.Context, guid uuid.UUID) error {
	err := s.repo.InsertUser(ctx, guid)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

func (s *service) IncreaseBalance(ctx context.Context, param models.BalanceIncreaseParams) error {
	userGuid, err := uuid.Parse(param.UserGUID.String())
	if err != nil {
		return fmt.Errorf("parsing user guid: %w", err)
	}

	queryParam := query.ChangeBalanceParams{
		UserGuid:     userGuid,
		OperationRef: param.OperationRef,
		Amount:       decimal.NewFromFloat(param.Amount),
	}

	err = s.repo.ChangeBalance(ctx, queryParam)
	if err != nil {
		return fmt.Errorf("increasing user balance: %w", err)
	}

	return nil
}

func (s *service) ReduceBalance(ctx context.Context, param models.BalanceReduceParams) error {
	userGuid, err := uuid.Parse(param.UserGUID.String())
	if err != nil {
		return fmt.Errorf("parsing user guid: %w", err)
	}

	queryParam := query.ChangeBalanceParams{
		UserGuid:     userGuid,
		OperationRef: param.OperationRef,
		Amount:       decimal.NewFromFloat(param.Amount).Neg(),
	}

	err = s.repo.ChangeBalance(ctx, queryParam)
	if err != nil {
		return fmt.Errorf("reducing user balance: %w", err)
	}

	return nil
}

func (s *service) GetUserBalance(ctx context.Context, guid uuid.UUID) (models.UserBalance, error) {
	repoRes, err := s.repo.GetUserBalance(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserBalance{Amount: 0.00, GUID: strfmt.UUID(guid.String())}, nil
		}

		return models.UserBalance{}, fmt.Errorf("getting user balance: %w", err)
	}

	return models.UserBalance{Amount: repoRes.Column2, GUID: strfmt.UUID(guid.String())}, nil
}
