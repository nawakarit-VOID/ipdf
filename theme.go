// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

func (m MyTheme) Color(name fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if v == theme.VariantDark {
		//theme black
		switch name {
		//🎨 สีหลัก
		case theme.ColorNameBackground: //→ พื้นหลังหลักของแอพ สีพื้นหลังสุด *ถ้าไม่มีภาพขั้นกลาง
			return color.NRGBA{46, 46, 51, 255}

		case theme.ColorNameForeground: //→ สีตัวอักษร/ไอคอน
			return color.White

			//🔘 ปุ่ม
		case theme.ColorNameButton: //→ สีปุ่มปกติ
			return color.NRGBA{200, 200, 200, 50}

		case theme.ColorNamePressed: //→ ตอนกด***
			return color.NRGBA{191, 218, 255, 255}

		case theme.ColorNameHover: //→ ตอนเอาเมาส์ชี้
			return color.NRGBA{255, 255, 255, 50}

		case theme.ColorNameDisabledButton: //→ ปุ่มที่กดไม่ได้***
			return color.NRGBA{255, 0, 0, 255}

		// 🧠 สถานะทั่วไป
		case theme.ColorNameDisabled: // → สีของ element ที่ใช้ไม่ได้***
			return color.NRGBA64{255, 0, 0, 255}

		case theme.ColorNameFocus: // → ตอนถูกเลือก / โฟกัส***
			return color.NRGBA{50, 50, 50, 40}

		// 🌈 สีหลักของแอพ
		case theme.ColorNamePrimary: // → สีเด่น (progress bar / highlight / ปุ่มสำคัญ)
			return color.NRGBA{109, 191, 57, 200}

			//🧾 Input / UI
		case theme.ColorNameInputBackground: // → พื้นหลังช่อง input / select
			return color.NRGBA{200, 200, 200, 50}

		case theme.ColorNamePlaceHolder: // → ตัวอักษร placeholder***
			return color.NRGBA{255, 0, 0, 255}

		// 🪟 Layer / พื้นหลังพิเศษ
		case theme.ColorNameMenuBackground: // → เมนู (dropdown / popup)
			return color.NRGBA{255, 255, 255, 50}

		case theme.ColorNameOverlayBackground: // → dialog / overlay
			return color.NRGBA{46, 46, 51, 255}

		case theme.ColorNameShadow: // → เงา
			return color.NRGBA{255, 255, 255, 40}

			//⚠️ สถานะพิเศษ
		case theme.ColorNameError: // → error (แดง)
			return color.NRGBA{255, 0, 0, 255}

		case theme.ColorNameSuccess: // → success (เขียว)
			return color.NRGBA{0, 255, 0, 255}

		case theme.ColorNameWarning: // → warning (เหลือง/ส้ม)
			return color.NRGBA{255, 165, 255, 255}

		}
	} else {
		//theme white
		switch name {
		//🎨 สีหลัก
		case theme.ColorNameBackground: //→ พื้นหลังหลักของแอพ สีพื้นหลังสุด *ถ้าไม่มีภาพขั้นกลาง
			return color.NRGBA{255, 255, 255, 255}

		case theme.ColorNameForeground: //→ สีตัวอักษร/ไอคอน
			return color.Black

			//🔘 ปุ่ม
		case theme.ColorNameButton: //→ สีปุ่มปกติ
			return color.NRGBA{50, 50, 50, 50}

		case theme.ColorNamePressed: //→ ตอนกด***
			return color.NRGBA{191, 218, 255, 255}

		case theme.ColorNameHover: //→ ตอนเอาเมาส์ชี้
			return color.NRGBA{255, 255, 255, 50}

		case theme.ColorNameDisabledButton: //→ ปุ่มที่กดไม่ได้***
			return color.NRGBA{255, 0, 0, 255}

		// 🧠 สถานะทั่วไป
		case theme.ColorNameDisabled: // → สีของ element ที่ใช้ไม่ได้***
			return color.NRGBA64{255, 0, 0, 255}

		case theme.ColorNameFocus: // → ตอนถูกเลือก / โฟกัส***
			return color.NRGBA{50, 50, 50, 40}

		// 🌈 สีหลักของแอพ
		case theme.ColorNamePrimary: // → สีเด่น (progress bar / highlight / ปุ่มสำคัญ)
			return color.NRGBA{109, 191, 57, 200}

			//🧾 Input / UI
		case theme.ColorNameInputBackground: // → พื้นหลังช่อง input / select
			return color.NRGBA{50, 50, 50, 50}

		case theme.ColorNamePlaceHolder: // → ตัวอักษร placeholder***
			return color.NRGBA{255, 0, 0, 255}

		// 🪟 Layer / พื้นหลังพิเศษ
		case theme.ColorNameMenuBackground: // → เมนู (dropdown / popup)
			return color.NRGBA{255, 255, 255, 50}

		case theme.ColorNameOverlayBackground: // → dialog / overlay
			return color.NRGBA{230, 230, 230, 255}

		case theme.ColorNameShadow: // → เงา
			return color.NRGBA{255, 255, 255, 40}

			//⚠️ สถานะพิเศษ
		case theme.ColorNameError: // → error (แดง)
			return color.NRGBA{255, 0, 0, 255}

		case theme.ColorNameSuccess: // → success (เขียว)
			return color.NRGBA{0, 255, 0, 255}

		case theme.ColorNameWarning: // → warning (เหลือง/ส้ม)
			return color.NRGBA{255, 165, 255, 255}

		}
	}
	return theme.DefaultTheme().Color(name, v)
}

// ต้องมีครบ
func (m MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return myFont
	//return theme.DefaultTheme().Font(s)
}
func (m MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {

	// 📏 spacing / ระยะ
	case theme.SizeNamePadding: // → ระยะห่างทั่วไป (margin/padding ของ widget)
		return 4
	case theme.SizeNameSeparatorThickness: // → ความหนาเส้นคั่น
		return 0

	// 🖼️ ไอคอน / scrollbar
	case theme.SizeNameInlineIcon: // → ขนาด icon ในปุ่ม/ข้อความ /dialog
		return 14

	case theme.SizeNameScrollBar: // → ความกว้าง scrollbar ปกติ
		return 12
	case theme.SizeNameScrollBarSmall: // → scrollbar แบบเล็ก
		return 3

	// 🔤 ขนาดตัวอักษร
	case theme.SizeNameText: // → ข้อความปกติ
		return 14
	case theme.SizeNameHeadingText: // → หัวข้อใหญ่
		return 20
	case theme.SizeNameSubHeadingText: // → หัวข้อรอง
		return 16
	case theme.SizeNameCaptionText: // → ตัวเล็ก (caption/คำอธิบาย)
		return 12

	// 🧾 input
	case theme.SizeNameInputBorder: // → ความหนาขอบ input
		return 12

	}
	return theme.DefaultTheme().Size(name)
}
