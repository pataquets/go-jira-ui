package jiraui

import (
	"fmt"
	"regexp"

	ui "gopkg.in/gizak/termui.v2"
)

type Search struct {
	command     string
	directionUp bool
	re          *regexp.Regexp
}

type BaseListPage struct {
	uiList        *ScrollableList
	cachedResults []string
	isPopulated   bool
	ActiveSearch  Search
}

func (p *BaseListPage) SetSearch(searchCommand string) {
	if len(searchCommand) < 2 {
		// must be '/a' minimum
		return
	}
	direction := []byte(searchCommand)[0]
	regex := "(?i)" + string([]byte(searchCommand)[1:])
	s := new(Search)
	s.command = searchCommand
	if direction == '?' {
		s.directionUp = true
	} else if direction == '/' {
		s.directionUp = false
	} else {
		// bad command
		return
	}
	if re, err := regexp.Compile(regex); err != nil {
		// bad regex
		return
	} else {
		s.re = re
		p.ActiveSearch = *s
	}
}

func (p *BaseListPage) IsPopulated() bool {
	if len(p.cachedResults) > 0 || p.isPopulated {
		return true
	} else {
		return false
	}
}

func (p *BaseListPage) PreviousLine(n int) {
	p.uiList.CursorUpLines(n)
}

func (p *BaseListPage) NextLine(n int) {
	p.uiList.CursorDownLines(n)
}

func (p *BaseListPage) PreviousPara() {
	p.PreviousLine(5)
}

func (p *BaseListPage) NextPara() {
	p.NextLine(5)
}

func (p *BaseListPage) PreviousPage() {
	p.uiList.PageUp()
}

func (p *BaseListPage) NextPage() {
	p.uiList.PageDown()
}

func (p *BaseListPage) PageLines() int {
	return p.uiList.Height - 2
}

func (p *BaseListPage) TopOfPage() {
	p.uiList.Cursor = 0
	p.uiList.ScrollToTop()
}

func (p *BaseListPage) BottomOfPage() {
	p.uiList.Cursor = len(p.uiList.Items) - 1
	p.uiList.ScrollToBottom()
}

func (p *BaseListPage) Id() string {
	return fmt.Sprintf("BaseListPage(%p)", p)
}

func (p *BaseListPage) Update() {
	log.Debugf("BaseListPage.Update(): self:        %s (%p)", p.Id(), p)
	log.Debugf("BaseListPage.Update(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	ui.Render(p.uiList)
}

func (p *BaseListPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	changePage()
	q.Create()
}

func (p *BaseListPage) Create() {
	log.Debugf("BaseListPage.Create(): self:        %s (%p)", p.Id(), p)
	log.Debugf("BaseListPage.Create(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	ui.Clear()
	ls := NewScrollableList()
	p.uiList = ls
	p.cachedResults = make([]string, 0)
	ls.Items = p.cachedResults
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Updating, please wait"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
