package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

/*
	path = "./data.txt"						// output file path
	item = 1000								// total number of metrics / lines in output file
	metric = "metricName"					// base metricName. output file has only unique metricName with a number suffix
	types = []string{
		"counter",
		"gauge",
		"histogram",
		"summary",
	}
	labels = []string{						// each metric will have from 1 to 4 sets of labels
		"label1 value1",
		"label2 value2",
		"label3 value3",
		"label4 value4",
	}
	delimeter = ","							// separate name, type, labels and metric value
	space = " "								// separate a set of label values or metric values
	valueBound = 5000						// metric values are [0, valueBound)
	bounds = []float32{0.01, 0.5, 0.99}		// fixed quantile/buckets
*/

// generateData writes metrics in the following format to a text file:
// 		 name, type, label1 labelvalue1 , value1 value2 value3 value4 value5
// gauge and counter has only one value
func generateData() {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	for i := 0; i < item; i++ {

		mName := metric + strconv.Itoa(i)
		mType := types[rand.Intn(4)]
		labelSize := rand.Intn(len(labels)) + 1
		b := &strings.Builder{}
		writeNameTypeLabel(mName, mType, labelSize, b)

		// gauge or counter
		if mType == types[0] || mType == types[1] {
			b.WriteString(strconv.Itoa(rand.Intn(valueBound)))
			// histogram
		} else if mType == types[2] {

			count := 0
			buckets := make([]int, 3, 3)
			for i := range bounds {
				n := rand.Intn(valueBound)
				buckets[i] = n
				count += n
			}
			b.WriteString(strconv.Itoa(rand.Intn(valueBound))) // sum
			b.WriteString(space)
			b.WriteString(strconv.Itoa(count)) // count
			b.WriteString(space)
			for _, val := range buckets {
				b.WriteString(strconv.Itoa(val)) // individual bucket
				b.WriteString(space)
			}
			// summary
		} else {
			b.WriteString(strconv.Itoa(rand.Intn(valueBound))) // sum
			b.WriteString(space)
			b.WriteString(strconv.Itoa(rand.Intn(valueBound))) // count
			b.WriteString(space)
			for range bounds {
				b.WriteString(fmt.Sprintf("%f", rand.Float32())) // individual quantile
				b.WriteString(space)
			}
		}
		b.WriteString("\n")
		f.WriteString(b.String())
	}
}

func writeNameTypeLabel(mName, mType string, labelSize int, b *strings.Builder) {
	b.WriteString(mName)
	b.WriteString(delimeter)
	b.WriteString(mType)
	b.WriteString(delimeter)
	for i := 0; i < labelSize; i++ {
		b.WriteString(labels[i])
		b.WriteString(space)
	}
	b.WriteString(delimeter)
}
