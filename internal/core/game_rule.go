package core

func isOccupied(pos Position, fleet Fleet, sunkPositions []Position) bool {
	for _, sub := range fleet {
		if sub.Position == pos {
			return true
		}
	}
	for _, sp := range sunkPositions {
		if sp == pos {
			return true
		}
	}
	return false
}

func GetValidMoveDestination(move MoveAction, targetSub Submarine, friendlyFleet Fleet, sunkPositions []Position) (Position, bool) {
	dx, dy := move.Direction.ToVector()
	finalPos, err := targetSub.Position.Move(dx*move.Distance, dy*move.Distance)
	if err != nil {
		return targetSub.Position, false
	}

	if move.Distance == 2 {
		intermediatePos, err := targetSub.Position.Move(dx, dy)
		if err != nil || isOccupied(intermediatePos, friendlyFleet, sunkPositions) {
			return targetSub.Position, false
		}
	}

	if isOccupied(finalPos, friendlyFleet, sunkPositions) {
		return targetSub.Position, false
	}

	return finalPos, true
}
