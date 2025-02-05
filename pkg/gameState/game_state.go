package game_state

import "fmt"

type GameState struct {
	gameState [3][3]uint8
}

func NewGameState() *GameState {
	return &GameState{
		gameState: [3][3]uint8{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
	}
}

func (gs *GameState) Set(i uint8, player uint8) {
	switch i {
	case 1:
		gs.gameState[0][0] = player
	case 2:
		gs.gameState[0][1] = player
	case 3:
		gs.gameState[0][2] = player
	case 4:
		gs.gameState[1][0] = player
	case 5:
		gs.gameState[1][1] = player
	case 6:
		gs.gameState[1][2] = player
	case 7:
		gs.gameState[2][0] = player
	case 8:
		gs.gameState[2][1] = player
	case 9:
		gs.gameState[2][2] = player
	default:
		fmt.Println("invalid index")
	}
}

func (gs *GameState) GetGameState() [3][3]uint8 {
	return gs.gameState
}
