package vindinium

type IntrovertBot struct{}

func (b *IntrovertBot) Move(state *State) Direction {
	dist, prev := Distances(state.Game.Board, &PathfinderSettings{AvoidPlayers: true}, *state.Hero.Pos)

	dest := *state.Hero.Pos
	if closestTavern, tavernDistance := Closest(*state.Hero.Pos, state.Game.Taverns, dist); state.Hero.Life <= 30 || state.Hero.Life <= 80 && tavernDistance == 1 {
		dest = closestTavern
	} else if mines := state.Game.NotMyMines(); len(mines) > 0 {
		dest, _ = Closest(*state.Hero.Pos, mines, dist)
	} else {
		return "Stay"
	}

	next := NextStepTowards(*state.Hero.Pos, dest, prev)
	return DirectionOf(*state.Hero.Pos, next)
}
