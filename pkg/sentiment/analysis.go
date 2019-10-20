package sentiment

const (
	positive = "positive"
	negative = "negative"
)

type wordFrequency struct {
	word    string
	counter map[string]int
}

type classifier struct {
	dataset map[string][]string
	words   map[string]wordFrequency
}

func newClassifier() *classifier {
	c := new(classifier)
	c.dataset = map[string][]string{
		positive: {},
		negative: {},
	}
	c.words = map[string]wordFrequency{}
	return c
}

func (c *classifier) train(dataset map[string]string) {
	for sentence, class := range dataset {
		c.addSentence(sentence, class)
		words := tokenize(sentence)
		for _, w := range words {
			c.addWord(w, class)
		}
	}
}

func (c classifier) classify(sentence string) map[string]float64 {
	words := tokenize(sentence)
	posProb := c.probability(words, positive)
	negProb := c.probability(words, negative)
	return map[string]float64{
		positive: posProb,
		negative: negProb,
	}
}

func (c *classifier) addSentence(sentence, class string) {
	c.dataset[class] = append(c.dataset[class], sentence)
}

func (c *classifier) addWord(word, class string) {
	wf, ok := c.words[word]
	if !ok {
		wf = wordFrequency{word: word, counter: map[string]int{
			positive: 0,
			negative: 0,
		}}
	}
	wf.counter[class]++
	c.words[word] = wf
}

func (c classifier) priorProb(class string) float64 {
	return float64(len(c.dataset[class])) / float64(len(c.dataset[positive])+len(c.dataset[negative]))
}

func (c classifier) totalWordCount(class string) int {
	posCount := 0
	negCount := 0
	for _, wf := range c.words {
		posCount += wf.counter[positive]
		negCount += wf.counter[negative]
	}
	if class == positive {
		return posCount
	} else if class == negative {
		return negCount
	} else {
		return posCount + negCount
	}
}

func (c classifier) totalDistinctWordCount() int {
	posCount := 0
	negCount := 0
	for _, wf := range c.words {
		posCount += zeroOneTransform(wf.counter[positive])
		negCount += zeroOneTransform(wf.counter[negative])
	}
	return posCount + negCount
}

func (c classifier) probability(words []string, class string) float64 {
	prob := c.priorProb(class)
	for _, w := range words {
		count := 0
		if wf, ok := c.words[w]; ok {
			count = wf.counter[class]
		}
		prob *= (float64((count + 1)) / float64((c.totalWordCount(class) + c.totalDistinctWordCount())))
	}
	for _, w := range words {
		count := 0
		if wf, ok := c.words[w]; ok {
			count += (wf.counter[positive] + wf.counter[negative])
		}
		prob /= (float64((count + 1)) / float64((c.totalWordCount("") + c.totalDistinctWordCount())))
	}
	return prob
}
