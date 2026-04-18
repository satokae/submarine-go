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

func ResolveAttack(attack AttackAction, defendingFleet Fleet) (AttackOutcome, Fleet) {
	newFleet := make(Fleet, len(defendingFleet))
	copy(newFleet, defendingFleet)
	hitIndex := -1

	for i, sub := range newFleet {
		if sub.HP > 0 && sub.Position == attack.Position {
			hitIndex = i
			break
		}
	}

	if hitIndex != -1 {
		newFleet[hitIndex].HP--
		if newFleet[hitIndex].HP <= 0 {
			return HitAndSunk, newFleet
		}
		return Hit, newFleet
	}

	neighbors := GetNeighbors(attack.Position)
	isNear := false

	for _, sub := range defendingFleet {
		if sub.HP <= 0 {
			continue
		}

		for _, n := range neighbors {
			if sub.Position == n {
				isNear = true
				break
			}
		}
		if isNear {
			break
		}
	}

	if isNear {
		return HighWaves, newFleet
	}

	return Miss, newFleet
}
