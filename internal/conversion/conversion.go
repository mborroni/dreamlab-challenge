package conversion

import (
	"strconv"
	"strings"
)

func IPv4ToDecimal(ip string) (int64, error) {
	nodes := strings.Split(ip, ".")
	if len(nodes) != 4 {
		return 0, NotIPv4{}
	}
	var numbers []int64
	for _, node := range nodes {
		number, err := strconv.ParseInt(node, 10, 64)
		if err != nil {
			return 0, NotIPv4{}
		}
		numbers = append(numbers, number)
	}
	decimal := numbers[0] << 24
	decimal += numbers[1] << 16
	decimal += numbers[2] << 8
	decimal += numbers[3]
	return decimal, nil
}

func DecimalToIPv4(decimalIP int64) string {
	var output []string
	output = append(output, strconv.FormatInt((decimalIP>>24)&0xff, 10))
	output = append(output, strconv.FormatInt((decimalIP>>16)&0xff, 10))
	output = append(output, strconv.FormatInt((decimalIP>>8)&0xff, 10))
	output = append(output, strconv.FormatInt(decimalIP&0xff, 10))

	result := strings.Join(output, ".")
	return result
}

type NotIPv4 struct{}

func (NotIPv4) Error() string {
	return "not an IPv4 address"
}
