package address

import (
	"testing"

	"bytes"
	"fmt"
)

func TestNewRows(t *testing.T) {

	type Town struct {
		town     string
		townKana string
	}

	type Case struct {
		data     string
		zip7     string
		expected []Town
	}

	cases := []Case{
		{
			zip7: "4506246",
			data: `23105,"450  ","4506246","ｱｲﾁｹﾝ","ﾅｺﾞﾔｼﾅｶﾑﾗｸ","ﾒｲｴｷﾐｯﾄﾞﾗﾝﾄﾞｽｸｴｱ(ｺｳｿｳﾄｳ)(46ｶｲ)","愛知県","名古屋市中村区","名駅ミッドランドスクエア（高層棟）（４６階）",0,0,0,0,0,0`,
			expected: []Town{
				// (xx階) の場合は1行のみ
				{
					town:     "名駅ミッドランドスクエア高層棟46階",
					townKana: "メイエキミッドランドスクエアコウソウトウ46カイ",
				},
			},
		},

		// (全域) は消す
		{
			zip7: "0895865",
			data: `01649,"08958","0895865","ﾎｯｶｲﾄﾞｳ","ﾄｶﾁｸﾞﾝｳﾗﾎﾛﾁｮｳ","ｱﾂﾅｲ(ｾﾞﾝｲｷ)","北海道","十勝郡浦幌町","厚内（全域）",0,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "厚内",
					townKana: "アツナイ",
				},
			},
		},

		// (成田国際空港内) は消す
		{
			zip7: "2820031",
			data: `12347,"282  ","2820031","ﾁﾊﾞｹﾝ","ｶﾄﾘｸﾞﾝﾀｺﾏﾁ","ﾋﾄｸﾜﾀﾞ(ﾅﾘﾀｺｸｻｲｸｳｺｳﾅｲ)","千葉県","香取郡多古町","一鍬田（成田国際空港内）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "一鍬田",
					townKana: "ヒトクワダ",
				},
			},
		},

		// イッチョウメ が連続
		{
			zip7: "6028064",
			data: `26102,"602  ","6028064","ｷｮｳﾄﾌ","ｷｮｳﾄｼｶﾐｷﾞｮｳｸ","ｲｯﾁｮｳﾒ","京都府","京都市上京区","一町目（上長者町通堀川東入、東堀川通上長者町上る、東堀川通中",0,0,0,0,0,0
26102,"602  ","6028064","ｷｮｳﾄﾌ","ｷｮｳﾄｼｶﾐｷﾞｮｳｸ","ｲｯﾁｮｳﾒ","京都府","京都市上京区","立売通下る）",0,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "一町目",
					townKana: "イッチョウメ",
				},
				{
					town:     "一町目上長者町通堀川東入",
					townKana: "イッチョウメ",
				},
				{
					town:     "一町目東堀川通上長者町上る",
					townKana: "イッチョウメ",
				},
				{
					town:     "一町目東堀川通中立売通下る",
					townKana: "イッチョウメ",
				},
			},
		},

		// () 内の数字っぽいやつは除去される
		{
			zip7: "0482402",
			data: `01407,"04824","0482402","ﾎｯｶｲﾄﾞｳ","ﾖｲﾁｸﾞﾝﾆｷﾁｮｳ","ｵｵｴ(1ﾁｮｳﾒ､2ﾁｮｳﾒ<651､662､668ﾊﾞﾝﾁ>ｲｶﾞｲ､3ﾁｮｳﾒ5､1","北海道","余市郡仁木町","大江（１丁目、２丁目「６５１、６６２、６６８番地」以外、３丁目５、１",1,0,1,0,0,0
01407,"04824","0482402","ﾎｯｶｲﾄﾞｳ","ﾖｲﾁｸﾞﾝﾆｷﾁｮｳ","3-4､20､678､687ﾊﾞﾝﾁ)","北海道","余市郡仁木町","３−４、２０、６７８、６８７番地）",1,0,1,0,0,0`,
			expected: []Town{
				{
					town:     "大江",
					townKana: "オオエ",
				},
			},
		},
		{
			zip7: "0482331",
			data: `01407,"04823","0482331","ﾎｯｶｲﾄﾞｳ","ﾖｲﾁｸﾞﾝﾆｷﾁｮｳ","ｵｵｴ(2ﾁｮｳﾒ651､662､668ﾊﾞﾝﾁ､3ﾁｮｳﾒ103､118､","北海道","余市郡仁木町","大江（２丁目６５１、６６２、６６８番地、３丁目１０３、１１８、",1,0,1,0,0,0
01407,"04823","0482331","ﾎｯｶｲﾄﾞｳ","ﾖｲﾁｸﾞﾝﾆｷﾁｮｳ","210､254､267､372､444､469ﾊﾞﾝﾁ)","北海道","余市郡仁木町","２１０、２５４、２６７、３７２、４４４、４６９番地）",1,0,1,0,0,0`,
			expected: []Town{
				{
					town:     "大江",
					townKana: "オオエ",
				},
			},
		},
		{
			zip7: "0300924",
			data: `02201,"030  ","0300924","ｱｵﾓﾘｹﾝ","ｱｵﾓﾘｼ","ﾀｷｻﾜ(ｼﾓｶﾜﾗ190-1)","青森県","青森市","滝沢（下川原１９０−１）",1,1,0,0,0,0`,
			expected: []Town{
				{
					town:     "滝沢",
					townKana: "タキサワ",
				},
			},
		},

		// 地割
		{
			zip7: "0295505",
			data: `03366,"02955","0295505","ｲﾜﾃｹﾝ","ﾜｶﾞｸﾞﾝﾆｼﾜｶﾞﾏﾁ","ﾕﾓﾄ29ﾁﾜﾘ､ﾕﾓﾄ30ﾁﾜﾘ","岩手県","和賀郡西和賀町","湯本２９地割、湯本３０地割",0,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "湯本",
					townKana: "ユモト",
				},
			},
		},
		{
			zip7: "0287913",
			data: `03507,"02879","0287913","ｲﾜﾃｹﾝ","ｸﾉﾍｸﾞﾝﾋﾛﾉﾁｮｳ","ﾀﾈｲﾁﾀﾞｲ24ﾁﾜﾘ-ﾀﾞｲ25ﾁﾜﾘ(ﾐﾄﾞﾘｶﾞｵｶﾁｮｳ､ﾖｺﾃ)","岩手県","九戸郡洋野町","種市第２４地割〜第２５地割（緑ケ丘町、横手）",0,1,0,0,0,0`,
			expected: []Town{
				{
					town:     "種市",
					townKana: "タネイチ",
				},
			},
		},

		//
		// xx を除く  0330071 0330072 0285102 9800065 9960301 3842304 4280049 4400075 6511102 7201264 7983321
		//

		// 犬落瀬（内金矢、内山、岡沼、金沢、金矢、上淋代、木越、権現沢、四木、七百、下久保「１７４を除く」、下淋代、高森、通目木、坪毛沢「２５、６３７、６４１、６４３、６４７を除く」、中屋敷、沼久保、根古橋、堀切沢、南平、柳沢、大曲）
		{
			zip7: "0330071",
			data: `02405,"033  ","0330071","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","ｲﾇｵﾄｾ(ｳﾁｶﾅﾔ､ｳﾁﾔﾏ､ｵｶﾇﾏ､ｶﾅｻﾞﾜ､ｶﾅﾔ､ｶﾐｻﾋﾞｼﾛ､ｷｺｼ､ｺﾞﾝｹﾞﾝｻﾜ､","青森県","上北郡六戸町","犬落瀬（内金矢、内山、岡沼、金沢、金矢、上淋代、木越、権現沢、",1,1,0,0,0,0
02405,"033  ","0330071","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","ｼｷ､ｼﾁﾋｬｸ､ｼﾓｸﾎﾞ<174ｦﾉｿﾞｸ>､ｼﾓｻﾋﾞｼﾛ､ﾀｶﾓﾘ､ﾂﾞﾒｷ､ﾂﾎﾞｹｻﾞﾜ<2","青森県","上北郡六戸町","四木、七百、下久保「１７４を除く」、下淋代、高森、通目木、坪毛沢「２",1,1,0,0,0,0
02405,"033  ","0330071","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","5､637､641､643､647ｦﾉｿﾞｸ>､ﾅｶﾔｼｷ､ﾇﾏｸﾎﾞ､ﾈｺﾊｼ､ﾎﾘｷﾘ","青森県","上北郡六戸町","５、６３７、６４１、６４３、６４７を除く」、中屋敷、沼久保、根古橋、堀切",1,1,0,0,0,0
02405,"033  ","0330071","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","ｻﾜ､ﾐﾅﾐﾀｲ､ﾔﾅｷﾞｻﾜ､ｵｵﾏｶﾞﾘ)","青森県","上北郡六戸町","沢、南平、柳沢、大曲）",1,1,0,0,0,0`,
			expected: []Town{
				{
					town:     "犬落瀬",
					townKana: "イヌオトセ",
				},
				{
					town:     "犬落瀬内金矢",
					townKana: "イヌオトセウチカナヤ",
				},
				{
					town:     "犬落瀬内山",
					townKana: "イヌオトセウチヤマ",
				},
				{
					town:     "犬落瀬岡沼",
					townKana: "イヌオトセオカヌマ",
				},
				{
					town:     "犬落瀬金沢",
					townKana: "イヌオトセカナザワ",
				},
				{
					town:     "犬落瀬金矢",
					townKana: "イヌオトセカナヤ",
				},
				{
					town:     "犬落瀬上淋代",
					townKana: "イヌオトセカミサビシロ",
				},
				{
					town:     "犬落瀬木越",
					townKana: "イヌオトセキコシ",
				},
				{
					town:     "犬落瀬権現沢",
					townKana: "イヌオトセゴンゲンサワ",
				},
				{
					town:     "犬落瀬四木",
					townKana: "イヌオトセシキ",
				},
				{
					town:     "犬落瀬七百",
					townKana: "イヌオトセシチヒャク",
				},
				{
					town:     "犬落瀬下久保",
					townKana: "イヌオトセシモクボ",
				},
				{
					town:     "犬落瀬下淋代",
					townKana: "イヌオトセシモサビシロ",
				},
				{
					town:     "犬落瀬高森",
					townKana: "イヌオトセタカモリ",
				},
				{
					town:     "犬落瀬通目木",
					townKana: "イヌオトセヅメキ",
				},
				{
					town:     "犬落瀬坪毛沢",
					townKana: "イヌオトセツボケザワ",
				},
				{
					town:     "犬落瀬中屋敷",
					townKana: "イヌオトセナカヤシキ",
				},
				{
					town:     "犬落瀬沼久保",
					townKana: "イヌオトセヌマクボ",
				},
				{
					town:     "犬落瀬根古橋",
					townKana: "イヌオトセネコハシ",
				},
				{
					town:     "犬落瀬堀切沢",
					townKana: "イヌオトセホリキリサワ",
				},
				{
					town:     "犬落瀬南平",
					townKana: "イヌオトセミナミタイ",
				},
				{
					town:     "犬落瀬柳沢",
					townKana: "イヌオトセヤナギサワ",
				},
				{
					town:     "犬落瀬大曲",
					townKana: "イヌオトセオオマガリ",
				},
			},
		},

		{
			zip7: "9800065",
			data: `04101,"980  ","9800065","ﾐﾔｷﾞｹﾝ","ｾﾝﾀﾞｲｼｱｵﾊﾞｸ","ﾂﾁﾄｲ(1ﾁｮｳﾒ<11ｦﾉｿﾞｸ>)","宮城県","仙台市青葉区","土樋（１丁目「１１を除く」）",0,0,1,0,0,0`,
			expected: []Town{
				{
					town:     "土樋",
					townKana: "ツチトイ",
				},
			},
		},
		// 折茂（
		// 		今熊「２１３〜２３４、２４０、２４７、２６２、２６６、２７５、２７７、２８０、２９５、１１９９、１２０６、１５０４を除く」、
		// 		大原、
		// 		沖山、
		// 		上折茂「１−１３、７１−１９２を除く」
		// ）
		{
			zip7: "0330072",
			data: `02405,"033  ","0330072","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","ｵﾘﾓ(ｲﾏｸﾏ<213-234､240､247､262､266､27","青森県","上北郡六戸町","折茂（今熊「２１３〜２３４、２４０、２４７、２６２、２６６、２７",1,1,0,0,0,0
02405,"033  ","0330072","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","5､277､280､295､1199､1206､1504ｦﾉｿﾞｸ>､","青森県","上北郡六戸町","５、２７７、２８０、２９５、１１９９、１２０６、１５０４を除く」、",1,1,0,0,0,0
02405,"033  ","0330072","ｱｵﾓﾘｹﾝ","ｶﾐｷﾀｸﾞﾝﾛｸﾉﾍﾏﾁ","ｵｵﾊﾗ､ｵｷﾔﾏ､ｶﾐｵﾘﾓ<1-13､71-192ｦﾉｿﾞｸ>)","青森県","上北郡六戸町","大原、沖山、上折茂「１−１３、７１−１９２を除く」）",1,1,0,0,0,0
`,
			expected: []Town{
				{
					town:     "折茂",
					townKana: "オリモ",
				},
				{
					town:     "折茂今熊",
					townKana: "オリモイマクマ",
				},
				{
					town:     "折茂大原",
					townKana: "オリモオオハラ",
				},
				{
					town:     "折茂沖山",
					townKana: "オリモオキヤマ",
				},
				{
					town:     "折茂上折茂",
					townKana: "オリモカミオリモ",
				},
			},
		},
		{
			zip7: "0282504",
			data: `03202,"02825","0282504","ｲﾜﾃｹﾝ","ﾐﾔｺｼ","ﾊｺｲｼ(ﾀﾞｲ2ﾁﾜﾘ<70-136>-ﾀﾞｲ4ﾁﾜﾘ<3-11>)","岩手県","宮古市","箱石（第２地割「７０〜１３６」〜第４地割「３〜１１」）",1,1,0,0,0,0`,
			expected: []Town{
				{
					town:     "箱石",
					townKana: "ハコイシ",
				},
			},
		},
		{
			zip7: "0285102",
			data: `03302,"02851","0285102","ｲﾜﾃｹﾝ","ｲﾜﾃｸﾞﾝｸｽﾞﾏｷﾏﾁ","ｸｽﾞﾏｷ(ﾀﾞｲ40ﾁﾜﾘ<57ﾊﾞﾝﾁ125､176ｦﾉｿﾞｸ>-ﾀﾞｲ45","岩手県","岩手郡葛巻町","葛巻（第４０地割「５７番地１２５、１７６を除く」〜第４５",1,1,0,0,0,0
03302,"02851","0285102","ｲﾜﾃｹﾝ","ｲﾜﾃｸﾞﾝｸｽﾞﾏｷﾏﾁ","ﾁﾜﾘ)","岩手県","岩手郡葛巻町","地割）",1,1,0,0,0,0`,
			expected: []Town{
				{
					town:     "葛巻",
					townKana: "クズマキ",
				},
			},
		},
		{
			zip7: "9996652",
			data: `06203,"99976","9996652","ﾔﾏｶﾞﾀｹﾝ","ﾂﾙｵｶｼ","ｿｴｶﾞﾜ(ﾜﾀﾄｻﾞﾜ<ﾀｹﾉｺｻﾞﾜｵﾝｾﾝ>)","山形県","鶴岡市","添川（渡戸沢「筍沢温泉」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "添川",
					townKana: "ソエガワ",
				},
				{
					town:     "添川渡戸沢",
					townKana: "ソエガワワタトザワ",
				},
			},
		},
		// 南山（４３０番地以上「１７７０−１〜２、１８６２−４２、１９２３−５を除く」、大谷地、折渡、鍵金野、金山、滝ノ沢、豊牧、沼の台、肘折、平林）
		{
			zip7: "9960301",
			data: `06365,"99602","9960301","ﾔﾏｶﾞﾀｹﾝ","ﾓｶﾞﾐｸﾞﾝｵｵｸﾗﾑﾗ","ﾐﾅﾐﾔﾏ(430ﾊﾞﾝﾁｲｼﾞｮｳ<1770-1-2､1862-42､","山形県","最上郡大蔵村","南山（４３０番地以上「１７７０−１〜２、１８６２−４２、",1,1,0,0,0,0
06365,"99602","9960301","ﾔﾏｶﾞﾀｹﾝ","ﾓｶﾞﾐｸﾞﾝｵｵｸﾗﾑﾗ","1923-5ｦﾉｿﾞｸ>､ｵｵﾔﾁ､ｵﾘﾜﾀﾘ､ｶﾝｶﾈﾉ､ｷﾝｻﾞﾝ､ﾀｷﾉｻﾜ､ﾄﾖﾏｷ､ﾇﾏﾉﾀﾞｲ､","山形県","最上郡大蔵村","１９２３−５を除く」、大谷地、折渡、鍵金野、金山、滝ノ沢、豊牧、沼の台、",1,1,0,0,0,0
06365,"99602","9960301","ﾔﾏｶﾞﾀｹﾝ","ﾓｶﾞﾐｸﾞﾝｵｵｸﾗﾑﾗ","ﾋｼﾞｵﾘ､ﾋﾗﾊﾞﾔｼ)","山形県","最上郡大蔵村","肘折、平林）",1,1,0,0,0,0
`,
			expected: []Town{
				{
					town:     "南山",
					townKana: "ミナミヤマ",
				},
				{
					town:     "南山大谷地",
					townKana: "ミナミヤマオオヤチ",
				},
				{
					town:     "南山折渡",
					townKana: "ミナミヤマオリワタリ",
				},
				{
					town:     "南山鍵金野",
					townKana: "ミナミヤマカンカネノ",
				},
				{
					town:     "南山金山",
					townKana: "ミナミヤマキンザン",
				},
				{
					town:     "南山滝ノ沢",
					townKana: "ミナミヤマタキノサワ",
				},
				{
					town:     "南山豊牧",
					townKana: "ミナミヤマトヨマキ",
				},
				{
					town:     "南山沼の台",
					townKana: "ミナミヤマヌマノダイ",
				},
				{
					town:     "南山肘折",
					townKana: "ミナミヤマヒジオリ",
				},
				{
					town:     "南山平林",
					townKana: "ミナミヤマヒラバヤシ",
				},
			},
		},
		{
			zip7: "3771405",
			data: `10425,"37714","3771405","ｸﾞﾝﾏｹﾝ","ｱｶﾞﾂﾏｸﾞﾝﾂﾏｺﾞｲﾑﾗ","ｶﾝﾊﾞﾗ(ﾓﾛｼｺ<ｱｻﾏｴﾝ>)","群馬県","吾妻郡嬬恋村","鎌原（モロシコ「浅間園」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "鎌原",
					townKana: "カンバラ",
				},
				{
					town:     "鎌原モロシコ",
					townKana: "カンバラモロシコ",
				},
			},
		},

		// TODO: カナのネストされた() の中、なんで処理できてるんだ？？
		// これデータおかしい。他のデータだと、 カナの中のカッコは <> になっている。
		{
			zip7: "3703321",
			data: `10429,"37033","3703321","ｸﾞﾝﾏｹﾝ","ｱｶﾞﾂﾏｸﾞﾝﾋｶﾞｼｱｶﾞﾂﾏﾏﾁ","ｲｽﾞﾐｻﾜ(ｴﾎﾞｼ(ﾊﾙﾅｺﾊﾝ)､ｴﾎﾞｼｺｸﾕｳﾘﾝ77ﾘﾝﾊﾝ)","群馬県","吾妻郡東吾妻町","泉沢（烏帽子「榛名湖畔」、烏帽子国有林７７林班）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "泉沢",
					townKana: "イズミサワ",
				},
				{
					town:     "泉沢烏帽子",
					townKana: "イズミサワエボシハルナコハン",
				},
				{
					town:     "泉沢烏帽子国有林77林班",
					townKana: "イズミサワエボシコクユウリン77リンハン",
				},
			},
		},
		{
			zip7: "3703311",
			data: `10429,"37033","3703311","ｸﾞﾝﾏｹﾝ","ｱｶﾞﾂﾏｸﾞﾝﾋｶﾞｼｱｶﾞﾂﾏﾏﾁ","ｵｶｻﾞｷ(ｴﾎﾞｼ<ﾊﾙﾅｺﾊﾝ>)","群馬県","吾妻郡東吾妻町","岡崎（烏帽子「榛名湖畔」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "岡崎",
					townKana: "オカザキ",
				},
				{
					town:     "岡崎烏帽子",
					townKana: "オカザキエボシ",
				},
			},
		},

		{
			zip7: "3703322",
			data: `10429,"37033","3703322","ｸﾞﾝﾏｹﾝ","ｱｶﾞﾂﾏｸﾞﾝﾋｶﾞｼｱｶﾞﾂﾏﾏﾁ","ｶﾜﾄﾞ(ｴﾎﾞｼ<ﾊﾙﾅｺﾊﾝ>)","群馬県","吾妻郡東吾妻町","川戸（烏帽子「榛名湖畔」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "川戸",
					townKana: "カワド",
				},
				{
					town:     "川戸烏帽子",
					townKana: "カワドエボシ",
				},
			},
		},

		{
			zip7: "3862211",
			data: `20207,"38622","3862211","ﾅｶﾞﾉｹﾝ","ｽｻﾞｶｼ","ﾆﾚｲﾏﾁ(3153-1-3153-1100<ﾐﾈﾉﾊﾗ>)","長野県","須坂市","仁礼町（３１５３−１〜３１５３−１１００「峰の原」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "仁礼町",
					townKana: "ニレイマチ",
				},
			},
		},

		// 茂田井（１〜５００「２１１番地を除く」「古町」、２５２７〜２５２９「土遠」）

		{
			zip7: "3842304",
			data: `20324,"38423","3842304","ﾅｶﾞﾉｹﾝ","ｷﾀｻｸｸﾞﾝﾀﾃｼﾅﾏﾁ","ﾓﾀｲ(1-500<211ﾊﾞﾝﾁｦﾉｿﾞｸ><ﾌﾙﾏﾁ>､2527-2529","長野県","北佐久郡立科町","茂田井（１〜５００「２１１番地を除く」「古町」、２５２７〜２５２９",1,0,0,0,0,0
20324,"38423","3842304","ﾅｶﾞﾉｹﾝ","ｷﾀｻｸｸﾞﾝﾀﾃｼﾅﾏﾁ","<ﾄﾞﾄｵ>)","長野県","北佐久郡立科町","「土遠」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "茂田井",
					townKana: "モタイ",
				},
			},
		},

		// 牧之原（２５０〜３４３番地「２５５、２５６、２５８、２５９、２６２、２７６、２９４〜３００、３０２〜３０４番地を除く」）
		{
			zip7: "4280049",
			data: `22209,"428  ","4280049","ｼｽﾞｵｶｹﾝ","ｼﾏﾀﾞｼ","ﾏｷﾉﾊﾗ(250-343ﾊﾞﾝﾁ<255､256､258､259､262､","静岡県","島田市","牧之原（２５０〜３４３番地「２５５、２５６、２５８、２５９、２６２、",1,0,0,0,0,0
22209,"428  ","4280049","ｼｽﾞｵｶｹﾝ","ｼﾏﾀﾞｼ","276､294-300､302-304ﾊﾞﾝﾁｦﾉｿﾞｸ>)","静岡県","島田市","２７６、２９４〜３００、３０２〜３０４番地を除く」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "牧之原",
					townKana: "マキノハラ",
				},
			},
		},

		// TODO: FIXME これデータがおかしい@@@@@@@@@
		// 山田町下谷上（大上谷、修法ケ原、中一里山「９番地の４、１２番地を除く」長尾山、再度公園）
		{
			zip7: "6511102",
			data: `28109,"65111","6511102","ﾋｮｳｺﾞｹﾝ","ｺｳﾍﾞｼｷﾀｸ","ﾔﾏﾀﾞﾁｮｳｼﾓﾀﾆｶﾞﾐ(ｵｵｶﾐﾀﾞﾆ､ｼｭｳﾎｳｶﾞﾊﾗ､ﾅｶｲﾁﾘﾔﾏ<9ﾊﾞﾝﾁﾉ4､12ﾊﾞﾝﾁｦﾉｿﾞｸ>ﾅｶﾞ","兵庫県","神戸市北区","山田町下谷上（大上谷、修法ケ原、中一里山「９番地の４、１２番地を除く」長",1,1,0,0,0,0
28109,"65111","6511102","ﾋｮｳｺﾞｹﾝ","ｺｳﾍﾞｼｷﾀｸ","ｵﾔﾏ､ﾌﾀﾀﾋﾞｺｳｴﾝ)","兵庫県","神戸市北区","尾山、再度公園）",1,1,0,0,0,0`,
			expected: []Town{
				{
					town:     "山田町下谷上",
					townKana: "ヤマダチョウシモタニガミ",
				},
				{
					town:     "山田町下谷上大上谷",
					townKana: "ヤマダチョウシモタニガミオオカミダニ",
				},
				{
					town:     "山田町下谷上修法ケ原",
					townKana: "ヤマダチョウシモタニガミシュウホウガハラ",
				},
				{
					town:     "山田町下谷上中一里山長尾山",
					townKana: "ヤマダチョウシモタニガミナカイチリヤマナガオヤマ",
				},
				{
					town:     "山田町下谷上再度公園",
					townKana: "ヤマダチョウシモタニガミフタタビコウエン",
				},
			},
		},
		{
			zip7: "6650808",
			data: `28214,"665  ","6650808","ﾋｮｳｺﾞｹﾝ","ﾀｶﾗﾂﾞｶｼ","ｷﾘﾊﾀ(ﾅｶﾞｵｻﾝ<ｿﾉﾀ>)","兵庫県","宝塚市","切畑（長尾山「その他」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "切畑",
					townKana: "キリハタ",
				},
				{
					town:     "切畑長尾山",
					townKana: "キリハタナガオサン",
				},
			},
		},
		{
			zip7: "6302168",
			data: `29201,"63021","6302168","ﾅﾗｹﾝ","ﾅﾗｼ","ﾎﾞﾀﾞｲｾﾝﾁｮｳ(173-257ﾊﾞﾝﾁ<ﾊﾁﾌﾞｾﾄｳｹﾞ>)","奈良県","奈良市","菩提山町（１７３〜２５７番地「鉢伏峠」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "菩提山町",
					townKana: "ボダイセンチョウ",
				},
			},
		},
		{
			zip7: "7200845",
			data: `34207,"720  ","7200845","ﾋﾛｼﾏｹﾝ","ﾌｸﾔﾏｼ","ｱｼﾀﾞﾁｮｳﾌｸﾀﾞ(376-10<ｾｲﾎｳｼﾞ>)","広島県","福山市","芦田町福田（３７６−１０「聖宝寺」）",1,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "芦田町福田",
					townKana: "アシダチョウフクダ",
				},
			},
		},
		// 町名と() 内が同じ場合は（）内削除
		{
			zip7: "6560514",
			data: `28224,"65605","6560514","ﾋｮｳｺﾞｹﾝ","ﾐﾅﾐｱﾜｼﾞｼ","ｶｼｭｳ(ｶｼｭｳ)","兵庫県","南あわじ市","賀集（賀集）",0,0,0,0,0,0`,
			expected: []Town{
				{
					town:     "賀集",
					townKana: "カシュウ",
				},
			},
		},
	}

	for x, c := range cases {
		t.Run(c.zip7, func(t *testing.T) {
			reader := bytes.NewReader([]byte(c.data))
			r := NewReader(reader)
			cols, _ := r.Read()

			rows := NewRows(cols)

			if len(rows) != len(c.expected) {
				fmt.Println(rows)
				t.Errorf("#%d: zip:%s want '%d', got '%d'\n", x, c.zip7, len(c.expected), len(rows))
			} else {
				for i, row := range rows {
					if row.Town != c.expected[i].town {
						t.Errorf("#%d: zip:%s want '%s', got '%s'\n", x, c.zip7, c.expected[i].town, row.Town)
					}
					if row.TownKana != c.expected[i].townKana {
						t.Errorf("#%d: zip:%s want '%s', got '%s'\n", x, c.zip7, c.expected[i].townKana, row.TownKana)
					}
				}
			}
		})
	}
}