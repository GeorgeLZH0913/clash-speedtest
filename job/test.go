package job

import (
	"sync"

	"github.com/starudream/clash-speedtest/api/clash"
	"github.com/starudream/clash-speedtest/api/common"
)

type Result struct {
	Proxy *clash.Proxy

	Ip      string
	Country string
	City    string
	Lat     string
	Lon     string

	threads   int
	total     *common.DownloadResult
	downloads []*common.DownloadResult
	mu        sync.Mutex

	delays *DelayResult
}

func (t *Result) SetDownload(i int, v *common.DownloadResult) {
	if v == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.downloads[i] = v
	t.total.TotalSize += v.TotalSize
	t.total.ConnTime += v.ConnTime
	t.total.RespTime += v.RespTime
}

func (t *Result) SetDelay(v *DelayResult) {
	if v == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.delays.Results = v.Results
	t.delays.Average = v.Average
	t.delays.StdDev = v.StdDev
	t.delays.SuccessCount = v.SuccessCount
	t.delays.FailureCount = v.FailureCount
}

func (t *Task) Test(proxy *clash.Proxy) (*Result, error) {
	err := t.clash.SetGlobalProxy(proxy.Name)
	if err != nil {
		return nil, err
	}

	result := &Result{
		Proxy:     proxy,
		threads:   t.Threads,
		total:     &common.DownloadResult{},
		downloads: make([]*common.DownloadResult, t.Threads),
	
		delays:    &DelayResult{},
	}

	err = t.Down(result)
	if err != nil {
		return nil, err
	}

	err = t.TestDelay(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
