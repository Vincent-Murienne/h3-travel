package utils

func ValidateCardNumber(number string) bool {
	sum := 0
	alt := false
	n := len(number)

	for i := n - 1; i >= 0; i-- {
		c := number[i]
		if c < '0' || c > '9' {
			return false
		}
		digit := int(c - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}
	return sum%10 == 0
}
