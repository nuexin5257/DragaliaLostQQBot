package gachaBot

import (
	"fmt"
	"github.com/Wall-ee/chinese2digits/chinese2digits"
	"iotqq-plugins-demo/Go/cards"
	"iotqq-plugins-demo/Go/plugin"
	"iotqq-plugins-demo/Go/random"
	"iotqq-plugins-demo/Go/summon"
	"iotqq-plugins-demo/Go/userData"
	"iotqq-plugins-demo/Go/util"
	"regexp"
	"strconv"
)

func init() {
	plugin.FactoryInstance.RegisterPlugin(&gachaBot{8})
}

var regex = regexp.MustCompile("^(.*)连$")
var triggerMap = map[string]int{
	"十连": 10,
	"百连": 100,
	"千连": 1000,
	"万连": 10000,
}

type gachaBot struct {
	priority int //[0~1000)
}

func (g *gachaBot) IsTrigger(req *plugin.Request) (res bool, vNext bool) {
	if util.KeyWordTrigger(req.Content, "抽卡") || util.KeyWordTrigger(req.Content, "单抽") {
		return true, false
	}

	if value, ok := triggerMap[req.Content]; ok {
		req.ExtraInfo = value
		return true, false
	}

	if res := regex.FindStringSubmatch(req.Content); len(res) > 1 {
		numStr := res[1]
		res := chinese2digits.TakeNumberFromString(numStr)
		fmt.Printf("%+v,%t", res, res)
		numDigit := res.(map[string]interface{})["replacedText"].(string)
		if num, ok := strconv.Atoi(numDigit); ok == nil && num >= 1 {
			req.ExtraInfo = num
			return true, false
		} else {
			fmt.Println(ok)
		}
	}
	return false, true
}

func (g *gachaBot) Process(req *plugin.Request) []*plugin.Result {

	if req.ExtraInfo != nil {
		var res []*plugin.Result
		num := req.ExtraInfo.(int)
		for _, i2 := range SummonALot(req.Udid, num, summon.GetMultiSummon(num)) {
			res = append(res, &plugin.Result{
				Content: i2.string,
				Pic:     i2.Image,
			})
		}
		return res
	} else {
		user := userData.GetUser(req.Udid)
		if user.RebornEggNumber >= 1 || user.Udid == 570966274 {
			card := cards.NotGachaPoolCardMgr.PickUpOne().PickUpByStar(0)[0].PickOne()
			res := new(summon.SummonRecord)
			res.Card = append(res.Card, summon.SummonCard{
				Card:     card,
				New:      false,
				StackNum: 1,
			})
			if !util.IntContain(card.ID, user.CardIndex) {
				res.Card[0].New = true
				user.CardIndex = append(user.CardIndex, card.ID)
			}
			user.RebornEggNumber--
			userData.SaveUserByUDID(req.Udid)
			img := res.ImageFormatV2(user.RebornEggNumber, user.Water)
			return []*plugin.Result{{Pic: img, Content: "\n每一世的轮回都为了找到你~转生券消耗一张" + random.RandomGetSuffix()}}
		}
		if user.SummonCardNum >= 1 {
			res := summon.OneSummon(user)
			user.SummonCardNum--
			userData.SaveUserByUDID(req.Udid)
			img := res.ImageFormatV2(user.SummonCardNum, user.Water)
			return []*plugin.Result{{Pic: img}}
		} else {
			return []*plugin.Result{{Content: fmt.Sprintf("%s召唤券不够了%s", req.NickName, random.RandomGetSuffix())}}
		}
	}
}

func (g *gachaBot) Priority() int {
	return g.priority
}
