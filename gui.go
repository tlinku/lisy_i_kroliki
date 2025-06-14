package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type GUI struct {
	app         fyne.App
	window      fyne.Window
	world       *World
	simulation  *Simulation
	gridWidget  *widget.RichText
	chartWidget *fyne.Container
	chartImage  *widget.Icon
	widthEntry  *widget.Entry
	heightEntry *widget.Entry
	foxEntry    *widget.Entry
	rabbitEntry *widget.Entry
	grassEntry  *widget.Entry
	startButton *widget.Button
	resetButton *widget.Button
	stepButton  *widget.Button
	turnLabel   *widget.Label
	statsLabel  *widget.Label
	turnData    []float64
	foxData     []float64
	rabbitData  []float64
	grassData   []float64
}

type Simulation struct {
	world   *World
	running bool
	ticker  *time.Ticker
	stopCh  chan bool
}

func NewGUI() *GUI {
	a := app.New()
	a.SetIcon(nil)

	w := a.NewWindow("Lisy i Króliki")
	w.Resize(fyne.NewSize(1600, 1000))
	w.CenterOnScreen()

	gui := &GUI{
		app:    a,
		window: w,
	}

	gui.initializeComponents()
	gui.setupLayout()
	gui.resetToDefaults()

	return gui
}

func (g *GUI) initializeComponents() {
	g.widthEntry = widget.NewEntry()
	g.widthEntry.SetText("20")
	g.heightEntry = widget.NewEntry()
	g.heightEntry.SetText("15")
	g.foxEntry = widget.NewEntry()
	g.foxEntry.SetText("5")
	g.rabbitEntry = widget.NewEntry()
	g.rabbitEntry.SetText("15")
	g.grassEntry = widget.NewEntry()
	g.grassEntry.SetText("50")
	g.startButton = widget.NewButton("▶ Start", g.toggleSimulation)
	g.resetButton = widget.NewButton("🔄 Reset", g.resetSimulation)
	g.stepButton = widget.NewButton("⏯ Krok", g.stepSimulation)
	g.turnLabel = widget.NewLabel("Tura: 0")
	g.statsLabel = widget.NewLabel("Populacja:\n🦊 Lisy: 0\n🐰 Króliki: 0\n🌱 Trawa: 0")
	g.gridWidget = widget.NewRichText()
	g.chartImage = widget.NewIcon(nil)
	g.chartWidget = container.NewVBox(
		widget.NewLabel("Wykres Populacji"),
		g.chartImage,
	)
}

func (g *GUI) setupLayout() {
	settingsForm := container.NewVBox(
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Szerokość:", g.widthEntry),
			widget.NewFormItem("Wysokość:", g.heightEntry),
			widget.NewFormItem("Lisy:", g.foxEntry),
			widget.NewFormItem("Króliki:", g.rabbitEntry),
			widget.NewFormItem("Trawa:", g.grassEntry),
		),
	)
	controlsBox := container.NewVBox(
		widget.NewSeparator(),
		g.startButton,
		g.stepButton,
		g.resetButton,
		widget.NewSeparator(),
		g.turnLabel,
		g.statsLabel,
	)

	gridScroll := container.NewScroll(g.gridWidget)
	gridScroll.SetMinSize(fyne.NewSize(500, 400))

	gridCard := container.NewBorder(
		nil, nil, nil,
		gridScroll,
	)

	leftPanel := container.NewVBox(
		container.NewBorder(
			nil, nil, nil,
			settingsForm,
		),
		controlsBox,
	)
	rightPanel := container.NewBorder(
		widget.NewLabel(""),
		nil, nil, nil,
		g.chartImage,
	)
	rightSplit := container.NewHSplit(
		gridCard,
		rightPanel,
	)
	rightSplit.SetOffset(0.4)
	content := container.NewHSplit(
		leftPanel,
		rightSplit,
	)
	content.SetOffset(0.25)

	g.window.SetContent(content)
}

func (g *GUI) createWorld() {
	width, _ := strconv.Atoi(g.widthEntry.Text)
	height, _ := strconv.Atoi(g.heightEntry.Text)
	foxCount, _ := strconv.Atoi(g.foxEntry.Text)
	rabbitCount, _ := strconv.Atoi(g.rabbitEntry.Text)
	grassCount, _ := strconv.Atoi(g.grassEntry.Text)
	if width < 5 || width > 50 {
		width = 20
	}
	if height < 5 || height > 50 {
		height = 15
	}
	if foxCount < 0 || foxCount > 50 {
		foxCount = 5
	}
	if rabbitCount < 0 || rabbitCount > 100 {
		rabbitCount = 15
	}
	if grassCount < 0 || grassCount > 200 {
		grassCount = 50
	}

	g.world = NewWorld(width, height)
	g.world.PopulateRandomly(foxCount, rabbitCount, grassCount)

	g.simulation = &Simulation{
		world:   g.world,
		running: false,
		stopCh:  make(chan bool),
	}
	g.turnData = []float64{}
	g.foxData = []float64{}
	g.rabbitData = []float64{}
	g.grassData = []float64{}

	g.updateDisplay()
	g.updateChart()
}

func (g *GUI) toggleSimulation() {
	if g.simulation == nil {
		return
	}

	if g.simulation.running {
		g.pauseSimulation()
	} else {
		g.startSimulation()
	}
}

