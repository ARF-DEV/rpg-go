package game

func getNeighbors(curPos Pos, level *Level) []Pos {
	neighbours := []Pos{}
	left := Pos{curPos.X - 1, curPos.Y}
	right := Pos{curPos.X + 1, curPos.Y}
	up := Pos{curPos.X, curPos.Y - 1}
	down := Pos{curPos.X, curPos.Y + 1}

	if level.GetTile(int(left.X), int(left.Y)).CanWalk() {
		neighbours = append(neighbours, left)

	}
	if level.GetTile(int(right.X), int(right.Y)).CanWalk() {
		neighbours = append(neighbours, right)
	}

	if level.GetTile(int(up.X), int(up.Y)).CanWalk() {
		neighbours = append(neighbours, up)
	}
	if level.GetTile(int(down.X), int(down.Y)).CanWalk() {
		neighbours = append(neighbours, down)
	}

	return neighbours
}
