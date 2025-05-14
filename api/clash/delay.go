package clash

import (
	"fmt"
)

type DelayResponse struct {
	DelayResponse int `json:"delay"`
}

func (c *Client) GetProxyDelay(name string) (*DelayResponse, error) {
    resp, err := c.R().
        SetQueryParams(map[string]string{
            "timeout": "5000",
            "url":     "http://www.gstatic.com/generate_204",
        }).
        SetResult(&DelayResponse{}).
        Get(c.Addr + "/proxies/" + name + "/delay")
    
    if err != nil {
        return nil, err
    }

    if !resp.IsSuccess() {
        return nil, fmt.Errorf("response status %s", resp.Status())
    }

    return resp.Result().(*DelayResponse), nil
}
