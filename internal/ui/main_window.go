package ui

import (
	"fmt"
	"log"
	"todo_list/internal/data"

	"github.com/jroimartin/gocui"
)

type ViewData interface{}

type TodoListViewData struct {
	curLine int
}

var _ ViewData = (*TodoListViewData)(nil)

type MainWindow struct {
	repository data.Repository
	viewDatas  map[*gocui.View]ViewData
}

func NewMainWindow() *MainWindow {
	return &MainWindow{
		repository: data.CreateRepository(),
		viewDatas:  map[*gocui.View]ViewData{},
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
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorYellow
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

	g.SetCurrentView(VIEW_TITLE_TODOLIST)
	return nil

}
func (m *MainWindow) keyBindding(g *gocui.Gui) error {
	// g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
	// 	currentView := g.CurrentView()
	// 	if currentView == nil {
	// 		g.SetCurrentView(VIEW_TITLE_TODOLIST)
	// 		return nil
	// 	}
	// 	currentViewName := currentView.Name()
	// 	viewNames := []string{VIEW_TITLE_TODOLIST /* , VIEW_TITLE_DETAILS */}
	// 	index := slices.IndexFunc(viewNames, func(name string) bool { return name == currentViewName })
	// 	if index+1 >= len(viewNames) {
	// 		g.SetCurrentView(viewNames[0])
	// 	} else {
	// 		g.SetCurrentView(viewNames[index+1])
	// 	}

	// 	return nil
	// })
	g.SetKeybinding(VIEW_TITLE_TODOLIST, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			ox, oy := v.Origin()
			cx, cy := v.Cursor()
			if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}

			_, ncy := v.Cursor()
			if line, err := v.Line(ncy); err == nil {
				if view, err := g.View(VIEW_TITLE_DETAILS); err == nil && view != nil {
					view.Clear()
					fmt.Fprintln(view, line)
				}
			}

		}
		return nil
	})
	g.SetKeybinding(VIEW_TITLE_TODOLIST, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			cx, cy := v.Cursor()
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}

			_, ncy := v.Cursor()
			if line, err := v.Line(ncy); err == nil {
				if view, err := g.View(VIEW_TITLE_DETAILS); err == nil && view != nil {
					view.Clear()
					fmt.Fprintln(view, line)
				}
			}
		}
		return nil
	})
	return nil
}
