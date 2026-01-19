package game

import (
	"math/rand"
)

// Représente un tuyau
type Pipe struct {
	X      float64 // Position X
	GapY   float64 // Position Y du gap
	Passed bool    // Si le joueur a déjà passé ce tuyau
}

// Créer un nouveau tuyau à la position X donnée (le gap est aléatoire)
func NewPipe(x float64) *Pipe {
	return &Pipe{
		X:      x,
		GapY:   randomGapY(),
		Passed: false,
	}
}

// Génère une position Y aléatoire pour le gap d'un tuyau
func randomGapY() float64 {
	margin := 60.0
	minY := margin
	maxY := float64(ScreenHeight) - margin - PipeGapHeight
	return minY + rand.Float64()*(maxY-minY)
}

// Déplace le tuyau vers la gauche
func (p *Pipe) Update() {
	p.X -= PipeSpeed
}

// Vérifie si le tuyau est complètement sorti de l'écran
func (p *Pipe) IsOffScreen() bool {
	return p.X+PipeWidth < 0
}

// Repositionne le tuyau à droite de l'écran avec un nouveau gap
func (p *Pipe) Reset(maxX float64) {
	p.X = maxX + PipeSpacing
	p.GapY = randomGapY()
	p.Passed = false
}

// Vérifie si l'oiseau entre en collision avec ce tuyau, renvoie true si c'est la cas
func (p *Pipe) CheckCollision(birdLeft, birdRight, birdTop, birdBottom float64) bool {
	pipeLeft := p.X
	pipeRight := p.X + PipeWidth

	// Tuyau du haut
	topPipeTop := 0.0
	topPipeBottom := p.GapY

	// Tuyau du bas
	bottomPipeTop := p.GapY + PipeGapHeight
	bottomPipeBottom := float64(ScreenHeight)

	// Collision avec le tuyau du haut
	if birdRight > pipeLeft && birdLeft < pipeRight && birdTop < topPipeBottom && birdBottom > topPipeTop {
		return true
	}

	// Collision avec le tuyau du bas
	if birdRight > pipeLeft && birdLeft < pipeRight && birdTop < bottomPipeBottom && birdBottom > bottomPipeTop {
		return true
	}

	return false
}
