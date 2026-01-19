package game

// Représente le personnage qu'on joue
type Bird struct {
	X, Y float64 // Position X et Y
	Vy   float64 // Vitesse verticale
	W    float64 // Largeur
	H    float64 // Hauteur
}

// Créer un nouvel oiseau à la position de départ
func NewBird() *Bird {
	return &Bird{
		X: BirdStartX,
		Y: float64(ScreenHeight / 2),
		W: BirdWidth,
		H: BirdHeight,
	}
}

// Fait sauter l'oiseau
func (b *Bird) Jump() {
	b.Vy = JumpVelocity
}

// Met à jour la position de l'oiseau avec la gravité
func (b *Bird) Update() {
	b.Vy += Gravity
	b.Y += b.Vy
}

// Vérifie les collisions avec les bords de l'écran
// Renvoie true si l'oiseau touche le sol ou le plafond
func (b *Bird) CheckBounds() bool {
	if b.Y+b.H/2 >= float64(ScreenHeight) {
		b.Y = float64(ScreenHeight) - b.H/2
		return true // Touche le sol
	}
	if b.Y-b.H/2 <= 0 {
		b.Y = b.H / 2
		b.Vy = 0
	}
	return false
}

// Remet l'oiseau à sa position initiale
func (b *Bird) Reset() {
	b.X = BirdStartX
	b.Y = float64(ScreenHeight / 2)
	b.Vy = 0
}

// Renvoie les coordonnées de collision de l'oiseau
func (b *Bird) GetBounds() (left, right, top, bottom float64) {
	left = b.X - b.W/2
	right = b.X + b.W/2
	top = b.Y - b.H/2
	bottom = b.Y + b.H/2
	return
}
