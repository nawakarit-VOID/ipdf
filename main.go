package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
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
	"time"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"image/color"

	"github.com/jung-kurt/gofpdf"
	"github.com/nfnt/resize"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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

// ฟังก์ชันสำหรับคำนวณสัญลักษณ์ความเร็วตามเปอร์เซ็นต์การใช้ CPU
func SpeedSymbol(pvcpus float64) string {
	switch {
	case pvcpus >= 86.7:
		return "🛰️"
	case pvcpus >= 75.1:
		return "🚀"
	case pvcpus >= 62.6:
		return "✈️"
	case pvcpus >= 51.0:
		return "🚔"
	case pvcpus >= 37.6:
		return "🚍"
	case pvcpus >= 25.1:
		return "🏍️"
	case pvcpus >= 12.6:
		return "🚖"
	default:
		return "🐌"
	}
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

// ฟังก์ชันสำหรับอัปเดตข้อความแสดงความเร็ว CPU//////////////////////////////////////////////////////////////////////////////////////
func main() {

	a := app.NewWithID("com.nawakarit.oneimage")
	a.SetIcon(resourceIconPng) //ต้องปิดตอนเขียนโค้ด แก้โค้ด หรือทดสอบรัน
	w := a.NewWindow("Image → PDF Ultra Fast")
	w.SetIcon(resourceIconPng) //ต้องปิดตอนเขียนโค้ด

	//หน้าหลักของแอปพลิเคชัน โดยใช้ Fyne ในการสร้าง
	title := canvas.NewText(
		"🤖 Image → PDF Ultra Fast 🤖",
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	)
	title.TextSize = 34
	title.TextStyle = fyne.TextStyle{Bold: true}

	// สร้าง progress bar และ label สำหรับแสดงสถานะการทำงาน
	progress := widget.NewProgressBar()
	progress.SetValue(0)

	// สร้าง label สำหรับแสดงสถานะการทำงานของไฟล์ภาพที่ถูกโหลดเข้ามา
	status := widget.NewLabel("No images")

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
	fileListContainer.SetMinSize(fyne.NewSize(100, 250))

	// การตั้งค่าเริ่มต้นของการใช้ CPU และการแสดงสัญลักษณ์ความเร็วตามเปอร์เซ็นต์การใช้ CPU////////////////////////////////////////////////
	maxCPU := runtime.NumCPU()       //จำนวน CPU สูงสุดของเครื่องที่สามารถใช้ได้ (เช่น 4, 8, 16 cores)
	pmcpu := 100.0 / float64(maxCPU) //เปอร์เซ็นต์การใช้ CPU ต่อ 1 core (100% หารด้วยจำนวน core สูงสุดของเครื่อง)
	pvcpu := float64(100)            //เปอร์เซ็นต์การใช้ CPU เริ่มต้นที่ 100% (ใช้ทุก core ที่มี)
	symbol := SpeedSymbol(100)       //แสดงสัญลักษณ์ความเร็วตามเปอร์เซ็นต์การใช้ CPU เริ่มต้น

	// สร้างข้อความแสดงความเร็ว CPU โดยใช้ canvas.NewText เพื่อให้สามารถปรับแต่งสีและขนาดได้มากขึ้น
	cpuLabel := canvas.NewText(
		fmt.Sprintf("CPU Speed x%.1f %s ( %.0f%% / cores ) %s", float64(maxCPU), symbol, pvcpu, symbol),
		color.RGBA{R: 255, G: 255, B: 255, A: 255})
	cpuLabel.TextSize = 25                          // ปรับขนาดตามต้องการ
	cpuLabel.TextStyle = fyne.TextStyle{Bold: true} // ถ้าต้องการตัวหนา

	//slider///////////////////////////////////////////////////////////////////////////////////////////////////////////////
	cpuSlider := widget.NewSlider(1, float64(maxCPU)) // สร้าง slider สำหรับเลือกจำนวน CPU ที่จะใช้ โดยมีค่าตั้งแต่ 1 ถึงจำนวน CPU สูงสุดของเครื่อง
	cpuSlider.Step = 1                                //ใช้เฉพาะจำนวนเต็ม เพราะ workers และ parallelism ต้องเป็นจำนวนเต็ม
	cpuSlider.Value = float64(maxCPU)                 //ตั้งค่าเริ่มต้นของ slider ให้เป็นจำนวน CPU สูงสุด (ใช้ทุก core ที่มี)
	cpuSlider.OnChanged = func(v float64) {           //เมื่อ slider ถูกเปลี่ยนค่า จะคำนวณเปอร์เซ็นต์การใช้ CPU ใหม่และอัปเดตข้อความใน cpuLabel ตามค่าที่เลือก
		pvcpus := pmcpu * v
		symbol := SpeedSymbol(pvcpus) //แสดงสัญลักษณ์ความเร็วตามเปอร์เซ็นต์การใช้ CPU เริ่มต้น
		cpuLabel.Text = fmt.Sprintf("CPU Speed x%.1f %s ( %.0f%% / cores ) %s", v, symbol, pvcpus, symbol)
		cpuLabel.Refresh()
	}

	//ปุ่มเลือกโฟลเดอร์///////////////////////////////////////////////////////////////////////////////////////////////////////////
	selectBtn := widget.NewButton("📂 Select Folder", func() {

		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {

			//if uri == nil {
			//	return
			//}
			if err != nil || uri == nil {
				return
			}

			files = nil

			/*list, _ := os.ReadDir(uri.Path())

			if err != nil {
				dialog.ShowError(err, w)
				return
			}*/

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
	selectBtn.Importance = widget.HighImportance //ตั้งค่าความสำคัญของปุ่มเป็น High เพื่อให้มีสีและดูโดดเด่นมากขึ้น

	//ปุ่มเคลียร์รายการภาพที่โหลดเข้ามา
	clearBtn := widget.NewButton("🗑 Clear Images", func() {

		files = nil
		fileStatus = nil

		progress.SetValue(0)

		status.SetText("No images")

	})
	clearBtn.Importance = widget.DangerImportance //ตั้งค่าความสำคัญของปุ่มเป็น Danger เพื่อให้มีสีแดงและดูโดดเด่นมากขึ้น

	//ปุ่มเริ่มแปลง
	convertBtn := widget.NewButton("🔀 Convert to PDF", func() {

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

			cores := int(cpuSlider.Value)

			runtime.GOMAXPROCS(cores)

			go startPipeline(progress, status, cores, path, fileList)

		}, w)

		save.Resize(fyne.NewSize(700, 600))
		save.SetFileName("output.pdf")
		save.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))

		save.Show()

	})
	convertBtn.Importance = widget.SuccessImportance //ตั้งค่าความสำคัญของปุ่มเป็น Success เพื่อให้มีสีเขียวและดูโดดเด่นมากขึ้น

	//ปุ่มแสดงข้อมูลเกี่ยวกับแอปพลิเคชัน
	aboutBtn := widget.NewButton("About", func() {

		title := widget.NewLabelWithStyle("Image to PDF", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		version := widget.NewLabel("Version 18.0.1.0")
		info := widget.NewLabel("🤖 Go v1.26.1 ")
		date := widget.NewLabel("Done. 2 april 2026")
		by := widget.NewLabel("by.เจช์ (วัดดงหมี)")
		cc := widget.NewLabel("========..ความเร็วระดับนั้น..เราเรียกว่าพื้นฐาน..========")
		dot := widget.NewLabel("🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀🚀")
		copy := widget.NewLabel("Copyright © 6365")

		content := container.NewVBox(
			title,
			version,
			widget.NewSeparator(),
			container.NewCenter(info),
			container.NewCenter(date),
			container.NewCenter(by),
			container.NewCenter(cc),
			container.NewCenter(dot),
			container.NewCenter(copy),
		)

		dialog.ShowCustom(
			"About",
			"Close",
			content,
			w,
		)

	})

	//จัดวาง UI
	cpuCard := widget.NewCard(
		"⚙ CPU Settings",
		"CPU ทุกคอร์ทำงานตาม % ที่กำหนด (ไม่รวมกับที่ใช้อยู่ในปัจจุบัน)",
		container.NewVBox(
			container.NewCenter(cpuLabel),
			cpuSlider,
		),
	)

	inputCard := widget.NewCard(
		"📂 Input Images",
		"เลือกโฟลเดอร์รูปภาพ (เรียงลำดับภาพในแฟ้มก่อน รับไฟล์ .jpg, .jpeg, .png, .webp, .bmp, .tiff)",
		container.NewVBox(
			selectBtn,
			clearBtn,
		),
	)

	convertCard := widget.NewCard(
		"🔄️ Convert",
		"เริ่มแปลงภาพเป็น PDF (ขนาดลดลง 30%-50%) (คุณภาพ 85%)",
		container.NewVBox(
			convertBtn,
		),
	)

	progressCard := widget.NewCard(
		"📊 Progress",
		"",
		container.NewVBox(
			progress,
			status,
			fileListContainer,
		),
	)

	ui := container.NewVBox(

		container.NewCenter(title),
		cpuCard,
		inputCard,
		convertCard,
		progressCard,
		aboutBtn,
	)

	w.SetContent(container.NewPadded(ui))

	w.Resize(fyne.NewSize(800, 700))

	w.ShowAndRun()
}

// การเริ่มต้น pipeline การแปลงภาพเป็น PDF โดยใช้ workers หลายตัวในการประมวลผลภาพ
func startPipeline(
	progress *widget.ProgressBar,
	status *widget.Label,
	cores int,
	output string,
	fileList *widget.List,
) {

	start := time.Now()
	total := len(files)

	jobs := make(chan Job, cores*4)
	decoded := make(chan Img, cores*4)
	resized := make(chan Img, cores*4)
	encoded := make(chan Encoded, cores*4)

	decodeWorkers := cores * 2 //workers ที่ทำหน้าที่ decode พร้อมกันในเวลาเดียวกัน (ใช้มากกว่า cores เพื่อให้แน่ใจว่า CPU ไม่ว่างระหว่างรอ I/O)
	resizeWorkers := cores     //workers ที่ทำหน้าที่ resize พร้อมกันในเวลาเดียวกัน (ใช้เท่ากับ cores เพราะเป็นงานที่ใช้ CPU เป็นหลัก)
	encodeWorkers := cores     //workers ที่ทำหน้าที่ encode พร้อมกันในเวลาเดียวกัน (ใช้เท่ากับ cores เพราะเป็นงานที่ใช้ CPU เป็นหลัก)
	// เพราะเขาใช้ goroutine โดยใช้ go funcแยกกัน 3 กลุ่ม คือ decode, resize, encode ซึ่งแต่ละกลุ่มมีจำนวน workers
	// ที่แตกต่างกันตามลักษณะงานที่ทำ โดยใช้ sync.WaitGroup เพื่อรอให้ทุก worker ในแต่ละกลุ่มทำงานเสร็จ ก่อนที่จะปิด channel และเริ่มเขียน PDF

	var wgDecode sync.WaitGroup
	var wgResize sync.WaitGroup
	var wgEncode sync.WaitGroup
	// ---------- decode workers ----------
	for i := 0; i < decodeWorkers; i++ {

		wgDecode.Add(1)

		go func() {

			//สร้าง worker สำหรับการ decode ภาพจากไฟล์ไปเป็น image.Image
			// โดยใช้ goroutine และ sync.WaitGroup เพื่อรอให้ทุก worker ทำงานเสร็จ
			defer wgDecode.Done()

			for j := range jobs {

				f, err := os.Open(j.path)
				if err != nil {
					fmt.Println("❌ decode fail:", j.path, err)
					updateStatus(j.index, "❌ decode error", fileList)
					continue
				}

				img, _, err := image.Decode(f)
				f.Close()
				if err != nil {
					continue
				}
				decoded <- Img{
					index: j.index,
					img:   img,
				}
				updateStatus(j.index, "🔀 decoding", fileList)
			}
		}()
	}

	// ---------- resize workers ----------
	for i := 0; i < resizeWorkers; i++ {
		wgResize.Add(1)

		go func() {

			defer wgResize.Done()

			for im := range decoded {

				b := im.img.Bounds()

				if b.Dx() > 2480 {

					im.img = resize.Resize(
						2480,
						0,
						im.img,
						resize.Bilinear,
					)

				}

				resized <- im
				updateStatus(im.index, "↔️ resizing", fileList)
			}

		}()
	}

	// ---------- encode workers ----------
	for i := 0; i < encodeWorkers; i++ {

		wgEncode.Add(1)

		go func() {

			defer wgEncode.Done()

			for im := range resized {

				buf := jpegPool.Get().(*bytes.Buffer)

				buf.Reset()

				jpeg.Encode(buf, im.img, &jpeg.Options{
					Quality: 85,
				})

				b := im.img.Bounds()

				encoded <- Encoded{
					index: im.index,
					buf:   buf,
					w:     float64(b.Dx()) * 0.264583,
					h:     float64(b.Dy()) * 0.264583,
				}

				updateStatus(im.index, "🔄 encoding", fileList)

			}

		}()
	}

	// ---------- feed jobs ----------
	go func() {

		for i, f := range files {

			jobs <- Job{
				index: i,
				path:  f,
			}

		}

		close(jobs)

	}()

	// ---------- close channels ----------

	go func() {
		wgDecode.Wait()
		close(decoded)
	}()

	go func() {
		wgResize.Wait()
		close(resized)
	}()

	go func() {
		wgEncode.Wait()
		close(encoded)
	}()

	writePDF(encoded, total, progress, status, start, output, fileList)

}

// ฟังก์ชันสำหรับการเขียนไฟล์ PDF โดยรับข้อมูลภาพที่ถูก encode แล้วจาก channel
// และจัดเรียงตามลำดับ index เพื่อให้ภาพอยู่ในลำดับที่ถูกต้องใน PDF จากนั้นใช้ gofpdf
// ในการสร้าง PDF และเพิ่มภาพลงไปทีละหน้า พร้อมอัปเดต progress bar และสถานะการทำงานใน UI

func writePDF(
	in <-chan Encoded,
	total int,
	progress *widget.ProgressBar,
	status *widget.Label,
	start time.Time,
	output string,
	fileList *widget.List,
) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	pageW, pageH := pdf.GetPageSize()

	buffer := map[int]Encoded{}

	next := 0
	done := 0

	for r := range in {

		buffer[r.index] = r

		for {

			res, ok := buffer[next]

			if !ok {
				break
			}

			delete(buffer, next)

			name := fmt.Sprintf("img%d", next)

			opt := gofpdf.ImageOptions{
				ImageType: "JPG",
			}

			pdf.RegisterImageOptionsReader(
				name,
				opt,
				bytes.NewReader(res.buf.Bytes()),
			)

			scale := pageW / res.w
			if res.h*scale > pageH {
				scale = pageH / res.h
			}

			w := res.w * scale
			h := res.h * scale

			x := (pageW - w) / 2
			y := (pageH - h) / 2

			pdf.AddPage()

			pdf.ImageOptions(name, x, y, w, h, false, opt, 0, "")

			jpegPool.Put(res.buf)

			next++
			done++

			updateStatus(next-1, "✅ done", fileList)

			elapsed := time.Since(start).Seconds()

			speed := float64(done) / elapsed

			per := elapsed / float64(done)

			eta := per * float64(total-done)

			fyne.Do(func() {

				progress.SetValue(float64(done) / float64(total))

				status.SetText(fmt.Sprintf(
					"🔀 %d / %d images   🎨 %.1f img/s   ⏱ ETA %.1fs",
					done,
					total,
					speed,
					eta,
				))

			})

		}
	}

	pdf.OutputFileAndClose(output)

	fyne.Do(func() {

		elapsed := time.Since(start).Seconds()

		status.SetText(fmt.Sprintf(
			"🤖✅✅ - - Done! - - ✅✅🤖 - - - ▶️ %d images in %.1fs 🕒 (%.1f img/s)",
			total,
			elapsed,
			float64(total)/elapsed,
		))
	})
}
