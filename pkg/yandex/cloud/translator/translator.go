package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	API       = "https://translate.api.cloud.yandex.net/translate/v2/translate"
	languages = map[string]string{
		"af":     "af",     // Африкаанс
		"am":     "am",     // Амхарский
		"ar":     "ar",     // Арабский
		"az":     "az",     // Азербайджанский
		"ba":     "ba",     // Башкирский
		"be":     "be",     // Белорусский
		"bg":     "bg",     // Болгарский
		"bn":     "bn",     // Бенгальский
		"bs":     "bs",     // Боснийский
		"ca":     "ca",     // Каталанский
		"ceb":    "ceb",    // Себуанский
		"cs":     "cs",     // Чешский
		"cv":     "cv",     // Чувашский
		"cy":     "cy",     // Валлийский
		"da":     "da",     // Датский
		"de":     "de",     // Немецкий
		"el":     "el",     // Греческий
		"emj":    "emj",    // Эмодзи
		"en":     "en",     // Английский
		"eo":     "eo",     // Эсперанто
		"es":     "es",     // Испанский
		"et":     "et",     // Эстонский
		"eu":     "eu",     // Баскский
		"fa":     "fa",     // Персидский
		"fi":     "fi",     // Финский
		"fr":     "fr",     // Французский
		"ga":     "ga",     // Ирландский
		"gd":     "gd",     // Шотландский (гэльский)
		"gl":     "gl",     // Галисийский
		"gu":     "gu",     // Гуджарати
		"he":     "he",     // Иврит
		"hi":     "hi",     // Хинди
		"hr":     "hr",     // Хорватский
		"ht":     "ht",     // Гаитянский
		"hu":     "hu",     // Венгерский
		"hy":     "hy",     // Армянский
		"id":     "id",     // Индонезийский
		"is":     "is",     // Исландский
		"it":     "it",     // Итальянский
		"ja":     "ja",     // Японский
		"jv":     "jv",     // Яванский
		"ka":     "ka",     // Грузинский
		"kazlat": "kazlat", // Казахский (латиница)
		"kk":     "kk",     // Казахский
		"km":     "km",     // Кхмерский
		"kn":     "kn",     // Каннада
		"ko":     "ko",     // Корейский
		"ky":     "ky",     // Киргизский
		"la":     "la",     // Латынь
		"lb":     "lb",     // Люксембургский
		"lo":     "lo",     // Лаосский
		"lt":     "lt",     // Литовский
		"lv":     "lv",     // Латышский
		"mg":     "mg",     // Малагасийский
		"mhr":    "mhr",    // Марийский
		"mi":     "mi",     // Маори
		"mk":     "mk",     // Македонский
		"ml":     "ml",     // Малаялам
		"mn":     "mn",     // Монгольский
		"mr":     "mr",     // Маратхи
		"mrj":    "mrj",    // Горномарийский
		"ms":     "ms",     // Малайский
		"mt":     "mt",     // Мальтийский
		"my":     "my",     // Бирманский
		"ne":     "ne",     // Непальский
		"nl":     "nl",     // Нидерландский
		"no":     "no",     // Норвежский
		"os":     "os",     // Осетинский
		"pa":     "pa",     // Панджаби
		"pap":    "pap",    // Папьяменто
		"pl":     "pl",     // Польский
		"pt":     "pt",     // Португальский
		// "pt":     "pt-BR",   // Португальский (бразильский)
		"ro":     "ro",     // Румынский
		"ru":     "ru",     // Русский
		"sah":    "sah",    // Якутский
		"si":     "si",     // Сингальский
		"sk":     "sk",     // Словацкий
		"sl":     "sl",     // Словенский
		"sq":     "sq",     // Албанский
		"sr":     "sr",     // Сербский
		"su":     "su",     // Сунданский
		"sv":     "sv",     // Шведский
		"sw":     "sw",     // Суахили
		"ta":     "ta",     // Тамильский
		"te":     "te",     // Телугу
		"tg":     "tg",     // Таджикский
		"th":     "th",     // Тайский
		"tl":     "tl",     // Тагальский
		"tr":     "tr",     // Турецкий
		"tt":     "tt",     // Татарский
		"udm":    "udm",    // Удмуртский
		"uk":     "uk",     // Украинский
		"ur":     "ur",     // Урду
		"uz":     "uz",     // Узбекский
		"uzbcyr": "uzbcyr", // Узбекский (кириллица)
		"vi":     "vi",     // Вьетнамский
		"xh":     "xh",     // Коса
		"yi":     "yi",     // Идиш
		"zh":     "zh",     // Китайский
		"zu":     "zu",     // Зулу
	}
)

type Translator struct {
	api      string
	apiKey   string
	folderID string
}

func New(apiKey string) *Translator {
	t := &Translator{
		api:    API,
		apiKey: apiKey,
	}

	return t
}

type request struct {
	SourceLanguageCode string   `json:"sourceLanguageCode"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
	Format             string   `json:"format"`
	Texts              []string `json:"texts"`
	FolderId           string   `json:"folderId"`
}

type responseJson struct {
	Translations []struct {
		Text                 string `json:"text"`
		DetectedLanguageCode string `json:"detectedLanguageCode"`
	} `json:"translations"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func ParseLang(l string) string {
	if strings.Contains(l, "-") {
		_l := strings.Split(l, "-")
		l = _l[0]
	}

	if l, ok := languages[l]; ok {
		return l
	}

	return ""
}

func (t *Translator) Translate(sourceLang, targetLang string, text string) (string, error) {

	skURL, _ := url.Parse(t.api)
	values := request{
		SourceLanguageCode: sourceLang,
		TargetLanguageCode: targetLang,
		Format:             "PLAIN_TEXT",
		Texts:              []string{text},
	}
	if t.folderID != "" {
		values.FolderId = t.folderID
	}

	body, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("failed marshal params: %w", err)
	}

	data := bytes.NewBuffer(body)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, skURL.String(), data)
	if err != nil {
		return "", err
	}
	if t.apiKey != "" {
		req.Header.Add("Authorization", "Api-Key "+t.apiKey)
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var res = responseJson{}
	err = json.Unmarshal(bBody, &res)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("yandex return status code `%v`, error: %s", response.StatusCode, res.ErrorMessage)
	}

	textResult := ""
	for _, t := range res.Translations {
		textResult += t.Text
	}

	return textResult, nil
}
