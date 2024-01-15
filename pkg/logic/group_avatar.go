package logic

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

type GroupAvatarGenerator struct {
	waitGroup sync.WaitGroup
	dir       string
	prefix    string
	outName   string
}

func (g *GroupAvatarGenerator) Generate(urls []string) (string, error) {
	g.waitGroup.Add(len(urls))
	localPaths := g.multiGetImages(urls)
	g.waitGroup.Wait()
	return g.Compose(localPaths, 6)
}

func (g *GroupAvatarGenerator) multiGetImages(urls []string) (localPaths []string) {
	now := time.Now().UnixMilli()
	for i, url := range urls {
		localPath := fmt.Sprintf("%s/%s_%d_%d.png", g.dir, g.prefix, now, i)
		localPaths = append(localPaths, localPath)
		go g.download(url, localPath)
	}
	return
}

func (g *GroupAvatarGenerator) download(url string, fileName string) {
	defer g.waitGroup.Done()
	out, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	client := http.Client{Timeout: 2 * time.Second}
	resp, errResp := client.Get(url)
	if errResp != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	pix, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return
	}
	_, err = io.Copy(out, bytes.NewBuffer(pix))

}

func (g *GroupAvatarGenerator) Compose(paths []string, gap int) (string, error) {
	total := len(paths)
	imageSize := 80
	imageRow := int(math.Sqrt(float64(len(paths)))) // 根据传入图片数量判断行列数量
	isFull := true
	if imageRow*imageRow < total {
		isFull = false
	}
	if !isFull {
		imageRow = imageRow + 1
	}
	imageCol := imageRow
	actualRow := total / imageCol
	if (total % imageCol) > 0 {
		actualRow += 1
	}

	backImg := imaging.New(
		imageSize*imageRow+gap*(imageCol+1),
		imageSize*imageCol+gap*(imageRow+1),
		color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	)
	startRowOffset := 0
	if actualRow < imageRow {
		startRowOffset = (imageSize + gap) / 2
	}
	tempRow, tempCol := 0, 0
	rowOffset := startRowOffset
	for tempRow < imageRow {
		rowOffset += gap
		actualCol := total - tempRow*imageCol
		if actualCol > imageCol {
			actualCol = imageCol
		}
		colOffset := (imageCol - actualCol) * imageSize / 2
		for tempCol < imageCol {
			colOffset += gap
			imgIndex := imageCol*tempRow + tempCol
			if imgIndex >= len(paths) {
				break
			}
			resizeImg, err := imaging.Open(paths[imgIndex])
			if err != nil {
				fmt.Println(err)
			}
			resizeImg = imaging.Resize(resizeImg, imageSize, imageSize, imaging.Lanczos) // 加了模糊操作
			backImg = imaging.Paste(backImg, resizeImg, image.Pt(tempCol*imageSize+colOffset, rowOffset))
			tempCol += 1
		}
		tempRow += 1
		rowOffset += imageSize
		tempCol = 0
	}
	imagePath := fmt.Sprintf("%s/%s", g.dir, g.outName)
	err := imaging.Save(backImg, imagePath)
	if err != nil {
		fmt.Println(err)
	}
	return imagePath, err
}

func NewGroupAvatarGenerator(dir, prefix, outName string) *GroupAvatarGenerator {
	return &GroupAvatarGenerator{
		waitGroup: sync.WaitGroup{},
		dir:       dir,
		prefix:    prefix,
		outName:   outName,
	}
}
