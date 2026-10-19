package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
)

const (
	driverLocationKeyPrefix = "driver:loc:"
	driverLocationTTL       = 60 * time.Second
)

var ErrLocationNotFound = errors.New("driver location not found")

type ILocationRepository interface {
	SetDriverLocation(ctx context.Context, driverID uuid.UUID, loc *domain.DriverLocation) error
	GetDriverLocation(ctx context.Context, driverID uuid.UUID) (*domain.DriverLocation, error)
}

type LocationRepository struct {
	Client *redis.Client
}

func driverLocationKey(driverID uuid.UUID) string {
	return driverLocationKeyPrefix + driverID.String()
}

func (r *LocationRepository) SetDriverLocation(ctx context.Context, driverID uuid.UUID, loc *domain.DriverLocation) error {
	key := driverLocationKey(driverID)
	fields := map[string]interface{}{
		"lat":        strconv.FormatFloat(loc.Lat, 'f', -1, 64),
		"lng":        strconv.FormatFloat(loc.Lng, 'f', -1, 64),
		"updated_at": loc.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if loc.Heading != nil {
		fields["heading"] = strconv.FormatFloat(*loc.Heading, 'f', -1, 64)
	}
	if loc.Speed != nil {
		fields["speed"] = strconv.FormatFloat(*loc.Speed, 'f', -1, 64)
	}

	pipe := r.Client.Pipeline()
	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, driverLocationTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *LocationRepository) GetDriverLocation(ctx context.Context, driverID uuid.UUID) (*domain.DriverLocation, error) {
	values, err := r.Client.HGetAll(ctx, driverLocationKey(driverID)).Result()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, ErrLocationNotFound
	}

	loc, err := parseDriverLocation(values)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func parseDriverLocation(values map[string]string) (*domain.DriverLocation, error) {
	lat, err := strconv.ParseFloat(values["lat"], 64)
	if err != nil {
		return nil, fmt.Errorf("parse lat: %w", err)
	}
	lng, err := strconv.ParseFloat(values["lng"], 64)
	if err != nil {
		return nil, fmt.Errorf("parse lng: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, values["updated_at"])
	if err != nil {
		return nil, fmt.Errorf("parse updated_at: %w", err)
	}

	loc := &domain.DriverLocation{
		Lat:       lat,
		Lng:       lng,
		UpdatedAt: updatedAt,
	}
	if heading, ok := values["heading"]; ok && heading != "" {
		parsed, err := strconv.ParseFloat(heading, 64)
		if err != nil {
			return nil, fmt.Errorf("parse heading: %w", err)
		}
		loc.Heading = &parsed
	}
	if speed, ok := values["speed"]; ok && speed != "" {
		parsed, err := strconv.ParseFloat(speed, 64)
		if err != nil {
			return nil, fmt.Errorf("parse speed: %w", err)
		}
		loc.Speed = &parsed
	}
	return loc, nil
}
