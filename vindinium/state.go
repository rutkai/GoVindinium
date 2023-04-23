package vindinium

type State struct {
	Game    *Game  `json:"game"`
	Hero    *Hero  `json:"hero"`
	Token   string `json:"token"`
	ViewUrl string `json:"viewUrl"`
	PlayUrl string `json:"PlayUrl"`
}

func (s *State) Init() {
	s.Game.State = s
	s.Game.Hero = s.Hero
	s.Game.Token = s.Token
	s.Game.ViewUrl = s.ViewUrl
	s.Game.PlayUrl = s.PlayUrl

	s.Game.Board.parseTiles()
	s.Game.populateMines()
	s.Game.populateTaverns()
}
