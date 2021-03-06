package collectorBot

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"iotqq-plugins-demo/Go/cards"
	"iotqq-plugins-demo/Go/plugin"
	"iotqq-plugins-demo/Go/summon"
	"iotqq-plugins-demo/Go/userData"
	"iotqq-plugins-demo/Go/util"
	"sort"
	"strings"
)

func init() {
	plugin.FactoryInstance.RegisterPlugin(&collectorBot{4})
}

type collectorBot struct {
	priority int //[0~1000)
}

func (c *collectorBot) IsTrigger(req *plugin.Request) (res bool, vNext bool) {
	f := plugin.NewCommonPrefixTriggerFunc("图鉴")
	return f(req)
}

func (c *collectorBot) Process(req *plugin.Request) []*plugin.Result {
	args := strings.Split(req.Content, " ")
	flagSet := flag.NewFlagSet("图鉴", 0)
	var star int
	var cardType int
	var isOwn int
	var buf bytes.Buffer
	flagSet.SetOutput(&buf)
	flagSet.IntVar(&star, "s", 5, "star 5/4/3")
	flagSet.IntVar(&cardType, "t", 0, "type 0-all/1-character/2-dragon")
	flagSet.IntVar(&isOwn, "o", 0, "isOwn 0-all/1-own/2-not own")

	flagSet.Parse(args[1:])
	help := ""
	fmt.Println(star, cardType, isOwn)
	if !util.IntContain(star, []int{3, 4, 5}) || !util.IntContain(cardType, []int{0, 1, 2}) ||
		!util.IntContain(isOwn, []int{0, 1, 2}) {
		help += "参数错误"
		flagSet.Usage()
	}
	help += buf.String()
	if help == "" {
		cardIndex := userData.GetUser(req.Udid).CardIndex
		img := newCollectionImage(star, cardType, isOwn, cardIndex)
		return []*plugin.Result{{Pic: img}}
	} else {
		return []*plugin.Result{{Content: help}}
	}
}

type printCard struct {
	cards.Card
	isOwn bool
}

func newCollectionImage(star, cardType, isOwn int, cardIndex []int) *image.RGBA {
	var pc printCards
	for _, card := range cards.Cards {
		if card.Star != star {
			continue
		}
		if card.CardType != cardType && 0 != cardType {
			continue
		}
		own := util.IntContain(card.ID, cardIndex)
		if isOwn != 0 && (isOwn != 2 || own) && (isOwn != 1 || !own) {
			continue
		}
		pc = append(pc, printCard{card, own})
	}

	rowNum := 6
	bg := image.NewRGBA(image.Rect(0, 0, 80*rowNum, 20+80*((len(pc)+5)/rowNum)-1))
	util.Clear(bg, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	mask := summon.GetBlackMask(80, 80)

	sort.Sort(pc)

	row, col := 0, 0
	for _, card := range pc {
		cardImg := summon.GetCardImage(card.IconUrl)
		if cardImg.Bounds().Dx() != 80 {
			cardImg = resize.Resize(80, 80, cardImg, resize.Lanczos3)
		}
		dp := image.Point{X: 80 * row, Y: 80 * col}
		col += (row + 1) / rowNum
		row = (row + 1) % rowNum
		draw.Draw(bg, image.Rectangle{Min: dp, Max: dp.Add(cardImg.Bounds().Max)}, cardImg, image.Point{}, draw.Over)
		if !card.isOwn {
			draw.Draw(bg, image.Rectangle{Min: dp, Max: dp.Add(cardImg.Bounds().Max)}, mask, image.Point{}, draw.Over)
		}
	}

	//context:=gg.NewContextForRGBA(bg)
	//context := gg.NewContext(400, 400)
	//context.DrawString("test",100,100)
	//context.SetRGB(0, 0, 0)
	//context.Fill()

	return bg
}

type printCards []printCard

func (p printCards) Len() int {
	return len(p)
}

func (p printCards) Less(i, j int) bool {
	if p[i].CardType != p[j].CardType {
		return p[i].CardType < p[j].CardType
	}
	return p[i].ID < p[j].ID
}

func (p printCards) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (c *collectorBot) Priority() int {
	return c.priority
}
