package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/signintech/gopdf"
)

func main() {
	// 1. Открываем или создаём файл
	f, err := excelize.OpenFile("report.xlsx")
	if err != nil {
		f = excelize.NewFile()
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// 2. Определяем базовое имя листа и проверяем его наличие
	baseName := "СлужебнаяИнформация"
	sheetName := baseName
	sheets := f.GetSheetList()

	// Проверяем, существует ли уже лист с таким именем
	sheetExists := false
	for _, name := range sheets {
		if strings.EqualFold(name, baseName) { // Регистронезависимая проверка
			sheetExists = true
			break
		}
	}

	if sheetExists {
		// Если лист существует, генерируем новое имя с порядковым номером
		for counter := 1; ; counter++ {
			candidateName := fmt.Sprintf("%s_%d", baseName, counter)
			found := false
			for _, name := range sheets {
				if strings.EqualFold(name, candidateName) {
					found = true
					break
				}
			}
			if !found {
				sheetName = candidateName
				break
			}
		}
	}

	// 3. Добавляем новый лист и помещаем его в начало
	if _, err := f.NewSheet(sheetName); err != nil {
		log.Fatalf("Ошибка добавления листа: %v", err)
	}

	// Перемещаем лист на первую позицию (логика для пустой и непустой книги)
	sheets = f.GetSheetList()
	if len(sheets) > 1 {
		// Если в книге уже были листы, первый в списке — старый первый лист
		firstSheetName := sheets[0]
		if strings.EqualFold(firstSheetName, sheetName) {
			firstSheetName = sheets[1] // корректировка, если новый лист уже на первом месте
		}
		if err := f.MoveSheet(sheetName, firstSheetName); err != nil {
			log.Fatalf("Ошибка перемещения листа: %v", err)
		}
	}
	// Если книга была пустой, новый лист и так станет первым

	// 4. Делаем лист активным
	index, _ := f.GetSheetIndex(sheetName)
	f.SetActiveSheet(index)

	// 5. Заполняем лист данными и применяем стили
	// Устанавливаем ширину столбцов A и B
	f.SetColWidth(sheetName, "B", "B", 40)
	f.SetColWidth(sheetName, "F", "F", 70)

	// Создаем стиль для границы.
	// Поддерживаемые типы (Type): 1 - тонкая, 2 - средняя, 5 - жирная, 6 - двойная и др.
	styleID, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "bottom", // Направление: top, bottom, left, right
				Color: "000000", // Цвет в формате HEX (Черный)
				Style: 5,        // Стиль линии: 2 — средняя)
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Применяем стиль к ячейке.
	// Граница "bottom" у ячейки B2 визуально станет линией между B2 и B3.
	err = f.SetCellStyle(sheetName, "F1", "F1", styleID)
	if err != nil {
		log.Fatal(err)
	}
	//
	f.SetCellValue(sheetName, "F1", "ТУТ ЧТО-ТО НАПИСАЛИ")
	f.SetCellValue(sheetName, "F2", "тут написали какой-то пункт и текст")

	customFill := excelize.Fill{
		Type:    "pattern",
		Color:   []string{"F2F2F2"}, // Светло-серый фон
		Pattern: 1,
	}

	// Создаём стиль для заголовков (жирный, цветной шрифт)
	headerStyleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "003366", // Тёмно-синий цвет в формате RRGGBB
		},
		Fill: customFill,
	})
	if err != nil {
		log.Fatalf("Ошибка создания стиля: %v", err)
	}

	// Создаём стиль для обычных данных
	dataStyleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  14,
			Color: "000000", // Чёрный цвет
		},
	})
	if err != nil {
		log.Fatalf("Ошибка создания стиля данных: %v", err)
	}

	// Создаём стиль для пункта
	punktStyleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "000000", // Чёрный цвет
		},
	})
	if err != nil {
		log.Fatalf("Ошибка создания стиля для пункта: %v", err)
	}
	/*
		// Создаём стиль для заливки
		fillStyleID, err := f.NewStyle(&excelize.Style{
			Fill: customFill,
		})
		if err != nil {
			log.Fatalf("Ошибка создания стиля для заливки: %v", err)
		}
	*/
	// Записываем заголовки с датой
	now := time.Now().Format("2006-01-02 15:04:05")
	f.SetCellValue(sheetName, "B6", "Файл скачал:")
	f.SetCellValue(sheetName, "C6", "petrov_p")
	f.SetCellValue(sheetName, "B7", "Дата скачивания:")
	f.SetCellValue(sheetName, "C7", now)

	// Применяем стили к ячейкам
	f.SetCellStyle(sheetName, "B6", "B7", headerStyleID) // Заголовки в столбце И
	f.SetCellStyle(sheetName, "C6", "C7", dataStyleID)   // Данные в столбце С

	f.SetCellStyle(sheetName, "F2", "F2", punktStyleID) // Данные в столбце F2

	//	f.SetCellStyle(sheetName, "A1", "F60", fillStyleID) // Данные в столбце A1:F60

	// 6. Добавляем мета-информацию
	customProps := []excelize.CustomProperty{
		{Name: "Файл загружен пользователем", Value: "petrov_p"},
		{Name: "DownloadedBy", Value: "Иванов Иван Иванович"},
		{Name: "Файл загружен в", Value: now},
		{Name: "SourceFile", Value: "report.xlsx"},
	}

	for _, prop := range customProps {
		if err := f.SetCustomProps(prop); err != nil {
			log.Fatalf("Ошибка установки свойства %s: %v", prop.Name, err)
		}
	}

	// 7. Защита книги от структурных изменений
	if err := f.ProtectWorkbook(&excelize.WorkbookProtectionOptions{
		Password:      "strong_password_123",
		LockStructure: true,
	}); err != nil {
		log.Fatalf("Ошибка защиты книги: %v", err)
	}

	// 8. Защита всех листов от редактирования (включая новый служебный)
	for _, name := range f.GetSheetList() {
		if err := f.ProtectSheet(name, &excelize.SheetProtectionOptions{
			Password:            "strong_password_123",
			SelectLockedCells:   true,  // Разрешаем только выделять защищённые ячейки
			SelectUnlockedCells: false, // Запрещаем выделять незащищённые
		}); err != nil {
			log.Fatalf("Ошибка защиты листа '%s': %v", name, err)
		}
	}

	// 9. Сохраняем файл (без шифрования, чтобы он открывался)
	if err := f.SaveAs("report_protected.xlsx"); err != nil {
		log.Fatalf("Ошибка сохранения: %v", err)
	}

	fmt.Printf("Файл успешно обработан и сохранён как report_protected.xlsx\n")

	// PDF

	originalFile := "original.pdf"  // Ваш существующий файл
	tempMetaFile := "temp_meta.pdf" // Временный файл со служебной инфой
	outputFile := "result.pdf"      // Итоговый объединенный файл

	// --- ШАГ 1: Генерируем временный PDF с одной служебной страницей ---
	createMetaPage(tempMetaFile)
	// Гарантируем удаление временного файла после работы программы
	defer os.Remove(tempMetaFile)

	// --- ШАГ 2: Объединяем существующий файл и служебную страницу ---
	// pdfcpu.Merge принимает массив путей к файлам и склеивает их по порядку
	filesToMerge := []string{tempMetaFile, originalFile}

	err = api.MergeCreateFile(filesToMerge, outputFile, false, nil)
	if err != nil {
		log.Fatalf("Ошибка при объединении PDF файлов: %v", err)
	}

	fmt.Printf("Страница успешно добавлена! Итоговый файл: %s\n", outputFile)

}

// Функция генерации служебной страницы (аналогична предыдущему шагу)
func createMetaPage(filename string) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Важно: arial.ttf должен лежать рядом или укажите к нему путь
	err := pdf.AddTTFFont("Arial", "./arial.ttf")
	if err != nil {
		log.Fatalf("Шрифт не найден: %v", err)
	}

	// Заголовок
	_ = pdf.SetFont("Arial", "", 16)
	err = pdf.Cell(&gopdf.Rect{W: 200, H: 20}, "СЛУЖЕБНАЯ ИНФОРМАЦИЯ")
	if err != nil {
		log.Fatalf("Ошибка записи: %v", err)
	}
	pdf.Br(25)

	// Контент
	_ = pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(120, 120, 120)

	metaData := []string{
		fmt.Sprintf("Дата добавления: %s", time.Now().Format("2006-01-02 15:04:05")),
		"Статус документа: Архивы / Проверено",
		"Системный тег: META-APPEND-CONFIDENTIAL",
	}

	for _, line := range metaData {
		_ = pdf.Cell(nil, line)
		pdf.Br(18)
	}

	// Сохраняем временный файл
	err = pdf.WritePdf(filename)
	if err != nil {
		log.Fatalf("Не удалось создать временный файл: %v", err)
	}
}
