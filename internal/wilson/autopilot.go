package wilson

// Autopilot plays Wilson for --demo, the attract mode the catalog preview is
// recorded from. It looks ahead in each lane and shoots, changes lane or
// jumps to get past whatever is coming, and drifts towards pick-ups when the
// way is clear. It is meant to look like play, not to be perfect.

func harmful(t ObjectType) bool {
	return t == TypeObstacle || t == TypeCar || t == TypeEnemy
}

func wanted(t ObjectType) bool {
	return t == TypeGun || t == TypeAmmo || t == TypePowerUp || t == TypeEnableJump
}

func unwanted(t ObjectType) bool {
	return t == TypePowerDown || t == TypeDisableJump
}

// clear reports whether nothing harmful in lane reaches Wilson within
// `ahead` columns of his position.
func (s *GameState) clear(lane int, ahead float64) bool {
	for _, o := range s.Objects {
		if o.Lane != lane || !harmful(o.Type) {
			continue
		}
		w := 3.0
		if o.Type == TypeCar {
			w = 6.0
		}
		if o.X >= ChickenX-2-w && o.X <= ChickenX+w+ahead {
			return false
		}
	}
	return true
}

// shy reports whether a pick-up Wilson would rather miss lies just ahead.
func (s *GameState) shy(lane int, ahead float64) bool {
	for _, o := range s.Objects {
		if o.Lane == lane && unwanted(o.Type) && o.X >= ChickenX-1 && o.X <= ChickenX+ahead {
			return true
		}
	}
	return false
}

// Autopilot makes one decision; call it once per frame before Update.
func (s *GameState) Autopilot() {
	if !s.Started || s.GameOver || s.Paused || s.Jumping || s.InvincibleTimer > 0 {
		return
	}
	near := (s.Speed+1.2)*4 + 3 // columns covered in a few frames
	far := near * 2.2
	lane := s.ChickenLane

	if !s.clear(lane, near) {
		if s.HasGun && s.Ammo > 0 && !s.IsShooting {
			s.Shoot()
			return
		}
		for _, d := range []int{-1, 1} {
			l := lane + d
			if l >= 0 && l < Lanes && s.clear(l, far) {
				s.move(d)
				return
			}
		}
		for _, d := range []int{-1, 1} {
			l := lane + d
			if l >= 0 && l < Lanes && s.clear(l, near) {
				s.move(d)
				return
			}
		}
		s.Jump()
		return
	}

	// The way is clear: step aside from a bad pick-up, or go for a good one.
	if s.shy(lane, near) {
		for _, d := range []int{-1, 1} {
			l := lane + d
			if l >= 0 && l < Lanes && s.clear(l, far) && !s.shy(l, near) {
				s.move(d)
				return
			}
		}
	}
	for _, o := range s.Objects {
		if !wanted(o.Type) || o.Lane == lane || o.X < ChickenX || o.X > ChickenX+far {
			continue
		}
		d := 1
		if o.Lane < lane {
			d = -1
		}
		if s.clear(lane+d, far) {
			s.move(d)
		}
		return
	}
}

func (s *GameState) move(d int) {
	if d < 0 {
		s.MoveUp()
	} else {
		s.MoveDown()
	}
}
