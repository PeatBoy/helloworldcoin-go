package vo

/*
 @author king 409060350@qq.com
*/

type ActiveMinerRequest struct {
}
type ActiveMinerResponse struct {
}

type DeactiveMinerRequest struct {
}
type DeactiveMinerResponse struct {
}

type IsMinerActiveRequest struct {
}
type IsMinerActiveResponse struct {
	MinerInActiveState bool `json:"minerInActiveState"`
}
type GetMaxBlockHeightRequest struct {
}
type GetMaxBlockHeightResponse struct {
	MaxBlockHeight uint64 `json:"maxBlockHeight"`
}
type SetMaxBlockHeightRequest struct {
	MaxBlockHeight uint64 `json:"maxBlockHeight"`
}
type SetMaxBlockHeightResponse struct {
}
