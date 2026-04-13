// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
)

// การเริ่มต้น pipeline การแปลงภาพเป็น PDF โดยใช้ workers หลายตัวในการประมวลผลภาพ
func StartPipeline(
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
					Quality: 100,
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

	WritePDF(encoded, total, progress, status, start, output, fileList)

}
