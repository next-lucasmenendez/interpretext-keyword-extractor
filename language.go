package gotagger

import (
	"os"
	"path"
	"errors"
	"bufio"
)

var supported = map[string][]string{
	"es": []string{ "a", "al", "algo", "algunas", "algunos", "ante", "antes", "como", "con", "contra", "cual", "cuando",
		"de", "del", "desde", "donde", "durante", "e", "el", "ella", "ellas", "ellos", "en", "entre", "era", "erais",
		"eran", "eras", "eres", "es", "esa", "esas", "ese", "eso", "esos", "esta", "estaba", "estabais", "estaban",
		"estabas", "estad", "estada", "estadas", "estado", "estados", "estamos", "estando", "estar", "estaremos",
		"estará", "estarán", "estarás", "estaré", "estaréis", "estaría", "estaríais", "estaríamos", "estarían",
		"estarías", "estas","este", "estemos", "esto", "estos", "estoy", "estuve", "estuviera", "estuvierais",
		"estuvieran", "estuvieras","estuvieron", "estuviese", "estuvieseis", "estuviesen", "estuvieses", "estuvimos",
		"estuviste", "estuvisteis", "estuviéramos", "estuviésemos", "estuvo", "está", "estábamos", "estáis", "están",
		"estás", "esté", "estéis", "estén", "estés", "fue", "fuera", "fuerais", "fueran", "fueras", "fueron", "fuese",
		"fueseis", "fuesen", "fueses","fui", "fuimos", "fuiste", "fuisteis", "fuéramos", "fuésemos", "ha", "habida",
		"habidas", "habido", "habidos", "habiendo", "habremos", "habrá", "habrán", "habrás", "habré", "habréis",
		"habría", "habríais", "habríamos", "habrían", "habrías", "habéis", "había", "habíais", "habíamos", "habían",
		"habías", "han", "has", "hasta", "hay", "haya", "hayamos", "hayan", "hayas", "hayáis", "he", "hemos", "hube",
		"hubiera", "hubierais", "hubieran", "hubieras", "hubieron", "hubiese", "hubieseis", "hubiesen", "hubieses",
		"hubimos", "hubiste", "hubisteis", "hubiéramos", "hubiésemos", "hubo", "la", "las", "le", "les", "lo", "los",
		"me", "mi", "mis", "mucho", "muchos", "muy", "más", "mí", "mía", "mías", "mío", "míos", "nada", "ni", "no",
		"nos", "nosotras", "nosotros", "nuestra", "nuestras", "nuestro", "nuestros", "o", "os", "otra", "otras", "otro",
		"otros", "para", "pero", "poco", "por", "porque", "que", "quien", "quienes", "qué", "se", "sea", "seamos",
		"sean", "seas", "seremos", "será", "serán", "serás", "seré", "seréis", "sería", "seríais", "seríamos", "serían",
		"serías", "seáis", "sido", "siendo", "sin", "sobre", "sois", "somos", "son", "soy", "su", "sus", "suya",
		"suyas", "suyo", "suyos", "sí", "también", "tanto", "te", "tendremos", "tendrá", "tendrán", "tendrás", "tendré",
		"tendréis", "tendría", "tendríais", "tendríamos", "tendrían", "tendrías", "tened", "tenemos", "tenga",
		"tengamos", "tengan", "tengas", "tengo", "tengáis", "tenida", "tenidas", "tenido", "tenidos", "teniendo",
		"tenéis", "tenía", "teníais", "teníamos", "tenían", "tenías", "ti","tiene", "tienen", "tienes", "todo", "todos",
		"tu", "tus", "tuve", "tuviera", "tuvierais", "tuvieran", "tuvieras", "tuvieron", "tuviese", "tuvieseis",
		"tuviesen", "tuvieses", "tuvimos", "tuviste", "tuvisteis", "tuviéramos", "tuviésemos", "tuvo", "tuya", "tuyas",
		"tuyo", "tuyos", "tú", "un", "una", "uno", "unos", "vosotras", "vosotros", "vuestra", "vuestras", "vuestro",
		"vuestros", "y", "ya", "yo", "él", "éramos",
	},
	"en": []string{ "a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are",
		"aren't", "as", "at", "be", "because", "been", "before", "being", "below", "between", "both", "but", "by",
		"can't", "cannot", "could", "couldn't", "did", "didn't", "do", "does", "doesn't", "doing", "don't", "down",
		"during", "each", "few", "for", "from", "further", "had", "hadn't", "has", "hasn't", "have", "haven't",
		"having", "he", "he'd", "he'll", "he's", "her", "here", "here's", "hers", "herself", "him", "himself", "his",
		"how", "how's", "i", "i'd", "i'll", "i'm", "i've", "if", "in", "into", "is", "isn't", "it", "it's", "its",
		"itself", "let's", "me", "more", "most", "mustn't", "my", "myself", "no", "nor", "not", "of", "off", "on",
		"once", "only", "or", "other", "ought", "our", "ours", "ourselves", "out", "over", "own", "same", "shan't",
		"she", "she'd", "she'll", "she's", "should", "shouldn't", "so", "some", "such", "than", "that", "that's", "the",
		"their", "theirs", "them", "themselves", "then", "there", "there's", "these", "they", "they'd", "they'll",
		"they're", "they've", "this", "those", "through", "to", "too", "under", "until", "up", "very", "was", "wasn't",
		"we", "we'd", "we'll", "we're", "we've", "were", "weren't", "what", "what's", "when", "when's", "where",
		"where's", "which", "while", "who", "who's", "whom", "why", "why's", "with", "won't", "would", "wouldn't",
		"you", "you'd", "you'll", "you're", "you've", "your", "yours", "yourself", "yourselves",
	},
}

type language struct {
	code string
	stopwords []string
}

func loadLanguage(code string) (language, error) {
	var err error
	var stopwords []string

	if env := os.Getenv("STOPWORDS"); env != "" {
		var f string = path.Join(env, code)

		var fd *os.File
		if fd, err = os.Open(f); err != nil {
			return language{}, err
		}
		defer fd.Close()

		var scanner *bufio.Scanner = bufio.NewScanner(fd)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			var w string = scanner.Text()
			if len(w) > 0 {
				stopwords = append(stopwords, w)
			}
		}
	} else {
		var ok bool
		if stopwords, ok = supported[code]; !ok {
			return language{}, errors.New("language not supported")
		}
	}

	return language{ code , stopwords}, nil
}