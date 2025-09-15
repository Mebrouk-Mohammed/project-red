package main

import "fmt"

type Personnage struct {
	Name      string
	life      int
	Maxlife   int
	Strenght  int
	Money     int
	Inventory []string
}

func (p *Personnage) AfficherStatut() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("life: %d/%d\n", p.life, p.Maxlife)
}
func (p *Personnage) AfficherInventaire() {
	fmt.Println("Inventory:")
	if len(p.Inventory) == 0 {
		fmt.Println("  (vide)")
	}
	for _, item := range p.Inventory {
		fmt.Printf("  - %s\n", item)
	}
}

func (p *Personnage) PrendreDegats(damage int) {
	p.life -= damage
	if p.life < 0 {
		p.life = 0
	}
	fmt.Printf("%s a pris %d points de dégâts. life restante: %d/%d\n", p.Name, damage, p.life, p.Maxlife)
	if p.life == 0 {
		fmt.Printf("%s est mort!\n", p.Name)
	}
}
func (p *Personnage) Soigner(heal int) {
	p.life += heal
	if p.life > p.Maxlife {
		p.life = p.Maxlife
	}
	fmt.Printf("%s a été soigné de %d points. life actuelle: %d/%d\n", p.Name, heal, p.life, p.Maxlife)
}
