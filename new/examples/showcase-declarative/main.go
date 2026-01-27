package main

import p "preveltekit"

func main() {
	println("main: starting")
	p.Hydrate(&App{
		CurrentPage: p.New("basics"),
	},
		p.WithChild("/basics", &Basics{
			Count:    p.New(0),
			Name:     p.New(""),
			Message:  p.New(""),
			DarkMode: p.New(false),
			Agreed:   p.New(false),
			Score:    p.New(0),
		}),
		p.WithChild("/components", &Components{
			Message:      p.New(""),
			ClickCount:   p.New(0),
			CardTitle:    p.New(""),
			AlertType:    p.New(""),
			AlertMessage: p.New(""),
		}),
		p.WithChild("/lists", &Lists{
			Items:     p.NewList[string](),
			NewItem:   p.New(""),
			ItemCount: p.New(0),
		}),
		p.WithChild("/routing", &Routing{
			CurrentTab:  p.New(""),
			CurrentStep: p.New(0),
		}),
		p.WithChild("/links", &Links{
			LastNavigation: p.New(""),
		}),
		p.WithChild("/fetch", &Fetch{
			Status:  p.New(""),
			RawData: p.New(""),
		}),
		p.WithChild("/storage", &Storage{
			Theme:  p.NewLocalStore("theme", ""),
			Notes:  p.New(""),
			Status: p.New(""),
		}),
		p.WithChild("/debounce", &Debounce{
			SearchInput:   p.New(""),
			SearchResult:  p.New(""),
			SearchCount:   p.New(0),
			ClickCount:    p.New(0),
			ThrottleCount: p.New(0),
			Status:        p.New(""),
		}),
		p.WithChild("/bitcoin", &Bitcoin{
			Price:      p.New(""),
			Symbol:     p.New(""),
			UpdateTime: p.New(""),
			Loading:    p.New(false),
			Error:      p.New(""),
		}),
		// Register nested components used by p.Comp()
		p.WithNestedComponent("Badge", func() p.Component {
			return &Badge{Label: p.New("")}
		}),
		p.WithNestedComponent("Card", func() p.Component {
			return &Card{Title: p.New("")}
		}),
		p.WithNestedComponent("Button", func() p.Component {
			return &Button{Label: p.New("")}
		}),
		p.WithNestedComponent("Alert", func() p.Component {
			return &Alert{Type: p.New(""), Message: p.New("")}
		}),
	)
}
