package services

import (
    "context"
    "wallet-app/internal/domain/models"
    "wallet-app/internal/domain/repository"

    "github.com/google/uuid"
)

type WalletService struct {
    repo *repository.WalletRepository
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
    return &WalletService{repo: repo}
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID, currency string) (*models.Wallet, error) {
    wallet := &models.Wallet{
        ID:       uuid.New(),
        UserID:   userID,
        Currency: currency,
        Balance:  0,
    }

    if err := s.repo.Create(ctx, wallet); err != nil {
        return nil, err
    }

    return wallet, nil
}

func (s *WalletService) ProcessTransaction(ctx context.Context, walletID uuid.UUID, amount float64, txType models.TransactionType) error {
    var finalAmount float64
    if txType == models.Withdraw {
        finalAmount = -amount
    } else {
        finalAmount = amount
    }

    return s.repo.UpdateBalance(ctx, walletID, finalAmount)
}