func (g *GUI) startSimulation() {
	if g.simulation == nil || g.simulation.running {
		return
	}

	g.simulation.running = true
	g.startButton.SetText("⏸ Pauza")

	g.simulation.ticker = time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			select {
			case <-g.simulation.ticker.C:
				if g.simulation.running {
					g.stepSimulation()
				}
			case <-g.simulation.stopCh:
				return
			}
		}
	}()
}

func (g *GUI) pauseSimulation() {
	if g.simulation == nil || !g.simulation.running {
		return
	}

	g.simulation.running = false
	g.startButton.SetText("▶ Start")

	if g.simulation.ticker != nil {
		g.simulation.ticker.Stop()
	}
}

func (g *GUI) resetSimulation() {
	g.pauseSimulation()
	g.createWorld()
}

func (g *GUI) stepSimulation() {
	if g.world == nil {
		return
	}

	g.world.Simulate()
	g.updateDisplay()
	g.updateChart()
	if g.world.IsExtinct() {
		g.pauseSimulation()
	}
}

func (g *GUI) updateDisplay() {
	if g.world == nil {
		return
	}
	gridText := ""
	for y := 0; y < g.world.Height; y++ {
		for x := 0; x < g.world.Width; x++ {
			if g.world.Grid[y][x] != nil {
				gridText += g.world.Grid[y][x].GetIcon() + " "
			} else {
				gridText += "⬜"
			}
		}
		gridText += "\n"
	}

	g.gridWidget.ParseMarkdown("```\n" + gridText + "```")
	stats := g.world.GetStatistics()
	g.turnLabel.SetText(fmt.Sprintf("Tura: %d", g.world.Turn))
	g.statsLabel.SetText(fmt.Sprintf("Populacja:\n🦊 Lisy: %d\n🐰 Króliki: %d\n🌱 Trawa: %d\nRazem: %d",
		stats["Fox"], stats["Rabbit"], stats["Grass"],
		stats["Fox"]+stats["Rabbit"]+stats["Grass"]))
	g.turnData = append(g.turnData, float64(g.world.Turn))
	g.foxData = append(g.foxData, float64(stats["Fox"]))
	g.rabbitData = append(g.rabbitData, float64(stats["Rabbit"]))
	g.grassData = append(g.grassData, float64(stats["Grass"]))
	if len(g.turnData) > 50 {
		g.turnData = g.turnData[1:]
		g.foxData = g.foxData[1:]
		g.rabbitData = g.rabbitData[1:]
		g.grassData = g.grassData[1:]
	}
}

func (g *GUI) updateChart() {
	if len(g.turnData) < 1 {
		g.chartImage.SetResource(nil)
		return
	}
	p := plot.New()
	p.Title.Text = "Populacja w czasie"
	p.X.Label.Text = "Tura"
	p.Y.Label.Text = "Liczba organizmów"

	foxPoints := make(plotter.XYs, len(g.turnData))
	rabbitPoints := make(plotter.XYs, len(g.turnData))
	grassPoints := make(plotter.XYs, len(g.turnData))

	for i := range g.turnData {
		foxPoints[i].X = g.turnData[i]
		foxPoints[i].Y = g.foxData[i]
		rabbitPoints[i].X = g.turnData[i]
		rabbitPoints[i].Y = g.rabbitData[i]
		grassPoints[i].X = g.turnData[i]
		grassPoints[i].Y = g.grassData[i]
	}

	if len(g.turnData) >= 2 {
		foxLine, _ := plotter.NewLine(foxPoints)
		foxLine.Color = color.RGBA{R: 255, G: 100, B: 0, A: 255}
		foxLine.Width = vg.Points(2)

		rabbitLine, _ := plotter.NewLine(rabbitPoints)
		rabbitLine.Color = color.RGBA{R: 139, G: 69, B: 19, A: 255}
		rabbitLine.Width = vg.Points(2)

		grassLine, _ := plotter.NewLine(grassPoints)
		grassLine.Color = color.RGBA{R: 0, G: 128, B: 0, A: 255}
		grassLine.Width = vg.Points(2)

		p.Add(foxLine, rabbitLine, grassLine)
		p.Legend.Add("🦊 Lisy", foxLine)
		p.Legend.Add("🐰 Króliki", rabbitLine)
		p.Legend.Add("🌱 Trawa", grassLine)
	} else {
		foxScatter, _ := plotter.NewScatter(foxPoints)
		foxScatter.Color = color.RGBA{R: 255, G: 100, B: 0, A: 255}

		rabbitScatter, _ := plotter.NewScatter(rabbitPoints)
		rabbitScatter.Color = color.RGBA{R: 139, G: 69, B: 19, A: 255}

		grassScatter, _ := plotter.NewScatter(grassPoints)
		grassScatter.Color = color.RGBA{R: 0, G: 128, B: 0, A: 255}

		p.Add(foxScatter, rabbitScatter, grassScatter)
		p.Legend.Add("🦊 Lisy", foxScatter)
		p.Legend.Add("🐰 Króliki", rabbitScatter)
		p.Legend.Add("🌱 Trawa", grassScatter)
	}
	img := vgimg.New(vg.Points(1200), vg.Points(900))
	dc := draw.New(img)
	p.Draw(dc)
	var buf bytes.Buffer
	png.Encode(&buf, img.Image())
	resource := fyne.NewStaticResource("chart.png", buf.Bytes())
	g.chartImage.SetResource(resource)
}

func (g *GUI) resetToDefaults() {
	g.createWorld()
}

func (g *GUI) Run() {
	g.window.ShowAndRun()
}
