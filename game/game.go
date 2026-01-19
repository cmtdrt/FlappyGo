package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

var (
	skyColor = color.RGBA{R: 135, G: 206, B: 235, A: 255}
)

// Représente l'instance du jeu
type Game struct {
	bird    *Bird
	pipes   []*Pipe
	score   int
	started bool
	dead    bool

	frame *ebiten.Image
}

// Crée une nouvelle instance du jeu
func NewGame() *Game {
	g := &Game{
		bird:  NewBird(),
		frame: ebiten.NewImage(ScreenWidth, ScreenHeight),
	}

	// Crée 3 tuyaux de départ
	for i := 0; i < 3; i++ {
		x := float64(ScreenWidth + i*PipeSpacing)
		g.pipes = append(g.pipes, NewPipe(x))
	}

	return g
}

// Met à jour l'état du jeu
func (g *Game) Update() error {
	// Appui sur ESPACE pour commencer ou relancer
	if !g.started {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.started = true
			g.bird.Jump()
		}
		return nil
	}
	// Si on est mort, la touche ESPACE relance le jeu
	if g.dead {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset()
		}
		return nil
	}

	// Saut avec ESPACE
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.bird.Jump()
	}

	// Met à jour la position de l'oiseau
	g.bird.Update()

	// Vérifie si l'oiseau touche le sol ou le plafond
	if g.bird.CheckBounds() {
		g.dead = true
	}

	// Met à jour la position des tuyaux
	g.updatePipes()

	// Vérifie si l'oiseau entre en collision avec les tuyaux
	g.checkCollisions()

	return nil
}

// Met à jour tous les tuyaux et gère le recyclage
func (g *Game) updatePipes() {
	for _, p := range g.pipes {
		p.Update()

		// Incrémente le score quand on passe un tuyau
		if !p.Passed {
			birdLeft, _, _, _ := g.bird.GetBounds()
			if p.X+PipeWidth < birdLeft {
				p.Passed = true
				g.score++
			}
		}
	}

	// Recycle les tuyaux qui sont sortis de l'écran
	for _, p := range g.pipes {
		if p.IsOffScreen() {
			maxX := 0.0
			for _, other := range g.pipes {
				if other.X > maxX {
					maxX = other.X
				}
			}
			p.Reset(maxX)
		}
	}
}

// Vérifie les collisions entre l'oiseau et les tuyaux
func (g *Game) checkCollisions() {
	birdLeft, birdRight, birdTop, birdBottom := g.bird.GetBounds()

	for _, p := range g.pipes {
		if p.CheckCollision(birdLeft, birdRight, birdTop, birdBottom) {
			g.dead = true
			return
		}
	}
}

// Dessine le jeu à l'écran
func (g *Game) Draw(screen *ebiten.Image) {
	// Dessine la scène dans un buffer interne pour pouvoir appliquer un effet global quand on est mort.
	dst := g.frame
	dst.Clear()

	// Fond
	dst.Fill(skyColor)

	// Tuyaux
	for _, p := range g.pipes {
		// Tuyau du haut
		ebitenutil.DrawRect(dst, p.X, 0, PipeWidth, p.GapY, color.RGBA{R: 0, G: 200, B: 0, A: 255})
		// Tuyau du bas
		ebitenutil.DrawRect(dst, p.X, p.GapY+PipeGapHeight, PipeWidth, float64(ScreenHeight)-(p.GapY+PipeGapHeight), color.RGBA{R: 0, G: 200, B: 0, A: 255})
	}

	// Design de l'oiseau
	birdColor := color.RGBA{R: 255, G: 215, B: 0, A: 255}
	birdX := g.bird.X - g.bird.W/2
	birdY := g.bird.Y - g.bird.H/2
	ebitenutil.DrawRect(dst, birdX, birdY, g.bird.W, g.bird.H, birdColor)

	// Score
	scoreStr := fmt.Sprintf("%d", g.score)
	face := basicfont.Face7x13
	bounds := text.BoundString(face, scoreStr)

	// On centre le texte, puis on le grossit via une mise à l'échelle.
	scale := 4.0
	tmpW := bounds.Dx()
	tmpH := bounds.Dy()
	if tmpW < 1 {
		tmpW = 1
	}
	if tmpH < 1 {
		tmpH = 1
	}

	// Dessine le texte dans une image temporaire, puis scale cette image.
	tmp := ebiten.NewImage(tmpW, tmpH)
	tmp.Fill(color.Transparent)
	text.Draw(tmp, scoreStr, face, -bounds.Min.X, -bounds.Min.Y, color.White)

	w := float64(tmpW) * scale

	x := (float64(ScreenWidth) - w) / 2
	y := 60.0 // un peu en haut

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	dst.DrawImage(tmp, op)

	if !g.started {
		drawCenteredText(dst, "Appuie sur ESPACE pour commencer", float64(ScreenHeight)*0.35, 2.0, color.White)
	} else if g.dead {
		drawCenteredText(dst, "Perdu !", float64(ScreenHeight)*0.30, 3.0, color.White)
		drawCenteredText(dst, "Appuie sur ESPACE pour recommencer", float64(ScreenHeight)*0.30+45, 2.2, color.White)
	}

	// Si on est mort, on applique un effet de niveaux de gris au buffer, puis on le dessine à l'écran.
	// Sinon, on dessine le buffer tel quel.
	if g.dead {
		var cm colorm.ColorM
		cm.ChangeHSV(0, 0, 1) // saturation à 0 = noir et blanc
		colorm.DrawImage(screen, dst, cm, &colorm.DrawImageOptions{})
	} else {
		screen.DrawImage(dst, &ebiten.DrawImageOptions{})
	}
}

// Définit la taille logique de l'écran
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Réinitialise le jeu à son état initial
func (g *Game) reset() {
	g.bird.Reset()
	g.score = 0
	g.started = false
	g.dead = false

	g.pipes = g.pipes[:0]
	for i := 0; i < 3; i++ {
		x := float64(ScreenWidth + i*PipeSpacing)
		g.pipes = append(g.pipes, NewPipe(x))
	}
}

// Dessine un texte centré horizontalement, à la position y donnée, avec une échelle.
func drawCenteredText(dst *ebiten.Image, msg string, y float64, baseScale float64, clr color.Color) {
	face := basicfont.Face7x13
	bounds := text.BoundString(face, msg)

	w := bounds.Dx()
	h := bounds.Dy()
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	tmp := ebiten.NewImage(w, h)
	tmp.Fill(color.Transparent)
	text.Draw(tmp, msg, face, -bounds.Min.X, -bounds.Min.Y, clr)

	// Adapte l'échelle pour ne pas dépasser la largeur de l'écran (80% max).
	maxScaleFit := 0.8 * float64(ScreenWidth) / float64(w)
	scale := baseScale
	if scale > maxScaleFit {
		scale = maxScaleFit
	}
	if scale < 0.5 {
		scale = 0.5
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((float64(ScreenWidth)-float64(w)*scale)/2, y)
	dst.DrawImage(tmp, op)
}
