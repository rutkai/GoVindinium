package vindinium

import (
	"fmt"
	"strconv"
)

const (
	WALL = iota - 2
	AIR
	TAVERN

	AIR_TILE    = " "
	WALL_TILE   = "#"
	TAVERN_TILE = "["
	MINE_TILE   = "$"
	HERO_TILE   = "@"
)

var (
	AIM = map[Direction]*Position{
		"North": {-1, 0},
		"East":  {0, 1},
		"South": {1, 0},
		"West":  {0, -1},
	}
)

type Board struct {
	Size    int    `json:"size"`
	Tiles   string `json:"tiles"`
	Tileset [][]interface{}
}

type Position struct {
	X, Y int
}

func tileToInt(tiles string, index int) int {
	tile := []rune(tiles)[index]
	str, _ := strconv.Atoi(string(tile))

	return str
}

func (board *Board) parseTile(tile string) interface{} {
	switch string([]rune(tile)[0]) {
	case AIR_TILE:
		return AIR
	case WALL_TILE:
		return WALL
	case TAVERN_TILE:
		return TAVERN
	case MINE_TILE:
		char := string([]rune(tile)[1])
		id, _ := strconv.Atoi(char)
		return &MineTile{id}
	case HERO_TILE:
		char := string([]rune(tile)[1])
		id, _ := strconv.Atoi(char)
		return &HeroTile{id}
	default:
		return -3
	}
}

func (board *Board) parseTiles() {
	var vector [][]rune
	var matrix [][][]rune
	ts := make([][]interface{}, board.Size)

	for i := 0; i <= len(board.Tiles)-2; i = i + 2 {
		vector = append(vector, []rune(board.Tiles)[i:i+2])
	}

	for i := 0; i < len(vector); i = i + board.Size {
		matrix = append(matrix, vector[i:i+board.Size])
	}

	for xi, x := range matrix {
		innerList := make([]interface{}, board.Size)
		for xsi, xs := range x {
			innerList[xsi] = board.parseTile(string(xs))
		}
		ts[xi] = innerList
	}

	board.Tileset = ts
}

func (board *Board) Passable(loc Position) bool {
	return board.Tileset[loc.X][loc.Y] == AIR
}

func (board *Board) Wall(loc Position) bool {
	return board.Tileset[loc.X][loc.Y] == WALL
}

func (board *Board) Tavern(loc Position) bool {
	return board.Tileset[loc.X][loc.Y] == TAVERN
}

func (board *Board) Mine(loc Position) (bool, *MineTile) {
	switch tile := board.Tileset[loc.X][loc.Y]; tile.(type) {
	case *MineTile:
		return true, tile.(*MineTile)
	default:
		return false, nil
	}
}

func (board *Board) Hero(loc Position) (bool, *HeroTile) {
	switch tile := board.Tileset[loc.X][loc.Y]; tile.(type) {
	case *HeroTile:
		return true, tile.(*HeroTile)
	default:
		return false, nil
	}
}

func (board *Board) HasNeighbouringHero(loc Position, myself int) bool {
	// Direct neighbours
	if isHero, hero := board.Hero(*board.To(loc, "North")); isHero && hero.Id != myself {
		return true
	}
	if isHero, hero := board.Hero(*board.To(loc, "South")); isHero && hero.Id != myself {
		return true
	}
	if isHero, hero := board.Hero(*board.To(loc, "East")); isHero && hero.Id != myself {
		return true
	}
	if isHero, hero := board.Hero(*board.To(loc, "West")); isHero && hero.Id != myself {
		return true
	}

	// Neighbour corners
	if isHero, hero := board.Hero(*board.To(*board.To(loc, "North"), "East")); isHero && hero.Id != myself {
		return true
	}
	if isHero, hero := board.Hero(*board.To(*board.To(loc, "East"), "South")); isHero && hero.Id != myself {
		return true
	}
	if isHero, hero := board.Hero(*board.To(*board.To(loc, "South"), "West")); isHero && hero.Id != myself {
		return true
	}
	if isHero, hero := board.Hero(*board.To(*board.To(loc, "West"), "North")); isHero && hero.Id != myself {
		return true
	}

	return false
}

func (board *Board) To(loc Position, direction Direction) *Position {
	row := loc.X
	col := loc.Y
	dLoc := AIM[direction]
	nRow := row + dLoc.X
	if nRow < 0 {
		nRow = 0
	}
	if nRow > board.Size-1 {
		nRow = board.Size - 1
	}
	nCol := col + dLoc.Y
	if nCol < 0 {
		nCol = 0
	}
	if nCol > board.Size-1 {
		nCol = board.Size - 1
	}

	return &Position{nRow, nCol}
}

func DirectionOf(from Position, to Position) Direction {
	if AIM["North"].X+from.X == to.X && AIM["North"].Y+from.Y == to.Y {
		return "North"
	}
	if AIM["South"].X+from.X == to.X && AIM["South"].Y+from.Y == to.Y {
		return "South"
	}
	if AIM["West"].X+from.X == to.X && AIM["West"].Y+from.Y == to.Y {
		return "West"
	}
	if AIM["East"].X+from.X == to.X && AIM["East"].Y+from.Y == to.Y {
		return "East"
	}
	fmt.Println("ERROR: Invalid positions for direction calculation: ", from, to)
	return "Stay"
}

func (pos *Position) IsValid(size int) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X < size && pos.Y < size
}
