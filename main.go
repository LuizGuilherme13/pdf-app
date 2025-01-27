package main

import (
	"my_app/extract"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Extrator de PDF para XLSX")
	w.Resize(fyne.NewSize(600, 500))

	fileIcon := widget.NewFileIcon(nil)
	filePathURI := widget.NewLabel("...")

	//* Arquivo de Origem
	btnFile := widget.NewButton("Selecione o Arquivo", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.NewError(err, w).Show()
			}

			if reader != nil {
				filePathURI.SetText(reader.URI().String())
				reader.Close()
			}

		}, w)
	})
	filePath := container.NewHBox(fileIcon, filePathURI, btnFile)

	//* Local de destino
	destFileURI := widget.NewLabel("...")
	btnDestFile := widget.NewButton("Local de Destino", func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil {
				dialog.NewError(err, w).Show()
			}

			if lu != nil {
				destFileURI.SetText(lu.Path())

			}
		}, w)
	})
	destPath := container.NewHBox(fileIcon, destFileURI, btnDestFile)

	//* Botão Extrair
	btnExtract := widget.NewButton("Extrair", func() {
		if err := extract.Do(filePathURI.Text[7:], destFileURI.Text); err != nil {
			dialog.NewError(err, w).Show()
		}
	})

	box := container.NewVBox(filePath, destPath, btnExtract)

	w.SetContent(box)

	w.ShowAndRun()
}
