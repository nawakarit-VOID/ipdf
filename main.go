// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"bytes"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type FileStatus struct {
	Name   string
	Status string
}

var fileStatus []FileStatus

type Job struct {
	index int
	path  string
}

type Img struct {
	index int
	img   image.Image
}

type Encoded struct {
	index int
	buf   *bytes.Buffer
	w     float64
	h     float64
}

var files []string  // สร้างตัวแปร global สำหรับเก็บรายชื่อไฟล์ภาพที่ถูกโหลดเข้ามา และสถานะการทำงานของแต่ละไฟล์
var lastScroll = -1 //ตัวแปรสำหรับเก็บ index ของไฟล์ที่ถูกเลื่อนดูล่าสุดใน list เพื่อให้ UI เลื่อนไปที่ไฟล์นั้นเมื่อมีการอัปเดตสถานะการทำงาน

var jpegPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func updateStatus(index int, text string, list *widget.List) {

	fyne.Do(func() {

		fileStatus[index].Status = text

		list.Refresh()

		// เลื่อนไปที่ไฟล์ที่กำลังทำงาน
		list.ScrollTo(index)
		lastScroll = index

	})

}

// โหลด icon
func loadIcon(size int) fyne.Resource {
	var file string

	switch {
	case size >= 512:
		file = "icons/icon-512.png" ///ที่อยู่
	case size >= 256:
		file = "icons/icon-256.png"
	case size >= 128:
		file = "icons/icon-128.png"
	default:
		file = "icons/icon-64.png"
	}

	data, _ := iconFS.ReadFile(file)
	return fyne.NewStaticResource(file, data)
}

//go:embed icons/*
var iconFS embed.FS

//   //go:embed assets/background/bgW.png
//var bgW []byte

//var resLight = fyne.NewStaticResource("bgW.png", bgW)

//go:embed assets/font/Itim-Regular.ttf
var fontItim []byte
var myFont = fyne.NewStaticResource("Itim-Regular.ttf", fontItim)

var overlayW = color.NRGBA{250, 0, 0, 80}
var overlayB = color.NRGBA{0, 0, 0, 80}

//go:embed assets/lang/English.json
var enJSON []byte

//go:embed assets/lang/THAI.json
var thJSON []byte

