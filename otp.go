package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key     string
	Created time.Time
}

type RetentionMap map[string]OTP

func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rm := make(RetentionMap)

	go rm.Retention(ctx, retentionPeriod)

	return rm
}

// generate a new OTP to return to user
func (rm RetentionMap) NewOTP() OTP {
	otp := OTP{
		Key:     uuid.NewString(),
		Created: time.Now(),
	}

	rm[otp.Key] = otp
	return otp
}

// check retention map to see if otp exists, delete if it does
func (rm RetentionMap) VerifyOTP(otp string) bool {
	// Verify OTP is existing
	if _, ok := rm[otp]; !ok {
		// otp does not exist
		return false
	}
	delete(rm, otp)
	return true
}

// every 500 ms check the retention map, and remove expired otps
func (rm RetentionMap) Retention(ctx context.Context, retentionPeriod time.Duration) {
	ticker := time.NewTicker(500 * time.Millisecond)

	select {
	case <-ticker.C:
		for _, otp := range rm {
			if otp.Created.Add(retentionPeriod).Before(time.Now()) {
				delete(rm, otp.Key)
			}
		}
	case <-ctx.Done():
		return
	}

}
