package source

// Structure d'une entité (joueur ou monstre)
type Entity struct {
	Name   string
	Health int
	Damage int
}

// Structure d'une arme
type Weapon struct {
	Name   string
	Damage int
}

// Inflige des dégâts à l'entité
func (e *Entity) TakeDamage(damage int) {
	e.Health -= damage
	if e.Health < 0 {
		e.Health = 0
	}
}

// Fonction d'attaque entre deux entités
func Attack(attacker *Entity, defender *Entity, weapon Weapon) {
	if defender.Health > 0 {
		defender.TakeDamage(weapon.Damage)
	}
}
