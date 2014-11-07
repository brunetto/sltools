package main

import (
	"log"
	"path/filepath"
	"strconv"
	
	"github.com/brunetto/gowut/gwu"
	"bitbucket.org/brunetto/slpp/sla"
)


func main() {
	var ( 
	err error
	)
	
	// Create and build a window
	win := gwu.NewWindow("slpp-ui", "UI for slpp")
	win.Style().SetFullWidth()
// 	win.SetHAlign(gwu.HA_CENTER)
	win.SetCellPadding(2)
	
	
	cboxPanel := gwu.NewHorizontalPanel()
	cboxLabel := gwu.NewLabel("")
	cbox := gwu.NewCheckBox("")
	cbox.AddEHandlerFunc(func(e gwu.Event) {
		if cbox.State() {
			cbox.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD)
		} else {
			cbox.Style().SetFontWeight(gwu.FONT_WEIGHT_NORMAL)
		}
		e.MarkDirty(cbox)
		cboxLabel.SetText("Display <property one>, state is " + strconv.FormatBool(cbox.State()))
		e.MarkDirty(cboxLabel)
	}, gwu.ETYPE_CLICK)
	cboxPanel.Add(cbox)
	cboxPanel.AddHSpace(3)
	cboxPanel.Add(cboxLabel)
	win.Add(cboxPanel)
	
	
	
	
	inFilePanel := gwu.NewHorizontalPanel()
	inFilePanel.Add(gwu.NewLabel("Input file:"))
	inFileBox := gwu.NewTextBox("../files/all_the_fishes.txt") // string in inFileBox.Text()
	// TODO 1: gray text, deleted if/when I enter the path that should be black
	// TODO 2: deactivate textbox when load
	// TODO 3: discard loaded -> activate textbox
	inFileBox.AddSyncOnETypes(gwu.ETYPE_KEY_UP)
	inFilePanel.Add(inFileBox)
	win.Add(inFilePanel)
	
	loadPanel := gwu.NewHorizontalPanel()
	loadedLabel := gwu.NewLabel("")
	populateButton := gwu.NewButton("Load")
	populateButton.AddEHandlerFunc(func(e gwu.Event) {
		inPath := filepath.Dir(inFileBox.Text())
		inFile := filepath.Base(inFileBox.Text())
		starMap := make(sla.StarMapType)
		starMap.Populate(inPath, inFile)
		loadedLabel.SetText("Done!")
		e.MarkDirty(loadedLabel)
		//TODO:  goroutine with chanel to know the parsing %
	}, gwu.ETYPE_CLICK)
	loadPanel.Add(populateButton)
	loadPanel.Add(loadedLabel)
	win.Add(loadPanel)
	
	type wrt struct {
		Str string
	}
	/*
	func (w wrt) Write (p []byte) (n int, err error) {
		w.Str = p.String()
		return len(p), nil
	}
	
	func (w wrt) Get () (string) {
		return w.Str
	}	
	
	w := &wrt
	
	countExchPanel := gwu.NewHorizontalPanel()
	countExchLabel := gwu.NewLabel("")
	countExchButton := gwu.NewButton("Count Exchanges")
	populateButton.AddEHandlerFunc(func(e gwu.Event) {
		starMap.ExecOnAll("CountExchanges")
		starMap.PrintExchStats(&w)
		loadedLabel.SetText(w.Get())
		e.MarkDirty(loadedLabel)
	}, gwu.ETYPE_CLICK)
	loadPanel.Add(populateButton)
	loadPanel.Add(loadedLabel)
	win.Add(loadPanel)
	*/

	
	
	
	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("interface", "localhost:8081")//"localhost:8081")
	server.SetText("Hola!!")
	server.AddWin(win)
	if err = server.Start(""); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}





