package gameState

type GameState struct {
	Counter int `json:"counter"`
	Turn    int `json:"turn"`
}

func NewGameState() *GameState {
	return &GameState{}
}

func (gs *GameState) GetGameState() int {
	return gs.Counter
}

func (gs *GameState) IncCounter(nextTurn int) {
	gs.Counter++
	gs.Turn = nextTurn
}

func (gs *GameState) DecCounter(nextTurn int) {
	gs.Counter--
	gs.Turn = nextTurn
}
