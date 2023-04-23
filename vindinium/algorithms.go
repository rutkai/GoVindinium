package vindinium

import (
	"fmt"
	"golang.org/x/exp/slices"
)

type PathfinderSettings struct {
	AvoidPlayers bool
}

func Distances(b *Board, sets *PathfinderSettings, start Position) (dist map[Position]int, prev map[Position]Position) {
	var explored []Position
	var queue []Position
	dist = make(map[Position]int)
	prev = make(map[Position]Position)

	queue = append(queue, start)
	explored = append(explored, start)
	dist[start] = 0

	_, myself := b.Hero(start)

	tryDirection := func(current Position, p Position) {
		if p.IsValid(b.Size) && slices.Index(explored, p) == -1 {
			if b.Wall(p) || b.Passable(p) && b.HasNeighbouringHero(p, myself.Id) && sets.AvoidPlayers {
				return
			}

			if b.Passable(p) {
				queue = append(queue, p)
			}

			explored = append(explored, p)
			dist[p] = dist[current] + 1
			prev[p] = current
		}
	}

	for len(queue) != 0 {
		current := queue[0]

		tryDirection(current, Position{current.X - 1, current.Y})
		tryDirection(current, Position{current.X, current.Y - 1})
		tryDirection(current, Position{current.X + 1, current.Y})
		tryDirection(current, Position{current.X, current.Y + 1})

		queue = queue[1:]
	}

	return
}

func NextStepTowards(start Position, dest Position, prev map[Position]Position) Position {
	if _, ok := prev[dest]; !ok {
		fmt.Println("ERROR: Destination is unreachable! Dest: ", dest)
		return start
	}

	for prev[dest] != start {
		dest = prev[dest]
	}

	return dest
}

func Closest(ref Position, list []Position, dist map[Position]int) (Position, int) {
	if len(list) == 0 {
		fmt.Println("ERROR: Closest calculation got empty list! Ref: ", ref)
		return ref, 0
	}

	minDistance := 10000
	closest := ref
	for _, p := range list {
		distance, ok := dist[p]
		if !ok {
			fmt.Println("WARNING: Position is unreachable! Pos: ", p)
			continue
		}
		if distance < minDistance {
			minDistance = distance
			closest = p
		}
	}

	if minDistance == 10000 {
		fmt.Println("WARNING: There is no path to get to any of the destinations (closest calc)!")
	}
	return closest, minDistance
}
