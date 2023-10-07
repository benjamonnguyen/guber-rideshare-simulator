package msgmodel

type LocationMsg struct {
	UserId string `json:"userId"`
	Coord  [2]int `json:"coord"`
}

type ServerMsg struct {
	ClientId string `json:"clientId"`
	Msg      string `json:"msg"`
}
