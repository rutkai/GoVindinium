package vindinium

type Game struct {
	State    *State
	Board    *Board  `json:"board"`
	Heroes   []*Hero `json:"heroes"`
	Id       string  `json:"id"`
	Finished bool    `json:"finished"`
	Turn     int     `json:"turn"`
	MaxTurns int     `json:"maxTurns"`
	Hero     *Hero   `json:"hero"`
	Token    string  `json:"token"`
	Crashed  bool    `json:"crashed"`
	ViewUrl  string  `json:"viewUrl"`
	PlayUrl  string  `json:"PlayUrl"`
	Mines    map[Position]*MineTile
	Taverns  []Position
}

func NewGame(state *State) (game *Game) {
	game = state.Game
	game.State = state
	game.Board.parseTiles()
	game.populateMines()
	game.populateTaverns()

	return
}

func (game *Game) populateMines() {
	game.Mines = make(map[Position]*MineTile)

	for x := 0; x < game.Board.Size; x++ {
		for y := 0; y < game.Board.Size; y++ {
			if isMine, mine := game.Board.Mine(Position{x, y}); isMine {
				game.Mines[Position{x, y}] = mine
			}
		}
	}

	return
}

func (game *Game) populateTaverns() {
	game.Taverns = []Position{}

	for x := 0; x < game.Board.Size; x++ {
		for y := 0; y < game.Board.Size; y++ {
			if game.Board.Tavern(Position{x, y}) {
				game.Taverns = append(game.Taverns, Position{x, y})
			}
		}
	}

	return
}

func (game *Game) NotMyMines() []Position {
	var mines []Position

	for k, v := range game.Mines {
		if v.HeroId != game.Hero.Id {
			mines = append(mines, k)
		}
	}

	return mines
}
