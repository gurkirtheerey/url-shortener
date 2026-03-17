package main

// Base62 encoding converts a numeric ID into a short, URL-safe string.
// It uses 62 characters (a-z, A-Z, 0-9), similar to how binary uses 2 digits
// and decimal uses 10. This gives us short codes like "b", "ba", "dnh" instead
// of raw numbers like 1, 62, 12345.
//
// Because every database ID is unique, every base62 code is guaranteed unique.
// No need to check for collisions or generate random strings.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func encodeBase62(id int64) string {
	if id == 0 {
		return string(charset[0])
	}

	var result []byte
	for id > 0 {
		// id%62 gives us the index into our charset (the "remainder")
		// id/62 shifts to the next digit, just like id/10 would in decimal
		// We prepend each character so the most significant digit comes first
		result = append([]byte{charset[id%62]}, result...)
		id /= 62
	}
	return string(result)
}
