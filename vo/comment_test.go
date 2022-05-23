package vo

import (
	"encoding/json"
	"testing"
)

func TestComment_Next(t *testing.T) {
	s := "{\"code\":0,\"data\":{\"page\":{\"num\":1,\"size\":20,\"count\":346},\"replies\":[{\"rpid\":113675144992,\"oid\":981746036,\"mid\":5847253,\"like\":0,\"ctime\":1653036963,\"content\":{\"message\":\"[星星眼][拥抱][给心心][给心心]\"}},{\"rpid\":113674936736,\"oid\":981746036,\"mid\":386313929,\"like\":0,\"ctime\":1653036889,\"content\":{\"message\":\"🤤🤤\"}},{\"rpid\":113674863504,\"oid\":981746036,\"mid\":39597993,\"like\":3,\"ctime\":1653036812,\"content\":{\"message\":\"我靠我居然在520被一个女生告白了？？一直陪我聊天，情商好高还特别 可爱，终于能体验到恋爱的滋味了，反正没事干，不如和我一起复制粘贴做白日梦[doge]\"}},{\"rpid\":113674795520,\"oid\":981746036,\"mid\":1632325152,\"like\":1,\"ctime\":1653036795,\"content\":{\"message\":\"[喜欢][喜欢][喜欢]\"}},{\"rpid\":113674654976,\"oid\":981746036,\"mid\":114985474,\"like\":0,\"ctime\":1653036700,\"content\":{\"message\":\"大润发不相信爱情\"}},{\"rpid\":113674497360,\"oid\":981746036,\"mid\":64656691,\"like\":0,\"ctime\":1653036582,\"content\":{\"message\":\"貌似素颜更好看[哦呼]\"}},{\"rpid\":113674012320,\"oid\":981746036,\"mid\":440546563,\"like\":0,\"ctime\":1653036261,\"content\":{\"message\":\"以后还不知道谁这么有福气把UP娶走。\"}},{\"rpid\":113673629328,\"oid\":981746036,\"mid\":85685453,\"like\":0,\"ctime\":1653036027,\"content\":{\"message\":\"怎么那么像阿玛粽\"}},{\"rpid\":113673620912,\"oid\":981746036,\"mid\":298554328,\"like\":1,\"ctime\":1653036011,\"content\":{\"message\":\"[辣眼睛]\"}},{\"rpid\":113673544768,\"oid\":981746036,\"mid\":67461904,\"like\":0,\"ctime\":1653035969,\"content\":{\"message\":\"转圈的时候让我想到了大 妈[笑哭]\"}},{\"rpid\":113673458960,\"oid\":981746036,\"mid\":15988021,\"like\":0,\"ctime\":1653035948,\"content\":{\"message\":\"两眼一黑了\"}},{\"rpid\":113673419712,\"oid\":981746036,\"mid\":26142767,\"like\":0,\"ctime\":1653035936,\"content\":{\"message\":\"总感觉穿了高跟鞋小腿怪怪的[思考]\"}},{\"rpid\":113673520128,\"oid\":981746036,\"mid\":333180693,\"like\":0,\"ctime\":1653035920,\"content\":{\"message\":\"哇！男友！无中生有\"}},{\"rpid\":113673072944,\"oid\":981746036,\"mid\":24078673,\"like\":0,\"ctime\":1653035674,\"content\":{\"message\":\"这不黑料里面的女主角吗\"}},{\"rpid\":113672961632,\"oid\":981746036,\"mid\":17841850,\"like\":0,\"ctime\":1653035573,\"content\":{\"message\":\"没得\"}},{\"rpid\":113672763376,\"oid\":981746036,\"mid\":8542997,\"like\":0,\"ctime\":1653035492,\"content\":{\"message\":\"省流：有男友了\"}},{\"rpid\":113672577456,\"oid\":981746036,\"mid\":396610668,\"like\":0,\"ctime\":1653035330,\"content\":{\"message\":\"我们的人生永远不会有交集 你是路上引人注目的漂亮女生 是被宠爱的小女孩 是感情里的佼佼者 是不被定义的浪漫 我是云南的 云南怒江的\"}},{\"rpid\":113672174336,\"oid\":981746036,\"mid\":27809534,\"like\":1,\"ctime\":1653035093,\"content\":{\"message\":\"[汤圆]\"}},{\"rpid\":113672023232,\"oid\":981746036,\"mid\":252646508,\"like\":0,\"ctime\":1653034964,\"content\":{\"message\":\"我还以为你在隔离酒店拍的哈哈哈\"}},{\"rpid\":113671729296,\"oid\":981746036,\"mid\":627358238,\"like\":0,\"ctime\":1653034829,\"content\":{\"message\":\"江南style经久不衰哇\"}}]},\"mid\":0,\"task\":{\"TaskType\":0,\"Payload\":\"\"}}"
	var cmt MainComment
	err := json.Unmarshal([]byte(s), &cmt)
	if err != nil {
		t.Error(err)
	}
	cmt.Next()
}
