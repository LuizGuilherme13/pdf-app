package extract

import (
	"bufio"
	"fmt"
	"my_app/models"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func Do(filePath, dest string) error {
	err := api.ExtractContentFile(filePath, dest, nil, nil)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dest)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		n1 := ext(files[i].Name())
		n2 := ext(files[j].Name())

		return n1 < n2
	})

	re := regexp.MustCompile(`\((.*?)\)\s*Tj`)
	decoder := charmap.ISO8859_1.NewDecoder()

	viagem := models.Viagem{}
	viagens := []models.Viagem{}

	for i, f := range files {
		lines := []string{}

		//* LIMPANDO LINHAS ****************************************************
		{
			file, err := os.Open(filepath.Join(dest, f.Name()))
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				matches := re.FindAllStringSubmatch(scanner.Text(), -1)

				for _, m := range matches {
					res, _, err := transform.String(decoder, m[1])
					if err != nil {

						return err
					}

					lines = append(lines, res)
				}

			}

			if err := scanner.Err(); err != nil {
				return err
			}

			if err := file.Close(); err != nil {
				return err
			}
		}

		//* SALVANDO .TXT ******************************************************
		{
			file, err := os.OpenFile(filepath.Join(dest, f.Name()), os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			defer file.Close()

			writer := bufio.NewWriter(file)

			for _, line := range lines {
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					return err
				}
			}

			if err := writer.Flush(); err != nil {
				return err
			}

			if err := file.Close(); err != nil {
				return err
			}
		}

		//* MONTANDO VIAGENS ***************************************************
		{
			file, err := os.Open(filepath.Join(dest, f.Name()))
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				if scanner.Text() == "Viagem" {
					scanner.Scan()
					viagem.Nr = scanner.Text()
				}

				if scanner.Text() == "DESPESAS DA VIAGEM" {
				Exit:
					for scanner.Scan() {
						if scanner.Text() == "DETALHAMENTO COMBUSTÍVEL - VEÍCULO" || scanner.Text() == "OBSERVAÇÃO" {
							if viagem.PedagioDespesa == "" {
								viagem.PedagioDespesa = "0,00"
							}
							break Exit
						}

						if scanner.Text() == `1179 / PEDAGIO DIVERSOS` {
							for range 3 {
								scanner.Scan()
							}

							viagem.PedagioDespesa = scanner.Text()
							break Exit
						}
					}
				}

				if scanner.Text() == `\(+\) Base Comissão` {
					scanner.Scan()

					viagem.BaseComissao = scanner.Text()
				}

				if scanner.Text() == `\(-\) Despesas Desc. Comissão` {
					scanner.Scan()

					viagem.DescBaseComissao = scanner.Text()
				}

				if scanner.Text() == `\(+\) Pedágio` {
					scanner.Scan()

					viagem.PedagioFrete = scanner.Text()
				}
			}
		}

		if (i+1)%2 == 0 {
			viagens = append(viagens, viagem)
			viagem = models.Viagem{}
		}
	}

	fmt.Println("Total de viagens:", len(viagens))

	//* Criando arquivo XLSX
	{
		f := excelize.NewFile()
		defer f.Close()

		sheetName := "VIAGENS"
		f.SetSheetName(f.GetSheetName(0), sheetName)

		headers := []string{"nr_viagem", "base_comissao", "desc_base_comissao", "pedagio_frete", "pedagio_despesa"}

		for i, h := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue(sheetName, cell, h)
		}

		for i, v := range viagens {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), v.Nr)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), v.BaseComissao)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), v.DescBaseComissao)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), v.PedagioFrete)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), v.PedagioDespesa)
		}

		f.AddTable(sheetName, &excelize.Table{Range: fmt.Sprintf("A1:E%d", len(viagens))})

		err = f.SaveAs(filepath.Join(dest, "viagens.xlsx"))
		if err != nil {
			return err
		}
	}

	return nil
}

func ext(fileName string) int {
	re := regexp.MustCompile(`Content_page_(\d+)\.txt`)

	matches := re.FindStringSubmatch(fileName)

	if len(matches) > 1 {
		num, err := strconv.Atoi(matches[1])
		if err == nil {
			return num
		}
	}

	return 0
}
