// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"bytes"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jung-kurt/gofpdf"
)

// ฟังก์ชันสำหรับการเขียนไฟล์ PDF โดยรับข้อมูลภาพที่ถูก encode แล้วจาก channel
// และจัดเรียงตามลำดับ index เพื่อให้ภาพอยู่ในลำดับที่ถูกต้องใน PDF จากนั้นใช้ gofpdf
// ในการสร้าง PDF และเพิ่มภาพลงไปทีละหน้า พร้อมอัปเดต progress bar และสถานะการทำงานใน UI

func WritePDF(
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

				status.SetText(fmt.Sprintf(
					"🔀 %d / %d images   🎨 %.1f img/s   ⏱ ETA %.1fs",
					done,
					total,
					speed,
					eta,
				))

			})

		}

		totalSteps := float64(total) + 1 // +1 สำหรับ savePDF

		progress.SetValue(float64(done) / float64(totalSteps))
	}

	pdf.OutputFileAndClose(output)
	progress.SetValue(1.0)
	fyne.Do(func() {

		elapsed := time.Since(start).Seconds()

		status.SetText(fmt.Sprintf(
			"✅ Done --> %d images in %.1fs 🕒 (%.1f img/s)",
			total,
			elapsed,
			float64(total)/elapsed,
		))

	})
}
