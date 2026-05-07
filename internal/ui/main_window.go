package ui

import (
	"fmt"
	"log"
	"slices"
	"todo_list/internal/data"

	"github.com/jroimartin/gocui"
)

type ViewData interface{}

type TodoListViewData struct {
	curLine int
}

var _ ViewData = (*TodoListViewData)(nil)

type State struct {
	currentItem int
	currentView string
}
type MainWindow struct {
	repository data.Repository
	viewDatas  map[*gocui.View]ViewData
	State
}

func NewMainWindow() *MainWindow {
	return &MainWindow{
		repository: data.CreateRepository(),
		viewDatas:  map[*gocui.View]ViewData{},
		State:      State{currentItem: 0, currentView: VIEW_TITLE_TODOLIST},
	}
}

func (m *MainWindow) ShowAndRun() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.Cursor = true
	g.SetManagerFunc(m.layout)
	m.setCurrentView(g, VIEW_TITLE_TODOLIST)
	m.keyBindding(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

const (
	VIEW_TITLE_TODOLIST = "TodoList"
	VIEW_TITLE_DETAILS  = "Details"
	VIEW_TITLE_HELP     = "Help"
	VIEW_TITLE_INPUT    = "Input"
)

func (m *MainWindow) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	oneOfX := maxX / 10
	// oneOfY := maxY / 10
	{
		if v, err := g.SetView(VIEW_TITLE_TODOLIST, 0, 0, maxX-oneOfX*4, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Title = VIEW_TITLE_TODOLIST
			v.Autoscroll = true
			v.BgColor = gocui.ColorDefault // 默认背景色
			v.FgColor = gocui.ColorWhite
			v.SelBgColor = gocui.ColorWhite
			v.SelFgColor = gocui.ColorBlack
			v.Highlight = true
			v.Wrap = true

			if todoList, err := m.repository.GetTodos(); err == nil {
				for _, todo := range todoList {
					fmt.Fprintln(v, todo.Content)
				}
			}
		}
	}
	{
		if v, err := g.SetView(VIEW_TITLE_DETAILS, maxX-oneOfX*4, 0, maxX-oneOfX*2, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = VIEW_TITLE_DETAILS
			fmt.Fprintln(v, "detail view")
		}
	}
	{
		if v, err := g.SetView(VIEW_TITLE_HELP, maxX-oneOfX*2, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = VIEW_TITLE_HELP
		}
	}
	m.setCurrentView(g, m.currentView)

	return nil

}
func (m *MainWindow) keyBindding(g *gocui.Gui) error {
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		currentView := g.CurrentView()
		if currentView == nil {
			m.setCurrentView(g, VIEW_TITLE_TODOLIST)
			return nil
		}
		currentViewName := currentView.Name()
		viewNames := []string{VIEW_TITLE_TODOLIST, VIEW_TITLE_DETAILS}
		index := slices.IndexFunc(viewNames, func(name string) bool { return name == currentViewName })
		if index == -1 {
			return nil
		}
		if index+1 >= len(viewNames) {
			m.setCurrentView(g, viewNames[0])
		} else {
			m.setCurrentView(g, viewNames[index+1])
		}

		return nil
	})
	g.SetKeybinding(VIEW_TITLE_TODOLIST, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			_, oy := v.Origin()
			cx, cy := v.Cursor()
			if m.currentItem-1 <= 0 {
				v.SetCursor(cx, oy)
				m.currentItem = 0
			} else {
				v.SetCursor(cx, cy-1)
				m.currentItem--
			}
			// if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			// 	if err := v.SetOrigin(ox, oy-1); err != nil {
			// 		return err
			// 	}
			// }

			// _, ncy := v.Cursor()
			// if line, err := v.Line(ncy); err == nil {
			cx, cy = v.Cursor()
			if view, err := g.View(VIEW_TITLE_DETAILS); err == nil && view != nil {
				view.Clear()
				fmt.Fprintln(view, cx, cy, m.currentItem)
			}
			// }

		}
		return nil
	})
	g.SetKeybinding(VIEW_TITLE_TODOLIST, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			// _, oy := v.Origin()
			cx, cy := v.Cursor()
			if m.currentItem+1 >= m.repository.Size() {
				return nil
				// v.SetCursor(cx, oy)
				// m.currentItem = 0
			} else {
				v.SetCursor(cx, cy+1)
				m.currentItem++
			}

			// _, ncy := v.Cursor()
			// if line, err := v.Line(ncy); err == nil {
			cx, cy = v.Cursor()
			if view, err := g.View(VIEW_TITLE_DETAILS); err == nil && view != nil {
				view.Clear()
				fmt.Fprintln(view, cx, cy, m.currentItem)
			}
			// }
		}
		return nil
	})
	g.SetKeybinding(VIEW_TITLE_TODOLIST, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		maxX, maxY := g.Size()
		halfX := maxX / 2
		halfY := maxY / 2
		if v, err := g.SetView(VIEW_TITLE_INPUT, halfX/2, halfY/2, maxX-halfX/2, halfY/2+2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = VIEW_TITLE_INPUT
			v.Editable = true
			v.Frame = true
			v.Clear()
			m.setCurrentView(g, VIEW_TITLE_INPUT)
			g.SetKeybinding(VIEW_TITLE_INPUT, gocui.KeyCtrlQ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				m.setCurrentView(g, VIEW_TITLE_TODOLIST)
				g.DeleteView(VIEW_TITLE_INPUT)
				return nil
			})
			g.SetKeybinding(VIEW_TITLE_INPUT, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				content := v.Buffer()
				// fmt.Println(content)
				m.repository.CreateAndAddTodo(content, false)

				if v, err := m.setCurrentView(g, VIEW_TITLE_TODOLIST); err == nil {
					fmt.Fprintln(v, content)
				}
				g.DeleteView(VIEW_TITLE_INPUT)
				return nil
			})
		}
		return nil
	})
	return nil
}
func (m *MainWindow) setCurrentView(g *gocui.Gui, view string) (*gocui.View, error) {
	v, err := g.SetCurrentView(view)
	if err != nil {
		return v, err
	}
	m.currentView = view
	if v, err := g.View(VIEW_TITLE_HELP); err == nil {
		v.Clear()
		for _, line := range m.helpText(view) {
			fmt.Fprintln(v, line)
		}
	}
	return v, nil
}
func (m *MainWindow) helpText(view string) []string {
	texts := []string{"global:", "\ttab: switch view", ""}
	switch view {
	case VIEW_TITLE_TODOLIST:
		texts = append(texts, fmt.Sprintf("%v:", VIEW_TITLE_TODOLIST),
			"\t↑: move up to select item",
			"\t↓: move down to select item",
			"\t⮐: add new todo",
		)
	case VIEW_TITLE_DETAILS:
		texts = append(texts, fmt.Sprintf("%v:", VIEW_TITLE_DETAILS), "\tenter: modify")
	case VIEW_TITLE_INPUT:
		texts = append(texts, fmt.Sprintf("%v:", VIEW_TITLE_INPUT), "\tctrl+q: close view", "\tenter: close view and keep content")

	}
	return texts
}
