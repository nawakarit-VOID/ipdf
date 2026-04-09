# 🤖 Image → PDF Ultra Fast

> Convert images to PDF with blazing speed using Go + Fyne GUI

---

## 🚀 Features

* ⚡ Ultra-fast multi-core image processing
* 🧠 Smart pipeline (decode → resize → encode)
* 🖼️ Supports multiple image formats:

  * JPG / JPEG
  * PNG
  * WEBP
  * BMP
  * TIFF
* 📉 Automatic compression (30%–50% smaller)
* 📊 Real-time progress + ETA
* 🎛️ Adjustable CPU usage
* 🖥️ Clean GUI (built with Fyne)

---

## 📦 Download

👉 Download the latest version from **Releases**

---

## 🛠️ Build from Source

### Requirements

* Go 1.26.1
* Linux (recommended)

### Install dependencies

```bash
go mod tidy
```

### Run

```bash
go run main.go
```

### Build

```bash
go build -o image2pdf
```

---

## 🧠 Supported Formats (Important)

Go requires explicit image decoders.

Make sure these are included:

```go
import (
    _ "image/jpeg"
    _ "image/png"
    _ "image/gif"

    _ "golang.org/x/image/webp"
    _ "golang.org/x/image/bmp"
    _ "golang.org/x/image/tiff"
)
```

---

## ⚙️ How It Works

Pipeline architecture:

```
[Files]
   ↓
Decode (I/O heavy)
   ↓
Resize (CPU)
   ↓
Encode (CPU)
   ↓
PDF Writer
```

* Uses goroutines for parallel processing
* Optimized for multi-core CPUs
* Buffered channels for smooth flow

---

## 🎮 Usage

1. Click **Select Folder**
2. Choose image directory
3. Adjust CPU slider (optional)
4. Click **Convert to PDF**
5. Choose save location

---

## ⚠️ Notes

* GUI apps require system libraries (OpenGL / X11)
* For best portability, use **AppImage**
* Large images are automatically resized -##- **and scaled to fit A4 page size** -##-

---

## 🧪 Performance

| CPU     | Speed       |
| ------- | ----------- |
| 4 cores | ~3–5 img/s  |
| 8 cores | ~6–10 img/s |

(depending on image size)

---

## 📁 Project Structure

```
.
├── main.go
├── go.mod
├── assets/
└── README.md
```

---

## 🧑‍💻 Tech Stack

* Go
* Fyne (GUI)
* gofpdf
* nfnt/resize

---

## ❤️ Contributing

Pull requests are welcome!

---

## 📄 License

MIT License

---

## 🙌 Credits

Built with ❤️ using Go
