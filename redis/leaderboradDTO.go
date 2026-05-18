package main

type LeaderBoardEntry struct {
	UserID string  `json:"user_id"`
	Score  float64 `json:"score"`
}

type GetLeaderBoardResp struct {
	Data []LeaderBoardEntry `json:"data"`
}

type UpdatScoreRequest struct {
	Score  int    `json:"score"`
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
}
