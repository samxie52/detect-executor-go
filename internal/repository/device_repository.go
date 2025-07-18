package repository

import (
	"context"
	"detect-executor-go/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	BaseRepository

	Create(ctx context.Context, device *model.Device) error
	GetByID(ctx context.Context, id uint) (*model.Device, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*model.Device, error)
	Update(ctx context.Context, device *model.Device) error
	Delete(ctx context.Context, id uint) error

	List(ctx context.Context, page int, size int) ([]model.Device, error)
	GetOnlineDevices(ctx context.Context) ([]model.Device, error)
	GetByStatus(ctx context.Context, status string) ([]model.Device, error)

	UpdateLastSeen(ctx context.Context, deviceID string) error
	UpdateStatus(ctx context.Context, deviceID string, status string) error
}

type deviceRepository struct {
	BaseRepository
}

func NewDeviceRepository(db *gorm.DB, redis *redis.Client) DeviceRepository {
	return &deviceRepository{
		BaseRepository: NewBaseRepository(db, redis),
	}
}

// Implement DeviceRepository methods
func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.GetDb().WithContext(ctx).Create(device).Error
}

func (r *deviceRepository) GetByID(ctx context.Context, id uint) (*model.Device, error) {
	var device model.Device
	err := r.GetDb().WithContext(ctx).First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) GetByDeviceID(ctx context.Context, deviceID string) (*model.Device, error) {
	var device model.Device
	err := r.GetDb().WithContext(ctx).Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) Update(ctx context.Context, device *model.Device) error {
	return r.GetDb().WithContext(ctx).Save(device).Error
}

func (r *deviceRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDb().WithContext(ctx).Delete(&model.Device{}, id).Error
}

func (r *deviceRepository) List(ctx context.Context, page int, size int) ([]model.Device, error) {
	var devices []model.Device
	offset := (page - 1) * size
	err := r.GetDb().WithContext(ctx).Offset(offset).Limit(size).Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) GetOnlineDevices(ctx context.Context) ([]model.Device, error) {
	var devices []model.Device
	err := r.GetDb().WithContext(ctx).Where("status = ?", "online").Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) GetByStatus(ctx context.Context, status string) ([]model.Device, error) {
	var devices []model.Device
	err := r.GetDb().WithContext(ctx).Where("status = ?", status).Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) UpdateLastSeen(ctx context.Context, deviceID string) error {
	return r.GetDb().WithContext(ctx).Model(&model.Device{}).Where("device_id = ?", deviceID).Update("last_seen", "NOW()").Error
}

func (r *deviceRepository) UpdateStatus(ctx context.Context, deviceID string, status string) error {
	return r.GetDb().WithContext(ctx).Model(&model.Device{}).Where("device_id = ?", deviceID).Update("status", status).Error
}
