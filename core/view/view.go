package view

import (
	"fmt"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var bounds = []float64{0, 10, 50, 100, 500, 1000}

var (
	headStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#dff0a3b9"))
	labelStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFF"))
	valueStyle         = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("86"))
	speedStyle         = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#8e8d8d"))
	speedDivisionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#a5d6ff"))
)

type speedDivision map[int]string

var speed_divisions_view = speedDivision{
	0:  "[                    ]",
	1:  "[■■                  ]",
	2:  "[■■■■                ]",
	3:  "[■■■■■■              ]",
	4:  "[■■■■■■■■            ]",
	5:  "[■■■■■■■■■■          ]",
	6:  "[■■■■■■■■■■■■        ]",
	7:  "[■■■■■■■■■■■■■■      ]",
	8:  "[■■■■■■■■■■■■■■■■    ]",
	9:  "[■■■■■■■■■■■■■■■■■■  ]",
	10: "[■■■■■■■■■■■■■■■■■■■■]",
}

func (s speedDivision) View(n int) string {
	v, exists := s[n]
	if !exists {
		return s[0]
	}

	return v
}

type tickMsg struct{}

type ViewModel struct {
	getCur func() float64
	getAvg func() float64
	tickD  time.Duration
	//
	avg    float64
	cur    float64
	width  int
	height int
}

func NewViewModel(
	getCurr, getAvg func() float64,
	tickD time.Duration,
) *ViewModel {
	return &ViewModel{
		getCur: getCurr,
		getAvg: getAvg,
		tickD:  tickD,
	}
}

func (m *ViewModel) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *ViewModel) Init() tea.Cmd {
	return m.tick()
}

func (m *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tickMsg:
		m.cur = m.getCur()
		m.avg = m.getAvg()

		return m, m.tick()

	case tea.KeyMsg:
		str := msg.String()
		if str == "q" || str == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *ViewModel) level(n float64) int {
	return sort.SearchFloat64s(bounds, n)
}

func (m *ViewModel) View() string {
	s := fmt.Sprintf("\n%s\n\n%s %s %s \n%s\n\n%s %s %s\n%s\n\n",
		headStyle.Render("Internet speed test"),

		labelStyle.Render("Current:"),
		valueStyle.Render(fmt.Sprintf("%.3f", m.cur)),
		speedStyle.Render("Mbit/s"),
		speedDivisionStyle.Render(speed_divisions_view.View(m.level(m.cur))),

		labelStyle.Render("Average:"),
		valueStyle.Render(fmt.Sprintf("%.3f", m.avg)),
		speedStyle.Render("Mbit/s"),
		speedDivisionStyle.Render(speed_divisions_view.View(m.level(m.avg))),
	)
	return s
}
