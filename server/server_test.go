package server

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	couponv1 "github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1"
)

// TestHighTrafficCouponIssuance simulates 1000 requests per second traffic condition to verify concurrency handling in the coupon issuance system.
func TestHighTrafficCouponIssuance(t *testing.T) {
	srv := NewCouponIssuanceServer()

	now := time.Now().UTC()
	startAt := now.Add(-1 * time.Hour)
	endAt := now.Add(24 * time.Hour)

	const couponLimit = 10000

	createCampReq := connect.NewRequest(&couponv1.CreateCampaignRequest{
		CouponLimit: couponLimit,
		Name:        "High Traffic Test Campaign",
		Description: "Campaign for testing high traffic conditions (1000 RPS)",
		StartAt:     timestamppb.New(startAt),
		EndAt:       timestamppb.New(endAt),
	})

	createCampResp, err := srv.CreateCampaign(context.Background(), createCampReq)
	require.NoError(t, err)
	require.NotNil(t, createCampResp)

	campId := createCampResp.Msg.Campaign.Id

	t.Run("High traffic condition simulation (1000 RPS)", func(t *testing.T) {
		const (
			requestCount  = 1000
			testDuration  = 10
			totalRequests = requestCount * testDuration
		)

		// Track metrics
		var (
			successCount      int64 = 0
			failCount         int64 = 0
			rateLimitErrors   int64 = 0
			otherErrors       int64 = 0
			responseTimeSum   int64 = 0
			maxResponseTime   int64 = 0
			minResponseTime   int64 = 0xFFFFFFFF
			issuedCouponCodes       = sync.Map{} // Thread-safe map for tracking issued coupon IDs
			duplicateIssues   int64 = 0
		)

		// Wait group to ensure all goroutines complete
		var wg sync.WaitGroup
		wg.Add(totalRequests)

		// Channel to control request rate
		requestTicker := time.NewTicker(time.Second / requestCount)
		defer requestTicker.Stop()

		// Start time for the test
		testStartTime := time.Now()

		// Launch goroutines for all requests
		for i := 0; i < totalRequests; i++ {
			go func(idx int) {
				defer wg.Done()

				// Wait for ticker to control request rate
				<-requestTicker.C

				startTime := time.Now()
				req := connect.NewRequest(&couponv1.IssueCouponRequest{
					CampaignId: campId,
				})

				// Add some context timeout to prevent hanging requests
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				resp, err := srv.IssueCoupon(ctx, req)
				endTime := time.Now()

				// Calculate response time in milliseconds
				responseTime := endTime.Sub(startTime).Milliseconds()

				// Update response time metrics
				atomic.AddInt64(&responseTimeSum, responseTime)

				// Update max/min response times with atomic operations
				for {
					currentMax := atomic.LoadInt64(&maxResponseTime)
					if responseTime <= currentMax {
						break
					}
					if atomic.CompareAndSwapInt64(&maxResponseTime, currentMax, responseTime) {
						break
					}
				}

				for {
					currentMin := atomic.LoadInt64(&minResponseTime)
					if responseTime >= currentMin {
						break
					}
					if atomic.CompareAndSwapInt64(&minResponseTime, currentMin, responseTime) {
						break
					}
				}

				if err == nil {
					atomic.AddInt64(&successCount, 1)

					// Check for duplicate coupon code issuance
					code := resp.Msg.Coupon.Code
					if _, loaded := issuedCouponCodes.LoadOrStore(code, true); loaded {
						atomic.AddInt64(&duplicateIssues, 1)
						t.Logf("Duplicate coupon code detected: %s", code)
					}
				} else {
					atomic.AddInt64(&failCount, 1)

					// Categorize errors
					if err.Error() == "coupon limit exceeded" {
						atomic.AddInt64(&rateLimitErrors, 1)
					} else {
						atomic.AddInt64(&otherErrors, 1)
						t.Logf("Error issuing coupon: %v", err)
					}
				}
			}(i)
		}

		// Wait for all requests to complete
		wg.Wait()

		// Calculate test duration
		actualTestDuration := time.Since(testStartTime).Seconds()
		actualRequestRate := float64(totalRequests) / actualTestDuration

		// Verify final campaign state
		getCampaignReq := connect.NewRequest(&couponv1.GetCampaignRequest{
			CampaignId: campId,
		})
		getCampaignResp, err := srv.GetCampaign(context.Background(), getCampaignReq)
		require.NoError(t, err)

		actualCouponCount := len(getCampaignResp.Msg.Campaign.Coupons)

		// Calculate average response time
		avgResponseTime := float64(responseTimeSum) / float64(totalRequests)

		// Log test results
		t.Logf("High Traffic Test Results:")
		t.Logf("- Test duration: %.2f seconds", actualTestDuration)
		t.Logf("- Actual request rate: %.2f requests/second", actualRequestRate)
		t.Logf("- Success count: %d", successCount)
		t.Logf("- Failure count: %d", failCount)
		t.Logf("- Rate limit errors: %d", rateLimitErrors)
		t.Logf("- Other errors: %d", otherErrors)
		t.Logf("- Actual coupon count in campaign: %d", actualCouponCount)
		t.Logf("- Average response time: %.2f ms", avgResponseTime)
		t.Logf("- Min response time: %d ms", minResponseTime)
		t.Logf("- Max response time: %d ms", maxResponseTime)

		// Test assertions
		assert.Equal(t, int64(0), duplicateIssues, "No duplicate coupon code should be issued")
		assert.Equal(t, int(successCount), actualCouponCount, "Number of successful issuance should match actual coupon count")

		// If campaign limit is hit, check that we have the right number of failures
		if successCount >= couponLimit {
			assert.Equal(t, couponLimit, int(successCount),
				"Number of successful coupon issuances should not exceed the campaign limit")
			assert.Equal(t, totalRequests-couponLimit, int(failCount),
				"Number of failures should equal total requests minus the campaign limit")
		}

		// Response time assertions (adjust thresholds as needed for a system)
		assert.LessOrEqual(t, avgResponseTime, 200.0, "Average response time should be under 200ms")
	})

	t.Run("Concurrent campaign creation", func(t *testing.T) {
		// Test to ensure campaign creation is also thread-safe
		const concurrentCampaigns = 50
		var wg sync.WaitGroup
		wg.Add(concurrentCampaigns)

		campaignIDs := make([]uint32, concurrentCampaigns)
		campaignCreationErrors := make([]error, concurrentCampaigns)

		// Create multiple campaigns concurrently
		for i := 0; i < concurrentCampaigns; i++ {
			go func(idx int) {
				defer wg.Done()

				req := connect.NewRequest(&couponv1.CreateCampaignRequest{
					CouponLimit: 100,
					Name:        fmt.Sprintf("Concurrent Campaign %d", idx),
					Description: fmt.Sprintf("Campaign created in concurrent test %d", idx),
					StartAt:     timestamppb.New(startAt),
					EndAt:       timestamppb.New(endAt),
				})

				resp, err := srv.CreateCampaign(context.Background(), req)
				if err == nil {
					campaignIDs[idx] = resp.Msg.Campaign.Id
				}
				campaignCreationErrors[idx] = err
			}(i)
		}

		wg.Wait()

		// Check campaign creation results
		successfulCreations := 0
		for i, err := range campaignCreationErrors {
			if err == nil {
				successfulCreations++
				assert.NotEmpty(t, campaignIDs[i], "Campaign ID should not be empty")
			} else {
				t.Logf("Campaign creation error: %v", err)
			}
		}

		assert.Equal(t, concurrentCampaigns, successfulCreations,
			"All campaign creation requests should succeed")

		// Verify each created campaign has a unique ID
		idMap := make(map[uint32]bool)
		duplicates := 0

		for _, id := range campaignIDs {
			if idMap[id] {
				duplicates++
			}
			idMap[id] = true
		}

		assert.Equal(t, 0, duplicates, "No duplicate campaign IDs should be generated")
	})
}
