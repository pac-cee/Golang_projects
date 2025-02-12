package repository

import (
    "context"
    "wallet-app/internal/domain/models"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type WalletRepository struct {
    db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
    return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
    return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
    var wallet models.Wallet
    if err := r.db.WithContext(ctx).First(&wallet, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var wallet models.Wallet
        if err := tx.WithContext(ctx).First(&wallet, "id = ?", id).Error; err != nil {
            return err
        }

        wallet.Balance += amount
        if wallet.Balance < 0 {
            return ErrInsufficientFunds
        }

        return tx.Save(&wallet).Error
    })
}