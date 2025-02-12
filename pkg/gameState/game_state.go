package gameState

type GameState struct {
	Counter int   `json:"counter"`
	Turn    uint8 `json:"turn"`
}

func NewGameState() *GameState {
	return &GameState{
		Counter: 0,
		Turn:    1,
	}
}

func (gs *GameState) GetGameState() int {
	return gs.Counter
}

func (gs *GameState) IncCounter() {
	gs.Counter++
	gs.Turn = 3 - gs.Turn
}

func (gs *GameState) DecCounter() {
	gs.Counter--
	gs.Turn = 3 - gs.Turn
}
