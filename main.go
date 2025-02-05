package main

import (
	"my_app/extract"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Extrator de PDF para XLSX")
	w.Resize(fyne.NewSize(500, 400))

	fileIcon := widget.NewFileIcon(nil)
	folderIcon := widget.NewIcon(theme.FolderIcon())

	//* Arquivo de Origem
	filePathURI := widget.NewLabel("...")
	filePathURI.Resize(fyne.NewSize(100, filePathURI.MinSize().Height))
	btnFile := widget.NewButton("Procurar", func() {
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

	v1 := container.NewVBox(
		container.NewHBox(widget.NewLabel("Arquivo PDF:"), layout.NewSpacer(), btnFile),
		container.NewHBox(fileIcon, filePathURI),
	)

	//* Local de destino
	destFileURI := widget.NewLabel("...")
	btnDestFile := widget.NewButton("Procurar", func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil {
				dialog.NewError(err, w).Show()
			}

			if lu != nil {
				destFileURI.SetText(lu.Path())

			}
		}, w)
	})

	v2 := container.NewVBox(
		container.NewHBox(widget.NewLabel("Local de Destino:"), layout.NewSpacer(), btnDestFile),
		container.NewHBox(folderIcon, destFileURI),
	)

	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()
	progressDesc := widget.NewLabel("Extraindo dados... Aguarde.")
	progressDesc.Hide()

	//* Botão Extrair
	btnExtract := widget.NewButton("Extrair", func() {
		progressBar.Show()
		progressDesc.Show()

		go func() {
			if err := extract.Do(filePathURI.Text[7:], destFileURI.Text); err != nil {
				dialog.NewError(err, w).Show()
			}
			progressBar.Hide()
			progressDesc.Hide()
			dialog.ShowInformation("Concluído", "Extração finalizada com sucesso!", w)
		}()
	})

	box := container.NewVBox(
		container.NewPadded(v1),
		container.NewPadded(widget.NewSeparator()),
		container.NewPadded(v2),
		container.NewPadded(widget.NewSeparator()),
		container.NewHBox(layout.NewSpacer(), btnExtract, layout.NewSpacer()),
		layout.NewSpacer(),
		container.NewPadded(
			container.NewVBox(
				progressBar,
				progressDesc,
			),
		),
		layout.NewSpacer(),
	)

	w.SetContent(box)

	w.ShowAndRun()
}