// ฟังก์ชันสำหรับอัปเดตข้อความแสดงความเร็ว CPU//////////////////////////////////////////////////////////////////////////////////////
func main() {

	a := app.NewWithID("com.nawakarit.ipdf")
	icon := loadIcon(64) //เอา data มาใช้
	w := a.NewWindow("ipdf")
	w.SetIcon(icon)
	a.Settings().SetTheme(&MyTheme{})

	//bg := canvas.NewImageFromResource(resLight)
	//bg.FillMode = canvas.ImageFillCover

	/*
		Stretch → ยืดเต็ม (อาจเบี้ยว)
		Contain → ไม่เบี้ยว แต่มีขอบ
		Cover → เต็มจอ + ไม่เบี้ยว
	*/
	// ============================================================================
	// เปลี่ยนภาษา
	// ============================================================================
	var en map[string]string
	var th map[string]string
	json.Unmarshal(enJSON, &en)
	json.Unmarshal(thJSON, &th)

	// create i18n
	tr := NewI18n(en, th)

	// 🔥 language select // เปลี่ยนภาษา
	langSelect := NewSelect(tr, []string{"en", "th"}, func(val string) {
		tr.SetLang(val)
	})
	langSelect.SetSelected("en")

	// UI

	// สร้าง progress bar และ label สำหรับแสดงสถานะการทำงาน
	progress := widget.NewProgressBar()
	progress.SetValue(0)

	status := NewLabel(tr, "No images")

	// ============================================================================
	// list widget
	// ============================================================================
	// สร้าง list widget สำหรับแสดงชื่อไฟล์ภาพและสถานะการทำงานของแต่ละไฟล์ โดยใช้ข้อมูลจาก fileStatus
	// ซึ่งเป็น slice ของ FileStatus struct ที่เก็บชื่อไฟล์และสถานะการทำงานของแต่ละไฟล์

	fileList := widget.NewList(

		func() int {
			return len(fileStatus)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {

			fs := fileStatus[i]

			o.(*widget.Label).SetText(
				fmt.Sprintf("%03d  %-25s %s",
					i+1,
					fs.Name,
					fs.Status,
				),
			)

		},
	)

	fileListContainer := container.NewVScroll(fileList)

	fileListContainer.SetMinSize(fyne.NewSize(250, 250))

	maxCPU := runtime.NumCPU() //จำนวน CPU สูงสุดของเครื่องที่สามารถใช้ได้ (เช่น 4, 8, 16 cores)

	// ============================================================================
	// เลือกแฟ้ม
	// ============================================================================
	selectBtn := NewButton(tr, "Select Folder", func() {

		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {

			//if uri == nil {
			//	return
			//}
			if err != nil || uri == nil {
				return
			}

			files = nil

			list, err := os.ReadDir(uri.Path())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			validExt := map[string]bool{
				".jpg": true, ".jpeg": true, ".png": true,
				".webp": true, ".bmp": true, ".tiff": true,
			}

			for _, f := range list {
				ext := strings.ToLower(filepath.Ext(f.Name()))
				if validExt[ext] {
					files = append(files, filepath.Join(uri.Path(), f.Name()))
				}
			}

			sort.Strings(files)

			fileStatus = nil

			fmt.Println("FILES COUNT:", len(files))
			for _, f := range files {
				fmt.Println("FILE:", f)

				fileStatus = append(fileStatus, FileStatus{
					Name:   filepath.Base(f),
					Status: "🎨 image",
				})

			}

			fileList.Refresh()
			fmt.Println("files count:", len(files))
			status.SetText(fmt.Sprintf("Loaded %d images", len(files)))

		}, w)

		fd.Resize(fyne.NewSize(800, 600))
		if l, err := storage.ListerForURI(storage.NewFileURI("/media")); err == nil {
			fd.SetLocation(l)
		} else if l, err := storage.ListerForURI(storage.NewFileURI("/mnt")); err == nil {
			fd.SetLocation(l)
		}
		fd.Show()

	})
	//selectBtn.Importance = widget.HighImportance //ตั้งค่าความสำคัญของปุ่มเป็น High เพื่อให้มีสีและดูโดดเด่นมากขึ้น

	//ปุ่มเคลียร์รายการภาพที่โหลดเข้ามา
	clearBtn := NewButton(tr, "Clear", func() {

		files = nil
		fileStatus = nil

		progress.SetValue(0)

		status.SetText("No images")

	})
	//clearBtn.Importance = widget.DangerImportance //ตั้งค่าความสำคัญของปุ่มเป็น Danger เพื่อให้มีสีแดงและดูโดดเด่นมากขึ้น

	//ปุ่มเริ่มแปลง
	convertBtn := NewButton(tr, "Convert", func() {

		if len(files) == 0 {
			status.SetText("No images")
			return
		}

		save := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {

			if uc == nil {
				return
			}

			path := uc.URI().Path()
			uc.Close()

			runtime.GOMAXPROCS(maxCPU)

			go StartPipeline(progress, status, maxCPU, path, fileList)

		}, w)

		save.Resize(fyne.NewSize(700, 600))
		save.SetFileName("output.pdf")
		save.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))

		save.Show()

	})

	abbtn := widget.NewButton("!", func() {
		dialog.ShowInformation("about", "ipdf v1.2.0\nGolang + fyne\n\nโปรแกรมแปลงรูปภาพเป็น pdf (เร็วมาก กินแรมน้อย)\nจริงๆสร้างมาเพื่ออ่านการ์ตูน รองรับ 1 - หลายพันรูป\nคู่มือการใช้ รวมถึงเวอร์ชันใหม่ๆเข้าไปดูได้ในที่ github ด้านล่าง\nคลิกเข้าไปที่แท็บ Repositories แล้วหาชื่อโปรแกรม\n\nBy nawakarit - เจช์ (วัดดงหมี)\nhttps://github.com/nawakarit-VOID\n© 2026", w)
	})

	// ============================================================================
	// จัดวาง UI
	// ============================================================================

	labelt := NewLabel(tr, "We believe we are the fastest.")
	labelt.Alignment = fyne.TextAlignCenter

	label := NewLabel(tr, "Arrange the images in the folder first.\nSupports .jpg, .jpeg, .png, .webp, .bmp, and .tiff files.")
	label.Alignment = fyne.TextAlignCenter

	selectBtn1 := container.NewGridWrap(fyne.NewSize(150, 35), selectBtn)
	clearBtn1 := container.NewGridWrap(fyne.NewSize(150, 35), clearBtn)
	convertBtn1 := container.NewGridWrap(fyne.NewSize(150, 35), convertBtn)

	input1 := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewCenter(container.NewHBox(selectBtn1, clearBtn1, convertBtn1)))

	prog := container.NewGridWrap(fyne.NewSize(395, 35), progress)
	TR := container.NewGridWrap(fyne.NewSize(49, 35), langSelect)
	abbtn1 := container.NewGridWrap(fyne.NewSize(10, 35), abbtn)

	ProgressTR := container.NewBorder(
		nil, nil, nil, nil,
		container.NewCenter(container.NewHBox(prog, TR, abbtn1)))

	top := container.NewVBox(

		labelt, label, input1, ProgressTR,
		container.NewCenter(status),
	)

	content := container.NewBorder(
		top,               //บน
		nil,               //ล่าง
		nil,               // ซ้าย
		nil,               //ขวา
		fileListContainer, // กลาง
	)

	ui := container.NewStack(
		//bg, // พื้นหลัง
		//overlay,
		content, // UI ด้านบน
	)

	//w.SetContent(container.NewPadded(ui))
	w.SetContent(ui)
	w.Resize(fyne.NewSize(500, 700))
	//w.SetFixedSize(true)
	w.ShowAndRun()
}
