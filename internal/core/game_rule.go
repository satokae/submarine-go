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

func IsAttackPossible(attack AttackAction, friendlyFleet Fleet) bool {
	activeFriendlyFleet := Fleet{}
	for _, sub := range friendlyFleet {
		if sub.HP > 0 {
			activeFriendlyFleet = append(activeFriendlyFleet, sub)
		}
	}

	for _, sub := range activeFriendlyFleet {
		if sub.Position == attack.Position {
			return false
		}
	}

	neighbors := GetNeighbors(attack.Position)

	for _, n := range neighbors {
		for _, sub := range activeFriendlyFleet {
			if sub.Position == n {
				return true
			}
		}
	}

	return false
}
