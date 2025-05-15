package job

import (
	"math"
	"fmt"
	"sync"
	"time"
)

type DelayResult struct {
	Results      []int   // 延迟的所有测试结果(失败为-1)
	Average      float64 // 延迟的平均值
	StdDev       float64 // 延迟的标准差
	SuccessCount int     // 成功次数
	FailureCount int     // 失败次数
}

func (t *Task) TestDelay(result *Result) (err error) {
	stats, err := t.TestProxyDelay(result.Proxy.Name)

	if err != nil {
		fmt.Println("测试延迟失败:", err)
		return err
	}
	result.SetDelay(stats)

	fmt.Println("延迟测试结果:", stats.Results)
	fmt.Printf("平均延迟: %.2f ms\n", stats.Average)
	fmt.Printf("标准差: %.2f ms\n", stats.StdDev)
	fmt.Printf("成功次数: %d\n", stats.SuccessCount)
	fmt.Printf("失败次数: %d\n", stats.FailureCount)

	return nil
}

// TestProxyDelay 调用 GetProxyDelay 10 次，计算返回测试结果
func (t *Task) TestProxyDelay(name string) (*DelayResult, error) {
	const numTests = 10
	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		results     = make([]int, numTests)
		successData []int
	)

	wg.Add(numTests)
	for i := 0; i < numTests; i++ {
		go func(index int) {
			defer wg.Done()

			res, err := t.clash.GetProxyDelay(name)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results[index] = -1
				return
			}

			results[index] = res.DelayResponse
			successData = append(successData, res.DelayResponse)
			time.Sleep(2 * 100 * time.Millisecond)
		}(i)
	}
	wg.Wait()

	// 统计成功次数
	successCount := len(successData)
	failureCount := numTests - successCount

	var average, stdDev float64
	if successCount > 0 {
		// 计算平均值
		sum := 0
		for _, delay := range successData {
			sum += delay
		}
		average = float64(sum) / float64(successCount)

		// 计算标准差
		variance := 0.0
		for _, delay := range successData {
			diff := float64(delay) - average
			variance += diff * diff
		}
		variance /= float64(successCount)
		stdDev = math.Sqrt(variance)
	}

	return &DelayResult{
		Results:      results,
		Average:      average,
		StdDev:       stdDev,
		SuccessCount: successCount,
		FailureCount: failureCount,
	}, nil
}
