package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
)

type (
	errorData struct {
		NumSlotsBehind int64 `json:"numSlotsBehind"`
	}

	getHealthRpcError struct {
		Message string    `json:"message"`
		Data    errorData `json:"data"`
		Code    int64     `json:"code"`
	}

	getHealthResponse struct {
		Result string            `json:"result"`
		Error  getHealthRpcError `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#gethealth
func (c *RPCClient) GetHealth(ctx context.Context) (*string, *getHealthRpcError, error) {
	// only retrieve results of method

	body, err := c.rpcRequest(ctx, formatRPCRequest("getHealth", []interface{}{}))
	if err != nil {
		return nil, nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("getHealth response: %v", string(body))

	var resp getHealthResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, &resp.Error, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil, nil
}
