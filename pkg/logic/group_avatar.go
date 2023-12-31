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
	return g.compose(localPaths, 3)
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

func (g *GroupAvatarGenerator) compose(paths []string, gap int) (string, error) {
	imageSize := 80
	imageRow := int(math.Sqrt(float64(len(paths)))) // 根据传入图片数量判断行列数量
	imageCol := imageRow
	rowOffset := gap
	colOffset := gap
	backImg := imaging.New(
		imageSize*imageRow+gap*(imageCol+1),
		imageSize*imageCol+gap*(imageRow+1),
		color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	)
	tempRow, tempCol := 1, 1
	for tempRow < imageRow+1 {
		for tempCol < imageCol+1 {
			imgIndex := imageCol*(tempRow-1) + tempCol - 1
			resizeImg, err := imaging.Open(paths[imgIndex])
			if err != nil {
				fmt.Println(err)
			}
			resizeImg = imaging.Resize(resizeImg, imageSize, imageSize, imaging.Lanczos) // 加了模糊操作
			backImg = imaging.Paste(backImg, resizeImg, image.Pt((tempCol-1)*imageSize+colOffset, (tempRow-1)*imageSize+rowOffset))
			tempCol += 1
			colOffset += gap
		}
		tempRow += 1
		rowOffset += gap
		tempCol = 1
		colOffset = gap
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
